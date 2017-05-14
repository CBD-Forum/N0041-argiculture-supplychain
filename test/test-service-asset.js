var tape = require('tape');
var _test = require('tape-promise');
var path = require('path');
var test = _test(tape);

var util = require("./util.js");

var Constant = require("../server/constant.js");
var UserService = require(__dirname + "/../server/service/user.service.js");
var OrderService = require(__dirname + "/../server/service/order.service.js");
var AssetService = require(__dirname + "/../server/service/asset.service.js");

var Response = require("./response.js");

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

let token;
let user;

let financialToken;
let financial;

let currOrder

test('\n\n***** Asset Service Test: asset  *****\n\n', (t) => {
  var currAdmin = null;
  util.enroll("admin", "adminpw")
  .then((admin) => {
    t.pass("init chaincode success");
    currAdmin = admin;

    t.comment("测试农民账号1初始化")
    let registerReq = { body: {username: "farmer2", password: "farmer2", type: UserService.UserTypeFarmer}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试登录接口");
    let loginReq = { body: {username: "farmer2", password: "farmer2"}}
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

    t.comment("测试导入账单xinxi");
    let req = {user:user}
    let res = new Response();
    UserService.importOrderData(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("获取用户信息");
    let req = {user:user}
    let res = new Response();
    UserService.getUserInfo(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("开始测试获取可借款金额");
    let req = {user:user}
    let res = new Response();
    UserService.getFarmerLoanableMoney(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("获取所有未接的订单");
    let req = {user:user, query: {status: OrderService.OrderStatusInit}}
    let res = new Response();
    OrderService.getOrdersByStatus(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    let orderID;
    let data = result.result;
    orderID = data[0].OrderID;
    currOrder = data[0];

    t.comment("农民接单")
    let req = {user:user, body: {orderID: orderID}}
    let res = new Response();
    OrderService.matchOrderByFarmer(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("开始测试获取可借款金额");
    let req = {user:user}
    let res = new Response();
    UserService.getFarmerLoanableMoney(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民借款");
    let req = {user:user, body: {loan: result.result}}
    let res = new Response();
    UserService.loanByFarmer(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试金融端注册")
    let registerReq = { body: {username: "financial1", password: "financial1", type: UserService.UserTypeFinancial}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("金融端登录接口");
    let loginReq = { body: {username: "financial1", password: "financial1"}}
    let res = new Response();
    UserService.login(loginReq, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    financialTokenn = result.result.token;
    financial = result.result;

    t.comment("开始测试设置账号信息");
    let req = { user:financial, body: {token: token, name: "金融编号A", identity: "637687198904047682", location: "内蒙古"}};
    let res = new Response();
    console.log(req)
    UserService.setAccount(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("开始获取申请贷款列表");
    let req = {user:financial}
    let res = new Response();
    UserService.getLoanApplyList(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    let data = result.result;
    let farmerID = data[0].FarmerID;

    t.comment("获取该农民的贷款订单");
    let req = {user:financial, query: {farmerID: farmerID}}
    let res = new Response();
    UserService.getReceiveableOrdersByFarmer(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    let data = result.result;
    let farmerID = data[0].FarmerID;
    t.comment("审核通过");
    let req = {user:financial, body: {farmerID: farmerID}}
    let res = new Response();
    UserService.loanApprove(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民获取钱包");
    let req = {user:user}
    let res = new Response();
    UserService.getUserInfo(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民获取自己接的单子");
    let req = {user:user}
    let res = new Response();
    OrderService.getOrdersByFarmer(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment(JSON.stringify(currOrder))
    t.comment("农民检查下自己的资产是否够")
    let req = {user:user, query: {categoryID: currOrder.CategoryID}}
    let res = new Response();
    AssetService.getAssetByCategoryIDNOwner(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment(JSON.stringify(currOrder))
    t.comment("假设农民资产不够，需要再生产")
    let req = {user:user, body: {categoryID: currOrder.CategoryID, amount: currOrder.Amount, materialCost: 1, laborCost: 1}}
    let res = new Response();
    AssetService.createAsset(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民获取自己的所有资产")
    let req = {user:user}
    let res = new Response();
    AssetService.getFarmerAssetsByOwner(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民生产完成")
    let req = {user:user, body: {orderID: currOrder.OrderID}};
    let res = new Response();
    AssetService.produceOrderByFarmer(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民查看下当前订单进展")
    let req = {user:user, query: {orderID: currOrder.OrderID}};
    let res = new Response();
    OrderService.getOrder(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("农民获取自己的所有资产")
    let req = {user:user}
    let res = new Response();
    AssetService.getFarmerAssetsByOwner(req, res, null);
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