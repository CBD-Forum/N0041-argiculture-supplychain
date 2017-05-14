let http            = require('http');

let UserService     = require('./service/user.service.js');
let AssetService    = require('./service/asset.service.js');
let OrderService    = require('./service/order.service.js');
let PackageService    = require('./service/package.service.js');
let WarehouseService    = require('./service/warehouse.service.js');
let LogisticService    = require('./service/logistic.service.js');
let MerchandiseService  = require('./service/merchandise.service.js');
let BlockchainService  = require('./service/blockchain.service.js');

let TokenValidation = require(__dirname + '/middleware/tokenvalidation.js');

var start = function(app, redis, port) {
  app.get('/test', (req, res, next) => {
    res.status(200).send('test');
  });
  app.post('/user/types', UserService.getUserType);
  app.post('/user', UserService.register);
  app.post('/user/login', UserService.login);

  app.use('/session/', TokenValidation.validate);
  app.get('/session/test', (req, res, next) => {
    res.status(200).send('session test');
  });

  app.post('/session/user/account', UserService.setAccount);
  app.get('/session/user/account', UserService.getUserInfo);
  app.get('/session/user/money/history', UserService.getMoneyHistory);

  app.get('/session/user/loan', UserService.getFarmerLoanableMoney); // 农民获取自己的可借款金额
  app.post('/session/user/loan', UserService.loanByFarmer); // 农民借款
  app.get('/session/loanapplies', UserService.getLoanApplyList); // 金融端获取申请贷款列表
  app.post('/session/loan/approve', UserService.loanApprove); // 农民借款审核通过
  app.get('/session/orders/receivable', UserService.getReceiveableOrdersByFarmer); // 金融端查看农民的应收账单
  app.post('/session/user/orders', UserService.importOrderData);  // 农民导入订单

  app.post('/session/order', OrderService.createOrder);
  app.get('/session/order', OrderService.getOrder);
  app.get('/session/orders/farmer', OrderService.getOrdersByFarmer);  // 农民获取自己的订单
  app.get('/session/orders/owner', OrderService.getOrdersByOwner);    // 商户获取自己的所有订单
  app.get('/session/orders/status', OrderService.getOrdersByStatus);  // 获取各个阶段的订单，不同角色获取到不同的
  app.post('/session/order/match', OrderService.matchOrderByFarmer);  // 农民接单
  app.post('/session/order/pay', OrderService.payOrder);

  app.get('/session/asset/category', AssetService.getCategory);
  app.post('/session/asset', AssetService.createAsset);
  app.get('/session/asset', AssetService.getAsset);
  app.get('/session/assets', AssetService.getFarmerAssetsByOwner);
  app.post('/session/asset/produce', AssetService.produceOrderByFarmer); // 农民生产好了资产, 提交该订单到 待包装 状态

  app.post('/session/package', PackageService.packageOrder); // 包装商打包某个订单
  app.get('/session/packages/order', PackageService.getBigPackagesByOrder); // 包装商获取某个订单的包装
  app.get('/session/packages', PackageService.getBigPackagesByOwner); // 包装商获取自己的包装
  app.get('/session/packages/status', PackageService.getBigPackageByOrderStatus); // 不同角色查看流转到自己的包装

  app.post('/session/warehouse', WarehouseService.insertWarehouseStoreIn); // 仓库把包装入库
  app.get('/session/warehouses/order', WarehouseService.getWarehousesByOrder); // 获取某个订单的库存
  app.get('/session/warehouses', WarehouseService.getWarehousesByOwner); // 获取自己的所有仓库信息

  app.post('/session/logistic', LogisticService.insertLogistic); // 物流把物品发货
  app.get('/session/logistics', LogisticService.getLogisticsByOwner); // 物流查看自己的所有物流
  app.get('/session/logistics/order', LogisticService.getLogisticByOrder); // 查看某个订单的物流

  app.post('/session/merchandise', MerchandiseService.insertMechandise); // 商户收货
  app.get('/session/merchandises', MerchandiseService.getMerchandiseByOwner); // 商户查看自己的所有物品
  app.get('/session/merchandise/flow', MerchandiseService.getMerchandiseFlow); // 商品溯源
  app.post('/session/merchandise/buy', MerchandiseService.purchaseMechandise);  // 模拟用户购买商品

  app.get('/block/height', BlockchainService.getBlockHeight);
  app.get('/block/num', BlockchainService.getBlockByNumber);
  app.get('/block/hash', BlockchainService.getBlockByHash); 

  app.use(function (req, res, next) {
    let err = new Error('Not Found');
    err.status = 404;
    next(err);
  });

  app.use(function (err, req, res, next) {        // = development error handler, print stack trace
    // console.log('Error Handler -', req.url, err);
    let errorCode = err.status || 500;
    res.status(errorCode);
    if (req.bag) {
      req.bag.error = {msg: err.stack, status: errorCode};
      if (req.bag.error.status === 404) {
        req.bag.error.msg = 'Sorry, I cannot locate that file';
      }
    }
    //res.render('template/error', {bag: req.bag});
    res.send({'message':err});
  });

  let server = http.createServer(app).listen(port, function () {
    console.log('Server Up');
    console.log('INFO', 'Startup complete on port', server.address().port);
  });
  server.timeout = 2400000;
}

module.exports.start = start;