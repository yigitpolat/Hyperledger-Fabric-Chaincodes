package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
	"time"
)

type NotaryApp struct {
}

type User struct {
	UserID    string
	UserName  string
	UserSname string
	Amount    int
}

type UserList struct {
	UserIDs []string
}

type Asset struct {
	AssetID   string
	AssetType string
	UserID    string
}

type AssetList struct {
		AssetIDs []string
}

type UserResponse struct {
	UserID    string
	UserName  string
	UserSname string
	Amount    int
	AssetList []Asset
}

func (t *NotaryApp) Init(stub shim.ChaincodeStubInterface) peer.Response {
	var err error

	assetList := AssetList{AssetIDs: nil}

	assetListAsBytes, _ := json.Marshal(assetList)

	err = stub.PutState("assetlist", assetListAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Asset List is created.")

	userList := UserList{UserIDs: nil}

	userListAsBytes, _ := json.Marshal(userList)

	err = stub.PutState("userlist", userListAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("User List is created.")

	return shim.Success(nil)
}

func (t *NotaryApp) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "adduser" {
		return t.addUser(stub, args)
	} else if function == "deleteuser" {
		return t.deleteUser(stub, args)
	} else if function == "getuser" {
		return t.getUser(stub, args)
	} else if function == "getuserhistory" {
		return t.getUserHistory(stub, args)
	} else if function == "addasset" {
		return t.addAsset(stub, args)
	} else if function == "exchangeasset" {
		return t.exchangeAsset(stub, args)
	} else if function == "getasset" {
		return t.getAsset(stub, args)
	} else if function == "addAmount" {
		return t.Deposit(stub, args)
	} else if function == "getassethistory" {
		return t.getAssetHistory(stub, args)
	} else if function == "getallusers" {
		return t.getAllUsers(stub, args)
	} else if function == "getallassets" {
		return t.getAllAssets(stub, args)
	} else if function == "deleteasset" {
		return t.deleteAsset(stub, args)
	}
	return shim.Error("Invalid function name.")
}

func (t *NotaryApp) addUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("=============== Start Add User ===============")

	var UserID string
	var UserName string
	var UserSname string
	var Amount int
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	UserID = args[0]
	UserName = args[1]
	UserSname = args[2]
	Amount, err = strconv.Atoi(args[3])

	var user = User{UserID: UserID, UserName: UserName, UserSname: UserSname, Amount: Amount}

	userAsBytes, _ := json.Marshal(user)

	err = stub.PutState(UserID, userAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	userListAsBytes, err := stub.GetState("userlist")

	if userListAsBytes == nil {
		return shim.Error("Result is null")
	}

	var userList UserList
	err = json.Unmarshal(userListAsBytes, &userList)

	if err != nil {
		fmt.Println("There was an error: ", err)
	}

	userList.UserIDs = append(userList.UserIDs, UserID)

	userListAsBytes, _ = json.Marshal(userList)

	err = stub.PutState("userlist", userListAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("New user is created.")
	return shim.Success(nil)

}

func (t *NotaryApp) deleteUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Delete User ===============")

	if len(args) != 1 {
		return shim.Error("Invalid number of args")
	}

	UserID := args[0]

	err := stub.DelState(UserID)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("User with User ID: ", UserID, " has deleted.")

	fmt.Println("=============== End Delete User ===============")
	return shim.Success(nil)
}

func (t *NotaryApp) getUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var userID string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID of the user to getUser")
	}

	userID = args[0]

	// Get the state from the ledger
	UserAsBytes, err := stub.GetState(userID)
	AssetListAsBytes, err := stub.GetState("assetlist")

	assets := make([]Asset, 0)

	var assetIDs AssetList
	err = json.Unmarshal(AssetListAsBytes, &assetIDs)


	for i := 0; i< len(assetIDs.AssetIDs); i++{

		var tempAsset Asset
		tempAssetByte, err := stub.GetState(assetIDs.AssetIDs[i])
		
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(tempAssetByte, &tempAsset)

		if tempAsset.UserID == userID {
			assets = append(assets, tempAsset)
		}
	}

	if err != nil {
		jsonResp := "Failed to get state for " + userID
		return shim.Error(jsonResp)
	}

	if UserAsBytes == nil {
		jsonResp := "Null amount for " + userID
		return shim.Error(jsonResp)
	}

	var user UserResponse
	err = json.Unmarshal(UserAsBytes, &user)

	user.AssetList = assets

	UserResponseAsByte, err := json.Marshal(user)
	fmt.Println("User with User ID: ", userID, " has been listed.")

	return shim.Success(UserResponseAsByte)
}

