package main

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//MerchandiseFlow 物品溯源链
type MerchandiseFlow struct {
	Merchandise *Merchandise
	Logistic    *Logistic
	Warehouse   *WarehouseStoreIn
	BigPackage  *BigPackage
	FarmerAsset *FarmerAsset
}

type JutianChaincode struct {
}

func (t *JutianChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### jutain_cc Init ###########")
	return shim.Success(nil)

}

func (t *JutianChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("########### invoke function: " + function + " &&&###########")
	if function == "createAccount" {
		if len(args) != 4 {
			return wrapToPbResponse(nil, errors.New("len of args should be 4"))
		}
		accountID := args[0]
		userName := args[1]
		password := args[2]
		userType, err := strconv.Atoi(args[3])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[3]+" cannot convert to i, "+err.Error()))
		}
		return t.createAccount(stub, accountID, userName, password, userType)
	} else if function == "checkAccountExist" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		userName := args[0]
		return t.checkAccountExist(stub, userName)
	} else if function == "checkPassword" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		userName := args[0]
		password := args[1]
		return t.checkPassword(stub, userName, password)
	} else if function == "getAccount" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		accountID := args[0]
		return t.getAccount(stub, accountID)
	} else if function == "initUserInfo" {
		if len(args) != 5 {
			return wrapToPbResponse(nil, errors.New("len of args should be 5"))
		}
		ownerID := args[0]
		name := args[1]
		identity := args[2]
		location := args[3]
		userType, err := strconv.Atoi(args[4])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[4]+" cannot convert to i, "+err.Error()))
		}
		return t.initUserInfo(stub, ownerID, name, identity, location, userType)
	} else if function == "getUserInfo" {
		ownerID := args[0]
		return t.getUserInfo(stub, ownerID)
	} else if function == "changeTemperature" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		ownerID := args[0]
		temperature, err := strconv.Atoi(args[1])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[1]+" cannot convert to i, "+err.Error()))
		}
		return t.changeTemperature(stub, ownerID, temperature)
	} else if function == "getMoneyHistory" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getMoneyHistory(stub, ownerID)
	} else if function == "changeEquipmentStatus" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		ownerID := args[0]
		status := args[1]
		return t.changeEquipmentStatus(stub, ownerID, status)
	} else if function == "insertFarmerAsset" {
		if len(args) != 5 {
			return wrapToPbResponse(nil, errors.New("len of args should be 5"))
		}
		ownerID := args[0]
		categoryID := args[1]
		amount, err := strconv.Atoi(args[2])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[2]+" cannot convert to i, "+err.Error()))
		}
		materialUnitCost, err := strconv.Atoi(args[3])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[3]+" cannot convert to i, "+err.Error()))
		}
		laborUnitCost, err := strconv.Atoi(args[4])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[4]+" cannot convert to i, "+err.Error()))
		}
		return t.insertFarmerAsset(stub, ownerID, categoryID, amount, materialUnitCost, laborUnitCost)
	} else if function == "getFarmerAsset" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		assetID := args[0]
		return t.getFarmerAsset(stub, assetID)
	} else if function == "getAssetByCategoryIDNOwner" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		ownerID := args[0]
		categoryID := args[1]
		return t.getAssetByCategoryIDNOwner(stub, ownerID, categoryID)
	} else if function == "getFarmerAssetsByOwner" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getFarmerAssetsByOwner(stub, ownerID)
	} else if function == "importOrderData" {
		// 这里先处理给农户的信用加分
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.importFarmerOrderData(stub, ownerID)
	} else if function == "getFarmerLoanableMoney" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getFarmerLoanableMoney(stub, ownerID)
	} else if function == "loanByFarmer" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		ownerID := args[0]
		loan, err := strconv.Atoi(args[1])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[1]+" cannot convert to i, "+err.Error()))
		}
		return t.loanByFarmer(stub, ownerID, loan)
	} else if function == "getLoanApplyList" {
		return t.getLoanApplyList(stub)
	} else if function == "getReceiveableOrdersByFarmer" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getReceiveableOrdersByFarmer(stub, ownerID)
	} else if function == "loanApplyApprove" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		farmerID := args[0]
		return t.loanApplyApprove(stub, farmerID)
	} else if function == "insertOrder" {
		// 商户生成订单
		if len(args) != 7 {
			return wrapToPbResponse(nil, errors.New("len of args should be 7"))
		}
		ownerID := args[0]
		title := args[1]
		categoryID := args[2]
		destination := args[3]
		amount, err := strconv.Atoi(args[4])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[4]+" cannot convert to i, "+err.Error()))
		}
		cost, err := strconv.Atoi(args[5])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[5]+" cannot convert to i, "+err.Error()))
		}
		detail := args[6]
		return t.insertOrder(stub, ownerID, title, categoryID, destination, amount, cost, detail)
	} else if function == "getOrder" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		orderID := args[0]
		return t.getOrder(stub, orderID)
	} else if function == "getOrdersByOwner" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getOrdersByOwner(stub, ownerID)
	} else if function == "getOrdersByFarmer" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		farmerID := args[0]
		return t.getOrdersByFarmer(stub, farmerID)
	} else if function == "getOrdersByStatus" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		status, err := strconv.Atoi(args[0])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[0]+" cannot convert to i, "+err.Error()))
		}
		if status != OrderStatusDeliveried &&
			status != OrderStatusFinished &&
			status != OrderStatusInit &&
			status != OrderStatusMatched &&
			status != OrderStatusPackaged &&
			status != OrderStatusPackagedStoreIn &&
			status != OrderStatusProduced {
			return wrapToPbResponse(nil, errors.New(args[1]+" status error "))
		}
		return t.getOrdersByStatus(stub, status)
	} else if function == "payOrder" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		ownerID := args[0]
		orderID := args[1]
		return t.payOrder(stub, ownerID, orderID)
	} else if function == "matchOrderByFarmer" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		farmerID := args[0]
		orderID := args[1]
		return t.matchOrderByFarmer(stub, farmerID, orderID)
	} else if function == "produceOrderByFarmer" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		farmerID := args[0]
		orderID := args[1]
		return t.produceOrderByFarmer(stub, farmerID, orderID)
	} else if function == "packageOrder" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		packagerID := args[0]
		orderID := args[1]
		return t.packageOrder(stub, packagerID, orderID)
	} else if function == "getBigPackagesByOrder" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		packagerID := args[0]
		orderID := args[1]
		return t.getBigPackagesByOrder(stub, packagerID, orderID)
	} else if function == "getBigPackagesByOwner" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getBigPackagesByOwner(stub, ownerID)
	} else if function == "getBigPackageByOrderStatus" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		status, err := strconv.Atoi(args[0])
		if err != nil {
			return wrapToPbResponse(nil, errors.New(args[0]+" cannot convert to i, "+err.Error()))
		}
		if status != OrderStatusDeliveried &&
			status != OrderStatusFinished &&
			status != OrderStatusInit &&
			status != OrderStatusMatched &&
			status != OrderStatusPackaged &&
			status != OrderStatusPackagedStoreIn &&
			status != OrderStatusProduced {
			return wrapToPbResponse(nil, errors.New(args[1]+" status error "))
		}
		return t.getBigPackagesByOrderStatus(stub, status)
	} else if function == "insertWarehouseStoreIn" {
		if len(args) != 3 {
			return wrapToPbResponse(nil, errors.New("len of args should be 3"))
		}
		ownerID := args[0]
		bigPackageID := args[1]
		wName := args[2]
		return t.insertWarehouseStoreIn(stub, ownerID, bigPackageID, wName, nil)
	} else if function == "insertWarehouseStoreInList" {
		if len(args) < 3 {
			return wrapToPbResponse(nil, errors.New("len of args should be 3"))
		}
		ownerID := args[0]
		wName := args[1]
		bigPackageIDs := args[2:]
		return t.insertWarehouseStoreInList(stub, ownerID, wName, bigPackageIDs)
	} else if function == "getWarehouseByOrder" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		orderID := args[0]
		return t.getWarehouseStoreInByOrder(stub, orderID)
	} else if function == "getWarehouseByOwner" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getStoresByOwner(stub, ownerID)
	} else if function == "insertLogistic" {
		if len(args) != 3 {
			return wrapToPbResponse(nil, errors.New("len of args should be 3"))
		}

		ownerID := args[0]
		bigPackageID := args[1]
		truck := args[2]

		return t.insertLogistic(stub, ownerID, bigPackageID, truck, nil)
	} else if function == "insertLogisticList" {
		if len(args) < 3 {
			return wrapToPbResponse(nil, errors.New("len of args should be 3"))
		}

		ownerID := args[0]
		truck := args[1]
		bigPackageIDs := args[2:]
		return t.insertLogisticList(stub, ownerID, truck, bigPackageIDs)
	} else if function == "getLogisticByOwner" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getLogisticByOwner(stub, ownerID)
	} else if function == "getLogisticByOrder" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		orderID := args[0]
		return t.getLogisticByOrder(stub, orderID)
	} else if function == "insertMechandise" {
		if len(args) != 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		ownerID := args[0]
		bigPackageID := args[1]
		return t.insertMechandise(stub, ownerID, bigPackageID, nil)
	} else if function == "insertMechandiseList" {
		if len(args) < 2 {
			return wrapToPbResponse(nil, errors.New("len of args should be 2"))
		}
		ownerID := args[0]
		bigPackageIDs := args[1:]
		return t.insertMechandiseList(stub, ownerID, bigPackageIDs)
	} else if function == "getMerchandiseByOwner" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		ownerID := args[0]
		return t.getMerchandiseByOwner(stub, ownerID)
	} else if function == "getMerchandiseFlow" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		packageID := args[0]
		return t.getMerchandiseFlow(stub, packageID)
	} else if function == "purchasePackage" {
		if len(args) != 1 {
			return wrapToPbResponse(nil, errors.New("len of args should be 1"))
		}
		packageID := args[0]
		return t.purchasePackage(stub, packageID)
	}
	return shim.Error("Receive unknown invoke function name -" + function + "'")
}

