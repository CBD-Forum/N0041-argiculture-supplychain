var tape = require('tape');
var _test = require('tape-promise');
var path = require('path');
var test = _test(tape);

var util = require("./util.js");

var Constant = require("../server/constant.js");
var UserService = require(__dirname + "/../server/service/user.service.js");
var BlockchainService = require(__dirname + "/../server/service/blockchain.service.js");

var Response = require("./response.js");

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

let token;
let user;


test('\n\n***** Asset Service Test: BlockchainService  *****\n\n', (t) => {
  var currAdmin = null;
  util.enroll("admin", "adminpw")
  .then((admin) => {
    t.pass("init chaincode success");
    currAdmin = admin;

    t.comment("测试用户账号1初始化")
    let registerReq = { body: {username: "ttt", password: "ttt", nickname: "tttnickname"}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试用户账号2初始化")
    let registerReq = { body: {username: "test2", password: "test2", nickname: "tttnickname"}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("获取getBlockHeight");
    let req = {};
    let res = new Response();
    BlockchainService.getBlockHeight(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("根据height获取block");
    let req = {query: {height: result.result - 1}};
    let res = new Response();
    BlockchainService.getBlockByNumber(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));
    t.comment("结束")
    t.end();
  }).catch((err) => {
    t.fail(err.stack);
    t.end();
  });
})