func (t *NotaryApp) getUserHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var userID string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	userID = args[0]

	resultsIterator, err := stub.GetHistoryForKey(userID)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (t *NotaryApp) addAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("=============== Start Add User ===============")

	var AssetID string
	var AssetType string
	var UserID string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	AssetID = args[0]
	AssetType = args[1]
	UserID = args[2]

	userValueBytes, err := stub.GetState(UserID)

	if err != nil {
		jsonResp := "Failed to get state for " + UserID
		return shim.Error(jsonResp)
	}

	if userValueBytes == nil {
		jsonResp := "Null amount for " + UserID
		return shim.Error(jsonResp)
	}

	var asset = Asset{AssetID: AssetID, AssetType: AssetType, UserID: UserID}



	assetAsBytes, _ := json.Marshal(asset)

	err = stub.PutState(AssetID, assetAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	assetListAsBytes, err := stub.GetState("assetlist")

	if assetListAsBytes == nil {
		return shim.Error("Result is null")
	}

	var assetList AssetList
	err = json.Unmarshal(assetListAsBytes, &assetList)

	if err != nil {
		fmt.Println("There was an error: ", err)
	}

	assetList.AssetIDs = append(assetList.AssetIDs, AssetID)

	assetListAsBytes, _ = json.Marshal(assetList)

	err = stub.PutState("assetlist", assetListAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("New asset is created.")
	return shim.Success(nil)


}

func (t *NotaryApp) getAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var AssetID string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID of the asset to getAsset")
	}

	AssetID = args[0]

	// Get the state from the ledger
	AssetAsBytes, err := stub.GetState(AssetID)
	if err != nil {
		jsonResp := "Failed to get state for " + AssetID
		return shim.Error(jsonResp)
	}

	if AssetAsBytes == nil {
		jsonResp := "Null amount for " + AssetID
		return shim.Error(jsonResp)
	}

	return shim.Success(AssetAsBytes)
}

func (t *NotaryApp) getAssetHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var AssetID string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	AssetID = args[0]

	resultsIterator, err := stub.GetHistoryForKey(AssetID)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (t *NotaryApp) deleteAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Delete Asset ===============")

	if len(args) != 1 {
		return shim.Error("Invalid number of args")
	}

	AssetID := args[0]

	err := stub.DelState(AssetID)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Asset with asset ID: ", AssetID, " has deleted.")

	fmt.Println("=============== End Delete Asset ===============")
	return shim.Success(nil)
}

func (t *NotaryApp) getAllUsers(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	userListAsBytes, err := stub.GetState("userlist")

	if err != nil {
		jsonResp := "Failed to get state for user list"
		return shim.Error(jsonResp)
	}

	if userListAsBytes == nil {
		jsonResp := "Null amount for user list"
		return shim.Error(jsonResp)
	}

	var userList UserList
	err = json.Unmarshal(userListAsBytes, &userList)

	if err != nil {
		fmt.Println("There was an error:", err)
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for _, value := range userList.UserIDs {

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		userAsByte, err := stub.GetState(value)

		if err != nil {
			jsonResp := "Failed to get state for user"
			return shim.Error(jsonResp)
		}

		if userAsByte == nil {
			jsonResp := "Null amount for user"
			return shim.Error(jsonResp)
		}

		buffer.Write(userAsByte)

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

func (t *NotaryApp) getAllAssets(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	assetListAsBytes, err := stub.GetState("assetlist")

	if err != nil {
		jsonResp := "Failed to get state for asset list"
		return shim.Error(jsonResp)
	}

	if assetListAsBytes == nil {
		jsonResp := "Null amount for asset list"
		return shim.Error(jsonResp)
	}

	var assetList AssetList
	err = json.Unmarshal(assetListAsBytes, &assetList)

	if err != nil {
		fmt.Println("There was an error:", err)
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for _, value := range assetList.AssetIDs {

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		assetAsByte, err := stub.GetState(value)

		if err != nil {
			jsonResp := "Failed to get state for asset"
			return shim.Error(jsonResp)
		}

		if assetAsByte == nil {
			jsonResp := "Null amount for asset"
			return shim.Error(jsonResp)
		}

		buffer.Write(assetAsByte)

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

func (t *NotaryApp) Deposit(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Deposit ===============")

	if len(args) != 2 {
		return shim.Error("Invalid number of args")
	}

	UserID := args[0]
	depositAmountAsString := args[1]
	depositAmount, _ := strconv.Atoi(depositAmountAsString)

	UserByBytes, err := stub.GetState(UserID)

	if err != nil {
		return shim.Error(err.Error())
	}

	user := User{}
	err = json.Unmarshal(UserByBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	user.Amount += depositAmount
	UserByBytes, _ = json.Marshal(user)
	err = stub.PutState("ACCOUNT"+user.UserID, UserByBytes)

	fmt.Println(depositAmount, "has been added to account", UserID)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("=============== End Deposit ===============")
	return shim.Success(nil)
}

func (t *NotaryApp) exchangeAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("=============== Start Exchange ===============")

	if len(args) != 4 {
		return shim.Error("Invalid number of args")
	}

	User1ID := args[0]
	User2ID := args[1]
	AssetID := args[2]
	AmountAsString := args[3]
	Amount, _ := strconv.Atoi(AmountAsString)

	User1ByBytes, err := stub.GetState(User1ID)
	User2ByBytes, err := stub.GetState(User2ID)
	AssetAsBytes, err := stub.GetState(AssetID)

	if err != nil {
		return shim.Error(err.Error())
	}


	var user1 User
	err = json.Unmarshal(User1ByBytes, &user1)

	var user2 User
	err = json.Unmarshal(User2ByBytes, &user2)

	var asset Asset
	err = json.Unmarshal(AssetAsBytes, &asset)


	user1.Amount += Amount
	user2.Amount -= Amount

	asset.UserID = user2.UserID


	User1ByBytes, _ = json.Marshal(user1)
	User2ByBytes, _ = json.Marshal(user2)
	AssetAsBytes, _ = json.Marshal(asset)

	err = stub.PutState(user1.UserID, User1ByBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(user2.UserID, User2ByBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(asset.AssetID, AssetAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("=============== End of Exchange ===============")
	return shim.Success(nil)


}

func main() {
	err := shim.Start(new(NotaryApp))
	if err != nil {
		fmt.Printf("Error starting User Chaincode: %s", err)
	}
}