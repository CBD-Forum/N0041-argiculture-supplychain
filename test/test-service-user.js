var tape = require('tape');
var _test = require('tape-promise');
var path = require('path');
var test = _test(tape);

var util = require("./util.js");

var Constant = require("../server/constant.js");
var UserService = require(__dirname + "/../server/service/user.service.js");

var Response = require("./response.js");

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

let token;
let user;


test('\n\n***** Asset Service Test: user  *****\n\n', (t) => {
  var currAdmin = null;
  util.enroll("admin", "adminpw")
  .then((admin) => {
    t.pass("init chaincode success");
    currAdmin = admin;

    t.comment("测试用户账号1初始化")
    let registerReq = { body: {username: "test1", password: "test1", type: UserService.UserTypeFarmer}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试用户账号2初始化")
    let registerReq = { body: {username: "test2", password: "test2", type: UserService.UserTypeFarmer}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试登录接口");
    let loginReq = { body: {username: "test1", password: "test1"}}
    let res = new Response();
    UserService.login(loginReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    token = result.result.token;
    user = result.result;

    t.comment("开始测试设置账号信息");
    let req = { user:user, body: {token: token, name: "农民编号A", identity: "637687198904047682", location: "内蒙古"}};
    let res = new Response();
    UserService.setAccount(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("开始获取设置账号信息");
    let req = { user: user };
    let res = new Response();
    UserService.getUserInfo(req, res, null);
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