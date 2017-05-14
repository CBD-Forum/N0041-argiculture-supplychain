package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type LoanApply struct {
	ObjectType    string
	FarmerID      string
	FarmerName    string
	Money         int
	LoanableMoney int // 可放款额度
	CreateDate    int64
}

var loanApplyPrefix = "loanApply:"
var LoanApplyObjectType = "loanApply"
var loanApplyObjIdx = "objectType~farmerID"

func insertLoanApplyToDb(stub shim.ChaincodeStubInterface, l *LoanApply) error {
	l.ObjectType = LoanApplyObjectType
	var key = loanApplyPrefix + l.FarmerID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get LoanApply from state db, err: " + err.Error())
	} else if wBytes != nil {
		return errors.New("LoanApply already exists, " + l.FarmerID)
	}

	wBytes, err = json.Marshal(l)
	if err != nil {
		return errors.New("cannot convert LoanApply to bytes: " + err.Error())
	}

	err = stub.PutState(key, wBytes)
	if err != nil {
		return errors.New("cannot put LoanApply to state db, err: " + err.Error())
	}

	// create index
	value := []byte{0x00}
	objIdx, err := stub.CreateCompositeKey(loanApplyObjIdx, []string{LoanApplyObjectType, l.FarmerID})
	if err != nil {
		return errors.New("create LoanApply composite key err, " + err.Error())
	}
	stub.PutState(objIdx, value)

	return nil
}

func getLoanApplyByFarmerFromDb(stub shim.ChaincodeStubInterface, farmerID string) (*LoanApply, error) {
	l := &LoanApply{}

	key := loanApplyPrefix + farmerID
	wBytes, err := stub.GetState(key)
	if err != nil {
		return l, errors.New("cannot get LoanApply from state db, err: " + err.Error())
	} else if wBytes == nil {
		return l, nil
	}

	err = json.Unmarshal(wBytes, l)
	if err != nil {
		return l, errors.New("cannot convert json bytes to LoanApply, err: " + err.Error())
	}
	return l, nil
}

func delLoanApply(stub shim.ChaincodeStubInterface, farmerID string) error {
	key := loanApplyPrefix + farmerID

	err := stub.PutState(key, nil)
	if err != nil {
		return err
	}

	objIdx, err := stub.CreateCompositeKey(loanApplyObjIdx, []string{LoanApplyObjectType, farmerID})
	if err != nil {
		return errors.New("create LoanApply composite key err, " + err.Error())
	}

	err = stub.PutState(objIdx, nil)
	if err != nil {
		return err
	}
	return nil
}

func getLoanApplyListFromDb(stub shim.ChaincodeStubInterface) ([]*LoanApply, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(loanApplyObjIdx, []string{LoanApplyObjectType})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get Logistic LoanApply by ownerID, err:" + err.Error())
	}

	return _getLoanApplyFromIdxIterator(stub, iterator)
}

//从iterator获取列表
func _getLoanApplyFromIdxIterator(stub shim.ChaincodeStubInterface, iterator shim.StateQueryIteratorInterface) ([]*LoanApply, error) {
	ws := make([]*LoanApply, 0)
	for iterator.HasNext() {
		key, _, err := iterator.Next()
		if err != nil {
			return nil, errors.New("get key from LoanApply err, " + err.Error())
		}
		_, compositeKeys, err := stub.SplitCompositeKey(key)
		if err != nil {
			return nil, errors.New("canot get compositeKeys, err " + err.Error())
		}

		farmerID := compositeKeys[1]
		w, err := getLoanApplyByFarmerFromDb(stub, farmerID)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}
	return ws, nil
}
