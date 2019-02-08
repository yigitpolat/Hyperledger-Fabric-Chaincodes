package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//import format "fmt"

// Voter defined as struct
type Voter struct {
	NationalID       string `json:"CandidateID"`
	Name             string `json:"Name"`
	VotedCandidateID string `json:"VotedCandidateID"`
}

// Candidate defined as struct
type Candidate struct {
	CandidateID string `json:"CandidateID"`
	Name        string `json:"Name"`
	TotalVote   int    `json:"TotalVote"`
}

// SmartContract defined as struct
type SmartContract struct {
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data, so be careful to avoid a scenario where you
// inadvertently clobber your ledger's data!
func (smartcontract *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("=============== Start Init ===============")

	candidates := []Candidate{
		Candidate{CandidateID: "100", Name: "Finis Valorum", TotalVote: 0},
		Candidate{CandidateID: "200", Name: "Palpatine", TotalVote: 0},
		Candidate{CandidateID: "300", Name: "Bail Antilles", TotalVote: 0},
	}

	i := 0
	for i < len(candidates) {
		candidateAsBytes, _ := json.Marshal(candidates[i])
		err := stub.PutState("CANDIDATE"+candidates[i].CandidateID, candidateAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Print("Added candidate CANDIDATE", candidates[i].CandidateID, "Candidate name:", candidates[i].Name, "\n")
		i = i + 1
	}
	fmt.Println("=============== End Init  ===============")
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The 'set'
// method may create a new asset by specifying a new key-value pair.
func (smartcontract *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()
	fmt.Println()

	if function == "registerVoter" {
		return smartcontract.registerVoter(stub, args)
	} else if function == "queryVoter" {
		return smartcontract.queryVoter(stub, args)
	} else if function == "queryAllVoters" {
		return smartcontract.queryAllVoters(stub)
	} else if function == "queryCandidate" {
		return smartcontract.queryCandidate(stub, args)
	} else if function == "queryAllCandidates" {
		return smartcontract.queryAllCandidates(stub)
	} else if function == "addVote" {
		return smartcontract.addVote(stub, args)
	} else if function == "getHistory" {
		return smartcontract.getHistory(stub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (smartcontract *SmartContract) registerVoter(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Register Voter =============== ")

	if len(args) != 2 {
		return shim.Error("Invalid number of arguments.")
	}

	nationalID := args[0]
	name := args[1]

	voter := Voter{NationalID: nationalID, Name: name, VotedCandidateID: ""}
	voterAsBytes, _ := json.Marshal(voter)
	err := stub.PutState("VOTER"+nationalID, voterAsBytes)
	if err != nil {
		return shim.Error("Voter already registered")
	}
	fmt.Print("Added voter. VOTER", voter.NationalID, ", Voter info", voter)
	fmt.Println("=============== End Register Voter =============== ")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) queryVoter(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Query Voter =============== ")

	if len(args) != 1 {
		return shim.Error("Invalid number of arguments.")
	}

	nationalID := args[0]
	voterAsBytes, err := stub.GetState("VOTER" + nationalID)

	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "VoterID: VOTER" + nationalID + ", Voter info: " + string(voterAsBytes)
	fmt.Println(jsonResp)

	fmt.Println("=============== End Query Voter =============== ")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) queryAllVoters(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("=============== Start Query All Voters =============== ")

	startKey := "VOTER0"
	endKey := "VOTER99999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		jsonResp := "VoterID: " + queryResponse.Key + ", Voter info: " + string(queryResponse.Value)
		fmt.Println(jsonResp)
	}

	fmt.Println("=============== End Query All Voters =============== ")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) queryCandidate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Query Candidate =============== ")

	if len(args) != 1 {
		return shim.Error("Invalid number of arguments.")
	}

	candidateID := args[0]
	candidateAsBytes, err := stub.GetState("CANDIDATE" + candidateID)

	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "CandidateID: CANDIDATE" + candidateID + ", Candidate info: " + string(candidateAsBytes)
	fmt.Println(jsonResp)

	fmt.Println("=============== End Query Candidate =============== ")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) queryAllCandidates(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("=============== Start Query All Candidates =============== ")

	startKey := "CANDIDATE0"
	endKey := "CANDIDATE99999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		jsonResp := "CandidateID: " + queryResponse.Key + ", Candidate info: " + string(queryResponse.Value)
		fmt.Println(jsonResp)
	}

	fmt.Println("=============== End Query All Voters =============== ")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) addVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Add Vote =============== ")

	if len(args) != 2 {
		return shim.Error("Invalid number of arguments.")
	}

	voterNationalID := args[0]
	candidateID := args[1]

	voterAsBytes, err := stub.GetState("VOTER" + voterNationalID)
	candidateAsBytes, err := stub.GetState("CANDIDATE" + candidateID)

	if err != nil {
		return shim.Error(err.Error())
	}

	voter := Voter{}
	candidate := Candidate{}

	json.Unmarshal(voterAsBytes, &voter)
	json.Unmarshal(candidateAsBytes, &candidate)

	if voter.VotedCandidateID != "" {
		return shim.Error("Voter already voted a candidate")
	}

	voter.VotedCandidateID = candidateID
	candidate.TotalVote++

	voterByBytes, _ := json.Marshal(voter)
	candidateByBytes, _ := json.Marshal(candidate)

	err = stub.PutState("VOTER"+voter.NationalID, voterByBytes)
	err = stub.PutState("CANDIDATE"+candidate.CandidateID, candidateByBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("=============== End Add Vote =============== ")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) getHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Get History =============== ")

	candidateID := "CANDIDATE" + args[0]
	historyInterface, _ := stub.GetHistoryForKey(candidateID)
	historyInterface.Close()

	for historyInterface.HasNext() {
		queryResponse, err := historyInterface.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println(queryResponse)
	}

	fmt.Println("=============== End Get History =============== ")
	return shim.Success(nil)
}
