var hfc = require('fabric-client');
var util = require('util');
var fs = require('fs');
var path = require('path');

var testUtil = require(__dirname + '/util.js');
var utils = require('fabric-client/lib/utils.js');
var Orderer = require('fabric-client/lib/Orderer.js');

var the_user = null;

var logger = utils.getLogger('create-channel');

hfc.addConfigFile(path.join(__dirname, '../../config/config.json'));
var ORGS = hfc.getConfigSetting('config');

module.exports.createChannel = function() {
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

		// Acting as a client in org1 when creating the channel
		var name = ORGS[ORGS.currOrg].name;
		return hfc.newDefaultKeyValueStore({
			path: testUtil.storePathForOrg(name)
		}).then((store) => {
			client.setStateStore(store);
			return testUtil.getSubmitter(client, ORGS.currOrg);
		})
		.then((admin) => {
			logger.info('Successfully enrolled user \'admin\'');
			the_user = admin;

			// readin the envelope to send to the orderer
			data = fs.readFileSync(__dirname + '/../../chaincode/channel/' + ORGS.chaincode.channelfile);
			var request = {
				envelope : data
			};
			// send to orderer
			return chain.createChannel(request);
		}, (err) => {
			throw new Error('Failed to enroll user \'admin\'. ' + err);
		})
		.then((response) => {
			logger.debug(' response ::%j',response);

			if (response && response.status === 'SUCCESS') {
				logger.info('Successfully created the channel.');
				return sleep(5000);
			} else {
				throw new Error('Failed to create the channel. ');
			}
		}, (err) => {
			throw new Error('Failed to initialize the channel: ' + err.stack ? err.stack : err);
		})
		.then((nothing) => {
			logger.info('Successfully waited to make sure new channel was created.');
			resolve();
		}, (err) => {
			reject('Failed to sleep due to error: ' + err.stack ? err.stack : err);
		});
	})
};

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}
