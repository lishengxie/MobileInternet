package main

import (
	"SmartReview/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

func main(){
	reviewChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating smartReview chaincode: %v", err)
	}
	if err := reviewChaincode.Start(); err != nil {
		log.Panicf("Error starting smartReview chaincode: %v", err)
	}
}