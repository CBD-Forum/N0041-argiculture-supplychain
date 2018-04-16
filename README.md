这里不再更新。
fork到我的项目里了：https://github.com/jasoncodingnow/N0041-argiculture-supplychain

### 基于区块链的农业供应链解决方案demo

#### 说明
* 该demo作为演示，展示了如何在农业场景下使用区块链:
  * 农产品流转溯源
  * 资金流
  * 信息流
  * 金融风控

* 各个角色：
  * 商户： 管理订单、结算支付
  * 农户： 管理资产、管理订单
  * 中间环节： 包装商、仓库、物流
  * 金融： 管理贷款还款、风险


#### 结构说明

* basedocker: node以及npm dependencies 打包成Docker container，避免每次Docker运行都需要download dependencies
* chaincode: chaincode代码以及配套工具:
  * org define
  * channel data tool
  * chaincode source code
  * certificate files
  * docker config files
* config: 配置文档
* public: 前端web页面
* server: 后端服务代码
  * chaincode: chaincode tool for service
  * middleware: service middleware
  * model
  * service
  * server.js: 后端服务初始化入口
* test: 后端service单元测试代码
* pm2.sh: 如果不采用Docker部署，pm2部署也是可以
* app.js: 初始化backend service，连接chaincode


