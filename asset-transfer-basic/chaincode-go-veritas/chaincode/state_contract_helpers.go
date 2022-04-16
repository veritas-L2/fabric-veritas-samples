package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//TODO: handle errors

func InitStateContract(ctx contractapi.TransactionContextInterface) {
	ctx.GetStub().InvokeChaincode("state-contract", [][]byte{[]byte("InitStateContract")}, "l2")
}

func ReleaseStateContract(ctx contractapi.TransactionContextInterface) {
	ctx.GetStub().InvokeChaincode("state-contract", [][]byte{[]byte("ReleaseStateContract")}, "l2")
}

func PutState(ctx contractapi.TransactionContextInterface, key string, value []byte) error {
	queryArgs := [][]byte{[]byte("PutState"), []byte(key), value}

	response := ctx.GetStub().InvokeChaincode("state-contract", queryArgs, "l2")
	if response.Status != shim.OK {
		return fmt.Errorf("failed to put to rollup state. %s", response.Payload)
	} else {
		return nil
	}
}

func GetState(ctx contractapi.TransactionContextInterface, key string) ([]byte, error) {
	queryArgs := [][]byte{[]byte("GetState"), []byte(key)}

	response := ctx.GetStub().InvokeChaincode("state-contract", queryArgs, "l2")

	if response.Status != shim.OK {
		return nil, fmt.Errorf("failed to read from rollup state. %s", response.Payload)
	} else {
		return response.Payload, nil
	}
}

func DeleteState(ctx contractapi.TransactionContextInterface, key string) error {
	queryArgs := [][]byte{[]byte("DeleteState"), []byte(key)}

	response := ctx.GetStub().InvokeChaincode("state-contract", queryArgs, "l2")

	if response.Status != shim.OK {
		return fmt.Errorf("failed to delete from rollup state. %s", response.Payload)
	} else {
		return nil
	}
}
