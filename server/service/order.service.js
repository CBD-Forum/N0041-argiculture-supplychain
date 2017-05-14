'use strict';

let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

let Tool = require("../tool.js")
let UserService = require("./user.service.js");
let AssetService = require("./asset.service.js");

class OrderService {

  static _checkStatus(status) {
    status = parseInt(status);
    if (status != OrderService.OrderStatusInit &&
        status != OrderService.OrderStatusMatched &&
        status != OrderService.OrderStatusProduced &&
        status != OrderService.OrderStatusPackaged &&
        status != OrderService.OrderStatusPackagedStoreIn &&
        status != OrderService.OrderStatusDeliveried &&
        status != OrderService.OrderStatusFinished)
      return false;
    return true
  }

  static createOrder(req, res, next) {
    if (!req.body || !req.body.categoryID || !req.body.amount || !req.body.cost || !req.body.destination) {
      console.log(req.body)
      return res.status(501).send('categoryID or amount or cost or detail or destination not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许商户调用');

    let categoryID = req.body.categoryID;
    let amount = parseInt(req.body.amount);
    let cost = parseInt(req.body.cost);
    let destination = req.body.destination;

    if (!AssetService.checkCategoryExists(categoryID))
      return res.status(501).send('categoryID not exists');

    let ownerID = req.user.id;
    let args = [ownerID, "", categoryID, destination, amount, cost, ""];

    Chaincode.invoke("insertOrder", args, Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getOrder(req, res, next) {
    let orderID = req.query.orderID;
    if (!orderID)
      return res.status(501).send('orderID not provided');

    let args = [orderID];
    Chaincode.query("getOrder", args, Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getOrdersByFarmer(req, res, next) {
    let farmerId = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    Chaincode.query("getOrdersByFarmer", [farmerId], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getOrdersByOwner(req, res, next) {
    let ownerID = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许商户调用');
    Chaincode.query("getOrdersByOwner", [ownerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getOrdersByStatus(req, res, next) {
    let ownerID = req.user.id;
    let status = req.query.status;
    if (!OrderService._checkStatus(status))
      return res.status(501).send("status error");

    // 判断权限
    let currUser = req.user;
    if (currUser.type == UserService.UserTypeFarmer) 
      status = OrderService.OrderStatusInit;
    else if (currUser.type == UserService.UserTypePackager) 
      status = OrderService.OrderStatusProduced;
    else if (currUser.type == UserService.UserTypeWarehouse) 
      status = OrderService.OrderStatusPackaged;
    else if (currUser.type == UserService.UserTypeLogistic) 
      status = OrderService.OrderStatusPackagedStoreIn;

    Chaincode.query("getOrdersByStatus", [status], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static payOrder(req, res, next) {
    if (!req.body || !req.body.orderID ) {
      return res.status(501).send('orderID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许商户调用');

    let orderID = req.body.orderID;
    let ownerID = req.user.id;
    Chaincode.invoke("payOrder", [ownerID, orderID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static matchOrderByFarmer(req, res, next) {
    if (!req.body || !req.body.orderID ) {
      return res.status(501).send('orderID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    let orderID = req.body.orderID;
    let farmerID = req.user.id;
    Chaincode.invoke("matchOrderByFarmer", [farmerID, orderID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

}

OrderService.OrderStatusInit            = 1 //订单初始化
OrderService.OrderStatusMatched         = 2 //订单被接单
OrderService.OrderStatusProduced        = 3 //订单生产完成
OrderService.OrderStatusPackaged        = 4 //订单打包
OrderService.OrderStatusPackagedStoreIn = 5 //订单入库
OrderService.OrderStatusDeliveried      = 6 //订单发货
OrderService.OrderStatusFinished        = 7 //订单结束

module.exports = OrderService