// 用户密码
func (t *JutianChaincode) createAccount(stub shim.ChaincodeStubInterface, accountID, userName, password string, userType int) pb.Response {
	acc := &Account{
		accountID,
		userName,
		password,
		userType,
	}
	err := insertAccountToDb(stub, acc)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse("OK", nil)
}

func (t *JutianChaincode) checkAccountExist(stub shim.ChaincodeStubInterface, userName string) pb.Response {
	account, err := getAccountByUserNameFromDb(stub, userName)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if account.UserName == "" {
		return wrapToPbResponse(nil, errors.New("not found"))
	}
	return wrapStructToPbResponse("OK", nil)
}

func (t *JutianChaincode) getAccount(stub shim.ChaincodeStubInterface, accountID string) pb.Response {
	account, err := getAccountFromDb(stub, accountID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if account.UserName == "" {
		return wrapToPbResponse(nil, errors.New("not found"))
	}
	account.Password = ""
	return wrapStructToPbResponse(account, nil)
}

func (t *JutianChaincode) checkPassword(stub shim.ChaincodeStubInterface, userName, password string) pb.Response {
	account, err := getAccountByUserNameFromDb(stub, userName)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if account.UserName == "" {
		return wrapToPbResponse(nil, errors.New("not found"))
	}
	if account.Password != password {
		return wrapToPbResponse(nil, errors.New("password error"))
	}
	account.Password = ""
	return wrapStructToPbResponse(account, nil)
}

// 账户相关
func (t *JutianChaincode) initUserInfo(stub shim.ChaincodeStubInterface, ownerID, name, identity, location string, userType int) pb.Response {
	account := &UserInfo{
		accountObjectType,
		AccountTypeFarmer,
		ownerID,
		name,
		identity,
		location,
		userType,
		600,
	}
	err := insertUserInfoToDb(stub, account)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	money, err := getMoneyFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if money.OwnerID == "" {
		// 初始化金额为1000
		money := &AccountMoney{
			ownerID,
			1000,
			nil,
			0,
			MakeTimestamp(),
		}
		err = insertMoneyToDb(stub, money)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
	}

	balance, err := getAccountBalanceFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if balance.OwnerID == "" {
		balance.OwnerID = ownerID
		balance.Loan = 0
		insertAccountBalanceToDb(stub, balance)
	}

	return wrapStructToPbResponse(account, nil)
}

func (t *JutianChaincode) getUserInfo(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	account, err := getUserInfoFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	temperature, err := getTemperatureFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	equipment, err := getEquipmentFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	money, err := getMoneyFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	balance, err := getAccountBalanceFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if balance.OwnerID == "" {
		balance.OwnerID = ownerID
		balance.Loan = 0
		err = insertAccountBalanceToDb(stub, balance)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
	}
	loanApply, err := getLoanApplyByFarmerFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	loanableMoney, err := getFarmerLoanableMoney(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	loanApply.LoanableMoney = loanableMoney

	accountInfo := &AccountInfo{
		account,
		temperature,
		equipment,
		money,
		balance,
		loanApply,
	}
	return wrapStructToPbResponse(accountInfo, nil)
}

func (t *JutianChaincode) changeTemperature(stub shim.ChaincodeStubInterface, ownerID string, temperature int) pb.Response {
	err := insertTemperatureToDb(stub, ownerID, temperature)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(true, nil)
}

func (t *JutianChaincode) getTemperature(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	temperature, err := getTemperatureFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapToPbResponse([]byte(strconv.Itoa(temperature)), nil)
}

func (t *JutianChaincode) changeEquipmentStatus(stub shim.ChaincodeStubInterface, ownerID string, status string) pb.Response {
	if status != EquipmentStatusFault && status != EquipmentStatusException && status != EquipmentStatusGOOD {
		return wrapToPbResponse(nil, errors.New("status 错误: "+status))
	}
	err := insertEquipmentToDb(stub, ownerID, status)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(true, nil)
}

// 获取钱包历史
func (t *JutianChaincode) getMoneyHistory(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	moneys, err := getMoneyHistory(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	moneyHistories := make([]*AccountMoneyHistory, 0)
	for _, money := range moneys {
		moneyHistory := &AccountMoneyHistory{}
		moneyHistory.AccountMoney = money
		moneyHistory.Flows = make([]interface{}, 0)
		if money.FlowIDs != nil && len(money.FlowIDs) > 0 {
			for _, ID := range money.FlowIDs {
				if money.FlowType == AccountMoneyFlowTypeAsset {
					asset, err := getFarmerAssetFromDb(stub, ID)
					if err != nil {
						return wrapToPbResponse(nil, err)
					} else if asset.AssetID != "" {
						moneyHistory.Flows = append(moneyHistory.Flows, asset)
					}
				} else if money.FlowType == AccountMoneyFlowTypeOrder {
					order, err := getOrderFromDb(stub, ID)
					if err != nil {
						return wrapToPbResponse(nil, err)
					} else if order.OrderID != "" {
						moneyHistory.Flows = append(moneyHistory.Flows, order)
					}
				} else if money.FlowType == AccountMoneyFlowTypeBigPackage {
					p, err := getBigPackageFromDb(stub, ID)
					if err != nil {
						return wrapToPbResponse(nil, err)
					} else if p.BigPackageID != "" {
						moneyHistory.Flows = append(moneyHistory.Flows, p)
					}
				} else if money.FlowType == AccountMoneyFlowTypeWarehouse {
					w, err := getWarehouseStoreInFromDb(stub, ID)
					if err != nil {
						return wrapToPbResponse(nil, err)
					} else if w.BigPackageID != "" {
						moneyHistory.Flows = append(moneyHistory.Flows, w)
					}
				} else if money.FlowType == AccountMoneyFlowTypeLogistic {
					l, err := getLogisticFromDb(stub, ID)
					if err != nil {
						return wrapToPbResponse(nil, err)
					} else if l.BigPackageID != "" {
						moneyHistory.Flows = append(moneyHistory.Flows, l)
					}
				} else if money.FlowType == AccountMoneyFlowTypeMerchandise {
					m, err := getMerchandiseFromDb(stub, ID)
					if err != nil {
						return wrapToPbResponse(nil, err)
					} else if m.PackageID != "" {
						moneyHistory.Flows = append(moneyHistory.Flows, m)
					}
				}
			}
		}
		moneyHistories = append(moneyHistories, moneyHistory)
	}
	return wrapStructToPbResponse(moneyHistories, nil)
}

// 资产
func (t *JutianChaincode) insertFarmerAsset(stub shim.ChaincodeStubInterface, ownerID, categoryID string, amount, materialUnitCost, laborUnitCost int) pb.Response {
	//先判别该类别下是否已经有资产了
	existAsset, err := getAssetByCategoryNOwnFromDb(stub, categoryID, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	// if existAsset.AssetID != "" {
	// 	return wrapToPbResponse(nil, errors.New("category and owner already exists"))
	// }

	money, err := getMoneyFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	if money.OwnerID == "" {
		// 先初始化钱包
		money := &AccountMoney{
			ownerID,
			1000,
			nil,
			0,
			MakeTimestamp(),
		}
		err = insertMoneyToDb(stub, money)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
	}

	// 判断钱够不够
	if (amount * (materialUnitCost + laborUnitCost)) > money.Money {
		return wrapToPbResponse(nil, errors.New("dont have enough money, money is "+strconv.Itoa(money.Money)))
	}

	asset := &FarmerAsset{
		farmerAssetObjectType,
		GenerateRandom(32),
		ownerID,
		categoryID,
		amount,
		materialUnitCost,
		laborUnitCost,
		MakeTimestamp(),
		"",
		"",
	}
	if existAsset.AssetID != "" {
		asset.Amount += existAsset.Amount
	}
	err = insertFarmerAssetToDb(stub, asset)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	// 修改自己的钱
	money.Money -= amount * (materialUnitCost + laborUnitCost)
	money.FlowType = AccountMoneyFlowTypeAsset
	money.FlowIDs = []string{asset.AssetID}
	err = insertMoneyToDb(stub, money)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	return wrapStructToPbResponse(asset, nil)
}

// 获取单个资产
func (t *JutianChaincode) getFarmerAsset(stub shim.ChaincodeStubInterface, assetID string) pb.Response {
	asset, err := getFarmerAssetFromDb(stub, assetID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(asset, nil)
}

// 根据类别获取资产
func (t *JutianChaincode) getAssetByCategoryIDNOwner(stub shim.ChaincodeStubInterface, ownerID, categoryID string) pb.Response {
	existAsset, err := getAssetByCategoryNOwnFromDb(stub, categoryID, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if existAsset.AssetID == "" {
		return wrapToPbResponse(nil, errors.New("asset not exist"))
	}
	return wrapStructToPbResponse(existAsset, nil)
}

// 获取拥有者的所有资产
func (t *JutianChaincode) getFarmerAssetsByOwner(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	assets, err := getFarmerAssetsByOwnerFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(assets, nil)
}

// 给农民的信用加分，目前导入，加15分
func (t *JutianChaincode) importFarmerOrderData(stub shim.ChaincodeStubInterface, farmerID string) pb.Response {
	farmerUserInfo, err := getUserInfoFromDb(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	farmerUserInfo.CreditScore += 15
	err = insertUserInfoToDb(stub, farmerUserInfo)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	return wrapStructToPbResponse(farmerUserInfo, nil)
}

func (t *JutianChaincode) getFarmerLoanableMoney(stub shim.ChaincodeStubInterface, farmerID string) pb.Response {
	money, err := getFarmerLoanableMoney(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(money, nil)
}

// 农民借钱
// 改成写入借款申请，然后再核实放款
func (t *JutianChaincode) loanByFarmer(stub shim.ChaincodeStubInterface, ownerID string, loan int) pb.Response {
	money, err := getFarmerLoanableMoney(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if loan > money {
		return wrapToPbResponse(nil, errors.New("借款的金额太大了"))
	}

	// balance, err := getAccountBalanceFromDb(stub, ownerID)
	// if err != nil {
	// 	return wrapToPbResponse(nil, err)
	// } else if balance.OwnerID == "" {
	// 	balance.OwnerID = ownerID
	// 	balance.Loan = loan
	// } else {
	// 	balance.Loan += loan
	// }
	// err = insertAccountBalanceToDb(stub, balance)
	// if err != nil {
	// 	return wrapToPbResponse(nil, err)
	// }

	// userMoney, err := getMoneyFromDb(stub, ownerID)
	// if err != nil {
	// 	return wrapToPbResponse(nil, err)
	// }
	// userMoney.Money += loan
	// userMoney.FlowType = AccountMoneyFlowTypeLoan
	// err = insertMoneyToDb(stub, userMoney)
	// if err != nil {
	// 	return wrapToPbResponse(nil, err)
	// }
	// return wrapStructToPbResponse(userMoney, nil)

	farmerUserInfo, err := getUserInfoFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	loanApply := &LoanApply{
		"",
		ownerID,
		farmerUserInfo.RealName,
		loan,
		0,
		MakeTimestamp(),
	}
	err = insertLoanApplyToDb(stub, loanApply)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	loanableMoney, err := getFarmerLoanableMoney(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	loanApply.LoanableMoney = loanableMoney

	return wrapStructToPbResponse(loanApply, nil)
}

// 获取借款的列表
func (t *JutianChaincode) getLoanApplyList(stub shim.ChaincodeStubInterface) pb.Response {
	l, err := getLoanApplyListFromDb(stub)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	for _, loanApply := range l {
		loanableMoney, err := getFarmerLoanableMoney(stub, loanApply.FarmerID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
		loanApply.LoanableMoney = loanableMoney
	}
	return wrapStructToPbResponse(l, nil)
}

// 获取该农民有哪些应收单据
func (t *JutianChaincode) getReceiveableOrdersByFarmer(stub shim.ChaincodeStubInterface, farmerID string) pb.Response {
	o := make([]*Order, 0)
	orders, err := getOrdersByFarmerFromDb(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	for _, order := range orders {
		if order.Payed == false {
			o = append(o, order)
		}
	}
	return wrapStructToPbResponse(o, nil)
}

// 核实后放款
func (t *JutianChaincode) loanApplyApprove(stub shim.ChaincodeStubInterface, farmerID string) pb.Response {
	loanApply, err := getLoanApplyByFarmerFromDb(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if loanApply.FarmerID == "" {
		return wrapToPbResponse(nil, errors.New("loan apply not found, farmerID:"+farmerID))
	}
	balance, err := getAccountBalanceFromDb(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if balance.OwnerID == "" {
		balance.OwnerID = farmerID
		balance.Loan = loanApply.Money
	} else {
		balance.Loan += loanApply.Money
	}
	err = insertAccountBalanceToDb(stub, balance)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	userMoney, err := getMoneyFromDb(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	userMoney.Money += loanApply.Money
	userMoney.FlowType = AccountMoneyFlowTypeLoan
	err = insertMoneyToDb(stub, userMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	err = delLoanApply(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	return wrapStructToPbResponse(userMoney, nil)
}

// 农民获取自己可借款金额
func getFarmerLoanableMoney(stub shim.ChaincodeStubInterface, ownerID string) (int, error) {
	orders, err := getOrdersByFarmerFromDb(stub, ownerID)
	if err != nil {
		return 0, err
	}
	money := 0
	for _, order := range orders {
		if !order.Payed && order.FarmerID != "" { // 只要有单子还没付，就可以借
			money += order.Price
		}
	}

	farmerUserInfo, err := getUserInfoFromDb(stub, ownerID)
	if err != nil {
		return 0, err
	}
	r1 := float32(money)
	if farmerUserInfo.CreditScore < 500 {
		money = 0
	} else if farmerUserInfo.CreditScore < 680 {
		money = int(r1 * float32(0.4))
	} else if farmerUserInfo.CreditScore < 800 {
		money = int(r1 * float32(0.6))
	}

	balance, err := getAccountBalanceFromDb(stub, ownerID)
	if err != nil {
		return 0, err
	}

	loanApply, err := getLoanApplyByFarmerFromDb(stub, ownerID)
	if err != nil {
		return 0, err
	}

	result := money - balance.Loan - loanApply.Money
	return result, nil
}

//
// 订单相关
//

func (t *JutianChaincode) insertOrder(stub shim.ChaincodeStubInterface, ownerID, title, categoryID, destination string, amount, cost int, detail string) pb.Response {
	order := &Order{
		orderObjectType,
		GenerateRandom(32),
		ownerID,
		"",
		title,
		categoryID,
		destination,
		amount,
		cost,
		detail,
		OrderStatusInit,
		false,
		MakeTimestamp(),
	}

	err := insertOrderToDb(stub, order)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(order, nil)
}

// 获取单个order
func (t *JutianChaincode) getOrder(stub shim.ChaincodeStubInterface, orderID string) pb.Response {
	order, err := getOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if order.OrderID == "" {
		return wrapToPbResponse(nil, errors.New("cannot find order by id: "+orderID))
	}
	return wrapStructToPbResponse(order, nil)
}

// 根据owner获取订单
func (t *JutianChaincode) getOrdersByOwner(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	orders, err := getOrdersByOwnerFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(orders, nil)
}

// 根据farmer获取订单
func (t *JutianChaincode) getOrdersByFarmer(stub shim.ChaincodeStubInterface, farmerID string) pb.Response {
	orders, err := getOrdersByFarmerFromDb(stub, farmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(orders, nil)
}

// 根据status获取单子
func (t *JutianChaincode) getOrdersByStatus(stub shim.ChaincodeStubInterface, status int) pb.Response {
	orders, err := getOrdersByStatusFromDb(stub, status)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(orders, nil)
}

func (t *JutianChaincode) getBigPackagesByOrderStatus(stub shim.ChaincodeStubInterface, status int) pb.Response {
	orders, err := getOrdersByStatusFromDb(stub, status)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if status == OrderStatusPackaged {
		bigPackages := make([]*BigPackage, 0)
		if len(orders) == 0 {
			return wrapStructToPbResponse(bigPackages, nil)
		}
		for _, order := range orders {
			packages, err := getBigPackagesByOrderFromDb(stub, order.OrderID)
			if err != nil {
				return wrapToPbResponse(nil, err)
			}
			for _, currPackage := range packages {
				if !currPackage.Delivered {
					bigPackages = append(bigPackages, currPackage)
				}
			}
		}
		return wrapStructToPbResponse(bigPackages, nil)
	} else if status == OrderStatusPackagedStoreIn {
		bigPackages := make([]*WarehouseStoreIn, 0)
		if len(orders) == 0 {
			return wrapStructToPbResponse(bigPackages, nil)
		}
		for _, order := range orders {
			packages, err := getWareshouseStoreInsByOrderFromDb(stub, order.OrderID)
			if err != nil {
				return wrapToPbResponse(nil, err)
			}
			for _, currPackage := range packages {
				if !currPackage.Sent {
					bigPackages = append(bigPackages, currPackage)
				}
			}
		}
		return wrapStructToPbResponse(bigPackages, nil)
	} else if status == OrderStatusDeliveried {
		bigPackages := make([]*Logistic, 0)
		if len(orders) == 0 {
			return wrapStructToPbResponse(bigPackages, nil)
		}
		for _, order := range orders {
			packages, err := getLogisticByOrderFromDb(stub, order.OrderID)
			if err != nil {
				return wrapToPbResponse(nil, err)
			}
			for _, currPackage := range packages {
				if !currPackage.Delivered {
					bigPackages = append(bigPackages, currPackage)
				}
			}
		}
		return wrapStructToPbResponse(bigPackages, nil)
	}

	return wrapToPbResponse(nil, errors.New("不能调用这个接口，或者status错误"))
}

// 订单送到货之后，商户决定支付流程中的所有费用
func (t *JutianChaincode) payOrder(stub shim.ChaincodeStubInterface, ownerID, orderID string) pb.Response {
	order, err := getOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if order.OwnerID != ownerID {
		return wrapToPbResponse(nil, errors.New("this order is not yours"))
		// } else if order.Status != OrderStatusFinished {
		// 	return wrapToPbResponse(nil, errors.New("this order not finished yet"))
	} else if order.Payed {
		return wrapToPbResponse(nil, errors.New("this order already payed yet"))
	}

	merchantMoney, err := getMoneyFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	assetCost := order.Price

	bigPackages, err := getBigPackagesByOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	var warehouseCost = 0
	var packageCost = 0
	var logisticCost = 0

	lastWarehouse := &WarehouseStoreIn{}
	lastLogistic := &Logistic{}

	for _, bigPackage := range bigPackages {
		packageCost += bigPackage.Cost
		w, err := getWarehouseStoreInFromDb(stub, bigPackage.BigPackageID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		} else if w.BigPackageID != "" {
			lastWarehouse = w
			warehouseCost += w.Cost
		}
		l, err := getLogisticFromDb(stub, bigPackage.BigPackageID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		} else if l.BigPackageID != "" {
			lastLogistic = l
			logisticCost += l.Cost
		}
	}

	merchantMoney.Money -= assetCost
	merchantMoney.FlowIDs = []string{orderID}
	merchantMoney.FlowType = AccountMoneyFlowTypeOrder
	err = insertMoneyToDb(stub, merchantMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	farmerAccountMoney, err := getMoneyFromDb(stub, order.FarmerID)
	farmerAccountMoney.Money += assetCost
	farmerAccountMoney.FlowIDs = []string{orderID}
	farmerAccountMoney.FlowType = AccountMoneyFlowTypeOrder
	err = insertMoneyToDb(stub, farmerAccountMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	//判断是否有借款需要还
	farmerBalance, err := getAccountBalanceFromDb(stub, order.FarmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if farmerBalance.Loan > 0 {
		if farmerBalance.Loan < farmerAccountMoney.Money {
			farmerAccountMoney.Money -= farmerBalance.Loan
			farmerAccountMoney.FlowType = AccountMoneyFlowTypeRepayment

			farmerBalance.Loan = 0
		} else {
			farmerBalance.Loan -= farmerAccountMoney.Money

			farmerAccountMoney.Money = 0
			farmerAccountMoney.FlowType = AccountMoneyFlowTypeRepayment
		}
		err = insertMoneyToDb(stub, farmerAccountMoney)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
		err = insertAccountBalanceToDb(stub, farmerBalance)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
	}

	merchantMoney.Money -= packageCost
	merchantMoney.FlowIDs = make([]string, 0)
	merchantMoney.FlowType = AccountMoneyFlowTypeBigPackage
	packagerMoney, err := getMoneyFromDb(stub, bigPackages[0].OwnerID)
	packagerMoney.Money += packageCost
	packagerMoney.FlowIDs = make([]string, 0)
	for _, bigPackage := range bigPackages {
		packagerMoney.FlowIDs = append(packagerMoney.FlowIDs, bigPackage.BigPackageID)
		merchantMoney.FlowIDs = append(merchantMoney.FlowIDs, bigPackage.BigPackageID)
	}
	packagerMoney.FlowType = AccountMoneyFlowTypeBigPackage
	err = insertMoneyToDb(stub, packagerMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	err = insertMoneyToDb(stub, merchantMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	merchantMoney.Money -= warehouseCost
	merchantMoney.FlowIDs = make([]string, 0)
	merchantMoney.FlowType = AccountMoneyFlowTypeWarehouse
	warehouseMoney, err := getMoneyFromDb(stub, lastWarehouse.OwnerID)
	warehouseMoney.Money += warehouseCost
	warehouseMoney.FlowIDs = make([]string, 0)
	for _, bigPackage := range bigPackages {
		warehouseMoney.FlowIDs = append(warehouseMoney.FlowIDs, bigPackage.BigPackageID)
		merchantMoney.FlowIDs = append(merchantMoney.FlowIDs, bigPackage.BigPackageID)
	}
	warehouseMoney.FlowType = AccountMoneyFlowTypeWarehouse
	err = insertMoneyToDb(stub, warehouseMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	err = insertMoneyToDb(stub, merchantMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	merchantMoney.Money -= logisticCost
	merchantMoney.FlowIDs = make([]string, 0)
	merchantMoney.FlowType = AccountMoneyFlowTypeLogistic
	logisticMoney, err := getMoneyFromDb(stub, lastLogistic.OwnerID)
	logisticMoney.Money += logisticCost
	logisticMoney.FlowIDs = make([]string, 0)
	for _, bigPackage := range bigPackages {
		logisticMoney.FlowIDs = append(logisticMoney.FlowIDs, bigPackage.BigPackageID)
		merchantMoney.FlowIDs = append(merchantMoney.FlowIDs, bigPackage.BigPackageID)
	}
	logisticMoney.FlowType = AccountMoneyFlowTypeLogistic
	err = insertMoneyToDb(stub, logisticMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	err = insertMoneyToDb(stub, merchantMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	order.Payed = true
	err = updateOrderToDb(stub, order, false)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	// 订单完成，该农户信用加分
	userInfo, err := getUserInfoFromDb(stub, order.FarmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	userInfo.CreditScore += 15
	insertUserInfoToDb(stub, userInfo)

	return wrapToPbResponse([]byte("OK"), nil)
}

// 农民看到订单，接单
func (t *JutianChaincode) matchOrderByFarmer(stub shim.ChaincodeStubInterface, farmerID, orderID string) pb.Response {
	order, err := getOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if order.OrderID == "" {
		return wrapToPbResponse(nil, errors.New("cannot find order by id: "+orderID))
	}
	if order.FarmerID != "" {
		return wrapToPbResponse(nil, errors.New("order already had farmer: "+orderID))
	}
	if order.Status != OrderStatusInit {
		return wrapToPbResponse(nil, errors.New("order cannot match farmer right now for status is "+strconv.Itoa(order.Status)+",id "+orderID))
	}

	// 需要判断农民是否可以接单
	// farmerMoney, err := getMoneyFromDb(stub, farmerID)
	// if err != nil {
	// 	return wrapToPbResponse(nil, err)
	// }
	// needMoney := (farmerAssetMaterialUnitCost + farmerAssetLaborlUnitCost) * order.Amount
	// if farmerMoney.Money < needMoney {
	// 	money, _ := getFarmerLoanableMoney(stub, farmerID)
	// 	if (farmerMoney.Money + money) < needMoney {
	// 		return wrapToPbResponse(nil, errors.New("没有足够的资金和信用接单"))
	// 	}
	// }

	order.FarmerID = farmerID
	order.Status = OrderStatusMatched
	err = updateOrderToDb(stub, order, true)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(order, nil)
}

// 农民生产完了，修改该订单状态，传递给下一个环节
func (t *JutianChaincode) produceOrderByFarmer(stub shim.ChaincodeStubInterface, farmerID, orderID string) pb.Response {
	order, err := getOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if order.OrderID == "" {
		return wrapToPbResponse(nil, errors.New("cannot find order by id: "+orderID))
	}
	if order.FarmerID != farmerID {
		return wrapToPbResponse(nil, errors.New("order's farmerID is not "+farmerID+", by id: "+orderID))
	}
	if order.Status != OrderStatusMatched {
		return wrapToPbResponse(nil, errors.New("order's status is wrong, id: "+orderID))
	}

	currAsset, err := getAssetByCategoryNOwnFromDb(stub, order.CategoryID, order.FarmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if currAsset.Amount < order.Amount {
		return wrapToPbResponse(nil, errors.New("dont have enough asset"))
	}

	currAsset.Amount -= order.Amount
	err = updateFarmerAssetToDb(stub, currAsset)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	order.Status = OrderStatusProduced
	err = updateOrderToDb(stub, order, false)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse("", nil)
}

//
// 包装商
//

var packageUnitPerAmount = 5
var bigPackagePerAmount = 100 // 指的是大的包装里面有 100 KG的玉米，可能是20个小包
var unitAmountPerBigPackage = (bigPackagePerAmount / packageUnitPerAmount)

func (t *JutianChaincode) packageOrder(stub shim.ChaincodeStubInterface, ownerID, orderID string) pb.Response {
	// 1. 确认条件是否符合 order.status
	// 2. 一个一个做包装
	order, err := getOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if order.OrderID == "" {
		return wrapToPbResponse(nil, errors.New("cannot find order by id: "+orderID))
	}
	if order.Status != OrderStatusProduced {
		return wrapToPbResponse(nil, errors.New("order's status is wrong, id: "+orderID))
	}

	// 取出，该订单对应的，农民那里的资产
	asset, err := getAssetByCategoryNOwnFromDb(stub, order.CategoryID, order.FarmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if asset.AssetID == "" {
		return wrapToPbResponse(nil, errors.New("找不到对应的asset"))
	}

	// 确认 packageUnit的数量
	// var packageUnitNum = order.Amount / packageUnitPerAmount
	// if order.Amount%packageUnitPerAmount > 0 {
	// 	packageUnitNum += 1
	// }

	// 确认 big package
	var bigPackageNum = order.Amount / bigPackagePerAmount
	// if order.Amount%bigPackagePerAmount > 0 {
	// 	bigPackageNum += 1
	// }
	var lastBigPackageAmount = order.Amount - bigPackageNum*bigPackagePerAmount

	// 生成大包装和小的包装
	bigPackages := make([]*BigPackageWithUnits, 0)
	for i := 0; i < bigPackageNum; i++ {
		bigPackage := &BigPackage{
			bigPackageObjectType,
			GenerateRandom(32),
			ownerID,
			orderID,
			bigPackagePerAmount,
			bigPackageCost,
			false,
			MakeTimestamp(),
		}
		bigPackageWithUnit := &BigPackageWithUnits{}
		bigPackageWithUnit.BigPackage = bigPackage

		err := insertBigPackageToDb(stub, bigPackage)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}

		bigPackages = append(bigPackages, bigPackageWithUnit)

		// 小的包装
		for j := 0; j < unitAmountPerBigPackage; j++ {
			packageUnit := &PackageUnit{
				packgeUnitObjectType,
				GenerateRandom(32),
				bigPackage.BigPackageID,
				asset.AssetID,
				packageUnitPerAmount,
				MakeTimestamp(),
			}
			err := insertPackageToDb(stub, packageUnit)
			if err != nil {
				return wrapToPbResponse(nil, err)
			}
			bigPackageWithUnit.PackageUnits = append(bigPackageWithUnit.PackageUnits, packageUnit)
		}
	}

	// 最后的一个大包装
	lastBigPackageAmount = order.Amount % bigPackagePerAmount
	if lastBigPackageAmount > 0 {
		bigPackage := &BigPackage{
			bigPackageObjectType,
			GenerateRandom(32),
			ownerID,
			orderID,
			order.Amount % bigPackagePerAmount,
			bigPackageCost,
			false,
			MakeTimestamp(),
		}
		bigPackageWithUnit := &BigPackageWithUnits{}
		bigPackageWithUnit.BigPackage = bigPackage

		err := insertBigPackageToDb(stub, bigPackage)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}

		bigPackages = append(bigPackages, bigPackageWithUnit)

		// 计算最后的包装箱数量
		lastPackageUnitNum := bigPackage.Amount / packageUnitPerAmount
		for j := 0; j < lastPackageUnitNum; j++ {
			packageUnit := &PackageUnit{
				packgeUnitObjectType,
				GenerateRandom(32),
				bigPackage.BigPackageID,
				asset.AssetID,
				packageUnitPerAmount,
				MakeTimestamp(),
			}
			err := insertPackageToDb(stub, packageUnit)
			if err != nil {
				return wrapToPbResponse(nil, err)
			}
			bigPackageWithUnit.PackageUnits = append(bigPackageWithUnit.PackageUnits, packageUnit)
		}

		lastPackageUnitAmount := bigPackage.Amount % packageUnitPerAmount
		if lastPackageUnitAmount > 0 {
			//最后一个小袋
			lastPackageUnit := &PackageUnit{
				packgeUnitObjectType,
				GenerateRandom(32),
				bigPackage.BigPackageID,
				asset.AssetID,
				lastPackageUnitAmount,
				MakeTimestamp(),
			}
			err = insertPackageToDb(stub, lastPackageUnit)
			if err != nil {
				return wrapToPbResponse(nil, err)
			}
			bigPackageWithUnit.PackageUnits = append(bigPackageWithUnit.PackageUnits, lastPackageUnit)
		}
	}

	// 修改订单信息
	order.Status = OrderStatusPackaged
	err = updateOrderToDb(stub, order, false)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	return wrapStructToPbResponse(bigPackages, nil)
}

//根据order获取bigPackage列表
func (t *JutianChaincode) getBigPackagesByOrder(stub shim.ChaincodeStubInterface, packagerID, orderID string) pb.Response {
	bigPackages, err := getBigPackagesByOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	bigPackageWithUnits := make([]*BigPackageWithUnits, 0)
	// 把BigPackage和对应的PackageUnit放在一起
	for _, bigPackage := range bigPackages {
		packageUnits, err := getPackageUnitsByBigPackageFromDb(stub, bigPackage.BigPackageID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
		bigPackageWithUnit := &BigPackageWithUnits{bigPackage, packageUnits}
		bigPackageWithUnits = append(bigPackageWithUnits, bigPackageWithUnit)
	}

	return wrapStructToPbResponse(bigPackageWithUnits, nil)
}

func (t *JutianChaincode) getBigPackagesByOwner(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	bigPackages, err := getBigPackagesByOwnFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	bigPackageWithUnits := make([]*BigPackageWithUnits, 0)
	// 把BigPackage和对应的PackageUnit放在一起
	for _, bigPackage := range bigPackages {
		packageUnits, err := getPackageUnitsByBigPackageFromDb(stub, bigPackage.BigPackageID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
		bigPackageWithUnit := &BigPackageWithUnits{bigPackage, packageUnits}
		bigPackageWithUnits = append(bigPackageWithUnits, bigPackageWithUnit)
	}
	return wrapStructToPbResponse(bigPackageWithUnits, nil)
}

//
//仓库
//

//批量入库
func (t *JutianChaincode) insertWarehouseStoreInList(stub shim.ChaincodeStubInterface, ownerID, wName string, bigPackageIDs []string) pb.Response {
	r := pb.Response{}
	for _, ID := range bigPackageIDs {
		r = t.insertWarehouseStoreIn(stub, ownerID, ID, wName, bigPackageIDs)
	}
	return r
}

func (t *JutianChaincode) insertWarehouseStoreIn(stub shim.ChaincodeStubInterface, ownerID, bigPackageID, wName string, bigPackageIDs []string) pb.Response {
	bigPackage, err := getBigPackageFromDb(stub, bigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if bigPackage.BigPackageID == "" {
		return wrapToPbResponse(nil, errors.New("big package not found, id is "+bigPackageID))
	}

	w := &WarehouseStoreIn{
		warehouseStoreInObjectType,
		bigPackageID,
		ownerID,
		bigPackage.OrderID,
		MakeTimestamp(),
		warehouseCost,
		wName,
		false,
	}
	err = insertWarehouseStoreInToDb(stub, w)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	// 更新大包装的状态
	bigPackage.Delivered = true
	err = updatePackageToDb(stub, bigPackage)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	// 检查该订单下所有包装状态，如果全部入库，修改订单状态
	bigPackageMap := make(map[string]string)
	for _, ID := range bigPackageIDs {
		bigPackageMap[ID] = ID
	}
	packages, err := getBigPackagesByOrderFromDb(stub, bigPackage.OrderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	allStoreIn := true
	for _, currPackage := range packages {
		if currPackage.BigPackageID == bigPackage.BigPackageID {
			continue
		}
		if !currPackage.Delivered && bigPackageMap[currPackage.BigPackageID] == "" {
			allStoreIn = false
		}
	}
	if allStoreIn {
		order, err := getOrderFromDb(stub, bigPackage.OrderID)
		order.Status = OrderStatusPackagedStoreIn
		err = updateOrderToDb(stub, order, false)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
	}

	return wrapStructToPbResponse(w, nil)
}

// 获取该订单下的所有的库存信息
func (t *JutianChaincode) getWarehouseStoreInByOrder(stub shim.ChaincodeStubInterface, orderID string) pb.Response {
	storesIn, err := getWareshouseStoreInsByOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(storesIn, nil)
}

// 获取该用户的所有库存
func (t *JutianChaincode) getStoresByOwner(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	storesIn, err := getWareshouseStoreInsByOwnFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(storesIn, nil)
}

//
//物流
//

func (t *JutianChaincode) insertLogisticList(stub shim.ChaincodeStubInterface, ownerID, truck string, bigPackageIDs []string) pb.Response {
	r := pb.Response{}
	for _, ID := range bigPackageIDs {
		r = t.insertLogistic(stub, ownerID, ID, truck, bigPackageIDs)
	}
	return r
}

func (t *JutianChaincode) insertLogistic(stub shim.ChaincodeStubInterface, ownerID, bigPackageID, truck string, bigPackageIDs []string) pb.Response {
	//步骤比较长
	// 1. 从包装查询到订单信息
	// 2. 从订单查询到终点和农民账户信息
	// 3. 从农民信息查询到起点
	bigPackage, err := getBigPackageFromDb(stub, bigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if bigPackage.BigPackageID == "" {
		return wrapToPbResponse(nil, errors.New("cannot get big package from id: "+bigPackageID))
	}

	order, err := getOrderFromDb(stub, bigPackage.OrderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if order.OrderID == "" {
		return wrapToPbResponse(nil, errors.New("cannot get order from id: "+bigPackage.OrderID))
	}

	farmer, err := getUserInfoFromDb(stub, order.FarmerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if farmer.OwnerID == "" {
		return wrapToPbResponse(nil, errors.New("cannot get UserInfo from id: "+order.FarmerID))
	}

	l := &Logistic{
		logisticObjectType,
		bigPackageID,
		ownerID,
		order.OrderID,
		farmer.Location,
		order.Destination,
		logisticCost,
		MakeTimestamp(),
		false,
	}

	err = insertLogisticToDb(stub, l)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	// 修改库存状态
	w, err := getWarehouseStoreInFromDb(stub, bigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	w.Sent = true
	err = updateWarehouseStoreInToDb(stub, w)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	//判断是否该订单下所有大包装都已经物流了，如果是，修改订单状态
	bigPackageMap := make(map[string]string)
	for _, ID := range bigPackageIDs {
		bigPackageMap[ID] = ID
	}
	bigPackages, err := getBigPackagesByOrderFromDb(stub, bigPackage.OrderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	allDelivered := true
	for _, currPackage := range bigPackages {
		if currPackage.BigPackageID == bigPackage.BigPackageID {
			continue
		}
		p, err := getLogisticFromDb(stub, currPackage.BigPackageID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		} else if p.BigPackageID == "" && bigPackageMap[currPackage.BigPackageID] == "" {
			allDelivered = false
			break
		}
	}
	if allDelivered {
		order.Status = OrderStatusDeliveried
		err = updateOrderToDb(stub, order, false)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
	}

	return wrapStructToPbResponse(l, nil)
}

func (t *JutianChaincode) getLogisticByOwner(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	logistics, err := getLogisticByOwnFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(logistics, nil)
}

func (t *JutianChaincode) getLogisticByOrder(stub shim.ChaincodeStubInterface, orderID string) pb.Response {
	bigPackages, err := getBigPackagesByOrderFromDb(stub, orderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	logistics := make([]*Logistic, 0)
	for _, bigPackage := range bigPackages {
		l, err := getLogisticFromDb(stub, bigPackage.BigPackageID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
		logistics = append(logistics, l)
	}
	return wrapStructToPbResponse(logistics, nil)
}

//
// 商户终端
//
func (t *JutianChaincode) insertMechandiseList(stub shim.ChaincodeStubInterface, ownerID string, bigPackageIDs []string) pb.Response {
	r := pb.Response{}
	for _, ID := range bigPackageIDs {
		r = t.insertMechandise(stub, ownerID, ID, bigPackageIDs)
	}
	return r
}

func (t *JutianChaincode) insertMechandise(stub shim.ChaincodeStubInterface, ownerID, bigPackageID string, bigPackageIDs []string) pb.Response {
	bigPackage, err := getBigPackageFromDb(stub, bigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if bigPackage.BigPackageID == "" {
		return wrapToPbResponse(nil, errors.New("cannot find bigPackage by id "+bigPackageID))
	}

	// 先获取订单判断是否是我的订单
	order, err := getOrderFromDb(stub, bigPackage.OrderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if order.OrderID == "" {
		return wrapToPbResponse(nil, errors.New("cannot find order by id "+bigPackage.OrderID))
	} else if order.OwnerID != ownerID {
		return wrapToPbResponse(nil, errors.New("this order not yours "))
	}

	logistic, err := getLogisticFromDb(stub, bigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	} else if logistic.BigPackageID == "" {
		return wrapToPbResponse(nil, errors.New("cannot find logistic by id "+bigPackageID))
	}

	// 1.根据大包装找到对应的小包装
	// 2.写入Merchandise
	// 3.修改logistic
	// 4.修改订单，如果全部到货的话
	smallPackages, err := getPackageUnitsByBigPackageFromDb(stub, bigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	merchandises := make([]*Merchandise, 0)
	for _, p := range smallPackages {
		m := &Merchandise{
			merchandiseObjectType,
			p.PackageID,
			ownerID,
			order.CategoryID,
			p.Amount,
			merchandisePrice,
			false,
			MakeTimestamp(),
		}
		err = insertMerchandiseToDb(stub, m)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}
		merchandises = append(merchandises, m)
	}

	bigPackages, err := getBigPackagesByOrderFromDb(stub, bigPackage.OrderID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	logistic.Delivered = true
	err = updateLogisticToDb(stub, logistic)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	bigPackageMap := make(map[string]string)
	for _, ID := range bigPackageIDs {
		bigPackageMap[ID] = ID
	}
	allDelivered := true
	for _, currPackage := range bigPackages {
		if currPackage.BigPackageID == bigPackage.BigPackageID {
			continue
		}
		p, err := getLogisticFromDb(stub, currPackage.BigPackageID)
		if err != nil {
			return wrapToPbResponse(nil, err)
		} else if (p.BigPackageID == "" || !p.Delivered) && bigPackageMap[currPackage.BigPackageID] == "" {
			allDelivered = false
			break
		}
	}
	if allDelivered {
		order.Status = OrderStatusFinished
		err = updateOrderToDb(stub, order, false)
		if err != nil {
			return wrapToPbResponse(nil, err)
		}

		// 全部到货了，则付款
		return t.payOrder(stub, ownerID, order.OrderID)
	}

	return wrapStructToPbResponse(merchandises, nil)
}

func (t *JutianChaincode) getMerchandiseByOwner(stub shim.ChaincodeStubInterface, ownerID string) pb.Response {
	merchandises, err := getMerchandisesByOwnerFromDb(stub, ownerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	return wrapStructToPbResponse(merchandises, nil)
}

// 商品溯源
func (t *JutianChaincode) getMerchandiseFlow(stub shim.ChaincodeStubInterface, packageID string) pb.Response {
	merchandise, err := getMerchandiseFromDb(stub, packageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	packageUnit, err := getPackageFromDb(stub, packageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	bigPackage, err := getBigPackageFromDb(stub, packageUnit.BigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	logistic, err := getLogisticFromDb(stub, bigPackage.BigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	warehouse, err := getWarehouseStoreInFromDb(stub, bigPackage.BigPackageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	asset, err := getFarmerAssetFromDb(stub, packageUnit.AssetID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	asset.LaborUnitCost = 0
	asset.MaterialUnitCost = 0

	farmerUserInfo, err := getUserInfoFromDb(stub, asset.OwnerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	asset.FarmerName = farmerUserInfo.RealName
	asset.Location = farmerUserInfo.Location

	mFlow := &MerchandiseFlow{
		merchandise,
		logistic,
		warehouse,
		bigPackage,
		asset,
	}
	return wrapStructToPbResponse(mFlow, nil)
}

// 模拟用户购买
func (t *JutianChaincode) purchasePackage(stub shim.ChaincodeStubInterface, packageID string) pb.Response {
	merchandise, err := getMerchandiseFromDb(stub, packageID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}
	if merchandise.SoldOut {
		return wrapToPbResponse(nil, errors.New("merchandise already sold out"))
	}
	merchandise.SoldOut = true
	err = updateMerchandiseToDb(stub, merchandise)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	merchantMoney, err := getMoneyFromDb(stub, merchandise.OwnerID)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	merchantMoney.Money += merchandise.Price
	merchantMoney.FlowIDs = []string{merchandise.PackageID}
	merchantMoney.FlowType = AccountMoneyFlowTypeMerchandise
	err = insertMoneyToDb(stub, merchantMoney)
	if err != nil {
		return wrapToPbResponse(nil, err)
	}

	return wrapStructToPbResponse(merchandise, nil)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Start Jutian Chaincode Date " + time.Now().Format("2006-01-02 15:04:05"))
	err := shim.Start(new(JutianChaincode))
	if err != nil {
		fmt.Printf("Error starting jutian chaincode - %s", err)
	}
}
