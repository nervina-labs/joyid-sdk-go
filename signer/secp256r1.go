package signer

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256r1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/sha256"
	"github.com/nervina-labs/joyid-sdk-go/utils"
	"github.com/nervosnetwork/ckb-sdk-go/v2/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

const (
	secp256r1EmptyWitnessLockLen = 129 // unlock_mode + pubkey + signature
)

type WebAuthnMsg struct {
	AuthData   string
	ClientData string
}

func GenerateWebAuthnChallenge(tx *types.Transaction) (string, error) {
	txHash := tx.ComputeHash()
	msg := txHash.Bytes()

	witnesses := tx.Witnesses
	if len(witnesses) < 1 {
		return "", errors.New("first witness cannot be empty")
	}
	firstWitnessArgs, err := types.DeserializeWitnessArgs(witnesses[0])
	if err != nil {
		return "", errors.New("first witness must be WitnessArgs")
	}
	emptyWitness := types.WitnessArgs{
		Lock:       make([]byte, secp256r1EmptyWitnessLockLen),
		InputType:  firstWitnessArgs.InputType,
		OutputType: firstWitnessArgs.OutputType,
	}
	emptyWitnessBytes := emptyWitness.Serialize()

	bytesLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytesLen, uint64(len(emptyWitnessBytes)))
	msg = append(msg, bytesLen...)
	msg = append(msg, emptyWitnessBytes...)

	for i := 1; i < len(witnesses); i++ {
		bytes := tx.Witnesses[i]
		bytesLen := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytesLen, uint64(len(bytes)))
		msg = append(msg, bytesLen...)
		msg = append(msg, bytes...)
	}

	msgHash := blake2b.Blake256(msg)
	if err != nil {
		return "", err
	}
	challenge := make([]byte, 86)
	base64.StdEncoding.Encode(challenge, msgHash)
	return utils.BytesToHex(challenge), nil
}

func signSecp256r1Tx(tx *types.Transaction, key *secp256r1.Key, mode unlockMode, webAuthn *WebAuthnMsg) error {
	clientDataHash := sha256.Sha256([]byte(webAuthn.ClientData))
	signData := fmt.Sprintf("%s%s", webAuthn.AuthData, clientDataHash)
	signature := key.Sign(signData)
	_, pubkey := key.Pubkey()

	if len(tx.Witnesses) < 1 {
		return errors.New("first witness cannot be empty")
	}
	firstWitnessArgs, err := types.DeserializeWitnessArgs(tx.Witnesses[0])
	if err != nil {
		return errors.New("first witness must be WitnessArgs")
	}
	firstWitnessArgs.Lock = []byte(fmt.Sprintf("%x%s%s%s%s", mode, pubkey, signature, webAuthn.AuthData, webAuthn.ClientData))

	tx.Witnesses[0] = firstWitnessArgs.Serialize()
	return nil
}
