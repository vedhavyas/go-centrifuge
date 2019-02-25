// +build unit

package p2p

import (
	"context"
	"github.com/centrifuge/go-centrifuge/documents"
	"testing"

	"github.com/centrifuge/go-centrifuge/testingutils/documents"

	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/p2p/common"
	"github.com/centrifuge/go-centrifuge/testingutils/config"
	"github.com/golang/protobuf/proto"

	"github.com/centrifuge/go-centrifuge/protobufs/gen/go/protocol"
	libp2pPeer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-protocol"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/p2p"
	"github.com/centrifuge/go-centrifuge/identity"
	testingcommons "github.com/centrifuge/go-centrifuge/testingutils/commons"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/go-centrifuge/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessenger struct {
	mock.Mock
}

func (mm *MockMessenger) Init(id ...protocol.ID) {
	mm.Called(id)
}

func (mm *MockMessenger) SendMessage(ctx context.Context, p libp2pPeer.ID, pmes *protocolpb.P2PEnvelope, protoc protocol.ID) (*protocolpb.P2PEnvelope, error) {
	args := mm.Called(ctx, p, pmes, protoc)
	resp, _ := args.Get(0).(*protocolpb.P2PEnvelope)
	return resp, args.Error(1)
}

func TestGetSignatureForDocument_fail_connect(t *testing.T) {
	centrifugeId, err := identity.ToCentID(utils.RandomSlice(identity.CentIDLength))
	idService := getIDMocks(centrifugeId)
	m := &MockMessenger{}
	testClient := &peer{config: cfg, idService: idService, mes: m, disablePeerStore: true}

	coreDoc, err := testingdocuments.GenerateCoreDocumentModel()
	assert.NoError(t, err)
	c, err := cfg.GetConfig()
	assert.NoError(t, err)
	c = updateKeys(c)
	ctx := testingconfig.CreateAccountContext(t, c)
	assert.Nil(t, err, "centrifugeId not initialized correctly ")
	mm := documents.Model{}
	mm.

	_, err = p2pcommon.PrepareP2PEnvelope(ctx, c.GetNetworkID(), p2pcommon.MessageTypeRequestSignature, &p2ppb.SignatureRequest{Document: &coreDoc.Document	})
	assert.NoError(t, err, "signature request could not be created")

	m.On("SendMessage", ctx, mock.Anything, mock.Anything, p2pcommon.ProtocolForCID(centrifugeId)).Return(nil, errors.New("some error"))
	resp, err := testClient.getSignatureForDocument(ctx, coreDoc, centrifugeId)
	m.AssertExpectations(t)
	assert.Error(t, err, "must fail")
	assert.Nil(t, resp, "must be nil")
}

func TestGetSignatureForDocument_fail_version_check(t *testing.T) {
	centrifugeId, err := identity.ToCentID(utils.RandomSlice(identity.CentIDLength))
	assert.Nil(t, err, "centrifugeId not initialized correctly ")

	idService := getIDMocks(centrifugeId)
	m := &MockMessenger{}
	testClient := &peer{config: cfg, idService: idService, mes: m, disablePeerStore: true}
	coreDoc, err := testingdocuments.GenerateCoreDocumentModel()
	assert.NoError(t, err)
	c, err := cfg.GetConfig()
	assert.NoError(t, err)
	c = updateKeys(c)
	ctx := testingconfig.CreateAccountContext(t, c)

	_, err = p2pcommon.PrepareP2PEnvelope(ctx, c.GetNetworkID(), p2pcommon.MessageTypeRequestSignature, &p2ppb.SignatureRequest{Document: &coreDoc.Document})
	assert.NoError(t, err, "signature request could not be created")

	m.On("SendMessage", ctx, mock.Anything, mock.Anything, p2pcommon.ProtocolForCID(centrifugeId)).Return(testClient.createSignatureResp("", nil), nil)
	resp, err := testClient.getSignatureForDocument(ctx, coreDoc.Document, centrifugeId)
	m.AssertExpectations(t)
	assert.Error(t, err, "must fail")
	assert.Contains(t, err.Error(), "Incompatible version")
	assert.Nil(t, resp, "must be nil")
}

func TestGetSignatureForDocument_fail_centrifugeId(t *testing.T) {
	centrifugeId, err := identity.ToCentID(utils.RandomSlice(identity.CentIDLength))
	idService := getIDMocks(centrifugeId)
	m := &MockMessenger{}
	testClient := &peer{config: cfg, idService: idService, mes: m, disablePeerStore: true}
	coreDoc, err := testingdocuments.GenerateCoreDocumentModel()
	c, err := cfg.GetConfig()
	assert.NoError(t, err)
	c = updateKeys(c)
	ctx := testingconfig.CreateAccountContext(t, c)

	assert.Nil(t, err, "centrifugeId not initialized correctly ")

	_, err = p2pcommon.PrepareP2PEnvelope(ctx, c.GetNetworkID(), p2pcommon.MessageTypeRequestSignature, &p2ppb.SignatureRequest{Document: &coreDoc.Document})
	assert.NoError(t, err, "signature request could not be created")

	randomBytes := utils.RandomSlice(identity.CentIDLength)
	signature := &coredocumentpb.Signature{EntityId: randomBytes, PublicKey: utils.RandomSlice(32)}
	m.On("SendMessage", ctx, mock.Anything, mock.Anything, p2pcommon.ProtocolForCID(centrifugeId)).Return(testClient.createSignatureResp(version.GetVersion().String(), signature), nil)

	resp, err := testClient.getSignatureForDocument(ctx, coreDoc.Document, centrifugeId)

	m.AssertExpectations(t)
	assert.Nil(t, resp, "must be nil")
	assert.Error(t, err, "must not be nil")
	assert.Contains(t, err.Error(), "[5]provided bytes doesn't match centID")

}

func getIDMocks(centrifugeId identity.CentID) *testingcommons.MockIDService {
	idService := &testingcommons.MockIDService{}
	id := &testingcommons.MockID{}
	id.On("CurrentP2PKey").Return("5dsgvJGnvAfiR3K6HCBc4hcokSfmjj", nil)
	idService.On("LookupIdentityForID", centrifugeId).Return(id, nil)
	return idService
}

func (s *peer) createSignatureResp(centNodeVer string, signature *coredocumentpb.Signature) *protocolpb.P2PEnvelope {
	req, err := proto.Marshal(&p2ppb.SignatureResponse{Signature: signature})
	if err != nil {
		return nil
	}

	dataReq := &p2ppb.Envelope{
		Header: &p2ppb.Header{
			NodeVersion: centNodeVer,
			Type:        p2pcommon.MessageTypeRequestSignatureRep.String(),
		},
		Body: req,
	}

	reqB, err := proto.Marshal(dataReq)
	if err != nil {
		return nil
	}

	return &protocolpb.P2PEnvelope{Body: reqB}
}
