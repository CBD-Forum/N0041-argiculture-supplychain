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

var path = require('path');
var fs = require('fs-extra');
var os = require('os');

var jsrsa = require('jsrsasign');
var KEYUTIL = jsrsa.KEYUTIL;

var hfc = require('fabric-client');
var copService = require('fabric-ca-client/lib/FabricCAClientImpl.js');
var User = require('fabric-client/lib/User.js');
var CryptoSuite = require('fabric-client/lib/impl/CryptoSuite_ECDSA_AES.js');
var KeyStore = require('fabric-client/lib/impl/CryptoKeyStore.js');
var ecdsaKey = require('fabric-client/lib/impl/ecdsa/key.js');

module.exports.readFile = readFile;

hfc.addConfigFile(path.join(__dirname, '../../config/config.json'));
var ORGS = hfc.getConfigSetting('config');

// directory for file based KeyValueStore
module.exports.KVS = '/tmp/hfc-test-kvs';
module.exports.storePathForOrg = function(org) {
	return module.exports.KVS + '_' + org;
};

function getSubmitter(username, password, client, loadFromConfig, userOrg) {
	var caUrl = ORGS[userOrg].ca;

	return client.getUserContext(username)
	.then((user) => {
		return new Promise((resolve, reject) => {
			if (user && user.isEnrolled()) {
				console.log('Successfully loaded member from persistence');
				return resolve(user);
			}
			if (!loadFromConfig) {
				var cop = new copService(caUrl);

				var member;
				return cop.enroll({
					enrollmentID: username,
					enrollmentSecret: password
				}).then((enrollment) => {
					console.log('Successfully enrolled user \'' + username + '\'');

					member = new User(username, client);
					return member.setEnrollment(enrollment.key, enrollment.certificate, ORGS[userOrg].mspid);
				}).then(() => {
					return client.setUserContext(member);
				}).then(() => {
					return resolve(member);
				}).catch((err) => {
					reject('Failed to enroll and persist user. Error: ' + err.stack ? err.stack : err);
				});
			} else {
				// need to load private key and pre-enrolled certificate from files based on the MSP
				// config directory structure:
				// <config>
				//    \_ keystore
				//       \_ admin.pem  <<== this is the private key saved in PEM file
				//    \_ signcerts
				//       \_ admin.pem  <<== this is the signed certificate saved in PEM file

				// first load the private key and save in the BCCSP's key store
				var privKeyPEM = path.join(__dirname, '../chaincode/msp/local/keystore/admin.pem');
				var pemData, member;
				return readFile(privKeyPEM)
				.then((data) => {
					pemData = data;
					// default crypto suite uses $HOME/.hfc-key-store as key store
					var kspath = CryptoSuite.getDefaultKeyStorePath();
					var testKey;
					return new KeyStore({
						path: kspath
					});
				}).then((store) => {
					var rawKey = KEYUTIL.getKey(pemData.toString());
					testKey = new ecdsaKey(rawKey);
					return store.putKey(testKey);
				}).then((value) => {
					// next save the certificate in a serialized user enrollment in the state store
					var certPEM = path.join(__dirname, '../../chaincode/msp/local/signcerts/admin.pem');
					return readFile(certPEM);
				}).then((data) => {
					member = new User(username, client);
					return member.setEnrollment(testKey, data.toString(), ORGS[userOrg].mspid);
				}).then(() => {
					return client.setUserContext(member);
				}).then((user) => {
					return resolve(user);
				}).catch((err) => {
					reject(new Error('Failed to load key or certificate and save to local stores. ' + err));
				});
			}
		});
	});
}

function readFile(path) {
	return new Promise((resolve, reject) => {
		fs.readFile(path, (err, data) => {
			if (!!err)
				reject(new Error('Failed to read file ' + path + ' due to error: ' + err));
			else
				resolve(data);
		});
	});
}

module.exports.getSubmitter = function(client, loadFromConfig, org) {
	if (arguments.length < 2) throw new Error('"client" and "test" are both required parameters');

	var fromConfig, userOrg;
	if (typeof loadFromConfig === 'boolean') {
		fromConfig = loadFromConfig;
	} else {
		fromConfig = false;
	}

	if (typeof loadFromConfig === 'string') {
		userOrg = loadFromConfig;
	} else {
		if (typeof org === 'string') {
			userOrg = org;
		} else {
			userOrg = 'org1';
		}
	}

	return getSubmitter('admin', 'adminpw', client, fromConfig, userOrg);
};
