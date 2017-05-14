'use strict';

let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

let Tool = require("../tool.js")
let UserService = require("./user.service.js");
let AssetService = require("./asset.service.js");

class WarehouseService {

  static _getWarehouseShelf() {
    let num = 10 + parseInt(Math.random() * 20);
    return "A号仓" + num + "号货架";
  }
  
  static insertWarehouseStoreIn(req, res, next) {
    if (!req.body || !req.body.bigPackageIDs) {
      return res.status(501).send('bigPackageIDs not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeWarehouse) 
      return res.status(501).send('身份错误，这个API只允许仓库调用');

    let bigPackageIDs = req.body.bigPackageIDs;
    let wName = WarehouseService._getWarehouseShelf();

    let ownerID = req.user.id;
    let args = [ownerID, wName];
    for (let key in bigPackageIDs) {
      args.push(bigPackageIDs[key])
    }

    Chaincode.invoke("insertWarehouseStoreInList", args, Constant.admin)
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

  static getWarehousesByOrder(req, res, next) {
    let orderID = req.query.orderID;
    if (!orderID)
      return res.status(501).send('orderID not provided');
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeWarehouse &&
        currUser.type != UserService.UserTypeLogistic) 
      return res.status(501).send('身份错误，这个API只允许仓库调用');

    Chaincode.query("getWarehouseByOrder", [orderID], Constant.admin)
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

  static getWarehousesByOwner(req, res, next) {
    let ownerID = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeWarehouse) 
      return res.status(501).send('身份错误，这个API只允许仓库调用');
    
    Chaincode.query("getWarehouseByOwner", [ownerID], Constant.admin)
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

module.exports = WarehouseService