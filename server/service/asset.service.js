'use strict';

let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

let Tool = require("../tool.js")
let UserService = require("./user.service.js");

class AssetService {

  static _getCategory() {
    return [
      {categoryID: 1, name: "保玉1号"},
      {categoryID: 2, name: "苏玉31"},
      {categoryID: 3, name: "中糯1号"},
      {categoryID: 4, name: "吉祥1号"},
      {categoryID: 5, name: "蠡玉88"},
      {categoryID: 6, name: "铁研358"},
      {categoryID: 7, name: "金诚508"},
      {categoryID: 8, name: "农华101"}
    ];
  }

  static getCategoryById(categoryID) {
    if (AssetService.categoryMap == null || AssetService.categoryMap == undefined) {
      AssetService.categoryMap = new Map();

      let categorys = AssetService._getCategory();
      for (let key in categorys) {
        let category = categorys[key];
        AssetService.categoryMap.set(category.categoryID, category)
      }
    }
    return AssetService.categoryMap.get(categoryID)
  }

  static checkCategoryExists(categoryID) {
    let category = AssetService.getCategoryById(parseInt(categoryID));
    return category != null && category != undefined;
  }

  static getCategory(req, res, next) {
    res.status(200).send(AssetService._getCategory());
  }

  static createAsset(req, res, next) {
    if (!req.body || !req.body.categoryID || !req.body.amount || !req.body.materialCost || !req.body.laborCost) {
      return res.status(501).send('categoryID or amount or materialCost or laborCost not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    let categoryID = req.body.categoryID;
    let amount = parseInt(req.body.amount);
    let materialCost = parseInt(req.body.materialCost);
    let laborCost = parseInt(req.body.laborCost);

    if (!AssetService.checkCategoryExists(categoryID))
      return res.status(501).send('categoryID not exists');

    let ownerID = req.user.id;
    let args = [ownerID, categoryID, amount, materialCost, laborCost];

    Chaincode.invoke("insertFarmerAsset", args, Constant.admin)
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
/**
  @api {get} /session/asset getAsset
  @apiName AssetService
  @apiGroup getAsset
  @apiDescription 获取单个资产

  @apiParam {Number} assetID 资产ID，必须

  @apiSuccess {Number} 200 请求成功
  @apiError {Number} 500 服务器错误
  @apiError {Number} 501 assetID not provided
  @apiErrorExample Error-response:
    HTTP/1.1 501
    {
      "result":"assetID not provided",
    }

*/

  // 获取单个资产
  static getAsset(req, res, next) {
    if (!req.query || !req.query.assetID) {
      return res.status(501).send('assetID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    let assetID = req.query.assetID;
    let ownerID = req.user.id;
    Chaincode.query("getFarmerAsset", [assetID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      console.log(err);
      res.status(500).send(err);
    })
  }

  // 根据类别和用户，获取单个资产
  static getAssetByCategoryIDNOwner(req, res, next) {
    if (!req.query || !req.query.categoryID) {
      console.log(req.query)
      return res.status(501).send('categoryID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    let categoryID = req.query.categoryID;
    let ownerID = req.user.id;
    Chaincode.query("getAssetByCategoryIDNOwner", [ownerID, categoryID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      console.log(err);
      res.status(500).send(err);
    })
  }

  // 获取该用户的所有资产
  static getFarmerAssetsByOwner(req, res, next) {
    let ownerID = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    Chaincode.query("getFarmerAssetsByOwner", [ownerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      console.log(err);
      res.status(500).send(err);
    })
  }

/**
  @api {get} /session/asset/history getAssetHistory
  @apiName AssetService
  @apiGroup getAssetHistory
  @apiDescription 获取资产变更历史

  @apiParam {Number} assetID 资产ID，必须
  @apiSuccess {Number} 200 请求成功
  @apiError {Number} 500 服务器错误
  @apiError {Number} 501 tradingID not provided
*/

  static getAssetHistory(req, res, next) {
    if (!req.query || !req.query.assetID) {
      return res.status(501).send('assetID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    Chaincode.query("getAssetHistory", [req.query.assetID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      for (var i = 0; i < result.length; i++) {
        for (var j = 0; j < result[i].TradeLogs.length; j++) {
          if (result[i].TradeLogs[j].LogID == "")
            continue;
          var tradeLog = result[i].TradeLogs[j];
          AssetService.attachCategoryToAsset(result[i].TradeLogs[j])
          if (String(req.user.id) == String(tradeLog.OldOwnerID)) {
            result[i].TradeLogs[j].OldOwnerName = req.user.username;
          } else {
            result[i].TradeLogs[j].OldOwnerName = Tool.getName() + "**";
          }
          if (String(req.user.id) == String(tradeLog.NewOwnerID)) {
            result[i].TradeLogs[j].NewOwnerName = req.user.username;
          } else {
            result[i].TradeLogs[j].NewOwnerName = Tool.getName() + "**";
          }
        }
      }
      res.status(200).send(result);
    }).catch((err) => {
      console.log(err);
      res.status(500).send(err);
    })
  } 

  static attachCategoryToAsset(asset) {
    let categoryID = parseInt(asset.CategoryID);
    if (!isNaN(categoryID)) {
      let c = AssetService.getCategoryById(categoryID);
      asset.Category = c.name;
    }
  }

  static produceOrderByFarmer(req, res, next) {
    if (!req.body || !req.body.orderID ) {
      return res.status(501).send('orderID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');
    
    let orderID = req.body.orderID;
    let farmerID = req.user.id;
    Chaincode.invoke("produceOrderByFarmer", [farmerID, orderID], Constant.admin)
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

module.exports = AssetService