

'use strict';

let express         = require('express');
let session         = require('express-session');
let cookieParser    = require('cookie-parser');
let bodyParser      = require('body-parser');
let url             = require('url');
let cors            = require('cors');
let fs              = require('fs');
let path            = require('path');
let hfc             = require('hfc');
let http            = require('http');
let Redis           = require('ioredis');
let server          = require('./server/server.js');

var serve_static = require('serve-static');

let config = fs.readFileSync(__dirname + '/config/config.json');
config = JSON.parse(config)["config"];

let app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded());
app.options('*', cors());
app.use(cors());

var compression = require('compression');
app.use(compression()); //use compression 
app.use(serve_static(path.join(__dirname, 'public')));

process.env.APP_NAME = 'ASSET_TRADING';

let redis = new Redis(config.redis);
redis.on('error', function(err) {
  console.log(err);
})
redis.set('foo', 'bar');
redis.get('foo', function (err, result) {
  if (result == 'bar')
    console.log('redis connect success');
});

let Constant = require("./server/constant.js");
Constant.redis = redis;

require('cf-deployment-tracker-client').track();
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';
process.env.NODE_ENV = 'production';
process.env.GOPATH = path.resolve(__dirname, 'chaincode');


require(__dirname + "/server/chaincode/chaincode.js").enroll()
.then((admin) => {
  Constant.setAdmin(admin);
  console.log("SUCCESS");
  server.start(app, redis, config.port);
}).catch((err) => {
  console.log(err)
});

// require(__dirname + "/server/chaincode/create-channel.js").createChannel()
// .then(() => {
//   return require(__dirname + "/server/chaincode/join-channel.js").joinChannel();
// }).then(() => {
//   return require(__dirname + "/server/chaincode/install-chaincode.js").installChaincode();
// }).then(() => {
//   return require(__dirname + "/server/chaincode/instantiate-chaincode.js").instantiateChaincode();
// }).then(() => {
//   return require(__dirname + "/server/chaincode/chaincode.js").enroll();
// }).then((admin) => {
//   console.log("SUCCESS");
//   server.start(app, redis, config.port);
// }).catch((err) => {
//   console.log(err)
// });

