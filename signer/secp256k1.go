package signer

import (
	"encoding/binary"
	"errors"

	"github.com/nervina-labs/joyid-sdk-go/crypto/keccak"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256k1"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

const (
	// unlock_mode + pubkey + signature
	secp256k1EmptyWitnessLockLen = 86
)

func signSecp256k1Tx(tx *types.Transaction, key *secp256k1.Key, mode byte) error {
	buf := tx.ComputeHash().Bytes()
	witnesses := tx.Witnesses
	if len(witnesses) < 1 {
		return errors.New("first witness cannot be empty")
	}
	firstWitnessArgs, err := types.DeserializeWitnessArgs(witnesses[0])
	if err != nil {
		return errors.New("first witness must be WitnessArgs")
	}
	emptyWitness := types.WitnessArgs{
		Lock:       make([]byte, secp256k1EmptyWitnessLockLen),
		InputType:  firstWitnessArgs.InputType,
		OutputType: firstWitnessArgs.OutputType,
	}
	emptyWitnessBytes := emptyWitness.Serialize()

	bytesLen := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytesLen, uint64(len(emptyWitnessBytes)))
	buf = append(buf, bytesLen...)
	buf = append(buf, emptyWitnessBytes...)

	for i := 1; i < len(witnesses); i++ {
		bytes := tx.Witnesses[i]
		bytesLen := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytesLen, uint64(len(bytes)))
		buf = append(buf, bytesLen...)
		buf = append(buf, bytes...)
	}

	sighash := keccak.Keccak256(buf)
	if err != nil {
		return err
	}

	// personal hash, ethereum prefix  \u0019Ethereum Signed Message:\n32
	personalEthereumSignPrefix := [...]byte{
		0x19, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x20, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64,
		0x20, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x3a, 0x0a, 0x33, 0x32,
	}
	message := personalEthereumSignPrefix[:]
	message = append(message, sighash...)
	messageHash := keccak.Keccak256(message)

	signature := key.Sign(messageHash)
	pubkeyHash := key.PubkeyHash()

	witnessArgsLock := []byte{mode}
	witnessArgsLock = append(witnessArgsLock, pubkeyHash...)
	witnessArgsLock = append(witnessArgsLock, signature...)
	firstWitnessArgs.Lock = witnessArgsLock
	tx.Witnesses[0] = firstWitnessArgs.Serialize()
	return nil
}
