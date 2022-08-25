// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bind_test

import (
	"context"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/fff-chain/go-fff"
	"github.com/fff-chain/go-fff/accounts/abi"
	"github.com/fff-chain/go-fff/accounts/abi/bind"
	"github.com/fff-chain/go-fff/common"
	"github.com/fff-chain/go-fff/common/hexutil"
	"github.com/fff-chain/go-fff/core/types"
	"github.com/fff-chain/go-fff/crypto"
	"github.com/fff-chain/go-fff/rlp"
)

type mockCaller struct {
	codeAtBlockNumber         *big.Int
	callContractBlockNumber   *big.Int
	pendingCodeAtCalled       bool
	pendingCallContractCalled bool
}

func (mc *mockCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	mc.codeAtBlockNumber = blockNumber
	return []byte{1, 2, 3}, nil
}

func (mc *mockCaller) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	mc.callContractBlockNumber = blockNumber
	return nil, nil
}

func (mc *mockCaller) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	mc.pendingCodeAtCalled = true
	return nil, nil
}

func (mc *mockCaller) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	mc.pendingCallContractCalled = true
	return nil, nil
}
func TestPassingBlockNumber(t *testing.T) {

	mc := &mockCaller{}

	bc := bind.NewBoundContract(common.HexToAddress("0x0"), abi.ABI{
		Methods: map[string]abi.Method{
			"something": {
				Name:    "something",
				Outputs: abi.Arguments{},
			},
		},
	}, mc, nil, nil)

	blockNumber := big.NewInt(42)

	bc.Call(&bind.CallOpts{BlockNumber: blockNumber}, nil, "something")

	if mc.callContractBlockNumber != blockNumber {
		t.Fatalf("CallContract() was not passed the block number")
	}

	if mc.codeAtBlockNumber != blockNumber {
		t.Fatalf("CodeAt() was not passed the block number")
	}

	bc.Call(&bind.CallOpts{}, nil, "something")

	if mc.callContractBlockNumber != nil {
		t.Fatalf("CallContract() was passed a block number when it should not have been")
	}

	if mc.codeAtBlockNumber != nil {
		t.Fatalf("CodeAt() was passed a block number when it should not have been")
	}

	bc.Call(&bind.CallOpts{BlockNumber: blockNumber, Pending: true}, nil, "something")

	if !mc.pendingCallContractCalled {
		t.Fatalf("CallContract() was not passed the block number")
	}

	if !mc.pendingCodeAtCalled {
		t.Fatalf("CodeAt() was not passed the block number")
	}
}

const hexData = "0x000000000000000000000000376c47978271565f56deb45495afa69e59c16ab200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000158"

func TestUnpackIndexedStringTyLogIntoMap(t *testing.T) {
	hash := crypto.Keccak256Hash([]byte("testName"))
	topics := []common.Hash{
		common.HexToHash("0x0"),
		hash,
	}
	mockLog := newMockLog(topics, common.HexToHash("0x0"))

	abiString := `[{"anonymous":false,"inputs":[{"indexed":true,"name":"name","type":"string"},{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"memo","type":"bytes"}],"name":"received","type":"event"}]`
	parsedAbi, _ := abi.JSON(strings.NewReader(abiString))
	bc := bind.NewBoundContract(common.HexToAddress("0x0"), parsedAbi, nil, nil, nil)

	expectedReceivedMap := map[string]interface{}{
		"name":   hash,
		"sender": common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		"amount": big.NewInt(1),
		"memo":   []byte{88},
	}
	unpackAndCheck(t, bc, expectedReceivedMap, mockLog)
}

func TestUnpackIndexedSliceTyLogIntoMap(t *testing.T) {
	sliceBytes, err := rlp.EncodeToBytes([]string{"name1", "name2", "name3", "name4"})
	if err != nil {
		t.Fatal(err)
	}
	hash := crypto.Keccak256Hash(sliceBytes)
	topics := []common.Hash{
		common.HexToHash("0x0"),
		hash,
	}
	mockLog := newMockLog(topics, common.HexToHash("0x0"))

	abiString := `[{"anonymous":false,"inputs":[{"indexed":true,"name":"names","type":"string[]"},{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"memo","type":"bytes"}],"name":"received","type":"event"}]`
	parsedAbi, _ := abi.JSON(strings.NewReader(abiString))
	bc := bind.NewBoundContract(common.HexToAddress("0x0"), parsedAbi, nil, nil, nil)

	expectedReceivedMap := map[string]interface{}{
		"names":  hash,
		"sender": common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		"amount": big.NewInt(1),
		"memo":   []byte{88},
	}
	unpackAndCheck(t, bc, expectedReceivedMap, mockLog)
}

