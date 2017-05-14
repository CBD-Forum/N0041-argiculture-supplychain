'use strict';

let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

let Tool = require("../tool.js")
let UserService = require("./user.service.js");
let AssetService = require("./asset.service.js");

class LogisticService {

  static _getLogisticProviders() {
    return [
      "圆圆通物流",
      "申通通物流",
      "千里眼物流",
      "哒哒哒物流"
    ]
  }

  static _getLogisticProvider() {
    let num = (Math.random() * 100) % 4;
    return LogisticService._getLogisticProviders()[num];
  }

  static _getTruck() {
    let num = parseInt(Math.random() * 90) + 10;
    return "货车" + num + "号";
  }

  static insertLogistic(req, res, next) {
    if (!req.body || !req.body.bigPackageIDs) {
      return res.status(501).send('bigPackageIDs not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeLogistic) 
      return res.status(501).send('身份错误，这个API只允许物流调用');

    let bigPackageIDs = req.body.bigPackageIDs;
    let truck = LogisticService._getTruck();

    let ownerID = req.user.id;
    let args = [ownerID, truck];
    for (let key in bigPackageIDs) {
      args.push(bigPackageIDs[key])
    }

    Chaincode.invoke("insertLogisticList", args, Constant.admin)
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

  static getLogisticsByOwner(req, res, next) {
    let ownerID = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeLogistic &&
        currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许物流调用');

    Chaincode.query("getLogisticByOwner", [ownerID], Constant.admin)
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
  
  static getLogisticByOrder(req, res, next) {
    if (!req.query || !req.query.orderID) {
      return res.status(501).send('orderID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeLogistic &&
        currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许物流调用');
    
    let orderID = req.query.orderID;
    Chaincode.query("getLogisticByOrder", [orderID], Constant.admin)
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

module.exports = LogisticService