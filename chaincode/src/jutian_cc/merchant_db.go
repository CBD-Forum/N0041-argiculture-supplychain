package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Merchandise struct {
	ObjectType string
	PackageID  string
	OwnerID    string
	CategoryID string
	Amount     int
	Price      int
	SoldOut    bool
	CreateDate int64
}

var merchandisePrice = 50

var merchandiseObjectType = "merchandise"

var merchandisePrefix = "merchandise:"

var merchandiseObjIdx = "objectType~packageID"
var merchandiseOwnIdx = "ownerID~packageID"

func insertMerchandiseToDb(stub shim.ChaincodeStubInterface, m *Merchandise) error {
	var key = merchandisePrefix + m.PackageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get Merchandise from state db, err: " + err.Error())
	} else if wBytes != nil {
		return errors.New("Merchandise already exists, " + m.PackageID)
	}

	wBytes, err = json.Marshal(m)
	if err != nil {
		return errors.New("cannot convert Merchandise to bytes: " + err.Error())
	}

	err = stub.PutState(key, wBytes)
	if err != nil {
		return errors.New("cannot put Merchandise to state db, err: " + err.Error())
	}

	// create index
	value := []byte{0x00}

	objIdx, err := stub.CreateCompositeKey(merchandiseObjIdx, []string{merchandiseObjectType, m.PackageID})
	if err != nil {
		return errors.New("create Merchandise composite key err, " + err.Error())
	}
	stub.PutState(objIdx, value)

	ownIdx, err := stub.CreateCompositeKey(merchandiseOwnIdx, []string{m.OwnerID, m.PackageID})
	if err != nil {
		return errors.New("create Merchandise composite key err, " + err.Error())
	}
	stub.PutState(ownIdx, value)

	return nil
}

func updateMerchandiseToDb(stub shim.ChaincodeStubInterface, m *Merchandise) error {
	var key = merchandisePrefix + m.PackageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get Merchandise from state db, err: " + err.Error())
	} else if wBytes == nil {
		return errors.New("Merchandise not exists, " + m.PackageID)
	}

	wBytes, err = json.Marshal(m)
	if err != nil {
		return errors.New("cannot convert Merchandise to bytes: " + err.Error())
	}

	err = stub.PutState(key, wBytes)
	if err != nil {
		return errors.New("cannot put Merchandise to state db, err: " + err.Error())
	}

	return nil
}

func getMerchandiseFromDb(stub shim.ChaincodeStubInterface, packageID string) (*Merchandise, error) {
	m := &Merchandise{}

	var key = merchandisePrefix + packageID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return m, errors.New("cannot get Merchandise from state db, err: " + err.Error())
	} else if wBytes == nil {
		return m, nil
	}

	err = json.Unmarshal(wBytes, m)
	if err != nil {
		return m, errors.New("cannot convert json bytes to Merchandise, err: " + err.Error())
	}
	return m, nil
}

func getMerchandisesByOwnerFromDb(stub shim.ChaincodeStubInterface, ownerID string) ([]*Merchandise, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(merchandiseOwnIdx, []string{ownerID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get Merchandise CreateCompositeKey by ownerID, err:" + err.Error())
	}

	return _getMerchandisesFromIdxIterator(stub, iterator)
}

//从iterator获取 Merchandise 列表
func _getMerchandisesFromIdxIterator(stub shim.ChaincodeStubInterface, iterator shim.StateQueryIteratorInterface) ([]*Merchandise, error) {
	ws := make([]*Merchandise, 0)
	for iterator.HasNext() {
		key, _, err := iterator.Next()
		if err != nil {
			return nil, errors.New("get key from Merchandise err, " + err.Error())
		}
		_, compositeKeys, err := stub.SplitCompositeKey(key)
		if err != nil {
			return nil, errors.New("canot get compositeKeys, err " + err.Error())
		}

		packageID := compositeKeys[1]
		w, err := getMerchandiseFromDb(stub, packageID)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}
	return ws, nil
}
