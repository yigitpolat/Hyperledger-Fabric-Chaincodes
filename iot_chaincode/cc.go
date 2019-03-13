package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//SmartContract asd
type SmartContract struct {
}

//Value asd
type Value struct {
	SensorID string `json:"sensorID"`
	Temp     string `json:"temp"`
	Time     string `json:"time"`
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("Error when starting SmartContract", err)
	}
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Chaincode instantiated")
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "registerSensor" {
		return s.registerSensor(stub, args)
	} else if function == "addTemp" {
		return s.addTemp(stub, args)
	} else if function == "getHistory" {
		return s.getHistory(stub, args)
	}
	return shim.Success(nil)
}

func (s *SmartContract) registerSensor(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("invalid number of arguments")
	}
	sensorID := args[0]
	value := Value{SensorID: args[0], Temp: " ", Time: " "}
	valueAsBytes, _ := json.Marshal(value)
	stub.PutState(sensorID, valueAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) addTemp(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("invalid number of arguments")
	}
	sensorID := args[0]
	valueAsBytes, _ := stub.GetState(sensorID)
	value := Value{}
	json.Unmarshal(valueAsBytes, &value)
	fmt.Println(value)

	value.Temp = args[1]
	value.Time = args[2]
	fmt.Println(value)
	valueByBytes, _ := json.Marshal(value)
	stub.PutState(sensorID, valueByBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("invalid number of arguments")
	}
	sensorID := args[0]
	iterator, _ := stub.GetHistoryForKey(sensorID)
	iterator.Close()
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			shim.Error(err.Error())
		}
		fmt.Println(queryResponse)
	}
	return shim.Success(nil)
}
