package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Logistic struct {
	ObjectType    string
	BigPackageID  string
	OwnerID       string
	OrderID       string
	StartLocation string
	EndLocation   string
	Cost          int
	CreateDate    int64
	Delivered     bool
}

var logisticCost = 2

var logisticObjectType = "logistic"

var logisticPrefix = "logistic:"

var logisticOwnIdx = "ownerID~logisticID"
var logisticObjIdx = "objectType~logisticID"
var logisticOrderIdx = "orderID~logisticID"

func insertLogisticToDb(stub shim.ChaincodeStubInterface, l *Logistic) error {
	var key = logisticPrefix + l.BigPackageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get Logistic from state db, err: " + err.Error())
	} else if wBytes != nil {
		return errors.New("Logistic already exists, " + l.BigPackageID)
	}

	wBytes, err = json.Marshal(l)
	if err != nil {
		return errors.New("cannot convert Logistic to bytes: " + err.Error())
	}

	err = stub.PutState(key, wBytes)
	if err != nil {
		return errors.New("cannot put Logistic to state db, err: " + err.Error())
	}

	// create index
	value := []byte{0x00}

	ownIdx, err := stub.CreateCompositeKey(logisticOwnIdx, []string{l.OwnerID, l.BigPackageID})
	if err != nil {
		return errors.New("create Logistic composite key err, " + err.Error())
	}
	stub.PutState(ownIdx, value)

	objIdx, err := stub.CreateCompositeKey(logisticObjIdx, []string{logisticObjectType, l.BigPackageID})
	if err != nil {
		return errors.New("create Logistic composite key err, " + err.Error())
	}
	stub.PutState(objIdx, value)

	objOrdIdx, err := stub.CreateCompositeKey(logisticOrderIdx, []string{l.OrderID, l.BigPackageID})
	if err != nil {
		return errors.New("create Logistic composite key err, " + err.Error())
	}
	stub.PutState(objOrdIdx, value)

	return nil
}

func updateLogisticToDb(stub shim.ChaincodeStubInterface, l *Logistic) error {
	var key = logisticPrefix + l.BigPackageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get Logistic from state db, err: " + err.Error())
	} else if wBytes == nil {
		return errors.New("Logistic not exists, " + l.BigPackageID)
	}

	wBytes, err = json.Marshal(l)
	if err != nil {
		return errors.New("cannot convert Logistic to bytes: " + err.Error())
	}

	err = stub.PutState(key, wBytes)
	if err != nil {
		return errors.New("cannot put Logistic to state db, err: " + err.Error())
	}
	return nil
}

func getLogisticFromDb(stub shim.ChaincodeStubInterface, bigPackageID string) (*Logistic, error) {
	l := &Logistic{}

	var key = logisticPrefix + bigPackageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return l, errors.New("cannot get Logistic from state db, err: " + err.Error())
	} else if wBytes == nil {
		return l, nil
	}

	err = json.Unmarshal(wBytes, l)
	if err != nil {
		return l, errors.New("cannot convert json bytes to Logistic, err: " + err.Error())
	}
	return l, nil
}

func getLogisticByOwnFromDb(stub shim.ChaincodeStubInterface, ownerID string) ([]*Logistic, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(logisticOwnIdx, []string{ownerID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get Logistic CreateCompositeKey by ownerID, err:" + err.Error())
	}

	return _getLogisticsFromIdxIterator(stub, iterator)
}

func getLogisticByOrderFromDb(stub shim.ChaincodeStubInterface, orderID string) ([]*Logistic, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(logisticOrderIdx, []string{orderID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get Logistic CreateCompositeKey by ownerID, err:" + err.Error())
	}

	return _getLogisticsFromIdxIterator(stub, iterator)
}

//从iterator获取warehouse store in 列表
func _getLogisticsFromIdxIterator(stub shim.ChaincodeStubInterface, iterator shim.StateQueryIteratorInterface) ([]*Logistic, error) {
	ws := make([]*Logistic, 0)
	for iterator.HasNext() {
		key, _, err := iterator.Next()
		if err != nil {
			return nil, errors.New("get key from Logistic err, " + err.Error())
		}
		_, compositeKeys, err := stub.SplitCompositeKey(key)
		if err != nil {
			return nil, errors.New("canot get compositeKeys, err " + err.Error())
		}

		bigPackageID := compositeKeys[1]
		w, err := getLogisticFromDb(stub, bigPackageID)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}
	return ws, nil
}
