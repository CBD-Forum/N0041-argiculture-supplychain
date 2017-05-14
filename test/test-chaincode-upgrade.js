let fs              = require('fs');

var tape = require('tape');
var _test = require('tape-promise');
var path = require('path');
var test = _test(tape);

process.env.GOPATH = path.resolve(__dirname, '../chaincode');

let config = fs.readFileSync(__dirname + '/../config/config.json');
config = JSON.parse(config)["config"];

test('\n\n***** Asset Chaincode Test: restart chaincode *****\n\n', (t) => {
  return require(__dirname + "/../server/chaincode/install-chaincode.js").installChaincode()
  .then(() => {
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