var hfc = require('fabric-client');
var util = require('util');
var fs = require('fs');
var path = require('path');
let Redis           = require('ioredis');

var utils = require('fabric-client/lib/utils.js');
var Orderer = require('fabric-client/lib/Orderer.js');
var Peer = require('fabric-client/lib/Peer.js');
var EventHub = require('fabric-client/lib/EventHub.js');
var the_user = null;

var testUtil = require(__dirname + '/../server/chaincode/util.js');

var logger = utils.getLogger('test-util');

hfc.addConfigFile(path.join(__dirname, '../config/config.json'));
var ORGS = hfc.getConfigSetting('config');

let Constant = require(__dirname +'/../server/constant.js');
let redis = new Redis(ORGS.redis);
Constant.redis = redis;

module.exports.initChaincode = function () {
  return new Promise((resolve, reject) => {
    require(__dirname + "/../server/chaincode/create-channel.js").createChannel()
    .then(() => {
      return require(__dirname + "/../server/chaincode/join-channel.js").joinChannel();
    }).then(() => {
      return require(__dirname + "/../server/chaincode/install-chaincode.js").installChaincode();
    }).then(() => {
      return require(__dirname + "/../server/chaincode/instantiate-chaincode.js").instantiateChaincode();
    }).then(() => {
      console.log("SUCCESS")
      resolve();
    }).catch((err) => {
      console.log(err)
      reject();
    });
  })
    
}

module.exports.enroll = function (enrollId, enrollSecret) {
  return new Promise((resolve, reject) => {
    var client = new hfc();
    var chain = client.newChain(ORGS.chaincode.channel);

    var caRootsPath = ORGS.orderer.tls_cacerts;
    let data = fs.readFileSync(path.join(__dirname, "../" + caRootsPath));
    let caroots = Buffer.from(data).toString();

    for (let key in ORGS) {
      if (key == ORGS.currOrg && ORGS.hasOwnProperty(key) && typeof ORGS[key].peer1 !== 'undefined') {
        let data = fs.readFileSync(path.join(__dirname, '../' + ORGS[key].peer1['tls_cacerts']));
        let peer = new Peer(
          ORGS[key].peer1.requests,
          {
            pem: Buffer.from(data).toString(),
            'ssl-target-name-override': ORGS[key].peer1['server-hostname']
          }
        );
        chain.addPeer(peer);
      }
    }

    chain.addOrderer(
      new Orderer(
        ORGS.orderer.url,
        {
          'pem': caroots,
          'ssl-target-name-override': ORGS.orderer['server-hostname']
        }
      )
    );

    var name = ORGS[ORGS.currOrg].name;
    return hfc.newDefaultKeyValueStore({
      path: testUtil.storePathForOrg(name)
    }).then((store) => {
      client.setStateStore(store);
      return testUtil.getSubmitter(client, ORGS.currOrg);
    }).then((admin) => {
      the_user = admin;

      Constant.admin = admin;
      Constant.chain = chain;
      Constant.ORGS = ORGS;

      resolve(admin);
    }).catch((err) => {
      console.log(err.stack);
      reject(err);
    })
  })
}