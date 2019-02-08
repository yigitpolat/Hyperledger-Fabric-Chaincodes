package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Car struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Color string `json:"color"`
	Year  string `json:"year"`
	Owner string `json:"owner"`
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("error while starting SmartContract")
	}
}

func (sc *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	cars := []Car{
		Car{ID: "1", Make: "Tesla", Color: "Grey", Year: "2015", Owner: "yigit"},
		Car{ID: "2", Make: "Ferro", Color: "red", Year: "2000", Owner: "kemal"},
		Car{ID: "3", Make: "Merco", Color: "green", Year: "1995", Owner: "damla"},
		Car{ID: "4", Make: "BMW", Color: "blue", Year: "1500", Owner: "begum"},
	}

	i := 0
	for i < len(cars) {
		carAsByte, _ := json.Marshal(cars[i])
		stub.PutState("CAR"+cars[i].ID, carAsByte)
		fmt.Println("CAR"+cars[i].ID, "added, values: ", cars[i])
		i = i + 1
	}
	return shim.Success(nil)
}

func (sc *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "createCar" {
		sc.createCar(stub, args)
	} else if function == "queryAllCars" {
		sc.queryAllCars(stub)
	} else if function == "getCarHistory" {
		sc.getCarHistory(stub, args)
	} else if function == "deleteCar" {
		sc.deleteCar(stub, args)
	} else if function == "changeCarOwner" {
		sc.changeCarOwner(stub, args)
	}

	return shim.Success(nil)
}

func (sc *SmartContract) createCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 5 {
		shim.Error("Invalid number of arguments")
	}

	car := Car{ID: args[0], Make: args[1], Color: args[2], Year: args[3], Owner: args[4]}
	carAsBytes, _ := json.Marshal(car)
	stub.PutState("CAR"+args[0], carAsBytes)
	fmt.Println("CAR"+args[0], "added, values: ", car)

	return shim.Success(nil)
}

func (sc *SmartContract) queryAllCars(stub shim.ChaincodeStubInterface) peer.Response {
	startKey := "CAR0"
	endKey := "CAR99999"

	iterator, err := stub.GetStateByRange(startKey, endKey)

	if err != nil {
		shim.Error(err.Error())
	}

	iterator.Close()
	for iterator.HasNext() {
		queryResult, err := iterator.Next()
		if err != nil {
			shim.Error(err.Error())
		}
		fmt.Println("Key:", queryResult.Key, " Values:", string(queryResult.Value))
	}
	return shim.Success(nil)
}

func (sc *SmartContract) getCarHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		shim.Error("Invalid number of arguments")
	}

	iterator, err := stub.GetHistoryForKey("CAR" + args[0])

	if err != nil {
		shim.Error(err.Error())
	}

	iterator.Close()

	for iterator.HasNext() {
		queryResult, err := iterator.Next()
		if err != nil {
			shim.Error(err.Error())
		}
		fmt.Println(queryResult)
	}
	return shim.Success(nil)
}

func (sc *SmartContract) deleteCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		shim.Error("Invalid number of arguments")
	}

	err := stub.DelState("CAR" + args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (sc *SmartContract) changeCarOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		shim.Error("Invalid number of arguments")
	}

	carAsBytes, _ := stub.GetState("CAR" + args[0])
	car := Car{}
	json.Unmarshal(carAsBytes, &car)
	car.Owner = args[1]
	carAsBytes, _ = json.Marshal(car)
	stub.PutState("CAR"+args[0], carAsBytes)
	return shim.Success(nil)
}
