package ed25519

import (
	"encoding/base64"

	"github.com/CentrifugeInc/go-centrifuge/centrifuge/config"
	"github.com/CentrifugeInc/go-centrifuge/centrifuge/utils"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	mh "github.com/multiformats/go-multihash"
	"golang.org/x/crypto/ed25519"
)

var log = logging.Logger("ed25519")

const (
	PublicKey  = "PUBLIC KEY"
	PrivateKey = "PRIVATE KEY"
)

// GetPublicSigningKey returns the public key from the file
func GetPublicSigningKey(fileName string) (publicKey ed25519.PublicKey) {
	key, err := utils.ReadKeyFromPemFile(fileName, PublicKey)

	if err != nil {
		log.Fatal(err)
	}
	publicKey = ed25519.PublicKey(key)
	return
}

// GetPrivateSigningKey returns the private key from the file
func GetPrivateSigningKey(fileName string) (privateKey ed25519.PrivateKey) {
	key, err := utils.ReadKeyFromPemFile(fileName, PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	privateKey = ed25519.PrivateKey(key)
	return
}

// GetSigningKeyPairFromConfig returns the public and private key pair from the config
func GetSigningKeyPairFromConfig() (publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {
	pub, priv := config.Config.GetSigningKeyPair()
	publicKey = GetPublicSigningKey(pub)
	privateKey = GetPrivateSigningKey(priv)
	return
}

// GenerateSigningKeyPair generates ed25519 key pair
func GenerateSigningKeyPair() (publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// PublicKeyToP2PKey returns p2pId from the public key
func PublicKeyToP2PKey(publicKey [32]byte) (p2pId peer.ID, err error) {
	// Taken from peer.go#IDFromPublicKey#L189
	// TODO As soon as this is merged: https://github.com/libp2p/go-libp2p-kad-dht/pull/129 we can get rid of this function
	// and only do:
	// pk, err := crypto.UnmarshalEd25519PublicKey(publicKey[:])
	// pid, error := IDFromPublicKey(pk)
	pk, err := crypto.UnmarshalEd25519PublicKey(publicKey[:])
	bpk, err := pk.Bytes()
	hash, err := mh.Sum(bpk[:], mh.SHA2_256, -1)
	if err != nil {
		return "", err
	}

	p2pId = peer.ID(hash)
	return
}

// GetIDConfig reads the keys and ID from the config and returns a the Identity config
func GetIDConfig() (identityConfig *config.IdentityConfig, err error) {
	pk, pvk := GetSigningKeyPairFromConfig()
	decodedId, err := base64.StdEncoding.DecodeString(string(config.Config.GetIdentityId()))
	if err != nil {
		return nil, err
	}

	identityConfig = &config.IdentityConfig{
		ID:         decodedId,
		PublicKey:  pk,
		PrivateKey: pvk,
	}
	return
}