package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	putStateError = 1
	getStateError = 2
	jsonMarshalError = 3
	jsonUnMarshalError = 4
	distributePaperFailed = 5
	existsError = 6
	notExistsError = 7
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) int {
	reviewerSet := ReviewerSet{
		Reviewers: map[string]string{},
	}
	reviewerSetJSON, err := json.Marshal(reviewerSet)
	if err != nil {
		return jsonMarshalError
	}
	err = ctx.GetStub().PutState("reviewerset", reviewerSetJSON)
	if err != nil {
		return putStateError
	}

	authorSet := AuthorSet{
		Authors: map[string]string{},
	}
	authorSetJSON, err := json.Marshal(authorSet)
	if err != nil {
		return jsonMarshalError
	}
	err = ctx.GetStub().PutState("authorset", authorSetJSON)
	if err != nil {
		return putStateError
	}
	return 0
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}