package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 论文结构体
type Paper struct {
	Title        string   `json:"title"`
	ID           string   `json:"id"`
	AuthorList   []string `json:"authorlist"`
	ReviewerList []string `json:"reviewerlist"`
	ReviewList   map[string]Review `json:"reviewlist"`
}

func (s *SmartContract) AddtoPaperSet(ctx contractapi.TransactionContextInterface, title string, id string) error {
	paperSetJson, err := ctx.GetStub().GetState("paperset")
	if err != nil {
		return err
	}
	var paperSet PaperSet
	err = json.Unmarshal(paperSetJson, &paperSet)
	if err != nil {
		return err
	}

	if _, ok := paperSet.Papers[title]; ok {
		return fmt.Errorf("Paper already %s exists", title)
	}

	paperSet.Papers[title] = id
	newPaperSet := PaperSet{
		Papers: paperSet.Papers,
	}
	newPaperSetJSON, err := json.Marshal(newPaperSet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("paperset", newPaperSetJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) GetPaperID(ctx contractapi.TransactionContextInterface, title string) (string, error){
	paperSetJson, err := ctx.GetStub().GetState("paperset")
	if err != nil {
		return "",err
	}
	var paperSet PaperSet
	err = json.Unmarshal(paperSetJson, &paperSet)
	if err != nil {
		return "", err
	}

	if _, ok := paperSet.Papers[title]; !ok {
		return "", fmt.Errorf("Paper %s not exist", title)
	}
	return paperSet.Papers[title],nil
}

func (s *SmartContract) AddPaper(ctx contractapi.TransactionContextInterface, title string, id string, authorList string, keywords string) error {
	authorList = strings.Trim(authorList, " ")
	authorListArr := strings.Split(authorList, "/")

	keywords = strings.Trim(keywords, " ")
	keywordsArr := strings.Split(keywords, "/")

	err := s.AddtoPaperSet(ctx,title,id)
	if err != nil {
		return err
	}

	reviewerList, err := s.distributePaper(ctx, keywordsArr)
	if err != nil {
		return err
	}
	paper := Paper{
		Title:        title,
		ID:           id,
		AuthorList:   authorListArr,
		ReviewerList: reviewerList,
		ReviewList:   make(map[string]Review),
	}
	paperJSON, err := json.Marshal(paper)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, paperJSON)
	if err != nil {
		return err
	}

	var author *Author
	for _, each := range authorListArr {

		author, err = s.ReadAuthor(ctx,each)
		if err != nil{
			return err
		}

		newAuthor := Author{
			ID:             author.ID,
			Name:           author.Name,
			Passwd:         author.Passwd,
			Email:          author.Email,
			CommittedPaper: append(author.CommittedPaper, title),
		}
		newAuthorJSON, err := json.Marshal(newAuthor)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(author.ID, newAuthorJSON)
		if err != nil {
			return err
		}
	}

	var reviewer Reviewer
	for _, each := range reviewerList {
		reviewerJSON, err := ctx.GetStub().GetState(each)
		if err != nil {
			return err
		}
		err = json.Unmarshal(reviewerJSON, &reviewer)
		if err != nil {
			return err
		}
		newReviewer := Reviewer{
			ID:              reviewer.ID,
			Name:            reviewer.Name,
			Passwd:          reviewer.Passwd,
			Email:           reviewer.Email,
			ResearchTarget:  reviewer.ResearchTarget,
			ReviewedPaper:   reviewer.ReviewedPaper,
			UNReviewedPaper: append(reviewer.UNReviewedPaper, title),
		}
		newReviewerJSON, err := json.Marshal(newReviewer)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(reviewer.ID, newReviewerJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SmartContract) GetPaper(ctx contractapi.TransactionContextInterface, title string) (*Paper, error) {
	paperID,err := s.GetPaperID(ctx,title)
	paperJSON, err := ctx.GetStub().GetState(paperID)
	if err != nil {
		return nil, err
	}
	if paperJSON == nil {
		return nil, fmt.Errorf("Paper %s doesn't exist", title)
	}
	var paper Paper
	err = json.Unmarshal(paperJSON, &paper)
	if err != nil {
		return nil, err
	}
	return &paper, nil
}

func (s *SmartContract) distributePaper(ctx contractapi.TransactionContextInterface, keywords []string) ([]string, error) {
	//fmt.Println(keywords)
	reviewerSetJson, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return nil, err
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJson, &reviewerSet)
	if err != nil {
		return nil, err
	}
	var reviewerIDSet []string
	for _, value := range reviewerSet.Reviewers {
		reviewerIDSet = append(reviewerIDSet, value)
	}
	return reviewerIDSet, nil
}
