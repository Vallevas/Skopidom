// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BlockchainMetaData contains all meta data concerning the Blockchain contract.
var BlockchainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"itemId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"action\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"actor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"payload\",\"type\":\"string\"}],\"name\":\"AuditLogged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"getEvent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"itemId\",\"type\":\"uint256\"}],\"name\":\"getItemEventCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"itemId\",\"type\":\"uint256\"}],\"name\":\"getItemEvents\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalEventCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"itemId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"action\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"payload\",\"type\":\"string\"}],\"name\":\"logEvent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b505f5f81905550610d82806100225f395ff3fe608060405234801561000f575f5ffd5b5060043610610055575f3560e01c80630ed0f06b146100595780631bd7a2601461007757806336529d57146100a75780635c650d8d146100d75780636d1884e014610107575b5f5ffd5b61006161013c565b60405161006e91906105f2565b60405180910390f35b610091600480360381019061008c919061063d565b610144565b60405161009e919061071f565b60405180910390f35b6100c160048036038101906100bc919061063d565b6101ab565b6040516100ce91906105f2565b60405180910390f35b6100f160048036038101906100ec91906107a0565b6101c8565b6040516100fe91906105f2565b60405180910390f35b610121600480360381019061011c919061063d565b6103ee565b604051610133969594939291906108e0565b60405180910390f35b5f5f54905090565b606060025f8381526020019081526020015f2080548060200260200160405190810160405280929190818152602001828054801561019f57602002820191905f5260205f20905b81548152602001906001019080831161018b575b50505050509050919050565b5f60025f8381526020019081526020015f20805490509050919050565b5f5f5f8154809291906101da9061097a565b91905055505f5f5490505f6040518060c0016040528083815260200189815260200188888080601f0160208091040260200160405190810160405280939291908181526020018383808284375f81840152601f19601f8201169050808301925050505050505081526020013373ffffffffffffffffffffffffffffffffffffffff16815260200142815260200186868080601f0160208091040260200160405190810160405280939291908181526020018383808284375f81840152601f19601f8201169050808301925050505050505081525090508060015f8481526020019081526020015f205f820151815f01556020820151816001015560408201518160020190816102e99190610bfc565b506060820151816003015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506080820151816004015560a082015181600501908161034f9190610bfc565b5090505060025f8981526020019081526020015f2082908060018154018082558091505060019003905f5260205f20015f90919091909150553373ffffffffffffffffffffffffffffffffffffffff1688837fc160262fc35c7fe77e722bfdfc4a0b237ed57a1f260ba0f6b26b856b496dcd2f8a8a428b8b6040516103d8959493929190610d05565b60405180910390a4819250505095945050505050565b5f5f60605f5f60605f60015f8981526020019081526020015f206040518060c00160405290815f82015481526020016001820154815260200160028201805461043690610a1b565b80601f016020809104026020016040519081016040528092919081815260200182805461046290610a1b565b80156104ad5780601f10610484576101008083540402835291602001916104ad565b820191905f5260205f20905b81548152906001019060200180831161049057829003601f168201915b50505050508152602001600382015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016004820154815260200160058201805461052590610a1b565b80601f016020809104026020016040519081016040528092919081815260200182805461055190610a1b565b801561059c5780601f106105735761010080835404028352916020019161059c565b820191905f5260205f20905b81548152906001019060200180831161057f57829003601f168201915b5050505050815250509050805f015181602001518260400151836060015184608001518560a001519650965096509650965096505091939550919395565b5f819050919050565b6105ec816105da565b82525050565b5f6020820190506106055f8301846105e3565b92915050565b5f5ffd5b5f5ffd5b61061c816105da565b8114610626575f5ffd5b50565b5f8135905061063781610613565b92915050565b5f602082840312156106525761065161060b565b5b5f61065f84828501610629565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b61069a816105da565b82525050565b5f6106ab8383610691565b60208301905092915050565b5f602082019050919050565b5f6106cd82610668565b6106d78185610672565b93506106e283610682565b805f5b838110156107125781516106f988826106a0565b9750610704836106b7565b9250506001810190506106e5565b5085935050505092915050565b5f6020820190508181035f83015261073781846106c3565b905092915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f8401126107605761075f61073f565b5b8235905067ffffffffffffffff81111561077d5761077c610743565b5b60208301915083600182028301111561079957610798610747565b5b9250929050565b5f5f5f5f5f606086880312156107b9576107b861060b565b5b5f6107c688828901610629565b955050602086013567ffffffffffffffff8111156107e7576107e661060f565b5b6107f38882890161074b565b9450945050604086013567ffffffffffffffff8111156108165761081561060f565b5b6108228882890161074b565b92509250509295509295909350565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61087382610831565b61087d818561083b565b935061088d81856020860161084b565b61089681610859565b840191505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6108ca826108a1565b9050919050565b6108da816108c0565b82525050565b5f60c0820190506108f35f8301896105e3565b61090060208301886105e3565b81810360408301526109128187610869565b905061092160608301866108d1565b61092e60808301856105e3565b81810360a08301526109408184610869565b9050979650505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610984826105da565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036109b6576109b561094d565b5b600182019050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610a3257607f821691505b602082108103610a4557610a446109ee565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302610aa77fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610a6c565b610ab18683610a6c565b95508019841693508086168417925050509392505050565b5f819050919050565b5f610aec610ae7610ae2846105da565b610ac9565b6105da565b9050919050565b5f819050919050565b610b0583610ad2565b610b19610b1182610af3565b848454610a78565b825550505050565b5f5f905090565b610b30610b21565b610b3b818484610afc565b505050565b5f5b82811015610b6157610b565f828401610b28565b600181019050610b42565b505050565b601f821115610bb45782821115610bb357610b8081610a4b565b610b8983610a5d565b610b9285610a5d565b6020861015610b9f575f90505b808301610bae82840382610b40565b505050505b5b505050565b5f82821c905092915050565b5f610bd45f1984600802610bb9565b1980831691505092915050565b5f610bec8383610bc5565b9150826002028217905092915050565b610c0582610831565b67ffffffffffffffff811115610c1e57610c1d6109c1565b5b610c288254610a1b565b610c33828285610b66565b5f60209050601f831160018114610c64575f8415610c52578287015190505b610c5c8582610be1565b865550610cc3565b601f198416610c7286610a4b565b5f5b82811015610c9957848901518255600182019150602085019450602081019050610c74565b86831015610cb65784890151610cb2601f891682610bc5565b8355505b6001600288020188555050505b505050505050565b828183375f83830152505050565b5f610ce4838561083b565b9350610cf1838584610ccb565b610cfa83610859565b840190509392505050565b5f6060820190508181035f830152610d1e818789610cd9565b9050610d2d60208301866105e3565b8181036040830152610d40818486610cd9565b9050969550505050505056fea2646970667358221220a666d3617d74abfb1b42271c92918de6ea25a470ccde0dfd49353658c02c336464736f6c63430008230033",
}

