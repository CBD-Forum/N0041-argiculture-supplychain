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

var Response = require("./response.js");

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

let token;
let user;

let currOrder

let bigPackages;

test('\n\n***** Asset Service Test: WarehouseService  *****\n\n', (t) => {
  var currAdmin = null;
  util.enroll("admin", "adminpw")
  .then((admin) => {
    t.pass("init chaincode success");
    currAdmin = admin;

    t.comment("测试农民账号1初始化")
    let registerReq = { body: {username: "warehouse1", password: "warehouse1", type: UserService.UserTypeWarehouse}};
    let res = new Response();
    UserService.register(registerReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("测试登录接口");
    let loginReq = { body: {username: "warehouse1", password: "warehouse1"}}
    let res = new Response();
    UserService.login(loginReq, res, null);
    return res.promise;
  }).then((result) => {
    t.pass(JSON.stringify(result));

    token = result.result.token;
    user = result.result;

    t.comment("开始测试设置账号信息");
    let req = { user:user, body: {token: token, name: "仓库商编号A", identity: "637687198904047682", location: "内蒙古"}};
    let res = new Response();
    UserService.setAccount(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("获取所有的未入库的baozhuang")
    let req = {user:user}
    let res = new Response();
    PackageService.getBigPackageByOrderStatus(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    t.comment("获取所有未入库的订单");
    let req = {user:user, query: { status: OrderService.OrderStatusPackaged} }
    let res = new Response();
    OrderService.getOrdersByStatus(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));

    let orderID;
    let data = result.result;
    orderID = data[0].OrderID;
    currOrder = data[0];

    t.comment("获取订单下的所有包装信息")
    let req = {user:user, body: {orderID: currOrder.OrderID}}
    let res = new Response();
    PackageService.getBigPackagesByOrder(req, res, null);
    return res.promise;

  }).then((result) => {
    t.pass(JSON.stringify(result));
    bigPackages = result.result;

    t.comment("给第一个大包装入库")
    let bigPackage = bigPackages[0].BigPackage;
    let req = {user:user, body: {bigPackageID: bigPackage.BigPackageID, cost: 10}}
    let res = new Response();
    WarehouseService.insertWarehouseStoreIn(req, res, null);
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

  //   t.comment("给第2个大包装入库")
  //   let bigPackage = bigPackages[1].BigPackage;
  //   let req = {user:user, body: {bigPackageID: bigPackage.BigPackageID, cost: 10}}
  //   let res = new Response();
  //   WarehouseService.insertWarehouseStoreIn(req, res, null);
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

    t.comment("获取订单下的仓库包装信息")
    let req = {user:user, query: {orderID: currOrder.OrderID}}
    let res = new Response();
    WarehouseService.getWarehousesByOrder(req, res, null);
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