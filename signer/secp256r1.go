package signer

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256r1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/sha256"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

const (
	secp256r1EmptyWitnessLockLen = 129 // unlock_mode + pubkey + signature
)

func signSecp256r1Tx(tx *types.Transaction, key *secp256r1.Key, mode unlockMode) error {
	txHash := tx.ComputeHash()
	msg := txHash.Bytes()

	witnesses := tx.Witnesses
	if len(witnesses) < 1 {
		return errors.New("first witness cannot be empty")
	}
	firstWitnessArgs, err := types.DeserializeWitnessArgs(witnesses[0])
	if err != nil {
		return errors.New("first witness must be WitnessArgs")
	}
	firstWitness := types.WitnessArgs{
		Lock:       make([]byte, secp256r1EmptyWitnessLockLen),
		InputType:  firstWitnessArgs.InputType,
		OutputType: firstWitnessArgs.OutputType,
	}
	firstWitnessBytes := firstWitness.Serialize()

	bytesLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytesLen, uint64(len(firstWitnessBytes)))
	msg = append(msg, bytesLen...)
	msg = append(msg, firstWitnessBytes...)

	for i := 1; i < len(witnesses); i++ {
		bytes := tx.Witnesses[i]
		bytesLen := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytesLen, uint64(len(bytes)))
		msg = append(msg, bytesLen...)
		msg = append(msg, bytes...)
	}

	msgHash, err := blake2b.Blake256(msg)
	if err != nil {
		return err
	}
	challenge := make([]byte, 86)
	base64.StdEncoding.Encode(challenge, msgHash)

	authData := "49960de5880e8c687434170f6476605b8fe4aeb9a28632c7995cf3ba831d97630162f9fb77"
	clientData := fmt.Sprintf("7b2274797065223a22776562617574686e2e676574222c226368616c6c656e6765223a22%x222c226f726967696e223a22687474703a2f2f6c6f63616c686f73743a38303030222c2263726f73734f726967696e223a66616c73657d", challenge)

	clientDataHash := sha256.Sha256([]byte(clientData))
	signData := fmt.Sprintf("%s%s", authData, clientDataHash)
	signature := key.Sign(signData)
	_, pubkey := key.Pubkey()

	firstWitness.Lock = []byte(fmt.Sprintf("%x%s%s%s%s", mode, pubkey, signature, authData, clientData))

	signedWitnesses := [][]byte{firstWitness.Serialize()}
	tx.Witnesses = append(signedWitnesses, witnesses[1:]...)

	return nil
}
