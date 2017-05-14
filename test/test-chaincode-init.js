var tape = require('tape');
var _test = require('tape-promise');
var path = require('path');
var test = _test(tape);

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

test('\n\n***** Asset Chaincode Test: init chaincode *****\n\n', (t) => {
  require(__dirname + "/../server/chaincode/create-channel.js").createChannel()
  .then(() => {
    t.pass("create channel success")
    return require(__dirname + "/../server/chaincode/join-channel.js").joinChannel();
  }).then(() => {
    t.pass("join channel success")
    return require(__dirname + "/../server/chaincode/install-chaincode.js").installChaincode();
  }).then(() => {
    return sleep(5000);
  }).then(() => {
    t.pass("install chaincode success")
    return require(__dirname + "/../server/chaincode/instantiate-chaincode.js").instantiateChaincode();
  }).then(() => {
    t.pass("instantiate chaincode success")
    console.log("SUCCESS")
    t.end();
  }).catch((err) => {
    t.fail(err)
    console.log(err)
    t.end();
  });
});

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}