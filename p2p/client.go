package p2p

import (
	"context"
	"fmt"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"

	"github.com/centrifuge/go-centrifuge/documents"

	"github.com/centrifuge/go-centrifuge/p2p/common"

	"github.com/golang/protobuf/proto"

	"github.com/centrifuge/go-centrifuge/contextutil"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/errors"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/p2p"
	"github.com/centrifuge/go-centrifuge/centerrors"
	"github.com/centrifuge/go-centrifuge/code"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/version"
	libp2pPeer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

func (s *peer) SendAnchoredDocument(ctx context.Context, receiverID identity.CentID, in *p2ppb.AnchorDocumentRequest) (*p2ppb.AnchorDocumentResponse, error) {
	nc, err := s.config.GetConfig()
	if err != nil {
		return nil, err
	}

	peerCtx, cancel := context.WithTimeout(ctx, nc.GetP2PConnectionTimeout())
	defer cancel()

	tc, err := s.config.GetAccount(receiverID[:])
	if err == nil {
		// this is a local account
		h := s.handlerCreator()
		// the following context has to be different from the parent context since its initiating a local peer call
		localCtx, err := contextutil.New(peerCtx, tc)
		if err != nil {
			return nil, err
		}
		return h.SendAnchoredDocument(localCtx, in, receiverID[:])
	}

	id, err := s.idService.LookupIdentityForID(receiverID)
	if err != nil {
		return nil, err
	}

	// this is a remote account
	pid, err := s.getPeerID(id)
	if err != nil {
		return nil, err
	}

	envelope, err := p2pcommon.PrepareP2PEnvelope(ctx, nc.GetNetworkID(), p2pcommon.MessageTypeSendAnchoredDoc, in)
	if err != nil {
		return nil, err
	}

	recv, err := s.mes.SendMessage(
		ctx, pid,
		envelope,
		p2pcommon.ProtocolForCID(receiverID))
	if err != nil {
		return nil, err
	}

	recvEnvelope, err := p2pcommon.ResolveDataEnvelope(recv)
	if err != nil {
		return nil, err
	}

	// handle client error
	if p2pcommon.MessageTypeError.Equals(recvEnvelope.Header.Type) {
		return nil, convertClientError(recvEnvelope)
	}

	if !p2pcommon.MessageTypeSendAnchoredDocRep.Equals(recvEnvelope.Header.Type) {
		return nil, errors.New("the received send anchored document response is incorrect")
	}

	r := new(p2ppb.AnchorDocumentResponse)
	err = proto.Unmarshal(recvEnvelope.Body, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// OpenClient returns P2PServiceClient to contact the remote peer
func (s *peer) getPeerID(id identity.Identity) (libp2pPeer.ID, error) {
	lastB58Key, err := id.CurrentP2PKey()
	if err != nil {
		return "", errors.New("error fetching p2p key: %v", err)
	}
	target := fmt.Sprintf("/ipfs/%s", lastB58Key)
	log.Info("Opening connection to: %s", target)
	ipfsAddr, err := ma.NewMultiaddr(target)
	if err != nil {
		return "", err
	}

	pid, err := ipfsAddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		return "", err
	}

	peerID, err := libp2pPeer.IDB58Decode(pid)
	if err != nil {
		return "", err
	}

	if !s.disablePeerStore {
		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", pid))
		targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
		// We have a peer ID and a targetAddr so we add it to the peer store
		// so LibP2P knows how to contact it
		s.host.Peerstore().AddAddr(peerID, targetAddr, pstore.PermanentAddrTTL)
	}

	return peerID, nil
}

// getSignatureForDocument requests the target node to sign the document
func (s *peer) getSignatureForDocument(ctx context.Context, model documents.Model, cid identity.CentID) (*p2ppb.SignatureResponse, error) {
	nc, err := s.config.GetConfig()
	if err != nil {
		return nil, err
	}

	var resp *p2ppb.SignatureResponse
	var header *p2ppb.Header
	tc, err := s.config.GetAccount(cid[:])

	cd, err := model.PackCoreDocument()
	if err != nil {
		return nil, err
	}

	if err == nil {
		// this is a local account
		h := s.handlerCreator()
		// create a context with receiving account value
		localPeerCtx, err := contextutil.New(ctx, tc)
		if err != nil {
			return nil, err
		}

		resp, err = h.RequestDocumentSignature(localPeerCtx, &p2ppb.SignatureRequest{Document: &cd})
		if err != nil {
			return nil, err
		}
		header = &p2ppb.Header{NodeVersion: version.GetVersion().String()}
	} else {
		// this is a remote account
		id, err := s.idService.LookupIdentityForID(cid)
		if err != nil {
			return nil, err
		}

		receiverPeer, err := s.getPeerID(id)
		if err != nil {
			return nil, err
		}

		envelope, err := p2pcommon.PrepareP2PEnvelope(ctx, nc.GetNetworkID(), p2pcommon.MessageTypeRequestSignature, &p2ppb.SignatureRequest{Document: &cd})
		if err != nil {
			return nil, err
		}
		log.Infof("Requesting signature from %s\n", receiverPeer)
		recv, err := s.mes.SendMessage(ctx, receiverPeer, envelope, p2pcommon.ProtocolForCID(cid))
		if err != nil {
			return nil, err
		}
		recvEnvelope, err := p2pcommon.ResolveDataEnvelope(recv)
		if err != nil {
			return nil, err
		}
		// handle client error
		if p2pcommon.MessageTypeError.Equals(recvEnvelope.Header.Type) {
			return nil, convertClientError(recvEnvelope)
		}
		if !p2pcommon.MessageTypeRequestSignatureRep.Equals(recvEnvelope.Header.Type) {
			return nil, errors.New("the received request signature response is incorrect")
		}
		resp = new(p2ppb.SignatureResponse)
		err = proto.Unmarshal(recvEnvelope.Body, resp)
		if err != nil {
			return nil, err
		}
		header = recvEnvelope.Header
	}

	err = validateSignatureResp(s.idService, cid, cd.SigningRoot, header, resp)
	if err != nil {
		return nil, err
	}

	log.Infof("Signature successfully received from %s\n", cid)
	return resp, nil
}

type signatureResponseWrap struct {
	resp *p2ppb.SignatureResponse
	err  error
}

func (s *peer) getSignatureAsync(ctx context.Context, model documents.Model, id identity.CentID, out chan<- signatureResponseWrap) {
	resp, err := s.getSignatureForDocument(ctx, model, id)
	out <- signatureResponseWrap{
		resp: resp,
		err:  err,
	}
}

// GetSignaturesForDocument requests peer nodes for the signature, verifies them, and returns those signatures.
func (s *peer) GetSignaturesForDocument(ctx context.Context, model documents.Model) (signatures []*coredocumentpb.Signature, err error) {
	in := make(chan signatureResponseWrap)
	defer close(in)

	nc, err := s.config.GetConfig()
	if err != nil {
		return nil, err
	}

	self, err := contextutil.Self(ctx)
	if err != nil {
		return nil, errors.New("failed to get self ID")
	}

	cs, err := model.GetCollaborators(self.ID)
	if err != nil {
		return nil, errors.New("failed to get external collaborators")
	}

	var count int
	peerCtx, _ := context.WithTimeout(ctx, nc.GetP2PConnectionTimeout())
	for _, c := range cs {
		count++
		go s.getSignatureAsync(peerCtx, model, c, in)
	}

	var responses []signatureResponseWrap
	for i := 0; i < count; i++ {
		responses = append(responses, <-in)
	}

	for _, resp := range responses {
		if resp.err != nil {
			// this error is ignored since we would still anchor the document
			log.Error(resp.err)
			continue
		}

		signatures = append(signatures, resp.resp.Signature)
	}

	return signatures, nil
}

func convertClientError(recv *p2ppb.Envelope) error {
	resp := new(errorspb.Error)
	err := proto.Unmarshal(recv.Body, resp)
	if err != nil {
		return err
	}
	return errors.New(resp.Message)
}

func validateSignatureResp(
	identityService identity.Service,
	receiver identity.CentID,
	signingRoot []byte,
	header *p2ppb.Header,
	resp *p2ppb.SignatureResponse) error {

	compatible := version.CheckVersion(header.NodeVersion)
	if !compatible {
		return version.IncompatibleVersionError(header.NodeVersion)
	}

	err := identity.ValidateCentrifugeIDBytes(resp.Signature.EntityId, receiver)
	if err != nil {
		return centerrors.New(code.AuthenticationFailed, err.Error())
	}

	err = identityService.ValidateSignature(resp.Signature, signingRoot)
	if err != nil {
		return centerrors.New(code.AuthenticationFailed, "signature invalid")
	}
	return nil
}
