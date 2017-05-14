/**
 * Copyright 2017 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

// This is an end-to-end test that focuses on exercising all parts of the fabric APIs
// in a happy-path scenario
'use strict';

var path = require('path');
var fs = require('fs');
var util = require('util');

var hfc = require('fabric-client');
var utils = require('fabric-client/lib/utils.js');
var Peer = require('fabric-client/lib/Peer.js');
var Orderer = require('fabric-client/lib/Orderer.js');
var EventHub = require('fabric-client/lib/EventHub.js');
var testUtil = require(__dirname + '/util.js');

var Constant = require('../constant.js');

var logger = utils.getLogger('install-chaincode');

hfc.addConfigFile(path.join(__dirname, '../../config/config.json'));
var ORGS = hfc.getConfigSetting('config');
Constant.ORGS = ORGS;

var tx_id = null;
var nonce = null;
var the_user = null;
var allEventhubs = [];

module.exports.instantiateChaincode = function () {
	return new Promise((resolve, reject) => {
		// this is a transaction, will just use org1's identity to
		// submit the request
		var org = ORGS.currOrg;
		var client = new hfc();
		var chain = client.newChain(ORGS.chaincode.channel);

		// 放在Constant下
		Constant.chain = chain;
		Constant.client = client;

		var caRootsPath = ORGS.orderer.tls_cacerts;
		let data = fs.readFileSync(path.join(__dirname, '../../' + caRootsPath));
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

		var orgName = ORGS[org].name;

		var targets = [],
			eventhubs = [];
		// set up the chain to use each org's 'peer1' for
		// both requests and events
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

				let eh = new EventHub();
				eh.setPeerAddr(
					ORGS[key].peer1.events,
					{
						pem: Buffer.from(data).toString(),
						'ssl-target-name-override': ORGS[key].peer1['server-hostname']
					}
				);
				
				eventhubs.push(eh);
				allEventhubs.push(eh);
			}
		}

		var name = ORGS[ORGS.currOrg].name;
		return hfc.newDefaultKeyValueStore({
			path: testUtil.storePathForOrg(name)
		}).then((store) => {
			client.setStateStore(store);
			return testUtil.getSubmitter(client, org);
		}).then((admin) => {

			logger.info('Successfully enrolled user \'admin\'');
			the_user = admin;
			Constant.admin = admin;

			// read the config block from the orderer for the chain
			// and initialize the verify MSPs based on the participating
			// organizations
			return chain.initialize();
		}, (err) => {

			logger.info('Failed to enroll user \'admin\'. ' + err);
			throw new Error('Failed to enroll user \'admin\'. ' + err);

		}).then((success) => {

			nonce = utils.getNonce();
			tx_id = chain.buildTransactionID(nonce, the_user);

			// send proposal to endorser
			var request = {
				chaincodePath: ORGS.chaincode.chaincodePath,
	      chaincodeId: ORGS.chaincode.chaincodeId,
	      chaincodeVersion: ORGS.chaincode.chaincodeVersion,
				fcn: 'init',
				args: ['a', '100', 'b', '200'],
				chainId: ORGS.chaincode.channel,
				txId: tx_id,
				nonce: nonce
			};

			return chain.sendInstantiateProposal(request);

		}, (err) => {

			logger.error('Failed to initialize the chain');
			throw new Error('Failed to initialize the chain');

		}).then((results) => {

			var proposalResponses = results[0];

			var proposal = results[1];
			var header   = results[2];
			var all_good = true;
			for(var i in proposalResponses) {
				let one_good = false;
				if (proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200) {
					one_good = true;
					logger.info('instantiate proposal was good');
				} else {
					logger.error('instantiate proposal was bad');
				}
				all_good = all_good & one_good;
			}
			if (all_good) {
				logger.info(util.format('Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s", metadata - "%s", endorsement signature: %s', proposalResponses[0].response.status, proposalResponses[0].response.message, proposalResponses[0].response.payload, proposalResponses[0].endorsement.signature));
				var request = {
					proposalResponses: proposalResponses,
					proposal: proposal,
					header: header
				};

				// set the transaction listener and set a timeout of 30sec
				// if the transaction did not get committed within the timeout period,
				// fail the test
				var deployId = tx_id.toString();

				var eventPromises = [];
				eventhubs.forEach((eh) => {
					eh.connect();
					let txPromise = new Promise((_resolve, _reject) => {
						let handle = setTimeout(_reject, 30000);

						eh.registerTxEvent(deployId.toString(), (tx, code) => {
							logger.info('The chaincode instantiate transaction has been committed on peer '+ eh.ep._endpoint.addr);
							clearTimeout(handle);
							eh.unregisterTxEvent(deployId);
							eh.disconnect();

							if (code !== 'VALID') {
								logger.error('The chaincode instantiate transaction was invalid, code = ' + code);
								_reject();
							} else {
								logger.info('The chaincode instantiate transaction was valid.');
								_resolve();
							}
						});
					});

					eventPromises.push(txPromise);
				});

				var sendPromise = chain.sendTransaction(request);
				return Promise.all([sendPromise].concat(eventPromises))
				.then((results) => {
					logger.debug('Event promise all complete and testing complete');
					return results[0]; // the first returned value is from the 'sendPromise' which is from the 'sendTransaction()' call
				}).catch((err) => {
					logger.error('Failed to send instantiate transaction and get notifications within the timeout period.');
					throw new Error('Failed to send instantiate transaction and get notifications within the timeout period.');
				});
			} else {
				logger.error('Failed to send instantiate Proposal or receive valid response. Response null or status is not 200. exiting...');
				throw new Error('Failed to send instantiate Proposal or receive valid response. Response null or status is not 200. exiting...');
			}
		}, (err) => {
			logger.error('Failed to send instantiate proposal due to error: ' + err.stack ? err.stack : err);
			throw new Error('Failed to send instantiate proposal due to error: ' + err.stack ? err.stack : err);
		}).then((response) => {
			if (response.status === 'SUCCESS') {
				logger.info('Successfully sent transaction to the orderer.');
				resolve()
			} else {
				logger.error('Failed to order the transaction. Error code: ' + response.status);
				reject('Failed to order the transaction. Error code: ' + response.status);
			}
		}, (err) => {
			logger.error('Failed to send instantiate due to error: ' + err.stack ? err.stack : err);
			reject('Failed to send instantiate due to error: ' + err.stack ? err.stack : err);
		});
	})
};
