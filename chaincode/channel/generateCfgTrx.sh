#!/bin/bash

CHANNEL_NAME=$1
if [ -z "$1" ]; then
  echo "Setting channel to default name 'mychannel'"
  CHANNEL_NAME="mychannel"
fi

echo "Channel name - "$CHANNEL_NAME
echo

#Backup the original configtx.yaml
cp $GOPATH/src/github.com/hyperledger/fabric/common/configtx/tool/configtx.yaml $GOPATH/src/github.com/hyperledger/fabric/common/configtx/tool/configtx.yaml.orig
cp configtx.yaml $GOPATH/src/github.com/hyperledger/fabric/common/configtx/tool/configtx.yaml

cd $GOPATH/src/github.com/hyperledger/fabric/
echo "Building configtxgen"
make configtxgen

echo "Generating genesis block"
./build/bin/configtxgen -profile TwoOrgs -outputBlock orderer.block
mv orderer.block $GOPATH/src/datachain/jutian-backend/chaincode/channel/orderer.block

echo "Generating channel configuration transaction"
./build/bin/configtxgen -profile TwoOrgs -outputCreateChannelTx channel.tx -channelID $CHANNEL_NAME
mv channel.tx $GOPATH/src/datachain/jutian-backend/chaincode/channel/$CHANNEL_NAME

#reset configtx.yaml file to its original
cp common/configtx/tool/configtx.yaml.orig common/configtx/tool/configtx.yaml
rm common/configtx/tool/configtx.yaml.orig
