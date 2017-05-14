package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type WarehouseStoreIn struct {
	ObjectType    string
	BigPackageID  string
	OwnerID       string
	OrderID       string
	CreateDate    int64
	Cost          int
	WarehouseName string
	Sent          bool
}

var warehouseCost = 3

var warehouseStoreInObjectType = "warehouseStoreIn"
var warehousePrefix = "warehouseStoreIn:"
var warehouseOrderIdx = "orderID~bigPackageID"
var warehouseObjIdx = "objectType~bigPackageID"
var warehouseOwnIdx = "ownerID~bigPackageID"

func insertWarehouseStoreInToDb(stub shim.ChaincodeStubInterface, w *WarehouseStoreIn) error {
	w.CreateDate = MakeTimestamp()
	var key = warehousePrefix + w.BigPackageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get warehouse storein from state db, err: " + err.Error())
	} else if wBytes != nil {
		return errors.New("warehouse storein already exists, " + w.BigPackageID)
	}

	wBytes, err = json.Marshal(w)
	if err != nil {
		return errors.New("cannot convert warehouse storein to bytes: " + err.Error())
	}

	err = stub.PutState(key, wBytes)
	if err != nil {
		return errors.New("cannot put warehouse storein to state db, err: " + err.Error())
	}

	// create index
	value := []byte{0x00}

	objIdx, err := stub.CreateCompositeKey(warehouseObjIdx, []string{warehouseStoreInObjectType, w.BigPackageID})
	if err != nil {
		return errors.New("create warehouse storein composite key err, " + err.Error())
	}
	stub.PutState(objIdx, value)

	orderIdx, err := stub.CreateCompositeKey(warehouseOrderIdx, []string{w.OrderID, w.BigPackageID})
	if err != nil {
		return errors.New("create warehouse storein composite key err, " + err.Error())
	}
	stub.PutState(orderIdx, value)

	ownIdx, err := stub.CreateCompositeKey(warehouseOwnIdx, []string{w.OwnerID, w.BigPackageID})
	if err != nil {
		return errors.New("create warehouse storein composite key err, " + err.Error())
	}
	stub.PutState(ownIdx, value)

	return nil
}

func updateWarehouseStoreInToDb(stub shim.ChaincodeStubInterface, w *WarehouseStoreIn) error {
	var key = warehousePrefix + w.BigPackageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get warehouse storein from state db, err: " + err.Error())
	} else if wBytes == nil {
		return errors.New("warehouse storein not exists, " + w.BigPackageID)
	}

	wBytes, err = json.Marshal(w)
	if err != nil {
		return errors.New("cannot convert warehouse storein to bytes: " + err.Error())
	}

	err = stub.PutState(key, wBytes)
	if err != nil {
		return errors.New("cannot put warehouse storein to state db, err: " + err.Error())
	}
	return nil
}

func getWarehouseStoreInFromDb(stub shim.ChaincodeStubInterface, bigPackageID string) (*WarehouseStoreIn, error) {
	w := &WarehouseStoreIn{}

	var key = warehousePrefix + bigPackageID
	pBytes, err := stub.GetState(key)
	if err != nil {
		return w, errors.New("cannot get warehouse storein from state db, err: " + err.Error())
	} else if pBytes == nil {
		return w, nil
	}

	err = json.Unmarshal(pBytes, w)
	if err != nil {
		return w, errors.New("cannot convert json bytes to warehouse storein, err: " + err.Error())
	}
	return w, nil
}

func getWareshouseStoreInsByOrderFromDb(stub shim.ChaincodeStubInterface, orderID string) ([]*WarehouseStoreIn, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(warehouseOrderIdx, []string{orderID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get warehouse storein CreateCompositeKey by categoryID, err:" + err.Error())
	}

	return _getWarehouseStoreinsFromIdxIterator(stub, iterator)
}

func getWareshouseStoreInsByOwnFromDb(stub shim.ChaincodeStubInterface, ownerID string) ([]*WarehouseStoreIn, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(warehouseOwnIdx, []string{ownerID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get warehouse storein CreateCompositeKey by categoryID, err:" + err.Error())
	}

	return _getWarehouseStoreinsFromIdxIterator(stub, iterator)
}

//从iterator获取warehouse store in 列表
func _getWarehouseStoreinsFromIdxIterator(stub shim.ChaincodeStubInterface, iterator shim.StateQueryIteratorInterface) ([]*WarehouseStoreIn, error) {
	ws := make([]*WarehouseStoreIn, 0)
	for iterator.HasNext() {
		key, _, err := iterator.Next()
		if err != nil {
			return nil, errors.New("get key from bigPackages by farmerID err, " + err.Error())
		}
		_, compositeKeys, err := stub.SplitCompositeKey(key)
		if err != nil {
			return nil, errors.New("canot get compositeKeys, err " + err.Error())
		}

		bigPackageID := compositeKeys[1]
		w, err := getWarehouseStoreInFromDb(stub, bigPackageID)
		if err != nil {
			return nil, err
		}
		if w.BigPackageID != "" {
			ws = append(ws, w)
		}
	}
	return ws, nil
}
