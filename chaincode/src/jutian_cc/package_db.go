package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//PackageUnit 植物包装，最终给到消费者的就是这个单位
type PackageUnit struct {
	ObjectType   string
	PackageID    string
	BigPackageID string
	AssetID      string
	Amount       int
	CreateDate   int64
}

var packgeUnitObjectType = "packageUnit"
var packageBigIdx = "bigPackagID~packageID"
var packageObjIdx = "objectType~packageID"
var pakcageUnitPrefix = "packageUnit:"

// 目前的打包策略
// 小的package，每个为10
// 大的bigPackage，每个为100个小的Package

//BigPackage 物流过程中用到的包装单位
type BigPackage struct {
	ObjectType   string
	BigPackageID string
	OwnerID      string
	OrderID      string
	Amount       int
	Cost         int
	Delivered    bool
	CreateDate   int64
}

var bigPackageCost = 2

//BigPackageWithUnits 物流过程中用到的包装单位，包括小的包装
type BigPackageWithUnits struct {
	BigPackage   *BigPackage
	PackageUnits []*PackageUnit
}

var bigPackageObjectType = "bigPackage"
var bigPackageObjIdx = "objectType~bigPackageID"
var bigPackageOrderIdx = "orderID~bigPackageID"
var bigPackageOwnIdx = "OwnerID~bigPackageID"

var bigPackagePrefix = "bigPackage:"

func insertBigPackageToDb(stub shim.ChaincodeStubInterface, p *BigPackage) error {
	p.CreateDate = MakeTimestamp()
	var key = bigPackagePrefix + p.BigPackageID
	pBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get package from state db, err: " + err.Error())
	} else if pBytes != nil {
		return errors.New("package already exists, " + p.BigPackageID)
	}

	pBytes, err = json.Marshal(p)
	if err != nil {
		return errors.New("cannot convert package to bytes: " + err.Error())
	}

	err = stub.PutState(key, pBytes)
	if err != nil {
		return errors.New("cannot put package to state db, err: " + err.Error())
	}

	value := []byte{0x00}
	// 创建orderId的idx
	orderIdx, err := stub.CreateCompositeKey(bigPackageOrderIdx, []string{p.OrderID, p.BigPackageID})
	if err != nil {
		return errors.New("create big package composite key err, " + err.Error())
	}
	stub.PutState(orderIdx, value)

	objIdx, err := stub.CreateCompositeKey(bigPackageObjIdx, []string{bigPackageObjectType, p.BigPackageID})
	if err != nil {
		return errors.New("create big package composite key err, " + err.Error())
	}
	stub.PutState(objIdx, value)

	ownIdx, err := stub.CreateCompositeKey(bigPackageOwnIdx, []string{p.OwnerID, p.BigPackageID})
	if err != nil {
		return errors.New("create big package composite key err, " + err.Error())
	}
	stub.PutState(ownIdx, value)

	return nil
}

func updatePackageToDb(stub shim.ChaincodeStubInterface, p *BigPackage) error {
	var key = bigPackagePrefix + p.BigPackageID
	pBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get package from state db, err: " + err.Error())
	} else if pBytes == nil {
		return errors.New("package not exists, " + p.BigPackageID)
	}

	pBytes, err = json.Marshal(p)
	if err != nil {
		return errors.New("cannot convert package to bytes: " + err.Error())
	}

	err = stub.PutState(key, pBytes)
	if err != nil {
		return errors.New("cannot put package to state db, err: " + err.Error())
	}
	return nil
}

func getBigPackageFromDb(stub shim.ChaincodeStubInterface, bigPackageID string) (*BigPackage, error) {
	bigPackag := &BigPackage{}

	var key = bigPackagePrefix + bigPackageID
	pBytes, err := stub.GetState(key)
	if err != nil {
		return bigPackag, errors.New("cannot get package from state db, err: " + err.Error())
	} else if pBytes == nil {
		return bigPackag, nil
	}

	err = json.Unmarshal(pBytes, bigPackag)
	if err != nil {
		return bigPackag, errors.New("cannot convert json bytes to BigPackage, err: " + err.Error())
	}
	return bigPackag, nil
}

