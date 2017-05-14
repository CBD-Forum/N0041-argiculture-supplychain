var tape = require('tape');
var _test = require('tape-promise');
var path = require('path');
var test = _test(tape);

var util = require("./util.js");

var Constant = require("../server/constant.js");
var UserService = require(__dirname + "/../server/service/user.service.js");
var OrderService = require(__dirname + "/../server/service/order.service.js");

var Response = require("./response.js");

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

let token;
let user;

test('\n\n***** Asset Service Test: order  *****\n\n', (t) => {
  var currAdmin = null;
  util.enroll("admin", "adminpw")
  .then((admin) => {
    t.pass("init chaincode success");
    currAdmin = admin;

    t.comment("测试商户账号1初始化")
    let registerReq = { body: {username: "merchant1", password: "merchant1", type: UserService.UserTypeMerchant}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试登录接口");
    let loginReq = { body: {username: "merchant1", password: "merchant1"}}
    let res = new Response();
    UserService.login(loginReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    token = result.result.token;
    user = result.result;

    t.comment("测试发布订单");
    let req = {user:user, body: {categoryID: 1, amount: 60, cost: 500, destination: "上海"}}
    let res = new Response();
    OrderService.createOrder(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));
    let orderID = result.result.OrderID;

    t.comment("测试获取订单");
    let req = {user: user, query: {orderID: orderID}}
    let res = new Response();
    OrderService.getOrder(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试发布订单2");
    let req = {user:user, body: {categoryID: 1, amount: 95, cost: 600, destination: "北京"}}
    let res = new Response();
    OrderService.createOrder(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试获取自己发布的订单");
    let req = {user: user}
    let res = new Response();
    OrderService.getOrdersByOwner(req, res, null);
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