// BlockchainABI is the input ABI used to generate the binding from.
// Deprecated: Use BlockchainMetaData.ABI instead.
var BlockchainABI = BlockchainMetaData.ABI

// BlockchainBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BlockchainMetaData.Bin instead.
var BlockchainBin = BlockchainMetaData.Bin

// DeployBlockchain deploys a new Ethereum contract, binding an instance of Blockchain to it.
func DeployBlockchain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Blockchain, error) {
	parsed, err := BlockchainMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BlockchainBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Blockchain{BlockchainCaller: BlockchainCaller{contract: contract}, BlockchainTransactor: BlockchainTransactor{contract: contract}, BlockchainFilterer: BlockchainFilterer{contract: contract}}, nil
}

// Blockchain is an auto generated Go binding around an Ethereum contract.
type Blockchain struct {
	BlockchainCaller     // Read-only binding to the contract
	BlockchainTransactor // Write-only binding to the contract
	BlockchainFilterer   // Log filterer for contract events
}

// BlockchainCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlockchainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockchainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlockchainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockchainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlockchainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockchainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlockchainSession struct {
	Contract     *Blockchain       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlockchainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlockchainCallerSession struct {
	Contract *BlockchainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// BlockchainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlockchainTransactorSession struct {
	Contract     *BlockchainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BlockchainRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlockchainRaw struct {
	Contract *Blockchain // Generic contract binding to access the raw methods on
}

// BlockchainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlockchainCallerRaw struct {
	Contract *BlockchainCaller // Generic read-only contract binding to access the raw methods on
}

// BlockchainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlockchainTransactorRaw struct {
	Contract *BlockchainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlockchain creates a new instance of Blockchain, bound to a specific deployed contract.
func NewBlockchain(address common.Address, backend bind.ContractBackend) (*Blockchain, error) {
	contract, err := bindBlockchain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Blockchain{BlockchainCaller: BlockchainCaller{contract: contract}, BlockchainTransactor: BlockchainTransactor{contract: contract}, BlockchainFilterer: BlockchainFilterer{contract: contract}}, nil
}

// NewBlockchainCaller creates a new read-only instance of Blockchain, bound to a specific deployed contract.
func NewBlockchainCaller(address common.Address, caller bind.ContractCaller) (*BlockchainCaller, error) {
	contract, err := bindBlockchain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlockchainCaller{contract: contract}, nil
}

// NewBlockchainTransactor creates a new write-only instance of Blockchain, bound to a specific deployed contract.
func NewBlockchainTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockchainTransactor, error) {
	contract, err := bindBlockchain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlockchainTransactor{contract: contract}, nil
}

// NewBlockchainFilterer creates a new log filterer instance of Blockchain, bound to a specific deployed contract.
func NewBlockchainFilterer(address common.Address, filterer bind.ContractFilterer) (*BlockchainFilterer, error) {
	contract, err := bindBlockchain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlockchainFilterer{contract: contract}, nil
}

// bindBlockchain binds a generic wrapper to an already deployed contract.
func bindBlockchain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BlockchainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blockchain *BlockchainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blockchain.Contract.BlockchainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blockchain *BlockchainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockchain.Contract.BlockchainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blockchain *BlockchainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blockchain.Contract.BlockchainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blockchain *BlockchainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blockchain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blockchain *BlockchainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockchain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blockchain *BlockchainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blockchain.Contract.contract.Transact(opts, method, params...)
}

// GetEvent is a free data retrieval call binding the contract method 0x6d1884e0.
//
// Solidity: function getEvent(uint256 eventId) view returns(uint256, uint256, string, address, uint256, string)
func (_Blockchain *BlockchainCaller) GetEvent(opts *bind.CallOpts, eventId *big.Int) (*big.Int, *big.Int, string, common.Address, *big.Int, string, error) {
	var out []interface{}
	err := _Blockchain.contract.Call(opts, &out, "getEvent", eventId)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(string), *new(common.Address), *new(*big.Int), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	out5 := *abi.ConvertType(out[5], new(string)).(*string)

	return out0, out1, out2, out3, out4, out5, err

}

// GetEvent is a free data retrieval call binding the contract method 0x6d1884e0.
//
// Solidity: function getEvent(uint256 eventId) view returns(uint256, uint256, string, address, uint256, string)
func (_Blockchain *BlockchainSession) GetEvent(eventId *big.Int) (*big.Int, *big.Int, string, common.Address, *big.Int, string, error) {
	return _Blockchain.Contract.GetEvent(&_Blockchain.CallOpts, eventId)
}

// GetEvent is a free data retrieval call binding the contract method 0x6d1884e0.
//
// Solidity: function getEvent(uint256 eventId) view returns(uint256, uint256, string, address, uint256, string)
func (_Blockchain *BlockchainCallerSession) GetEvent(eventId *big.Int) (*big.Int, *big.Int, string, common.Address, *big.Int, string, error) {
	return _Blockchain.Contract.GetEvent(&_Blockchain.CallOpts, eventId)
}