func getBigPackagesByOrderFromDb(stub shim.ChaincodeStubInterface, orderID string) ([]*BigPackage, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(bigPackageOrderIdx, []string{orderID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get big package CreateCompositeKey by categoryID, err:" + err.Error())
	}

	return _getBigPackagesFromIdxIterator(stub, iterator)
}

func getBigPackagesByOwnFromDb(stub shim.ChaincodeStubInterface, ownerID string) ([]*BigPackage, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(bigPackageOwnIdx, []string{ownerID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get big package CreateCompositeKey by categoryID, err:" + err.Error())
	}

	return _getBigPackagesFromIdxIterator(stub, iterator)
}

func insertPackageToDb(stub shim.ChaincodeStubInterface, p *PackageUnit) error {
	p.CreateDate = MakeTimestamp()
	var key = pakcageUnitPrefix + p.PackageID
	pBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get package unit from state db, err: " + err.Error())
	} else if pBytes != nil {
		return errors.New("package unit already exists, " + p.PackageID)
	}

	pBytes, err = json.Marshal(p)
	if err != nil {
		return errors.New("cannot convert package unit to bytes: " + err.Error())
	}

	err = stub.PutState(key, pBytes)
	if err != nil {
		return errors.New("cannot put package unit to state db, err: " + err.Error())
	}

	value := []byte{0x00}
	// 创建orderId的idx
	bigIdx, err := stub.CreateCompositeKey(packageBigIdx, []string{p.BigPackageID, p.PackageID})
	if err != nil {
		return errors.New("create package composite key err, " + err.Error())
	}
	stub.PutState(bigIdx, value)

	objIdx, err := stub.CreateCompositeKey(packageObjIdx, []string{packgeUnitObjectType, p.PackageID})
	if err != nil {
		return errors.New("create package composite key err, " + err.Error())
	}
	stub.PutState(objIdx, value)

	return nil
}

func getPackageFromDb(stub shim.ChaincodeStubInterface, packageID string) (*PackageUnit, error) {
	packageUnit := &PackageUnit{}

	var key = pakcageUnitPrefix + packageID
	pBytes, err := stub.GetState(key)
	if err != nil {
		return packageUnit, errors.New("cannot get package unit from state db, err: " + err.Error())
	} else if pBytes == nil {
		return packageUnit, errors.New("cannot find package unit: " + packageID)
	}

	json.Unmarshal(pBytes, packageUnit)
	if err != nil {
		return packageUnit, errors.New("cannot convert json bytes to Package, err: " + err.Error())
	}
	return packageUnit, nil
}

func getPackageUnitsByBigPackageFromDb(stub shim.ChaincodeStubInterface, bigPackageID string) ([]*PackageUnit, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(packageBigIdx, []string{bigPackageID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get package units CreateCompositeKey by categoryID, err:" + err.Error())
	}

	return _getPackageUnitsFromIdxIterator(stub, iterator)
}

//从iterator获取packageUnit 列表
func _getPackageUnitsFromIdxIterator(stub shim.ChaincodeStubInterface, iterator shim.StateQueryIteratorInterface) ([]*PackageUnit, error) {
	packages := make([]*PackageUnit, 0)
	for iterator.HasNext() {
		key, _, err := iterator.Next()
		if err != nil {
			return nil, errors.New("get key from packages by farmerID err, " + err.Error())
		}
		_, compositeKeys, err := stub.SplitCompositeKey(key)
		if err != nil {
			return nil, errors.New("canot get compositeKeys, err " + err.Error())
		}

		packageID := compositeKeys[1]
		packageUnit, err := getPackageFromDb(stub, packageID)
		if err != nil {
			return nil, err
		}
		packages = append(packages, packageUnit)
	}
	return packages, nil
}

//从iterator获取big package 列表
func _getBigPackagesFromIdxIterator(stub shim.ChaincodeStubInterface, iterator shim.StateQueryIteratorInterface) ([]*BigPackage, error) {
	bigPackages := make([]*BigPackage, 0)
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
		bigPackage, err := getBigPackageFromDb(stub, bigPackageID)
		if err != nil {
			return nil, err
		}
		bigPackages = append(bigPackages, bigPackage)
	}
	return bigPackages, nil
}
