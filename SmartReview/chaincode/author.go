package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 论文作者结构体
type Author struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Passwd         string   `json:"passwd"`
	Email          string   `json:"email"`
	CommittedPaper []string `json:"paper"`
}

// 作者集合
type AuthorSet struct {
	Authors map[string]string `json:"authors"`
}

func (s *SmartContract) AddtoAuthorSet(ctx contractapi.TransactionContextInterface,name string, id string) int {
	authorSetJson, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return getStateError
	}
	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJson, &authorSet)
	if err != nil {
		return jsonUnMarshalError
	}

	authors := authorSet.Authors
	if _, ok := authors[name];ok{
		return existsError
	}
	authors[name] = id
	newAuthorSet := AuthorSet{
		Authors: authors,
	}
	newAuthorSetJSON, err := json.Marshal(newAuthorSet)
	if err != nil {
		return jsonMarshalError
	}
	err = ctx.GetStub().PutState("authorset", newAuthorSetJSON)
	if err != nil {
		return putStateError
	}
	return 0
}

func (s *SmartContract) CreateAuthor(ctx contractapi.TransactionContextInterface, name string, id string, passwd string, email string) int {
	err2 := s.AddtoAuthorSet(ctx,name,id)
	if err2!=0{
		return err2
	}
	author := Author{
		ID:             id,
		Name:           name,
		Passwd:         passwd,
		Email:          email,
		CommittedPaper: []string{},
	}
	authorJSON, err := json.Marshal(author)
	if err != nil {
		return jsonMarshalError
	}
	err = ctx.GetStub().PutState(id, authorJSON)
	if err != nil {
		return putStateError
	}

	return 0
}

func (s *SmartContract) GetAuthorID(ctx contractapi.TransactionContextInterface, name string)(string, int){
	authorSetJson, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return "", getStateError
	}
	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJson, &authorSet)
	if err != nil {
		return "",jsonUnMarshalError
	}
	if _,ok := authorSet.Authors[name];!ok{
		return "", notExistsError
	}
	return authorSet.Authors[name],0
}

func (s *SmartContract) ReadAuthor(ctx contractapi.TransactionContextInterface, name string) (*Author, int) {
	authorID,err2 := s.GetAuthorID(ctx,name)
	if err2!=0{
		return nil,err2
	}

	authorJSON, err := ctx.GetStub().GetState(authorID)
	if err != nil {
		return nil, getStateError
	}

	var author Author
	err = json.Unmarshal(authorJSON, &author)
	if err != nil {
		return nil, jsonUnMarshalError
	}
	return &author, 0
}

func (s *SmartContract) GetCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, int) {
	author,err := s.ReadAuthor(ctx,name)
	if err != 0 {
		return nil, err
	}
	var committedPaper []Paper
	for _, each := range author.CommittedPaper {
		paper,err := s.GetPaper(ctx,each)
		if err != 0 {
			return nil, err
		}
		committedPaper = append(committedPaper, *paper)
	}
	return committedPaper, 0
}

type comittedPaper struct {
	Name       string   `json:"name"`
	AuthorList []string `json:"authorlist"`
	Reviews    []Review `json:"reviews"`
}

func (s *SmartContract) AuthorCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]comittedPaper, int) {
	papers, err := s.GetCommittedPaper(ctx, name)
	if err != 0 {
		return nil, err
	}
	var res []comittedPaper
	for _, paper := range papers {
		reviews := []Review{}
		if len(paper.ReviewList) > 0{
			for _, review := range paper.ReviewList {
				reviews = append(reviews, Review{
					ReviewerID: review.ReviewerID,
					Content: review.Content,
				})
			}
		} 
		tmp := comittedPaper{
			Name:       paper.Title,
			AuthorList: paper.AuthorList,
			Reviews:    reviews,
		}
		res = append(res, tmp)
	}
	return res, 0
}

func (s *SmartContract) ValidateAuthor(ctx contractapi.TransactionContextInterface,name string, passwd string)(bool,int){
	author,err := s.ReadAuthor(ctx,name)
	if err != 0 {
		return false, err
	}
	if author == nil {
		return false, err
	}
	return author.Passwd == passwd, 0
}