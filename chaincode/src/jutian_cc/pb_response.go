package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type PbResponse struct {
	Success bool
	Data    interface{}
	Err     string
}

func wrapToPbResponse(data []byte, err error) pb.Response {
	if err != nil {
		response := PbResponse{
			false,
			nil,
			err.Error(),
		}
		bytes, err := StructToJSONBytes(response)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(bytes)
	}

	response := PbResponse{
		true,
		data,
		"",
	}
	bytes, err := StructToJSONBytes(response)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}

func wrapStructToPbResponse(data interface{}, err error) pb.Response {
	LogStruct(data)
	if err != nil {
		response := PbResponse{
			false,
			nil,
			err.Error(),
		}
		LogStruct(response)
		bytes, err := StructToJSONBytes(response)
		LogMessage(string(bytes))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(bytes)
	}

	response := PbResponse{
		true,
		data,
		"",
	}
	LogStruct(response)
	bytes, err := StructToJSONBytes(response)
	LogMessage(string(bytes))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bytes)
}
