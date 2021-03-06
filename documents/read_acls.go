package documents

import (
	"bytes"
	"context"
	"fmt"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/go-centrifuge/contextutil"
	"github.com/centrifuge/go-centrifuge/crypto"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/protobufs/gen/go/document"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/precise-proofs/proofs/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// initReadRules initiates the read rules for a given CoreDocumentModel.
// Collaborators are given Read_Sign action.
// if the rules are created already, this is a no-op.
// if collaborators are empty, it is a no-op
func (cd *CoreDocument) initReadRules(collaborators []identity.DID) {
	if len(cd.Document.Roles) > 0 && len(cd.Document.ReadRules) > 0 {
		return
	}

	if len(collaborators) < 1 {
		return
	}

	cd.addCollaboratorsToReadSignRules(collaborators)
}

// findRole calls OnRole for every role that matches the actions passed in
func findRole(cd coredocumentpb.CoreDocument, onRole func(rridx, ridx int, role *coredocumentpb.Role) bool, actions ...coredocumentpb.Action) bool {
	am := make(map[int32]struct{})
	for _, a := range actions {
		am[int32(a)] = struct{}{}
	}

	for i, rule := range cd.ReadRules {
		if _, ok := am[int32(rule.Action)]; !ok {
			continue
		}

		for j, rk := range rule.Roles {
			role, err := getRole(rk, cd.Roles)
			if err != nil {
				// seems like roles and rules are not in sync
				// skip to next one
				continue
			}

			if onRole(i, j, role) {
				return true
			}

		}
	}

	return false
}

// GetExternalCollaborators returns collaborators of a Document without the own centID.
func (cd *CoreDocument) GetExternalCollaborators(self identity.DID) ([][]byte, error) {
	var cs [][]byte
	for _, c := range cd.Document.Collaborators {
		id := identity.NewDIDFromBytes(c)
		if !self.Equal(id) {
			cs = append(cs, c)
		}
	}

	return cs, nil
}

// NFTOwnerCanRead checks if the nft owner/account can read the Document
func (cd *CoreDocument) NFTOwnerCanRead(tokenRegistry TokenRegistry, registry common.Address, tokenID []byte, account identity.DID) error {
	// check if the account can read the doc
	if cd.AccountCanRead(account) {
		return nil
	}

	// check if the nft is present in read rules
	found := findRole(cd.Document, func(_, _ int, role *coredocumentpb.Role) bool {
		_, found := isNFTInRole(role, registry, tokenID)
		return found
	}, coredocumentpb.Action_ACTION_READ)

	if !found {
		return errors.New("nft not found in the Document")
	}

	// get the owner of the NFT
	owner, err := tokenRegistry.OwnerOf(registry, tokenID)
	if err != nil {
		return errors.New("failed to get NFT owner: %v", err)
	}

	// TODO(ved): this will always fail until we roll out identity v2 with CentID type as common.Address
	if !bytes.Equal(owner.Bytes(), account[:]) {
		return errors.New("account (%v) not owner of the NFT", account.String())
	}

	return nil
}

// AccountCanRead validate if the core Document can be read by the account .
// Returns an error if not.
func (cd *CoreDocument) AccountCanRead(account identity.DID) bool {
	// loop though read rules, check all the rules
	return findRole(cd.Document, func(_, _ int, role *coredocumentpb.Role) bool {
		_, found := isAccountInRole(role, account)
		return found
	}, coredocumentpb.Action_ACTION_READ, coredocumentpb.Action_ACTION_READ_SIGN)
}

// addNFTToReadRules adds NFT token to the read rules of core Document.
func (cd *CoreDocument) addNFTToReadRules(registry common.Address, tokenID []byte) error {
	nft, err := ConstructNFT(registry, tokenID)
	if err != nil {
		return errors.New("failed to construct NFT: %v", err)
	}

	role := &coredocumentpb.Role{RoleKey: utils.RandomSlice(32)}
	role.Nfts = append(role.Nfts, nft)
	cd.addNewRule(role, coredocumentpb.Action_ACTION_READ)
	return cd.setSalts()
}

// AddNFT returns a new CoreDocument model with nft added to the Core Document. If grantReadAccess is true, the nft is added
// to the read rules.
func (cd *CoreDocument) AddNFT(grantReadAccess bool, registry common.Address, tokenID []byte) (*CoreDocument, error) {
	ncd, err := cd.PrepareNewVersion(nil, false)
	if err != nil {
		return nil, errors.New("failed to prepare new version: %v", err)
	}

	nft := getStoredNFT(ncd.Document.Nfts, registry.Bytes())
	if nft == nil {
		nft = new(coredocumentpb.NFT)
		// add 12 empty bytes
		eb := make([]byte, 12, 12)
		nft.RegistryId = append(registry.Bytes(), eb...)
		ncd.Document.Nfts = append(ncd.Document.Nfts, nft)
	}
	nft.TokenId = tokenID

	if grantReadAccess {
		err = ncd.addNFTToReadRules(registry, tokenID)
		if err != nil {
			return nil, err
		}
	}

	return ncd, ncd.setSalts()
}

// IsNFTMinted checks if the there is an NFT that is minted against this Document in the given registry.
func (cd *CoreDocument) IsNFTMinted(tokenRegistry TokenRegistry, registry common.Address) bool {
	nft := getStoredNFT(cd.Document.Nfts, registry.Bytes())
	if nft == nil {
		return false
	}

	_, err := tokenRegistry.OwnerOf(registry, nft.TokenId)
	return err == nil
}

// CreateNFTProofs generate proofs returns proofs for NFT minting.
func (cd *CoreDocument) CreateNFTProofs(
	docType string,
	account identity.DID,
	registry common.Address,
	tokenID []byte,
	nftUniqueProof, readAccessProof bool) (proofs []*proofspb.Proof, err error) {

	if len(cd.Document.DataRoot) != idSize {
		return nil, errors.New("data root is invalid")
	}

	var pfKeys []string
	if nftUniqueProof {
		pk, err := getNFTUniqueProofKey(cd.Document.Nfts, registry)
		if err != nil {
			return nil, err
		}

		pfKeys = append(pfKeys, pk)
	}

	if readAccessProof {
		pks, err := getReadAccessProofKeys(cd.Document, registry, tokenID)
		if err != nil {
			return nil, err
		}

		pfKeys = append(pfKeys, pks...)
	}

	signingRootProofHashes, err := cd.getSigningRootProofHashes()
	if err != nil {
		return nil, errors.New("failed to generate signing root proofs: %v", err)
	}

	cdTree, err := cd.documentTree(docType)
	if err != nil {
		return nil, errors.New("failed to generate core Document tree: %v", err)
	}

	proofs, missedProofs := generateProofs(cdTree, pfKeys, append([][]byte{cd.Document.DataRoot}, signingRootProofHashes...))
	if len(missedProofs) != 0 {
		return nil, errors.New("failed to create proofs for fields %v", missedProofs)
	}

	return proofs, nil
}

// ConstructNFT appends registry and tokenID to byte slice
func ConstructNFT(registry common.Address, tokenID []byte) ([]byte, error) {
	var nft []byte
	// first 20 bytes of registry
	nft = append(nft, registry.Bytes()...)

	// next 32 bytes of the tokenID
	nft = append(nft, tokenID...)

	if len(nft) != nftByteCount {
		return nil, errors.New("byte length mismatch")
	}

	return nft, nil
}

// isNFTInRole checks if the given nft(registry + token) is part of the core Document role.
// If found, returns the index of the nft in the role and true
func isNFTInRole(role *coredocumentpb.Role, registry common.Address, tokenID []byte) (nftIdx int, found bool) {
	enft, err := ConstructNFT(registry, tokenID)
	if err != nil {
		return nftIdx, false
	}

	for i, n := range role.Nfts {
		if bytes.Equal(n, enft) {
			return i, true
		}
	}

	return nftIdx, false
}

func getStoredNFT(nfts []*coredocumentpb.NFT, registry []byte) *coredocumentpb.NFT {
	for _, nft := range nfts {
		if bytes.Equal(nft.RegistryId[:20], registry) {
			return nft
		}
	}

	return nil
}

func getReadAccessProofKeys(cd coredocumentpb.CoreDocument, registry common.Address, tokenID []byte) (pks []string, err error) {
	var rridx int  // index of the read rules which contain the role
	var ridx int   // index of the role
	var nftIdx int // index of the NFT in the above role
	var rk []byte  // role key of the above role

	found := findRole(cd, func(i, j int, role *coredocumentpb.Role) bool {
		z, found := isNFTInRole(role, registry, tokenID)
		if found {
			rridx = i
			ridx = j
			rk = role.RoleKey
			nftIdx = z
		}

		return found
	}, coredocumentpb.Action_ACTION_READ)

	if !found {
		return nil, ErrNFTRoleMissing
	}

	return []string{
		fmt.Sprintf(CDTreePrefix+".read_rules[%d].roles[%d]", rridx, ridx),          // proof that a read rule exists with the nft role
		fmt.Sprintf(CDTreePrefix+".roles[%s].nfts[%d]", hexutil.Encode(rk), nftIdx), // proof that role with nft exists
		fmt.Sprintf(CDTreePrefix+".read_rules[%d].action", rridx),                   // proof that this read rule has read access
	}, nil
}

func getNFTUniqueProofKey(nfts []*coredocumentpb.NFT, registry common.Address) (pk string, err error) {
	nft := getStoredNFT(nfts, registry.Bytes())
	if nft == nil {
		return pk, errors.New("nft is missing from the Document")
	}

	key := hexutil.Encode(nft.RegistryId)
	return fmt.Sprintf(CDTreePrefix+".nfts[%s]", key), nil
}

func getRoleProofKey(roles []*coredocumentpb.Role, roleKey []byte, account identity.DID) (pk string, err error) {
	role, err := getRole(roleKey, roles)
	if err != nil {
		return pk, err
	}

	idx, found := isAccountInRole(role, account)
	if !found {
		return pk, ErrNFTRoleMissing
	}

	return fmt.Sprintf(CDTreePrefix+".roles[%s].collaborators[%d]", hexutil.Encode(role.RoleKey), idx), nil
}

// isAccountInRole returns the index of the collaborator and true if account is in the given role as collaborators.
func isAccountInRole(role *coredocumentpb.Role, account identity.DID) (idx int, found bool) {
	for i, id := range role.Collaborators {
		if bytes.Equal(id, account[:]) {
			return i, true
		}
	}

	return idx, false
}

func getRole(key []byte, roles []*coredocumentpb.Role) (*coredocumentpb.Role, error) {
	for _, role := range roles {
		if utils.IsSameByteSlice(role.RoleKey, key) {
			return role, nil
		}
	}

	return nil, errors.New("role %d not found", key)
}

// isATInRole checks if the given access token is part of the core document role.
func isATInRole(role *coredocumentpb.Role, tokenID []byte) (*coredocumentpb.AccessToken, error) {
	for _, a := range role.AccessTokens {
		if bytes.Equal(tokenID, a.Identifier) {
			return a, nil
		}
	}
	return nil, errors.New("access token not found")
}

// validateAT validates that given access token against its signature
func validateAT(publicKey []byte, token *coredocumentpb.AccessToken, requesterID []byte) error {
	// assemble token message from the token for validation
	reqID := identity.NewDIDFromByte(requesterID)
	granterID := identity.NewDIDFromByte(token.Granter)
	tm, err := assembleTokenMessage(token.Identifier, granterID, reqID, token.RoleIdentifier, token.DocumentIdentifier)
	if err != nil {
		return err
	}
	validated := crypto.VerifyMessage(publicKey, tm, token.Signature, crypto.CurveSecp256K1)
	if !validated {
		return errors.New("access token is invalid")
	}
	return nil
}

// ATOwnerCanRead checks that the owner AT can read the document requested
func (cd *CoreDocument) ATOwnerCanRead(tokenID, docID []byte, account identity.DID) (err error) {
	// check if the access token is present in read rules of the document indicated in the AT request
	var at *coredocumentpb.AccessToken
	findRole(cd.Document, func(_, _ int, role *coredocumentpb.Role) bool {
		at, err = isATInRole(role, tokenID)
		if err != nil {
			return false
		}

		return true
	}, coredocumentpb.Action_ACTION_READ)

	if err != nil {
		return err
	}

	// check if the requested document is the document indicated in the access token
	if !bytes.Equal(at.DocumentIdentifier, docID) {
		return errors.New("the document requested does not match the document to which the access token grants access")
	}
	// validate the access token
	// TODO: fetch public key from Ethereum chain
	return validateAT(at.Key, at, account[:])
}

// AddAccessToken adds the AccessToken to the read rules of the document
func (cd *CoreDocument) AddAccessToken(ctx context.Context, payload documentpb.AccessTokenParams) (*CoreDocument, error) {
	ncd, err := cd.PrepareNewVersion(nil, false)
	if err != nil {
		return nil, err
	}

	role := new(coredocumentpb.Role)
	role.RoleKey = utils.RandomSlice(32)
	at, err := assembleAccessToken(ctx, payload, role.RoleKey)
	if err != nil {
		return nil, errors.New("failed to construct access token: %v", err)
	}

	role.AccessTokens = append(role.AccessTokens, at)
	ncd.addNewRule(role, coredocumentpb.Action_ACTION_READ)
	return ncd, ncd.setSalts()
}

// assembleAccessToken assembles a Read Access Token from the payload received
func assembleAccessToken(ctx context.Context, payload documentpb.AccessTokenParams, roleKey []byte) (*coredocumentpb.AccessToken, error) {
	account, err := contextutil.Account(ctx)
	if err != nil {
		return nil, err
	}
	tokenIdentifier := utils.RandomSlice(32)
	id, err := account.GetIdentityID()
	if err != nil {
		return nil, err
	}
	granterID := identity.NewDIDFromByte(id)
	roleID := roleKey
	granteeID, err := identity.NewDIDFromString(payload.Grantee)
	if err != nil {
		return nil, err
	}
	// assemble access token message to be signed
	docID, err := hexutil.Decode(payload.DocumentIdentifier)
	if err != nil {
		return nil, err
	}

	tm, err := assembleTokenMessage(tokenIdentifier, granterID, granteeID, roleID[:], docID)
	if err != nil {
		return nil, err
	}

	// fetch key pair from identity
	// TODO: change to signing key pair once secp scheme is available
	sig, err := account.SignMsg(tm)
	if err != nil {
		return nil, err
	}

	keys, err := account.GetKeys()
	if err != nil {
		return nil, err
	}

	// assemble the access token, appending the signature and public keys
	at := &coredocumentpb.AccessToken{
		Identifier:         tokenIdentifier,
		Granter:            granterID[:],
		Grantee:            granteeID[:],
		RoleIdentifier:     roleID[:],
		DocumentIdentifier: docID,
		Signature:          sig.Signature,
		Key:                keys[identity.KeyPurposeSigning].PublicKey,
	}

	return at, nil
}

// assembleTokenMessage assembles a token message
func assembleTokenMessage(tokenIdentifier []byte, granterID identity.DID, granteeID identity.DID, roleID []byte, docID []byte) ([]byte, error) {
	ids := [][]byte{tokenIdentifier, roleID, docID}
	for _, id := range ids {
		if len(id) != idSize {
			return nil, errors.New("invalid identifier length")
		}
	}

	tm := append(tokenIdentifier, granterID[:]...)
	tm = append(tm, granteeID[:]...)
	tm = append(tm, roleID...)
	tm = append(tm, docID...)
	return tm, nil
}
