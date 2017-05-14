'use strict';

var path = require('path');
var fs = require('fs');
var util = require('util');

var hfc = require('fabric-client');

var utils = require('fabric-client/lib/utils.js');
var Peer = require('fabric-client/lib/Peer.js');
var Orderer = require('fabric-client/lib/Orderer.js');
var EventHub = require('fabric-client/lib/EventHub.js');

var logger = utils.getLogger('invoke-chaincode');

var Constant = require(__dirname + '/../constant.js');
var testUtil = require('./util.js');

hfc.addConfigFile(path.join(__dirname, '../../config/config.json'));
var ORGS = hfc.getConfigSetting('config');

var org = ORGS[ORGS.currOrg];

class Chaincode {
  static invoke(functionName, args, emrollUser) {
    return new Promise((resolve, reject) => {
      let invokeResult = null;

      var nonce = utils.getNonce();
      var eventhub;
      var eventhubs = [];
      var chain = Constant.chain;

      var tx_id = chain.buildTransactionID(nonce, emrollUser);

      var requestArgs = [];
      for(let key in args) {
        requestArgs.push(args[key] + "");
      }

      var request = {
        chainId: Constant.ORGS.chaincode.channel,
        chaincodeId: Constant.ORGS.chaincode.chaincodeId,
        chaincodeVersion: Constant.ORGS.chaincode.chaincodeVersion,
        fcn: functionName,
        args: requestArgs,
        txId: tx_id,
        nonce: nonce,
      };

      let data = fs.readFileSync(path.join(__dirname, '../../' + Constant.ORGS[Constant.ORGS.currOrg].peer1['tls_cacerts']));
      eventhub = new EventHub();
      eventhub.setPeerAddr(
        Constant.ORGS[Constant.ORGS.currOrg].peer1.events,
        {
          pem: Buffer.from(data).toString(),
          'ssl-target-name-override': Constant.ORGS[Constant.ORGS.currOrg].peer1['server-hostname']
        }
      );
      eventhubs.push(eventhub);

      chain.sendTransactionProposal(request).then((results) => {
        var proposalResponses = results[0];

        var proposal = results[1];
        var header   = results[2];
        var all_good = true;
        for(var i in proposalResponses) {
          let one_good = false;
          if (proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200) {
            one_good = true;
            logger.info('transaction proposal was good');
          } else {
            logger.error(util.format('transaction proposal was bad, %s', JSON.stringify(proposalResponses)));
          }
          all_good = all_good & one_good;
        }
        if (all_good) {
          logger.info(util.format('Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s", metadata - "%s"', proposalResponses[0].response.status, proposalResponses[0].response.message, proposalResponses[0].response.payload));
          var request = {
            proposalResponses: proposalResponses,
            proposal: proposal,
            header: header
          };
          if (invokeResult == null) {
            invokeResult = proposalResponses[0].response.payload;
          }
          // resolve(invokeResult.toString()); // the first returned value is from the 'sendPromise' which is from the 'sendTransaction()' call

          var deployId = tx_id.toString();

          var eventPromises = [];
          eventhubs.forEach((eh) => {
            eh.connect();
            let txPromise = new Promise((resolve, reject) => {
              let handle = setTimeout(reject, 30000);

              eh.registerTxEvent(deployId.toString(), (tx, code) => {
                clearTimeout(handle);
                eh.disconnect();
                eh.unregisterTxEvent(deployId);

                if (code !== 'VALID') {
                  logger.error('The balance transfer transaction was invalid, code = ' + code);
                  reject();
                } else {
                  logger.info('The balance transfer transaction has been committed on peer '+ eh.ep._endpoint.addr);
                  resolve();
                }
              });
            });

            eventPromises.push(txPromise);
          });

          var sendPromise = chain.sendTransaction(request);
          return Promise.all([sendPromise].concat(eventPromises))
          .then((results) => {
            logger.debug(' event promise all complete and testing complete');
            resolve(invokeResult.toString()); // the first returned value is from the 'sendPromise' which is from the 'sendTransaction()' call
          }).catch((err) => {
            logger.error('Failed to send transaction and get notifications within the timeout period.');
            reject('Failed to send transaction and get notifications within the timeout period.');
          });
        } else {
          logger.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
          reject('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
        }
      })
    })
  }

  static query(functionName, args, enrollUser) {
    return new Promise((resolve, reject) => {
      var chain = Constant.chain;

      var nonce = utils.getNonce();
      var tx_id = chain.buildTransactionID(nonce, enrollUser);

      var requestArgs = [];
      for(let key in args) {
        requestArgs.push(args[key] + "");
      }

      var request = {
        chainId: Constant.ORGS.chaincode.channel,
        chaincodeId: Constant.ORGS.chaincode.chaincodeId,
        chaincodeVersion: Constant.ORGS.chaincode.chaincodeVersion,
        txId: tx_id,
        nonce: nonce,
        fcn: functionName,
        args: requestArgs
      };

      chain.queryByChaincode(request)
      .then((response_payloads) => {
        var result = [];
        for(let i = 0; i < response_payloads.length; i++) {
          result.push(response_payloads[i].toString('utf8'))
        }
        resolve(result);
      }, (err) => {
        logger.error('error when query');
        logger.error(err.stack ? err.stack : err);
        reject(err);
      });
    })
  }

  static enroll() {
    return new Promise((resolve, reject) => {
      var client = new hfc();
      var chain = client.newChain(ORGS.chaincode.channel);

      var caRootsPath = ORGS.orderer.tls_cacerts;
      let data = fs.readFileSync(path.join(__dirname, "../../" + caRootsPath));
      let caroots = Buffer.from(data).toString();

      chain.addOrderer(
        new Orderer(
          ORGS.orderer.url,
          {
            'pem': caroots,
            'ssl-target-name-override': ORGS.orderer['server-hostname']
          }
        )
      );

      let priameyKeySetted = false;
      for (let key in ORGS) {
      if (key == ORGS.currOrg && ORGS.hasOwnProperty(key) && typeof ORGS[key].peer1 !== 'undefined') {
        let data = fs.readFileSync(path.join(__dirname, '../../' + ORGS[key].peer1['tls_cacerts']));
        let peer = new Peer(
          ORGS[key].peer1.requests,
          {
            pem: Buffer.from(data).toString(),
            'ssl-target-name-override': ORGS[key].peer1['server-hostname']
          }
        );
        chain.addPeer(peer);
        if (!priameyKeySetted) {
          chain.setPrimaryPeer(peer);
          priameyKeySetted = true;
        }

        let eh = new EventHub();
        eh.setPeerAddr(
          ORGS[key].peer1.events,
          {
            pem: Buffer.from(data).toString(),
            'ssl-target-name-override': ORGS[key].peer1['server-hostname']
          }
        );
      }
    }

      Constant.chain = chain;
      Constant.ORGS = ORGS;

      var name = ORGS[ORGS.currOrg].name;
      return hfc.newDefaultKeyValueStore({
        path: testUtil.storePathForOrg(name)
      }).then((store) => {
        client.setStateStore(store);
        return testUtil.getSubmitter(client, ORGS.currOrg);
      }).then((admin) => {
        resolve(admin);
      }).catch((err) => {
        console.log(err.stack);
        reject(err);
      })
    })
  }
}

module.exports = Chaincode;