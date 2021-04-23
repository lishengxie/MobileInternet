package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// 审稿人集合
type ReviewerSet struct {
	Reviewers map[string]string `json:"reviewers"`
}

// 作者集合
type AuthorSet struct {
	Authors map[string]string `json:"authors"`
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	reviewerSet := ReviewerSet{
		Reviewers: make(map[string]string),
	}
	reviewerSetJSON, err := json.Marshal(reviewerSet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("reviewerset", reviewerSetJSON)
	if err != nil {
		return err
	}

	authorSet := AuthorSet{
		Authors: make(map[string]string),
	}
	authorSetJSON, err := json.Marshal(authorSet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("authorset", authorSetJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}
