package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Order struct {
	ObjectType  string
	OrderID     string
	OwnerID     string
	FarmerID    string
	Title       string
	CategoryID  string
	Destination string
	Amount      int
	Price       int
	Detail      string
	Status      int
	Payed       bool
	CreateDate  int64
}

type FilterOrdersFunc func(order *Order) bool

const (
	OrderStatusInit            = 1 //订单初始化
	OrderStatusMatched         = 2 //订单被接单
	OrderStatusProduced        = 3 //订单生产完成
	OrderStatusPackaged        = 4 //订单打包
	OrderStatusPackagedStoreIn = 5 //订单入库
	OrderStatusDeliveried      = 6 //订单发货
	OrderStatusFinished        = 7
)

var orderObjectType = "order"
var orderPrefix = "order:"
var orderObjIdx = "objectType~orderID"
var orderOwnerIdx = "ownerID~orderID"
var orderFarmerIdx = "farmerID~orderID"
var orderCategoryIdx = "categoryID~orderID"

func insertOrderToDb(stub shim.ChaincodeStubInterface, order *Order) error {
	var key = orderPrefix + order.OrderID
	orderBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get order bytes from state db, err: " + err.Error())
	} else if orderBytes != nil {
		return errors.New("order already exists, id is " + order.OrderID)
	}

	orderBytes, err = json.Marshal(order)
	if err != nil {
		return errors.New("cannot convert order to bytes: " + err.Error())
	}

	err = stub.PutState(key, orderBytes)
	if err != nil {
		return errors.New("cannot put order to state db, err: " + err.Error())
	}

	value := []byte{0x00}

	objIdxKey, err := stub.CreateCompositeKey(orderObjIdx, []string{orderObjectType, order.OrderID})
	if err != nil {
		return errors.New("create order compsite key err, " + err.Error())
	}
	stub.PutState(objIdxKey, value)

	ownerKey, err := stub.CreateCompositeKey(orderOwnerIdx, []string{order.OwnerID, order.OrderID})
	if err != nil {
		return errors.New("create order compsite key err, " + err.Error())
	}
	stub.PutState(ownerKey, value)

	categoryIDKey, err := stub.CreateCompositeKey(orderCategoryIdx, []string{order.CategoryID, order.OrderID})
	if err != nil {
		return errors.New("create order compsite key err, " + err.Error())
	}
	stub.PutState(categoryIDKey, value)

	farmerKey, err := stub.CreateCompositeKey(orderFarmerIdx, []string{order.FarmerID, order.OrderID})
	if err != nil {
		return errors.New("create order compsite key err, " + err.Error())
	}
	stub.PutState(farmerKey, value)

	return nil
}

// 修改订单，不涉及到分类、农民、商户的修改
func updateOrderToDb(stub shim.ChaincodeStubInterface, order *Order, changeFarmer bool) error {
	var key = orderPrefix + order.OrderID
	orderBytes, err := stub.GetState(key)
	if err != nil {
		return errors.New("cannot get order bytes from state db, err: " + err.Error())
	} else if orderBytes == nil {
		return errors.New("order not exists, id is " + order.OrderID)
	}

	orderBytes, err = json.Marshal(order)
	if err != nil {
		return errors.New("cannot convert order to bytes: " + err.Error())
	}

	err = stub.PutState(key, orderBytes)
	if err != nil {
		return errors.New("cannot put order to state db, err: " + err.Error())
	}

	value := []byte{0x00}
	farmerKey, err := stub.CreateCompositeKey(orderFarmerIdx, []string{order.FarmerID, order.OrderID})
	if err != nil {
		return errors.New("create order compsite key err, " + err.Error())
	}
	stub.PutState(farmerKey, value)

	return nil
}

func getOrderFromDb(stub shim.ChaincodeStubInterface, orderID string) (*Order, error) {
	order := &Order{}

	var key = orderPrefix + orderID
	orderBytes, err := stub.GetState(key)
	if err != nil {
		return order, errors.New("cannot get order bytes from state db, err: " + err.Error())
	} else if orderBytes == nil {
		return order, nil
	}

	err = json.Unmarshal(orderBytes, order)
	if err != nil {
		return order, errors.New("cannot convert order bytes to order, err: " + err.Error())
	}
	return order, nil
}

// 获取所有的订单
func getOrdersFromDb(stub shim.ChaincodeStubInterface) ([]*Order, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(orderObjIdx, []string{orderObjectType})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get orders CreateCompositeKey by objectType, err:" + err.Error())
	}

	return _getOrdersFromIdxIterator(stub, iterator, nil)
}

// 获取某个商户的订单
func getOrdersByOwnerFromDb(stub shim.ChaincodeStubInterface, ownerID string) ([]*Order, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(orderOwnerIdx, []string{ownerID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get orders CreateCompositeKey by ownerID, err:" + err.Error())
	}

	return _getOrdersFromIdxIterator(stub, iterator, nil)
}

// 获取农民接的订单
func getOrdersByFarmerFromDb(stub shim.ChaincodeStubInterface, farmerID string) ([]*Order, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(orderFarmerIdx, []string{farmerID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get orders CreateCompositeKey by farmerID, err:" + err.Error())
	}

	return _getOrdersFromIdxIterator(stub, iterator, nil)
}

// 获取所有的单子，以Status
func getOrdersByStatusFromDb(stub shim.ChaincodeStubInterface, status int) ([]*Order, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(orderObjIdx, []string{orderObjectType})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get orders CreateCompositeKey by categoryID, err:" + err.Error())
	}

	var f = func(o *Order) bool {
		if o.Status == status {
			return true
		}
		return false
	}
	return _getOrdersFromIdxIterator(stub, iterator, f)
}

// 获取某个类别的订单
func getOrdersByCategoryIDFromDb(stub shim.ChaincodeStubInterface, categoryID string) ([]*Order, error) {
	iterator, err := stub.GetStateByPartialCompositeKey(orderCategoryIdx, []string{categoryID})
	defer iterator.Close()
	if err != nil {
		return nil, errors.New("cannot get orders CreateCompositeKey by categoryID, err:" + err.Error())
	}

	return _getOrdersFromIdxIterator(stub, iterator, nil)
}

//从iterator获取order 列表
func _getOrdersFromIdxIterator(stub shim.ChaincodeStubInterface, iterator shim.StateQueryIteratorInterface, f FilterOrdersFunc) ([]*Order, error) {
	orders := make([]*Order, 0)
	for iterator.HasNext() {
		key, _, err := iterator.Next()
		if err != nil {
			return nil, errors.New("get key from orders by farmerID err, " + err.Error())
		}
		_, compositeKeys, err := stub.SplitCompositeKey(key)
		if err != nil {
			return nil, errors.New("canot get compositeKeys, err " + err.Error())
		}

		orderID := compositeKeys[1]
		order, err := getOrderFromDb(stub, orderID)
		if err != nil {
			return nil, err
		}
		if f != nil {
			if f(order) {
				orders = append(orders, order)
			}
		} else {
			orders = append(orders, order)
		}
	}
	return orders, nil
}
