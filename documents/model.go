package documents

import (
	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/storage"
	"github.com/centrifuge/precise-proofs/proofs/proto"
	"github.com/ethereum/go-ethereum/common"
)

// Model is an interface to abstract away model specificness like invoice or purchaseOrder
// The interface can cast into the type specified by the model if required
// It should only handle protocol-level Document actions
type Model interface {
	storage.Model

	// ID returns the Document identifier
	ID() []byte

	// CurrentVersion returns the current version identifier of the Document
	CurrentVersion() []byte

	// PreviousVersion returns the previous version identifier of the Document
	PreviousVersion() []byte

	// NextVersion returns the next version identifier of the Document.
	NextVersion() []byte

	// PackCoreDocument packs the implementing Document into a core Document
	// Should only be called when the Document is about to be put on wire.
	PackCoreDocument() (coredocumentpb.CoreDocument, error)

	// UnpackCoreDocument takes a core Document protobuf and loads the data into the model.
	UnpackCoreDocument(cd coredocumentpb.CoreDocument) error

	// DocumentType returns the type of the Document
	DocumentType() string

	// CalculateDataRoot calculates the data root of the model.
	CalculateDataRoot() ([]byte, error)

	// CalculateSigningRoot calculates the signing root of the model.
	CalculateSigningRoot() ([]byte, error)

	// CalculateDocumentRoot returns the Document root of the model.
	CalculateDocumentRoot() ([]byte, error)

	// PreviousDocumentRoot returns the Document root of the previous version.
	PreviousDocumentRoot() []byte

	// AppendSignatures appends the signatures to the model.
	AppendSignatures(signatures ...*coredocumentpb.Signature)

	// Signatures returns a copy of the signatures on the Document
	Signatures() []coredocumentpb.Signature

	// CreateProofs creates precise-proofs for given fields
	CreateProofs(fields []string) (proofs []*proofspb.Proof, err error)

	// CreateNFTProofs creates NFT proofs for minting.
	CreateNFTProofs(
		account identity.DID,
		registry common.Address,
		tokenID []byte,
		nftUniqueProof, readAccessProof bool) (proofs []*proofspb.Proof, err error)

	// IsNFTMinted checks if there is any NFT minted for the registry given
	IsNFTMinted(tr TokenRegistry, registry common.Address) bool

	// AddNFT adds an NFT to the Document.
	// Note: The Document should be anchored after successfully adding the NFT.
	AddNFT(grantReadAccess bool, registry common.Address, tokenID []byte) error

	// GetCollaborators returns the collaborators of this Document.
	// filter ids should not be returned
	GetCollaborators(filterIDs ...identity.DID) ([]identity.DID, error)

	// AccountCanRead returns true if the account can read the document
	AccountCanRead(account identity.DID) bool

	// NFTOwnerCanRead returns error if the NFT cannot read the document.
	NFTOwnerCanRead(tokenRegistry TokenRegistry, registry common.Address, tokenID []byte, account identity.DID) error

	// ATOwnerCanRead returns error if the NFT cannot read the document.
	ATOwnerCanRead(tokenID, docID []byte, account identity.DID) (err error)
}

// TokenRegistry defines NFT related functions.
type TokenRegistry interface {
	// OwnerOf to retrieve owner of the tokenID
	OwnerOf(registry common.Address, tokenID []byte) (common.Address, error)
}
