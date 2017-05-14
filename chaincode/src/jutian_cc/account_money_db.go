package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type AccountMoney struct {
	OwnerID    string
	Money      int
	FlowIDs    []string
	FlowType   int
	CreateDate int64
}

type AccountMoneyHistory struct {
	AccountMoney *AccountMoney
	Flows        []interface{}
}

type AccountBalance struct {
	OwnerID string
	Loan    int
}

var accountBalancePrefix = "accountBalance:"

const (
	AccountMoneyFlowTypeOrder       = 1
	AccountMoneyFlowTypeBigPackage  = 2
	AccountMoneyFlowTypeWarehouse   = 3
	AccountMoneyFlowTypeLogistic    = 4
	AccountMoneyFlowTypeMerchandise = 5
	AccountMoneyFlowTypeLoan        = 6
	AccountMoneyFlowTypeRepayment   = 7
	AccountMoneyFlowTypeAsset       = 8
)

var moneyKeyPrefix = "usermoney:"

func insertMoneyToDb(stub shim.ChaincodeStubInterface, accountMoney *AccountMoney) error {
	var key = moneyKeyPrefix + accountMoney.OwnerID
	accountMoney.CreateDate = MakeTimestamp()
	moneyBytes, err := json.Marshal(accountMoney)
	if err != nil {
		return errors.New("cannot convert user money to json bytes, err " + err.Error())
	}

	err = stub.PutState(key, moneyBytes)
	if err != nil {
		return errors.New("cannot insert user money to state db, err " + err.Error())
	}

	return nil
}

func getMoneyFromDb(stub shim.ChaincodeStubInterface, ownerID string) (*AccountMoney, error) {
	var key = moneyKeyPrefix + ownerID
	var money = &AccountMoney{}
	moneyBytes, err := stub.GetState(key)
	if err != nil {
		return money, errors.New("cannot get account money, err " + err.Error())
	} else if moneyBytes == nil {
		return money, nil
	}

	err = json.Unmarshal(moneyBytes, money)
	if err != nil {
		return money, errors.New("cannot convert bytes to usermoney, err " + err.Error() + ", data is " + string(moneyBytes) + "key is " + key)
	}
	return money, nil
}

func getMoneyHistory(stub shim.ChaincodeStubInterface, ownerID string) ([]*AccountMoney, error) {
	moneyHistory := make([]*AccountMoney, 0)
	var key = moneyKeyPrefix + ownerID
	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return moneyHistory, errors.New("cannot get money history, err " + err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		_, historicValue, err := resultsIterator.Next()
		if err != nil {
			return moneyHistory, errors.New("cannot get money history from iterator, err " + err.Error())
		}

		var money AccountMoney
		err = json.Unmarshal(historicValue, &money)
		if err != nil {
			return moneyHistory, errors.New("cannot convert money history to UserMoney, err " + err.Error())
		}
		moneyHistory = append(moneyHistory, &money)
	}
	return moneyHistory, nil
}

func insertAccountBalanceToDb(stub shim.ChaincodeStubInterface, balance *AccountBalance) error {
	var key = accountBalancePrefix + balance.OwnerID
	moneyBytes, err := json.Marshal(balance)
	if err != nil {
		return errors.New("cannot convert AccountBalance to json bytes, err " + err.Error())
	}

	err = stub.PutState(key, moneyBytes)
	if err != nil {
		return errors.New("cannot insert AccountBalance to state db, err " + err.Error())
	}

	return nil
}

func getAccountBalanceFromDb(stub shim.ChaincodeStubInterface, ownerID string) (*AccountBalance, error) {
	balance := &AccountBalance{}
	var key = accountBalancePrefix + ownerID
	moneyBytes, err := stub.GetState(key)
	if err != nil {
		return balance, errors.New("cannot get AccountBalance from state db, err " + err.Error())
	} else if moneyBytes == nil {
		return balance, nil
	}

	err = json.Unmarshal(moneyBytes, balance)
	if err != nil {
		return balance, errors.New("cannot convert bytes to AccountBalance, err " + err.Error() + ", data is " + string(moneyBytes) + "key is " + key)
	}
	return balance, nil
}
