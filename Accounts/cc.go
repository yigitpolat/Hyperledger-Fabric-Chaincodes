package main

//import format "fmt"
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Account defined as struct
type Account struct {
	AccountNumber string `json:"AccountNumber"`
	FirstName     string `json:"FirstName"`
	Amount        int    `json:"Amount"`
	Bank          string `json:"Bank"`
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
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The 'set'
// method may create a new asset by specifying a new key-value pair.
func (smartcontract *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()
	fmt.Println()
	fmt.Println("##Invoke is running " + function + "##")

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return smartcontract.initLedger(stub)
	} else if function == "queryAccount" {
		return smartcontract.queryAccount(stub, args)
		//TODO CHECK IF ACCOUNT EXISTS BEFORE
		//TODO PLAY WITH TXID
	} else if function == "createAccount" {
		return smartcontract.createAccount(stub, args)
	} else if function == "queryAllAccounts" {
		return smartcontract.queryAllAccounts(stub)
	} else if function == "deleteAccount" {
		return smartcontract.deleteAccount(stub, args)
	} else if function == "deposit" {
		return smartcontract.deposit(stub, args)
	} else if function == "getHistory" {
		return smartcontract.getHistory(stub, args)
	}
	/*	} else if function == "deposit" {
			return smartcontract.deposit(stub, args)
		} else if function == "transfer" {
			return smartcontract.transfer(stub, args)
		}
	*/
	return shim.Error("Invalid Smart Contract function name.")
}

func (smartcontract *SmartContract) initLedger(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("=============== Start Init Ledger ===============")

	accounts := []Account{
		Account{AccountNumber: "1010", FirstName: "Elrond", Amount: 100, Bank: "akbank"},
		Account{AccountNumber: "2020", FirstName: "Arwen", Amount: 200, Bank: "teb"},
		Account{AccountNumber: "3030", FirstName: "Aragorn", Amount: 300, Bank: "isbank"},
		Account{AccountNumber: "4040", FirstName: "Legolas", Amount: 400, Bank: "finansbank"},
		Account{AccountNumber: "5050", FirstName: "Frodo", Amount: 500, Bank: "akbank"},
	}

	i := 0
	for i < len(accounts) {
		accountAsBytes, _ := json.Marshal(accounts[i])
		err := stub.PutState("ACCOUNT"+accounts[i].AccountNumber, accountAsBytes)

		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Print("Added account ACCOUNT", accounts[i].AccountNumber, "\n")
		i = i + 1
	}

	fmt.Println("=============== End Init Ledger ===============")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) queryAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Query Account ===============")

	if len(args) != 1 {
		return shim.Error("Invalid number of args")
	}

	accountID := args[0]
	accountByBytes, err := stub.GetState(accountID)

	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "AccountID: " + accountID + ", Account info: " + string(accountByBytes)
	fmt.Println(jsonResp)
	fmt.Println("=============== End Query Account ===============")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) createAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Create Account ===============")

	if len(args) != 4 {
		return shim.Error("Invalid number of args")
	}

	accountNumber := args[0]
	firstName := args[1]
	amount, err := strconv.Atoi(args[2])
	bank := args[4]

	if err != nil {
		return shim.Error(err.Error())
	}

	account := Account{AccountNumber: accountNumber, FirstName: firstName, Amount: amount, Bank: bank}

	accountAsBytes, err := json.Marshal(account)
	stub.PutState("ACCOUNT"+account.AccountNumber, accountAsBytes)
	fmt.Println("=============== End Create Account ===============")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) queryAllAccounts(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("=============== Start Query All Account ===============")

	startKey := "ACCOUNT0000"
	endKey := "ACCOUNT9999"

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
		jsonResp := "AccountID: " + queryResponse.Key + ", Account info: " + string(queryResponse.Value)
		fmt.Println(jsonResp)
	}
	fmt.Println("=============== End Query All Account ===============")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) deleteAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Delete Account ===============")

	if len(args) != 1 {
		return shim.Error("Invalid number of args")
	}

	accountID := args[0]

	err := stub.DelState(accountID)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Account with AccountID", accountID, "has deleted.")

	fmt.Println("=============== End Delete Account ===============")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) deposit(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Deposit ===============")

	if len(args) != 2 {
		return shim.Error("Invalid number of args")
	}

	accountID := args[0]
	depositAmountAsString := args[1]
	depositAmount, _ := strconv.Atoi(depositAmountAsString)

	accountByBytes, err := stub.GetState(accountID)

	if err != nil {
		return shim.Error(err.Error())
	}

	account := Account{}
	json.Unmarshal(accountByBytes, &account)
	account.Amount += depositAmount
	accountAsBytes, _ := json.Marshal(account)
	err = stub.PutState("ACCOUNT"+account.AccountNumber, accountAsBytes)

	fmt.Println(depositAmount, "has been added to account", accountID)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("=============== End Deposit ===============")
	return shim.Success(nil)
}

func (smartcontract *SmartContract) getHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	historyInterface, _ := stub.GetHistoryForKey(args[0])
	historyInterface.Close()

	for historyInterface.HasNext() {
		queryResponse, err := historyInterface.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println(queryResponse)
	}

	return shim.Success(nil)
}

/*

func (smartcontract *SmartContract) transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

}
*/