// GetItemEventCount is a free data retrieval call binding the contract method 0x36529d57.
//
// Solidity: function getItemEventCount(uint256 itemId) view returns(uint256)
func (_Blockchain *BlockchainCaller) GetItemEventCount(opts *bind.CallOpts, itemId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Blockchain.contract.Call(opts, &out, "getItemEventCount", itemId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetItemEventCount is a free data retrieval call binding the contract method 0x36529d57.
//
// Solidity: function getItemEventCount(uint256 itemId) view returns(uint256)
func (_Blockchain *BlockchainSession) GetItemEventCount(itemId *big.Int) (*big.Int, error) {
	return _Blockchain.Contract.GetItemEventCount(&_Blockchain.CallOpts, itemId)
}

// GetItemEventCount is a free data retrieval call binding the contract method 0x36529d57.
//
// Solidity: function getItemEventCount(uint256 itemId) view returns(uint256)
func (_Blockchain *BlockchainCallerSession) GetItemEventCount(itemId *big.Int) (*big.Int, error) {
	return _Blockchain.Contract.GetItemEventCount(&_Blockchain.CallOpts, itemId)
}

// GetItemEvents is a free data retrieval call binding the contract method 0x1bd7a260.
//
// Solidity: function getItemEvents(uint256 itemId) view returns(uint256[])
func (_Blockchain *BlockchainCaller) GetItemEvents(opts *bind.CallOpts, itemId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _Blockchain.contract.Call(opts, &out, "getItemEvents", itemId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetItemEvents is a free data retrieval call binding the contract method 0x1bd7a260.
//
// Solidity: function getItemEvents(uint256 itemId) view returns(uint256[])
func (_Blockchain *BlockchainSession) GetItemEvents(itemId *big.Int) ([]*big.Int, error) {
	return _Blockchain.Contract.GetItemEvents(&_Blockchain.CallOpts, itemId)
}

// GetItemEvents is a free data retrieval call binding the contract method 0x1bd7a260.
//
// Solidity: function getItemEvents(uint256 itemId) view returns(uint256[])
func (_Blockchain *BlockchainCallerSession) GetItemEvents(itemId *big.Int) ([]*big.Int, error) {
	return _Blockchain.Contract.GetItemEvents(&_Blockchain.CallOpts, itemId)
}

// GetTotalEventCount is a free data retrieval call binding the contract method 0x0ed0f06b.
//
// Solidity: function getTotalEventCount() view returns(uint256)
func (_Blockchain *BlockchainCaller) GetTotalEventCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Blockchain.contract.Call(opts, &out, "getTotalEventCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalEventCount is a free data retrieval call binding the contract method 0x0ed0f06b.
//
// Solidity: function getTotalEventCount() view returns(uint256)
func (_Blockchain *BlockchainSession) GetTotalEventCount() (*big.Int, error) {
	return _Blockchain.Contract.GetTotalEventCount(&_Blockchain.CallOpts)
}

// GetTotalEventCount is a free data retrieval call binding the contract method 0x0ed0f06b.
//
// Solidity: function getTotalEventCount() view returns(uint256)
func (_Blockchain *BlockchainCallerSession) GetTotalEventCount() (*big.Int, error) {
	return _Blockchain.Contract.GetTotalEventCount(&_Blockchain.CallOpts)
}

// LogEvent is a paid mutator transaction binding the contract method 0x5c650d8d.
//
// Solidity: function logEvent(uint256 itemId, string action, string payload) returns(uint256)
func (_Blockchain *BlockchainTransactor) LogEvent(opts *bind.TransactOpts, itemId *big.Int, action string, payload string) (*types.Transaction, error) {
	return _Blockchain.contract.Transact(opts, "logEvent", itemId, action, payload)
}

// LogEvent is a paid mutator transaction binding the contract method 0x5c650d8d.
//
// Solidity: function logEvent(uint256 itemId, string action, string payload) returns(uint256)
func (_Blockchain *BlockchainSession) LogEvent(itemId *big.Int, action string, payload string) (*types.Transaction, error) {
	return _Blockchain.Contract.LogEvent(&_Blockchain.TransactOpts, itemId, action, payload)
}

// LogEvent is a paid mutator transaction binding the contract method 0x5c650d8d.
//
// Solidity: function logEvent(uint256 itemId, string action, string payload) returns(uint256)
func (_Blockchain *BlockchainTransactorSession) LogEvent(itemId *big.Int, action string, payload string) (*types.Transaction, error) {
	return _Blockchain.Contract.LogEvent(&_Blockchain.TransactOpts, itemId, action, payload)
}

// BlockchainAuditLoggedIterator is returned from FilterAuditLogged and is used to iterate over the raw logs and unpacked data for AuditLogged events raised by the Blockchain contract.
type BlockchainAuditLoggedIterator struct {
	Event *BlockchainAuditLogged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BlockchainAuditLoggedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockchainAuditLogged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BlockchainAuditLogged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BlockchainAuditLoggedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockchainAuditLoggedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockchainAuditLogged represents a AuditLogged event raised by the Blockchain contract.
type BlockchainAuditLogged struct {
	EventId   *big.Int
	ItemId    *big.Int
	Action    string
	Actor     common.Address
	Timestamp *big.Int
	Payload   string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuditLogged is a free log retrieval operation binding the contract event 0xc160262fc35c7fe77e722bfdfc4a0b237ed57a1f260ba0f6b26b856b496dcd2f.
//
// Solidity: event AuditLogged(uint256 indexed eventId, uint256 indexed itemId, string action, address indexed actor, uint256 timestamp, string payload)
func (_Blockchain *BlockchainFilterer) FilterAuditLogged(opts *bind.FilterOpts, eventId []*big.Int, itemId []*big.Int, actor []common.Address) (*BlockchainAuditLoggedIterator, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}
	var itemIdRule []interface{}
	for _, itemIdItem := range itemId {
		itemIdRule = append(itemIdRule, itemIdItem)
	}

	var actorRule []interface{}
	for _, actorItem := range actor {
		actorRule = append(actorRule, actorItem)
	}

	logs, sub, err := _Blockchain.contract.FilterLogs(opts, "AuditLogged", eventIdRule, itemIdRule, actorRule)
	if err != nil {
		return nil, err
	}
	return &BlockchainAuditLoggedIterator{contract: _Blockchain.contract, event: "AuditLogged", logs: logs, sub: sub}, nil
}

// WatchAuditLogged is a free log subscription operation binding the contract event 0xc160262fc35c7fe77e722bfdfc4a0b237ed57a1f260ba0f6b26b856b496dcd2f.
//
// Solidity: event AuditLogged(uint256 indexed eventId, uint256 indexed itemId, string action, address indexed actor, uint256 timestamp, string payload)
func (_Blockchain *BlockchainFilterer) WatchAuditLogged(opts *bind.WatchOpts, sink chan<- *BlockchainAuditLogged, eventId []*big.Int, itemId []*big.Int, actor []common.Address) (event.Subscription, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}
	var itemIdRule []interface{}
	for _, itemIdItem := range itemId {
		itemIdRule = append(itemIdRule, itemIdItem)
	}

	var actorRule []interface{}
	for _, actorItem := range actor {
		actorRule = append(actorRule, actorItem)
	}

	logs, sub, err := _Blockchain.contract.WatchLogs(opts, "AuditLogged", eventIdRule, itemIdRule, actorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockchainAuditLogged)
				if err := _Blockchain.contract.UnpackLog(event, "AuditLogged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAuditLogged is a log parse operation binding the contract event 0xc160262fc35c7fe77e722bfdfc4a0b237ed57a1f260ba0f6b26b856b496dcd2f.
//
// Solidity: event AuditLogged(uint256 indexed eventId, uint256 indexed itemId, string action, address indexed actor, uint256 timestamp, string payload)
func (_Blockchain *BlockchainFilterer) ParseAuditLogged(log types.Log) (*BlockchainAuditLogged, error) {
	event := new(BlockchainAuditLogged)
	if err := _Blockchain.contract.UnpackLog(event, "AuditLogged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
