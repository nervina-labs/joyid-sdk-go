# joyid-sdk-go

The go sdk of [JoyID Lock Script Demo with dynamic link library joyid.so](https://github.com/nervina-labs/joyid-lib-demo) and the contract deployment can be seen in [wiki-testnet-deployment](https://github.com/nervina-labs/joyid-lib-demo/wiki/Testnet-Deployment)

## Install

```shell
go get github.com/nervina-labs/joyid-sdk-go
```

## Quick Start

Some transfer examples with JoyID lock script are provided in example module.

### JoyID native unlock

- **Secp256r1(WebAuthn)**

```go
// example/main.go
func NativeTransferWithR1() error
```

- **Secp256k1(Ethereum)**

```go
// example/main.go
func NativeTransferWithK1() error
```

### JoyID subkey unlock

- **Secp256r1(WebAuthn)**

```go
// example/main.go
func SubkeyTransferWithR1() error
```

- **Secp256k1(Ethereum)**

```go
// example/main.go
func SubkeyTransferWithK1() error
```

### Add subkey with native unlock

- **Secp256r1(WebAuthn)**

Add secp256r1 subkey to JoyID account with secp256r1 native unlock

```go
// example/main.go
func AddSecp256r1SubkeyWithNativeUnlock() error
```

- **Secp256k1(Ethereum)**

Add secp256k1 subkey to JoyID account with secp256k1 native unlock

```go
// example/main.go
func AddSecp256k1SubkeyWithNativeUnlock() error
```
