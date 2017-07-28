package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("storecc")

// StoreCC
type StoreCC struct {
}

// Init
func (t *StoreCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *StoreCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_cc0 Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "load" {
		// queries an entity state
		return t.load(stub, args)
	} else if function == "save" {
		return t.save(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

// Query callback representing the query of a chaincode
func (t *StoreCC) load(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := string(Avalbytes)
	logger.Infof("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

type saveRsp struct {
	Code int
	Hash string
	TxID string
}

// Query callback representing the query of a chaincode
func (t *StoreCC) save(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	key := args[0]
	value := args[1]
	err := stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error("get failed ")
	}
	var rsp saveRsp
	rsp.Code = 0
	rsp.TxID = stub.GetTxID()
	b, _ := json.Marshal(rsp)

	return shim.Success(b)
}

func main() {
	fmt.Println("******  -- main storecc -- *******")
	err := shim.Start(new(StoreCC))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
