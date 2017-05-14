var tape = require('tape');
var _test = require('tape-promise');
var path = require('path');
var test = _test(tape);

var util = require("./util.js");

var Constant = require("../server/constant.js");
var UserService = require(__dirname + "/../server/service/user.service.js");
var OrderService = require(__dirname + "/../server/service/order.service.js");
var AssetService = require(__dirname + "/../server/service/asset.service.js");
var PackageService = require(__dirname + "/../server/service/package.service.js");
var WarehouseService = require(__dirname + "/../server/service/warehouse.service.js");
var LogisticService = require(__dirname + "/../server/service/logistic.service.js");
var MerchandiseService = require(__dirname + "/../server/service/merchandise.service.js");

var Response = require("./response.js");

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

let token;
let user;

let currOrder

let bigPackages;

let merchandises;

let farmer;

test('\n\n***** Asset Service Test: MerchandiseService  *****\n\n', (t) => {
  var currAdmin = null;
  util.enroll("admin", "adminpw")
  .then((admin) => {
    t.pass("init chaincode success");
    currAdmin = admin;

    t.comment("测试农民账号1初始化")
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

    t.comment("开始测试设置账号信息");
    let req = { user:user, body: {token: token, name: "商户编号A", identity: "637687198904047682", location: "内蒙古"}};
    let res = new Response();
    UserService.setAccount(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("获取所有未到货的订单");
    let req = {user:user, query: { status: OrderService.OrderStatusDeliveried} }
    let res = new Response();
    OrderService.getOrdersByStatus(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    let orderID;
    let data = result.result;
    orderID = data[0].OrderID;
    currOrder = data[0];

    t.comment("获取订单下的所有物流")
    let req = {user:user, query: {orderID: currOrder.OrderID}}
    let res = new Response();
    LogisticService.getLogisticByOrder(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));
    bigPackages = result.result;

    t.comment("给第一个大包装到货")
    let bigPackage = bigPackages[0];
    let req = {user:user, body: {bigPackageID: bigPackage.BigPackageID, cost: 10}}
    let res = new Response();
    MerchandiseService.insertMechandise(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("查看下当前订单进展")
    let req = {user:user, query: {orderID: currOrder.OrderID}};
    let res = new Response();
    OrderService.getOrder(req, res, null);
    return res.promise;

  // }).then((result) => {
  //   t.pass(JSON.stringify(result));

  //   t.comment("给第2个大包装发货")
  //   let bigPackage = bigPackages[1];
  //   let req = {user:user, body: {bigPackageID: bigPackage.BigPackageID, cost: 10}}
  //   let res = new Response();
  //   MerchandiseService.insertMechandise(req, res, null);
  //   return res.promise;

  // }).then((result) => {
  //   t.pass(JSON.stringify(result));

  //   t.comment("查看下当前订单进展")
  //   let req = {user:user, query: {orderID: currOrder.OrderID}};
  //   let res = new Response();
  //   OrderService.getOrder(req, res, null);
  //   return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试农民登录接口");
    let loginReq = { body: {username: "farmer2", password: "farmer2"}}
    let res = new Response();
    UserService.login(loginReq, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民获取可借款的金额");
    farmer = result.result;
    let req = {user:farmer}
    let res = new Response();
    UserService.getFarmerLoanableMoney(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民借款");
    let loan = result.result;
    let req = {user:farmer, body: {loan: loan}}
    let res = new Response();
    UserService.loanByFarmer(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民获取可借款的金额");
    let req = {user:farmer}
    let res = new Response();
    UserService.getFarmerLoanableMoney(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民获取钱包");
    let req = {user:farmer}
    let res = new Response();
    UserService.getUserInfo(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("获取目前的该商户的所有商品")
    let req = {user:user}
    let res = new Response();
    MerchandiseService.getMerchandiseByOwner(req, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    merchandises = result.result;
    t.comment("获取物品溯源")
    let req = {user:user, query: {packageID: merchandises[0].PackageID}}
    let res = new Response();
    MerchandiseService.getMerchandiseFlow(req, res, null);
    return res.promise;

  // }).then((result) => {
  //   t.pass(JSON.stringify(result));

  //   t.comment("支付订单金额")
  //   let req = {user:user, body: {orderID: currOrder.OrderID}};
  //   let res = new Response();
  //   OrderService.payOrder(req, res, null);
  //   return res.promise;

  // }).then((result) => {
  //   t.pass(JSON.stringify(result));

  //   t.comment("查看自己钱包")
  //   let req = { user: user };
  //   let res = new Response();
  //   UserService.getAccount(req, res, null);
  //   return res.promise;
  // }).then((result) => {
  //   t.pass(JSON.stringify(result));

  //   t.comment("模拟用户购买")
  //   let req = { user: user, body: {packageID:  merchandises[0].PackageID} };
  //   let res = new Response();
  //   MerchandiseService.purchaseMechandise(req, res, null);
  //   return res.promise;

  // }).then((result) => {
  //   t.pass(JSON.stringify(result));

  //   t.comment("查看钱包历史")
  //   let req = { user: user };
  //   let res = new Response();
  //   UserService.getMoneyHistory(req, res, null);
  //   return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));
    t.comment("结束")
    t.end();
  }).catch((err) => {
    t.fail(err.stack);
    t.end();
  });
})