'use strict';

let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

let Tool = require("../tool.js")
let UserService = require("./user.service.js");
let AssetService = require("./asset.service.js");

class MerchandiseService {

  static insertMechandise(req, res, next) {
    if (!req.body || !req.body.bigPackageIDs) {
      return res.status(501).send('bigPackageIDs not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许商户调用');

    let bigPackageIDs = req.body.bigPackageIDs;

    let ownerID = req.user.id;
    let args = [ownerID];
    for (let key in bigPackageIDs) {
      args.push(bigPackageIDs[key])
    }

    Chaincode.invoke("insertMechandiseList", args, Constant.admin)
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

  static getMerchandiseByOwner(req, res, next) {
    let ownerID = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeMerchant) 
      return res.status(501).send('身份错误，这个API只允许商户调用');

    Chaincode.query("getMerchandiseByOwner", [ownerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        for (var i = result.Data.length - 1; i >= 0; i--) {
          AssetService.attachCategoryToAsset(result.Data[i]);
        }
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }
  
  static getMerchandiseFlow(req, res, next) {
    if (!req.query || !req.query.packageID) {
      return res.status(501).send('packageID not provided');
    }

    let packageID = req.query.packageID;

    Chaincode.query("getMerchandiseFlow", [packageID], Constant.admin)
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

  static purchaseMechandise(req, res, next) {
    if (!req.body || !req.body.packageID) {
      return res.status(501).send('packageID not provided');
    }
    let packageID = req.body.packageID;

    let ownerID = req.user.id;
    let args = [packageID];

    Chaincode.invoke("purchasePackage", args, Constant.admin)
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

module.exports = MerchandiseService