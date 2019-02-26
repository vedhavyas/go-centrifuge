package documents

import (
	"context"
	"time"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/p2p"
	"github.com/centrifuge/go-centrifuge/anchors"
	"github.com/centrifuge/go-centrifuge/contextutil"
	"github.com/centrifuge/go-centrifuge/crypto/secp256k1"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/utils"
)

// Config defines required methods required for the documents package.
type Config interface {
	GetNetworkID() uint32
	GetIdentityID() ([]byte, error)
	GetP2PConnectionTimeout() time.Duration
}

// Client defines methods that can be implemented by any type handling p2p communications.
type Client interface {

	// GetSignaturesForDocument gets the signatures for Document
	GetSignaturesForDocument(ctx context.Context, model Model) ([]*coredocumentpb.Signature, error)

	// after all signatures are collected the sender sends the Document including the signatures
	SendAnchoredDocument(ctx context.Context, receiverID identity.CentID, in *p2ppb.AnchorDocumentRequest) (*p2ppb.AnchorDocumentResponse, error)
}

// defaultProcessor implements AnchorProcessor interface
type defaultProcessor struct {
	identityService  identity.Service
	p2pClient        Client
	anchorRepository anchors.AnchorRepository
	config           Config
}

// DefaultProcessor returns the default implementation of CoreDocument AnchorProcessor
func DefaultProcessor(idService identity.Service, p2pClient Client, repository anchors.AnchorRepository, config Config) AnchorProcessor {
	return defaultProcessor{
		identityService:  idService,
		p2pClient:        p2pClient,
		anchorRepository: repository,
		config:           config,
	}
}

// Send sends the given defaultProcessor to the given recipient on the P2P layer
func (dp defaultProcessor) Send(ctx context.Context, cd coredocumentpb.CoreDocument, id identity.CentID) (err error) {
	log.Infof("sending Document %x to recipient %x", cd.DocumentIdentifier, id)
	c, cancel := context.WithTimeout(ctx, dp.config.GetP2PConnectionTimeout())
	defer cancel()

	resp, err := dp.p2pClient.SendAnchoredDocument(c, id, &p2ppb.AnchorDocumentRequest{Document: &cd})
	if err != nil || !resp.Accepted {
		return errors.New("failed to send Document to the node: %v", err)
	}

	log.Infof("Sent Document to %x\n", id)
	return nil
}

// PrepareForSignatureRequests gets the core Document from the model, and adds the node's own signature
func (dp defaultProcessor) PrepareForSignatureRequests(ctx context.Context, model Model) error {
	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	_, err = model.DataRoot()
	if err != nil {
		return err
	}

	// calculate the signing root
	sr, err := model.SigningRoot()
	if err != nil {
		return errors.New("failed to calculate signing root: %v", err)
	}

	model.AppendSignatures(identity.Sign(self, identity.KeyPurposeSigning, sr))
	return nil
}

// RequestSignatures gets the core Document from the model, validates pre signature requirements,
// collects signatures, and validates the signatures,
func (dp defaultProcessor) RequestSignatures(ctx context.Context, model Model) error {
	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	idKeys, ok := self.Keys[identity.KeyPurposeSigning]
	if !ok {
		return errors.New("missing keys for signing")
	}

	psv := PreSignatureRequestValidator(self.ID[:], idKeys.PrivateKey, idKeys.PublicKey)
	err = psv.Validate(nil, model)
	if err != nil {
		return errors.New("failed to validate model for signature request: %v", err)
	}

	sigs, err := dp.p2pClient.GetSignaturesForDocument(ctx, model)
	if err != nil {
		return errors.New("failed to collect signatures from the collaborators: %v", err)
	}

	model.AppendSignatures(signs...)
	return nil
}

// PrepareForAnchoring validates the signatures and generates the Document root
func (dp defaultProcessor) PrepareForAnchoring(model Model) error {
	psv := PostSignatureRequestValidator(dp.identityService)
	err := psv.Validate(nil, model)
	if err != nil {
		return errors.New("failed to validate signatures: %v", err)
	}

	return nil
}

// AnchorDocument validates the model, and anchors the Document
func (dp defaultProcessor) AnchorDocument(ctx context.Context, model Model) error {
	pav := PreAnchorValidator(dp.identityService)
	err := pav.Validate(nil, model)
	if err != nil {
		return errors.New("pre anchor validation failed: %v", err)
	}

	dr, err := model.DocumentRoot()
	if err != nil {
		return errors.New("failed to get Document root: %v", err)
	}

	rootHash, err := anchors.ToDocumentRoot(dr)
	if err != nil {
		return errors.New("failed to get Document root: %v", err)
	}

	anchorID, err := anchors.ToAnchorID(model.CurrentVersion())
	if err != nil {
		return errors.New("failed to get anchor ID: %v", err)
	}

	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	// generate message authentication code for the anchor call
	mac, err := secp256k1.SignEthereum(anchors.GenerateCommitHash(anchorID, self.ID, rootHash), self.Keys[identity.KeyPurposeEthMsgAuth].PrivateKey)
	if err != nil {
		return errors.New("failed to generate ethereum MAC: %v", err)
	}

	log.Infof("Anchoring Document with identifiers: [Document: %#x, current: %#x], rootHash: %#x", model.ID(), model.CurrentVersion(), dr)
	confirmations, err := dp.anchorRepository.CommitAnchor(ctx, anchorID, rootHash, self.ID, [][anchors.DocumentProofLength]byte{utils.RandomByte32()}, mac)
	if err != nil {
		return errors.New("failed to commit anchor: %v", err)
	}

	<-confirmations
	log.Infof("Anchored Document with identifiers: [Document: %#x, current: %#x], rootHash: %#x", model.ID(), model.CurrentVersion(), dr)
	return nil
}

// SendDocument does post anchor validations and sends the Document to collaborators
func (dp defaultProcessor) SendDocument(ctx context.Context, model Model) error {
	av := PostAnchoredValidator(dp.identityService, dp.anchorRepository)
	err := av.Validate(nil, model)
	if err != nil {
		return errors.New("post anchor validations failed: %v", err)
	}

	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	cs, err := model.GetCollaborators(self.ID)
	if err != nil {
		return errors.New("get external collaborators failed: %v", err)
	}

	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core Document: %v", err)
	}

	for _, c := range cs {
		erri := dp.Send(ctx, cd, c)
		if erri != nil {
			err = errors.AppendError(err, erri)
		}
	}

	return err
}
