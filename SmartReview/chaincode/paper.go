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
	ReviewList   []Review `json:"reviewlist"`
}

func (s *SmartContract) AddPaper(ctx contractapi.TransactionContextInterface, title string, id string, authorList string, keywords string) error {
	authorList = strings.Trim(authorList, " ")
	authorListArr := strings.Split(authorList, "/")

	keywords = strings.Trim(keywords, " ")
	keywordsArr := strings.Split(keywords, "/")

	exists, err := s.AssetExists(ctx, title)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s exists", title)
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
		ReviewList:   []Review{},
	}
	paperJSON, err := json.Marshal(paper)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(title, paperJSON)
	if err != nil {
		return err
	}

	var author Author
	var authorID string

	authorSetJSON, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return err
	}

	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJSON, &authorSet)
	if err != nil {
		return err
	}
	for _, each := range authorListArr {
		if _, ok := authorSet.Authors[each]; !ok {
			return fmt.Errorf("Author %s not exists", each)
		}
		authorID = authorSet.Authors[each]
		authorJSON, err := ctx.GetStub().GetState(authorID)
		if err != nil {
			return err
		}
		err = json.Unmarshal(authorJSON, &author)
		if err != nil {
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
	reviewerSetJSON, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return err
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJSON, &reviewerSet)
	if err != nil {
		return err
	}

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
	paperJSON, err := ctx.GetStub().GetState(title)
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
