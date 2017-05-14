package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"errors"
)

type Account struct {
	AccountID string
	UserName  string
	Password  string
	UserType  int
}

type UserInfo struct {
	ObjectType  string
	Type        int
	OwnerID     string
	RealName    string
	IdentityID  string
	Location    string
	UserType    int
	CreditScore int
}

const (
	UserTypeFarmer    = 1
	UserTypePackager  = 2
	UserTypeWarehouse = 3
	UserTypeLogistic  = 4
	UserTypeMerchant  = 5
)

//AccountInfo 账户信息
type AccountInfo struct {
	UserInfo    *UserInfo
	Temperature int
	Equipment   string
	Money       *AccountMoney
	Balance     *AccountBalance
	LoanApply   *LoanApply
}

const (
	AccountTypeFarmer = 1
)

var accountObjectType = "AccountObjectType"
var userInfoIDPrefix = "userInfo:"
var accountPrefix = "account-password:"

func insertAccountToDb(stub shim.ChaincodeStubInterface, acc *Account) error {
	key := accountPrefix + acc.AccountID
	accountBytes, err := json.Marshal(&acc)
	if err != nil {
		return errors.New("cannot convert account to json bytes: " + err.Error())
	}
	err = stub.PutState(key, accountBytes)
	if err != nil {
		return errors.New("cannot save account to state db: " + err.Error())
	}

	userNameKey := accountPrefix + acc.UserName
	err = stub.PutState(userNameKey, accountBytes)
	if err != nil {
		return errors.New("cannot save account to state db: " + err.Error())
	}

	return nil
}

func getAccountFromDb(stub shim.ChaincodeStubInterface, AccountID string) (*Account, error) {
	key := accountPrefix + AccountID
	acc := &Account{}

	accountBytes, err := stub.GetState(key)
	if err != nil {
		return acc, errors.New("cannot get account from state db, key is " + key + ", err:" + err.Error())
	} else if accountBytes == nil {
		return acc, nil
	}

	err = json.Unmarshal(accountBytes, acc)
	if err != nil {
		return acc, errors.New("cannot convert bytes to account: " + err.Error())
	}

	return acc, nil
}

func getAccountByUserNameFromDb(stub shim.ChaincodeStubInterface, userName string) (*Account, error) {
	key := accountPrefix + userName
	acc := &Account{}

	accountBytes, err := stub.GetState(key)
	if err != nil {
		return acc, errors.New("cannot get account from state db, key is " + key + ", err:" + err.Error())
	} else if accountBytes == nil {
		return acc, nil
	}

	err = json.Unmarshal(accountBytes, acc)
	if err != nil {
		return acc, errors.New("cannot convert bytes to account: " + err.Error())
	}

	return acc, nil
}

func insertUserInfoToDb(stub shim.ChaincodeStubInterface, acc *UserInfo) error {
	key := userInfoIDPrefix + acc.OwnerID
	accountBytes, err := json.Marshal(&acc)
	if err != nil {
		return errors.New("cannot convert account to json bytes: " + err.Error())
	}

	fmt.Println(acc)
	err = stub.PutState(key, accountBytes)
	if err != nil {
		return errors.New("cannot save account to state db: " + err.Error())
	}
	return nil
}

func getUserInfoFromDb(stub shim.ChaincodeStubInterface, ownerID string) (*UserInfo, error) {
	key := userInfoIDPrefix + ownerID
	acc := &UserInfo{}

	accountBytes, err := stub.GetState(key)
	if err != nil {
		return acc, errors.New("cannot get account from state db, key is " + key + ", err:" + err.Error())
	} else if accountBytes == nil {
		return acc, nil
	}

	err = json.Unmarshal(accountBytes, acc)
	if err != nil {
		return acc, errors.New("cannot convert bytes to account: " + err.Error())
	}

	return acc, nil
}

// 温度

var temperaturePrefix = "account_temperature:"

func insertTemperatureToDb(stub shim.ChaincodeStubInterface, accountID string, temperature int) error {
	var key = temperaturePrefix + accountID
	var temperatureString = strconv.Itoa(temperature)
	err := stub.PutState(key, []byte(temperatureString))
	if err != nil {
		return errors.New("cannnt put temperatore to state db: " + err.Error())
	}
	return nil
}

func getTemperatureFromDb(stub shim.ChaincodeStubInterface, accountID string) (int, error) {
	var key = temperaturePrefix + accountID
	resBytes, err := stub.GetState(key)
	if err != nil {
		return 0, errors.New("cannot get temperature from state db, err: " + err.Error())
	} else if resBytes == nil {
		stub.PutState(key, []byte("0"))
		return 0, nil
	}
	return strconv.Atoi(string(resBytes))
}

const (
	EquipmentStatusFault     = "故障"
	EquipmentStatusException = "异常"
	EquipmentStatusGOOD      = "良好"
)

var equipmentPrefix = "account_equipment:"

func insertEquipmentToDb(stub shim.ChaincodeStubInterface, accountID string, status string) error {
	var key = temperaturePrefix + accountID
	err := stub.PutState(key, []byte(status))
	if err != nil {
		return errors.New("cannnt put equipment status to state db: " + err.Error())
	}
	return nil
}

func getEquipmentFromDb(stub shim.ChaincodeStubInterface, accountID string) (string, error) {
	var key = temperaturePrefix + accountID
	resBytes, err := stub.GetState(key)
	if err != nil {
		return "未知", errors.New("cannot get equipment status from state db, err: " + err.Error())
	} else if resBytes == nil {
		return "未知", nil
	}
	return string(resBytes), nil
}
