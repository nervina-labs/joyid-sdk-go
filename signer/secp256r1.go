package signer

import (
	"encoding/base64"
	"encoding/binary"
	"errors"

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
	msgHashHex := utils.BytesToHex(msgHash)

	if err != nil {
		return "", err
	}
	challenge := make([]byte, 86)
	base64.RawURLEncoding.Encode(challenge, []byte(msgHashHex))
	return utils.BytesToHex(challenge), nil
}

func signSecp256r1Tx(tx *types.Transaction, key *secp256r1.Key, mode byte, webAuthn *WebAuthnMsg) error {
	clientDataBytes, err := utils.HexToBytes(webAuthn.ClientData)
	if err != nil {
		return errors.New("hex convert error")
	}
	clientDataHash := sha256.Sha256(clientDataBytes)
	authData, err := utils.HexToBytes(webAuthn.AuthData)
	if err != nil {
		return errors.New("hex convert error")
	}
	signData := authData
	signData = append(signData, clientDataHash...)
	signature := key.Sign(sha256.Sha256(signData))
	_, pubkey := key.Pubkey()

	if len(tx.Witnesses) < 1 {
		return errors.New("first witness cannot be empty")
	}
	firstWitnessArgs, err := types.DeserializeWitnessArgs(tx.Witnesses[0])
	if err != nil {
		return errors.New("first witness must be WitnessArgs")
	}
	witnessArgsLock := []byte{mode}
	witnessArgsLock = append(witnessArgsLock, pubkey...)
	witnessArgsLock = append(witnessArgsLock, signature...)
	witnessArgsLock = append(witnessArgsLock, authData...)
	witnessArgsLock = append(witnessArgsLock, clientDataBytes...)
	firstWitnessArgs.Lock = witnessArgsLock
	if err != nil {
		return errors.New("hex convert error")
	}
	tx.Witnesses[0] = firstWitnessArgs.Serialize()
	return nil
}
