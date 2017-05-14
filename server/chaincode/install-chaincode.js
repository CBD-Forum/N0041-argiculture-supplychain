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
var testUtil = require(__dirname + '/util.js');

var logger = utils.getLogger('install-chaincode');

hfc.addConfigFile(path.join(__dirname, '../../config/config.json'));
var ORGS = hfc.getConfigSetting('config');

var tx_id = null;
var nonce = null;
var the_user = null;

module.exports.installChaincode = function() {
	return new Promise((resolve, reject) => {
		var org = ORGS.currOrg;

		var client = new hfc();
		var chain = client.newChain(ORGS.chaincode.channel);

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

		let priameyKeySetted = false;
		var targets = [];
		for (let key in ORGS[org]) {
			if (ORGS[org].hasOwnProperty(key)) {
				if (key.indexOf('peer') === 0) {
					let data = fs.readFileSync(path.join(__dirname, '../../' + ORGS[org][key]['tls_cacerts']));
					let peer = new Peer(
						ORGS[org][key].requests,
						{
							pem: Buffer.from(data).toString(),
							'ssl-target-name-override': ORGS[org][key]['server-hostname']
						}
					);

					targets.push(peer);
					chain.addPeer(peer);
					if (!priameyKeySetted) {
						chain.setPrimaryPeer(peer);
						priameyKeySetted = true;
					}
				}
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

			nonce = utils.getNonce();
			tx_id = chain.buildTransactionID(nonce, the_user);

			// send proposal to endorser
			var request = {
				targets: targets,
				chaincodePath: ORGS.chaincode.chaincodePath,
	      chaincodeId: ORGS.chaincode.chaincodeId,
	      chaincodeVersion: ORGS.chaincode.chaincodeVersion,
				txId: tx_id,
				nonce: nonce
			};

			return chain.sendInstallProposal(request);
		},
		(err) => {
			logger.error('Failed to enroll user \'admin\'. ' + err);
			throw new Error('Failed to enroll user \'admin\'. ' + err);
		}).then((results) => {
			var proposalResponses = results[0];

			var proposal = results[1];
			var header   = results[2];
			var all_good = true;
			for(var i in proposalResponses) {
				let one_good = false;
				if (proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200) {
					one_good = true;
					logger.info('install proposal was good');
				} else {
					logger.error('install proposal was bad');
				}
				all_good = all_good & one_good;
			}
			if (all_good) {
				logger.info(util.format('Successfully sent install Proposal and received ProposalResponse: Status - %s', proposalResponses[0].response.status));
				resolve();
			} else {
				reject('Failed to send install Proposal or receive valid response. Response null or status is not 200. exiting...');
			}
		},
		(err) => {
			logger.error('Failed to send install proposal due to error: ' + err.stack ? err.stack : err);
			reject('Failed to send install proposal due to error: ' + err.stack ? err.stack : err);
		});
	})
}