func TestUnpackIndexedArrayTyLogIntoMap(t *testing.T) {
	arrBytes, err := rlp.EncodeToBytes([2]common.Address{common.HexToAddress("0x0"), common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2")})
	if err != nil {
		t.Fatal(err)
	}
	hash := crypto.Keccak256Hash(arrBytes)
	topics := []common.Hash{
		common.HexToHash("0x0"),
		hash,
	}
	mockLog := newMockLog(topics, common.HexToHash("0x0"))

	abiString := `[{"anonymous":false,"inputs":[{"indexed":true,"name":"addresses","type":"address[2]"},{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"memo","type":"bytes"}],"name":"received","type":"event"}]`
	parsedAbi, _ := abi.JSON(strings.NewReader(abiString))
	bc := bind.NewBoundContract(common.HexToAddress("0x0"), parsedAbi, nil, nil, nil)

	expectedReceivedMap := map[string]interface{}{
		"addresses": hash,
		"sender":    common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		"amount":    big.NewInt(1),
		"memo":      []byte{88},
	}
	unpackAndCheck(t, bc, expectedReceivedMap, mockLog)
}

func TestUnpackIndexedFuncTyLogIntoMap(t *testing.T) {
	mockAddress := common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2")
	addrBytes := mockAddress.Bytes()
	hash := crypto.Keccak256Hash([]byte("mockFunction(address,uint)"))
	functionSelector := hash[:4]
	functionTyBytes := append(addrBytes, functionSelector...)
	var functionTy [24]byte
	copy(functionTy[:], functionTyBytes[0:24])
	topics := []common.Hash{
		common.HexToHash("0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"),
		common.BytesToHash(functionTyBytes),
	}
	mockLog := newMockLog(topics, common.HexToHash("0x5c698f13940a2153440c6d19660878bc90219d9298fdcf37365aa8d88d40fc42"))
	abiString := `[{"anonymous":false,"inputs":[{"indexed":true,"name":"function","type":"function"},{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"memo","type":"bytes"}],"name":"received","type":"event"}]`
	parsedAbi, _ := abi.JSON(strings.NewReader(abiString))
	bc := bind.NewBoundContract(common.HexToAddress("0x0"), parsedAbi, nil, nil, nil)

	expectedReceivedMap := map[string]interface{}{
		"function": functionTy,
		"sender":   common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		"amount":   big.NewInt(1),
		"memo":     []byte{88},
	}
	unpackAndCheck(t, bc, expectedReceivedMap, mockLog)
}

func TestUnpackIndexedBytesTyLogIntoMap(t *testing.T) {
	bytes := []byte{1, 2, 3, 4, 5}
	hash := crypto.Keccak256Hash(bytes)
	topics := []common.Hash{
		common.HexToHash("0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"),
		hash,
	}
	mockLog := newMockLog(topics, common.HexToHash("0x5c698f13940a2153440c6d19660878bc90219d9298fdcf37365aa8d88d40fc42"))

	abiString := `[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "userCount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "count",
				"type": "uint256"
			}
		],
		"name": "DistributeBlockRewardFunc",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "user",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "count",
				"type": "uint256"
			}
		],
		"name": "UserStakeFFF",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "user",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "count",
				"type": "uint256"
			}
		],
		"name": "UserUnStakeFFF",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "batchTransfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "string",
				"name": "reason",
				"type": "string"
			}
		],
		"name": "batchTransferFailed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "reason",
				"type": "bytes"
			}
		],
		"name": "batchTransferLowerFailed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "validator",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "deprecatedDeposit",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address payable",
				"name": "validator",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "directTransfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address payable",
				"name": "validator",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "directTransferFail",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "message",
				"type": "string"
			}
		],
		"name": "failReasonWithStr",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "feeBurned",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "key",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "value",
				"type": "bytes"
			}
		],
		"name": "paramChange",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "systemTransfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint8",
				"name": "channelId",
				"type": "uint8"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "msgBytes",
				"type": "bytes"
			}
		],
		"name": "unexpectedPackage",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "validator",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "validatorDeposit",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "validator",
				"type": "address"
			}
		],
		"name": "validatorEmptyJailed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "validator",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "validatorFelony",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "validator",
				"type": "address"
			}
		],
		"name": "validatorJailed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "validator",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "validatorMisdemeanor",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [],
		"name": "validatorSetUpdated",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "BIND_CHANNELID",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "BURN_ADDRESS",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "BURN_RATIO_SCALE",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "CODE_OK",
		"outputs": [
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "CROSS_CHAIN_CONTRACT_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "DUSTY_INCOMING",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "ERROR_FAIL_CHECK_VALIDATORS",
		"outputs": [
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "ERROR_FAIL_DECODE",
		"outputs": [
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "ERROR_LEN_OF_VAL_MISMATCH",
		"outputs": [
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "ERROR_RELAYFEE_TOO_LARGE",
		"outputs": [
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "ERROR_UNKNOWN_PACKAGE_TYPE",
		"outputs": [
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "EXPIRE_TIME_SECOND_GAP",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "GOV_CHANNELID",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "GOV_HUB_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "GetAllStakeInfo",
		"outputs": [
			{
				"components": [
					{
						"internalType": "address",
						"name": "stakeAddress",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "stakeCount",
						"type": "uint256"
					}
				],
				"internalType": "struct BSCValidatorSet.StakeInfo[]",
				"name": "",
				"type": "tuple[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "INCENTIVIZE_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "INIT_BURN_RATIO",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "INIT_VALIDATORSET_BYTES",
		"outputs": [
			{
				"internalType": "bytes",
				"name": "",
				"type": "bytes"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "JAIL_MESSAGE_TYPE",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "LIGHT_CLIENT_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "MAX_NUM_OF_VALIDATORS",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "MAX_STAKE_FFF_COUNT",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "PRECISION",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "RELAYERHUB_CONTRACT_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "SLASH_CHANNELID",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "SLASH_CONTRACT_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "STAKING_CHANNELID",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "SYSTEM_REWARD_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "StakeFFF",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "TOKEN_HUB_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "TOKEN_MANAGER_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "TRANSFER_IN_CHANNELID",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "TRANSFER_OUT_CHANNELID",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "UnStakeFFF",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "VALIDATORS_UPDATE_MESSAGE_TYPE",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "VALIDATOR_CONTRACT_ADDR",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "alreadyInit",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "bscChainID",
		"outputs": [
			{
				"internalType": "uint16",
				"name": "",
				"type": "uint16"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "burnRatio",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "burnRatioInitialized",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "currStakeFFF",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "currStakePeople",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "currStakePeopleIndex",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "currentValidatorSet",
		"outputs": [
			{
				"internalType": "address",
				"name": "consensusAddress",
				"type": "address"
			},
			{
				"internalType": "address payable",
				"name": "feeAddress",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "BBCFeeAddress",
				"type": "address"
			},
			{
				"internalType": "uint64",
				"name": "votingPower",
				"type": "uint64"
			},
			{
				"internalType": "bool",
				"name": "jailed",
				"type": "bool"
			},
			{
				"internalType": "uint256",
				"name": "incoming",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "currentValidatorSetMap",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "valAddr",
				"type": "address"
			}
		],
		"name": "deposit",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "expireTimeSecondGap",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "validator",
				"type": "address"
			}
		],
		"name": "felony",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "validator",
				"type": "address"
			}
		],
		"name": "getIncoming",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getValidators",
		"outputs": [
			{
				"internalType": "address[]",
				"name": "",
				"type": "address[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint8",
				"name": "channelId",
				"type": "uint8"
			},
			{
				"internalType": "bytes",
				"name": "msgBytes",
				"type": "bytes"
			}
		],
		"name": "handleAckPackage",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint8",
				"name": "channelId",
				"type": "uint8"
			},
			{
				"internalType": "bytes",
				"name": "msgBytes",
				"type": "bytes"
			}
		],
		"name": "handleFailAckPackage",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			},
			{
				"internalType": "bytes",
				"name": "msgBytes",
				"type": "bytes"
			}
		],
		"name": "handleSynPackage",
		"outputs": [
			{
				"internalType": "bytes",
				"name": "responsePayload",
				"type": "bytes"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "init",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "maxStakeFFFCount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "maxStakePeopleCount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "miniStakeFFFCount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "validator",
				"type": "address"
			}
		],
		"name": "misdemeanor",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "myBalance",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "numOfJailed",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "stakeInfoIndexMap",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "stakeInfoMap",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalInComing",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "key",
				"type": "string"
			},
			{
				"internalType": "bytes",
				"name": "value",
				"type": "bytes"
			}
		],
		"name": "updateParam",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`
	parsedAbi, _ := abi.JSON(strings.NewReader(abiString))
	bc := bind.NewBoundContract(common.HexToAddress("0x0"), parsedAbi, nil, nil, nil)

	expectedReceivedMap := map[string]interface{}{
		"content": hash,
		"sender":  common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		"amount":  big.NewInt(1),
		"memo":    []byte{88},
	}
	unpackAndCheck(t, bc, expectedReceivedMap, mockLog)
}

func unpackAndCheck(t *testing.T, bc *bind.BoundContract, expected map[string]interface{}, mockLog types.Log) {
	received := make(map[string]interface{})
	if err := bc.UnpackLogIntoMap(received, "received", mockLog); err != nil {
		t.Error(err)
	}

	if len(received) != len(expected) {
		t.Fatalf("unpacked map length %v not equal expected length of %v", len(received), len(expected))
	}
	for name, elem := range expected {
		if !reflect.DeepEqual(elem, received[name]) {
			t.Errorf("field %v does not match expected, want %v, got %v", name, elem, received[name])
		}
	}
}

func newMockLog(topics []common.Hash, txHash common.Hash) types.Log {
	return types.Log{
		Address:     common.HexToAddress("0x0"),
		Topics:      topics,
		Data:        hexutil.MustDecode(hexData),
		BlockNumber: uint64(26),
		TxHash:      txHash,
		TxIndex:     111,
		BlockHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		Index:       7,
		Removed:     false,
	}
}
