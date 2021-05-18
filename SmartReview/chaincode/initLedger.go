package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
)

type SmartContract struct {
	contractapi.Contract
}

// 用户集合
type UserSet struct {
	Users map[string]string `json:"users"`
}

// 论文集合
type PaperSet struct {
	Papers map[string]string `json:"papers"`
}

type SimilarityPair struct {
	Pair  string	`json:"pair"`
	Score float64	`json:"score"`
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	userSet := UserSet{
		Users: make(map[string]string),
	}
	userSetJSON, err := json.Marshal(userSet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("userset", userSetJSON)
	if err != nil {
		return err
	}

	paperSet := PaperSet{
		Papers: make(map[string]string),
	}
	paperSetJSON, err := json.Marshal(paperSet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("paperset", paperSetJSON)
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

func (s *SmartContract) AddSimilarityPair(ctx contractapi.TransactionContextInterface, pair string, score string) error {
	scoreFloat,err := strconv.ParseFloat(score,64)
	if err != nil {
		return err
	}

	similarityPair := SimilarityPair{
		Pair: pair,
		Score: scoreFloat,
	}

	similarityPairJSON, err := json.Marshal(similarityPair)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(pair,similarityPairJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) GetSimilarity(ctx contractapi.TransactionContextInterface, pair string) (float64, error) {
	similarityPairJSON, err := ctx.GetStub().GetState(pair)
	if err != nil {
		return 0.0, err
	}

	var similarityPair SimilarityPair

	err = json.Unmarshal(similarityPairJSON, &similarityPair)
	if err != nil {
		return 0.0, err
	}

	return similarityPair.Score, nil
}
