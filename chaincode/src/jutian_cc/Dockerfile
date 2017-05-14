from hyperledger/fabric-ccenv
COPY . $GOPATH/src/build-chaincode/
WORKDIR $GOPATH

COPY ./vendor/github.com $GOPATH/src/github.com

RUN go get github.com/op/go-logging 
RUN go install build-chaincode && mv $GOPATH/bin/build-chaincode $GOPATH/bin/assets_unit_test