'use strict';

let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

let Tool = require("../tool.js")
let UserService = require("./user.service.js");
let AssetService = require("./asset.service.js");
let OrderService = require("./order.service.js");

class PackageService {

  static packageOrder(req, res, next) {
    if (!req.body || !req.body.orderID ) {
      return res.status(501).send('orderID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypePackager) 
      return res.status(501).send('身份错误，这个API只允许包装商调用');

    let orderID = req.body.orderID;
    let packagerID = req.user.id;
    Chaincode.invoke("packageOrder", [packagerID, orderID], Constant.admin)
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

  static getBigPackagesByOrder(req, res, next) {
    if (!req.body || !req.body.orderID ) {
      return res.status(501).send('orderID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypePackager &&
        currUser.type != UserService.UserTypeWarehouse &&
        currUser.type != UserService.UserTypeLogistic &&
        currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许包装商调用');
    
    let orderID = req.body.orderID;
    let packagerID = req.user.id;
    Chaincode.invoke("getBigPackagesByOrder", [packagerID, orderID], Constant.admin)
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

  static getBigPackagesByOwner(req, res, next) {
    let currUser = req.user;
    if (currUser.type != UserService.UserTypePackager)
      return res.status(501).send('身份错误，这个API只允许包装商调用');

    let packagerID = req.user.id;
    Chaincode.invoke("getBigPackagesByOwner", [packagerID], Constant.admin)
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

  static getBigPackageByOrderStatus(req, res, next) {
    let currUser = req.user;
    if (currUser.type != UserService.UserTypePackager &&
        currUser.type != UserService.UserTypeWarehouse &&
        currUser.type != UserService.UserTypeLogistic &&
        currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许包装商调用');

    let status = 0;
    // 判断权限
    if (currUser.type == UserService.UserTypeFarmer) 
      status = OrderService.OrderStatusInit;
    else if (currUser.type == UserService.UserTypePackager) 
      status = OrderService.OrderStatusProduced;
    else if (currUser.type == UserService.UserTypeWarehouse) 
      status = OrderService.OrderStatusPackaged;
    else if (currUser.type == UserService.UserTypeLogistic) 
      status = OrderService.OrderStatusPackagedStoreIn;
    else if (currUser.type == UserService.UserTypeMerchant) 
      status = OrderService.OrderStatusDeliveried;

    Chaincode.invoke("getBigPackageByOrderStatus", [status], Constant.admin)
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

module.exports = PackageService