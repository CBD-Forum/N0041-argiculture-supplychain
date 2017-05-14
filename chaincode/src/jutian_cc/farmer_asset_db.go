package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type FarmerAsset struct {
	ObjectType       string
	AssetID          string
	OwnerID          string
	CategoryID       string
	Amount           int
	MaterialUnitCost int
	LaborUnitCost    int
	CreateDate       int64
	FarmerName       string
	Location         string
}

var farmerAssetMaterialUnitCost = 2 // 预设的成本
var farmerAssetLaborlUnitCost = 1   // 预设的成本

var farmerAssetObjectType = "farmerAsset"

var farmerAssetPrefix = "farmerAsset:"

var farmerAssetObjectTypeIdx = "objectType~assetID"
var farmerAssetOwnerIdx = "ownerID~assetID"

var farmerAssetCatNOwnIdx = "ownerID~categoryID" // key是如此。保存的是assetID

func insertFarmerAssetToDb(stub shim.ChaincodeStubInterface, asset *FarmerAsset) error {
	asset.CreateDate = MakeTimestamp()
	if asset.AssetID == "" {
		asset.AssetID = GenerateRandom(32)
	}

	var key = farmerAssetPrefix + asset.AssetID
	assetJSONBytes, err := json.Marshal(asset)
	if err != nil {
		return errors.New("cannot convert farmer asset to json bytes, err is " + err.Error())
	}
	err = stub.PutState(key, assetJSONBytes)
	if err != nil {
		return errors.New("cannot put farmer asset to state, err is " + err.Error())
	}

	objectIdx, err := stub.CreateCompositeKey(farmerAssetObjectTypeIdx, []string{farmerAssetObjectType, asset.AssetID})
	if err != nil {
		return errors.New("create item composite key err, " + err.Error())
	}
	value := []byte{0x00}
	stub.PutState(objectIdx, value)

	ownIdx, err := stub.CreateCompositeKey(farmerAssetOwnerIdx, []string{asset.OwnerID, asset.AssetID})
	if err != nil {
		return errors.New("create item composite key err, " + err.Error())
	}
	stub.PutState(ownIdx, value)

	// 保存类别。目前每个人每个分类只有一种
	catOwnKey := asset.OwnerID + "~" + asset.CategoryID
	stub.PutState(catOwnKey, []byte(asset.AssetID))

	return nil
}

//只更新部分数据
func updateFarmerAssetToDb(stub shim.ChaincodeStubInterface, asset *FarmerAsset) error {
	var key = farmerAssetPrefix + asset.AssetID
	assetJSONBytes, err := json.Marshal(asset)
	if err != nil {
		return errors.New("cannot convert farmer asset to json bytes, err is " + err.Error())
	}
	err = stub.PutState(key, assetJSONBytes)
	if err != nil {
		return errors.New("cannot put farmer asset to state, err is " + err.Error())
	}
	return nil
}

func getAssetByCategoryNOwnFromDb(stub shim.ChaincodeStubInterface, categoryID, ownerID string) (*FarmerAsset, error) {
	asset := &FarmerAsset{}

	catOwnKey := ownerID + "~" + categoryID
	assetID, err := stub.GetState(catOwnKey)
	if err != nil {
		return asset, errors.New("cannot find asset by category and ownerId " + ownerID + "," + categoryID + ",err: " + err.Error())
	}
	if assetID == nil {
		return asset, nil
	}
	return getFarmerAssetFromDb(stub, string(assetID))
}

func getFarmerAssetFromDb(stub shim.ChaincodeStubInterface, assetID string) (*FarmerAsset, error) {
	farmerAsset := &FarmerAsset{}
	var key = farmerAssetPrefix + assetID
	assBytes, err := stub.GetState(key)
	if err != nil {
		return farmerAsset, errors.New("cannot get farmer asset from state, err is " + err.Error())
	} else if assBytes == nil {
		return farmerAsset, nil
	}

	err = json.Unmarshal(assBytes, farmerAsset)
	if err != nil {
		return farmerAsset, errors.New("cannot convert farmer asset json bytes to Asset, err is " + err.Error())
	}
	return farmerAsset, nil
}

// 获取某个农民的所有资产
func getFarmerAssetsByOwnerFromDb(stub shim.ChaincodeStubInterface, ownerID string) ([]*FarmerAsset, error) {
	ownIterator, err := stub.GetStateByPartialCompositeKey(farmerAssetOwnerIdx, []string{ownerID})
	defer ownIterator.Close()
	if err != nil {
		return nil, errors.New("cannot get assets by owner, err " + err.Error())
	}

	assets := make([]*FarmerAsset, 0)
	var i int
	for i = 0; ownIterator.HasNext(); i++ {
		ownerKey, _, err := ownIterator.Next()
		if err != nil {
			return assets, errors.New("get owner key from asset err, " + err.Error())
		}
		_, compositeKeys, err := stub.SplitCompositeKey(ownerKey)
		if err != nil {
			return assets, errors.New("canot get compositeKeys, err " + err.Error())
		}

		assetID := compositeKeys[1]
		asset, err := getFarmerAssetFromDb(stub, assetID)
		if err != nil {
			return assets, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}
