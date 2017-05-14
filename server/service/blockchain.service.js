'use strict';

let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

let Tool = require("../tool.js")
let UserService = require("./user.service.js");

class BlockchainService {
  static getInfo(req, res, next) {
    Constant.chain.queryInfo()
    .then((result) => {
      console.log(result);
      return res.status(200).send(result);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getBlockHeight(req, res, next) {
    Constant.chain.queryInfo()
    .then((result) => {
      console.log(result.height)
      return res.status(200).send(result.height.low + "");
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getBlockByNumber(req, res, next) {
    if (!req.query || !req.query.height) {
      return res.status(501).send('height not provided, payload is ' + JSON.stringify(req.body));
    }
    var height = parseInt(req.query.height);
    if (isNaN(height))
      return res.status(501).send('height err' + req.query.height);

    Constant.chain.queryBlock(height)
    .then((result) => {
      result.header.previous_hash = result.header.previous_hash.toString('base64');
      result.header.data_hash = result.header.data_hash.toString('base64');
      result.data.data[0] = result.data.data[0].toString('base64');
      // result.header.previous_hash.buffer = result.header.previous_hash.buffer.toString();
      // result.header.data_hash.buffer = result.header.data_hash.buffer.toString();
      // result.data.data[0].buffer = result.data.data[0].buffer.toString();
      for (var i = result.metadata.metadata.length - 1; i >= 0; i--) {
         result.metadata.metadata[i] = result.metadata.metadata[i].toString('base64');
       } 
      return res.status(200).send(result);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getBlockByHash(req, res, next) {
    if (!req.query || !req.query.hash) {
      return res.status(501).send('hash not provided, payload is ' + JSON.stringify(req.body));
    }
    var hash = req.query.hash;

    Constant.chain.queryBlockByHash(hash)
    .then((result) => {
      result.header.previous_hash = result.header.previous_hash.toString('base64');
      result.header.data_hash = result.header.data_hash.toString('base64');
      result.data.data[0] = result.data.data[0].toString('base64');
      for (var i = result.metadata.metadata.length - 1; i >= 0; i--) {
         result.metadata.metadata[i] = result.metadata.metadata[i].toString('base64');
       } 
      return res.status(200).send(result);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }
}

module.exports = BlockchainService