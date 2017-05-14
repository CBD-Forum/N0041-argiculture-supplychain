#!/bin/bash

docker pull hyperledger/fabric-orderer:x86_64-1.0.0-alpha
docker pull hyperledger/fabric-peer:x86_64-1.0.0-alpha
docker pull hyperledger/fabric-zookeeper:x86_64-1.0.0-alpha
docker pull hyperledger/fabric-couchdb:x86_64-1.0.0-alpha
docker pull hyperledger/fabric-kafka:x86_64-1.0.0-alpha
docker pull hyperledger/fabric-ca:x86_64-1.0.0-alpha
docker pull hyperledger/fabric-ccenv:x86_64-1.0.0-alpha
docker pull hyperledger/fabric-javaenv:x86_64-1.0.0-alpha

docker tag hyperledger/fabric-orderer:x86_64-1.0.0-alpha hyperledger/fabric-orderer:latest
docker tag hyperledger/fabric-peer:x86_64-1.0.0-alpha hyperledger/fabric-peer:latest
docker tag hyperledger/fabric-zookeeper:x86_64-1.0.0-alpha hyperledger/fabric-zookeeper:latest
docker tag hyperledger/fabric-couchdb:x86_64-1.0.0-alpha hyperledger/fabric-couchdb:latest
docker tag hyperledger/fabric-kafka:x86_64-1.0.0-alpha hyperledger/fabric-kafka:latest
docker tag hyperledger/fabric-ca:x86_64-1.0.0-alpha hyperledger/fabric-ca:latest
docker tag hyperledger/fabric-ccenv:x86_64-1.0.0-alpha hyperledger/fabric-ccenv:latest
docker tag hyperledger/fabric-javaenv:x86_64-1.0.0-alpha hyperledger/fabric-javaenv:latest
