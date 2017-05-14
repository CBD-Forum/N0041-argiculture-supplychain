/**
 * Copyright 2016 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an 'AS IS' BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

var util = require('util');
var path = require('path');
var fs = require('fs');
var grpc = require('grpc');

var hfc = require('fabric-client');
var utils = require('fabric-client/lib/utils.js');
var Peer = require('fabric-client/lib/Peer.js');
var Orderer = require('fabric-client/lib/Orderer.js');
var EventHub = require('fabric-client/lib/EventHub.js');

var testUtil = require(__dirname + '/util.js');

var logger = utils.getLogger('join-channel');

var the_user = null;
var tx_id = null;
var nonce = null;

hfc.addConfigFile(path.join(__dirname, '../../config/config.json'));
var ORGS = hfc.getConfigSetting('config');

var allEventhubs = [];

var _commonProto = grpc.load(path.join(__dirname, '../../node_modules/fabric-client/lib/protos/common/common.proto')).common;

module.exports.joinChannel = function() {
	return new Promise((resolve, reject) => {
		var org = ORGS.currOrg;

		var client = new hfc();
		var chain = client.newChain(ORGS.chaincode.channel);

		var orgName = ORGS[org].name;

		var targets = [],
			eventhubs = [];

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

		for (let key in ORGS[org]) {
			if (ORGS[org].hasOwnProperty(key)) {
				if (key.indexOf('peer') === 0) {
					data = fs.readFileSync(path.join(__dirname, "../../" + ORGS[org][key]['tls_cacerts']));
					targets.push(
						new Peer(
							ORGS[org][key].requests,
							{
								pem: Buffer.from(data).toString(),
								'ssl-target-name-override': ORGS[org][key]['server-hostname']
							}
						)
					);

					let eh = new EventHub();
					eh.setPeerAddr(
						ORGS[org][key].events,
						{
							pem: Buffer.from(data).toString(),
							'ssl-target-name-override': ORGS[org][key]['server-hostname']
						}
					);
					eventhubs.push(eh);
					allEventhubs.push(eh);
				}
			}
		}

		var name = ORGS[ORGS.currOrg].name;
		return hfc.newDefaultKeyValueStore({
			path: testUtil.storePathForOrg(name)
		}).then((store) => {
			client.setStateStore(store);
			return testUtil.getSubmitter(client, org);
		})
		.then((admin) => {
			logger.info('Successfully enrolled user \'admin\'');
			the_user = admin;

			nonce = utils.getNonce();
			tx_id = chain.buildTransactionID(nonce, the_user);
			var request = {
				targets : targets,
				txId : 	tx_id,
				nonce : nonce
			};

			var eventPromises = [];
			eventhubs.forEach((eh) => {
				eh.connect();
				let txPromise = new Promise((_resolve, _reject) => {
					let handle = setTimeout(_reject, 30000);

					eh.registerBlockEvent((block) => {
						eh.disconnect();
						clearTimeout(handle);

						// in real-world situations, a peer may have more than one channels so
						// we must check that this block came from the channel we asked the peer to join
						if(block.data.data.length === 1) {
							// Config block must only contain one transaction
							// var envelope = _commonProto.Envelope.decode(block.data.data[0]);
							// var payload = _commonProto.Payload.decode(envelope.payload);
							// var channel_header = _commonProto.ChannelHeader.decode(payload.header.channel_header);

							// if (channel_header.channel_id === ORGS.chaincode.channel) {
							// 	logger.info('The new channel has been successfully joined on peer '+ eh.ep._endpoint.addr);
							// 	_resolve();
							// }
							_resolve();
						}
					});
				});

				eventPromises.push(txPromise);
			});

			sendPromise = chain.joinChannel(request);
			return Promise.all([sendPromise].concat(eventPromises));
		}, (err) => {
			logger.error('Failed to enroll user \'admin\' due to error: ' + err.stack ? err.stack : err);
			throw new Error('Failed to enroll user \'admin\' due to error: ' + err.stack ? err.stack : err);
		}).then((results) => {
			logger.info(util.format('Join Channel R E S P O N S E : %j', results));

			if(results[0] && results[0][0] && results[0][0].response && results[0][0].response.status == 200) {
				logger.info(util.format('Successfully joined peers in organization %s to join the channel', orgName));
				resolve();
			} else {
				logger.error(' Failed to join channel');
				throw new Error('Failed to join channel');
			}
		}, (err) => {
			reject('Failed to join channel due to error: ' + err.stack ? err.stack : err);
		});
	})
}

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}