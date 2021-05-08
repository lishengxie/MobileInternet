package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 论文结构体
type Paper struct {
	Title        string            `json:"title"`
	ID           string            `json:"id"`
	KeyWords     []string          `json:"keywords"`
	AuthorList   []string          `json:"authorlist"`
	ReviewerList []string          `json:"reviewerlist"`
	ReviewList   map[string]Review `json:"reviewlist"`
	StorePath	 string		   `json:"storepath"`
}

type PaperInfo struct {
	Title      string            `json:"title"`
	KeyWords   []string          `json:"keywords"`
	AuthorList []string          `json:"authorlist"`
	ReviewList map[string]Review `json:"reviewlist"`
}

func (s *SmartContract) GetPaperSet(ctx contractapi.TransactionContextInterface) (*PaperSet, error) {
	paperSetJson, err := ctx.GetStub().GetState("paperset")
	if err != nil {
		return nil, err
	}
	var paperSet PaperSet
	err = json.Unmarshal(paperSetJson, &paperSet)
	if err != nil {
		return nil, err
	}

	return &paperSet, nil
}

func (s *SmartContract) AddtoPaperSet(ctx contractapi.TransactionContextInterface, title string, id string) error {
	paperSet, err := s.GetPaperSet(ctx)
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

func (s *SmartContract) UpdatePaperSet(ctx contractapi.TransactionContextInterface, origin_title string, new_title string) error {
	paperSet, err := s.GetPaperSet(ctx)
	if err != nil {
		return err
	}

	if _, ok := paperSet.Papers[origin_title]; !ok {
		return fmt.Errorf("Paper not exists", origin_title)
	}

	id := paperSet.Papers[origin_title]
	delete(paperSet.Papers, origin_title)
	paperSet.Papers[new_title] = id

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

func (s *SmartContract) CreatePaper(ctx contractapi.TransactionContextInterface, title string, id string, authorList string, keywords string, path string) error {
	authorList = strings.Trim(authorList, " ")
	authorsArr := strings.Split(authorList, "/")

	keywords = strings.Trim(keywords, " ")
	keywordsArr := strings.Split(keywords, "/")

	err := s.AddtoPaperSet(ctx, title, id)
	if err != nil {
		return err
	}

	reviewerIDList, err := s.distributePaper(ctx, keywordsArr)
	if err != nil {
		return err
	}

	var authorIDList []string
	for _, authorName := range authorsArr {
		authorID, err := s.GetAuthorID(ctx, authorName)
		if err != nil {
			return err
		}
		authorIDList = append(authorIDList, authorID)
	}

	paper := Paper{
		Title:        title,
		ID:           id,
		KeyWords:     keywordsArr,
		AuthorList:   authorIDList,
		ReviewerList: reviewerIDList,
		ReviewList:   make(map[string]Review),
		StorePath: 	  path,
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
	for _, each := range authorIDList {
		author, err = s.ReadAuthor(ctx, each)
		if err != nil {
			return err
		}

		newAuthor := Author{
			ID:             author.ID,
			Name:           author.Name,
			Passwd:         author.Passwd,
			Email:          author.Email,
			CommittedPaper: append(author.CommittedPaper, id),
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

	var reviewer *Reviewer
	for _, each := range reviewerIDList {
		reviewer, err = s.ReadReviewer(ctx, each)
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
			UNReviewedPaper: append(reviewer.UNReviewedPaper, id),
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

func (s *SmartContract) UpdatePaperInfo(ctx contractapi.TransactionContextInterface, origin_title string, new_title string, addedAuthorList string) error {
	err := s.UpdatePaperSet(ctx, origin_title, new_title)
	if err != nil {
		return err
	}

	addedAuthor := strings.Trim(addedAuthorList, " ")
	addedAuthorArr := strings.Split(addedAuthor, "/")

	var addedAuthorIDList []string
	for _, authorName := range addedAuthorArr {
		authorID, err := s.GetAuthorID(ctx, authorName)
		if err != nil {
			return err
		}
		addedAuthorIDList = append(addedAuthorIDList, authorID)
	}

	id, err := s.GetPaperID(ctx, origin_title)
	if err != nil {
		return err
	}
	paper, err := s.GetPaper(ctx, id)
	if err != nil {
		return err
	}

	newPaper := Paper{
		Title:        new_title,
		ID:           paper.ID,
		KeyWords:     paper.KeyWords,
		AuthorList:   append(paper.AuthorList, addedAuthorIDList...),
		ReviewList:   paper.ReviewList,
		ReviewerList: paper.ReviewerList,
		StorePath: paper.StorePath,
	}

	newPaperJSON, err := json.Marshal(newPaper)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(paper.ID, newPaperJSON)
	if err != nil {
		return err
	}

	var author *Author
	for _, each := range addedAuthorIDList {

		author, err = s.ReadAuthor(ctx, each)
		if err != nil {
			return err
		}

		newAuthor := Author{
			ID:             author.ID,
			Name:           author.Name,
			Passwd:         author.Passwd,
			Email:          author.Email,
			CommittedPaper: append(author.CommittedPaper, paper.ID),
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

	return nil
}

func (s *SmartContract) GetPaperID(ctx contractapi.TransactionContextInterface, title string) (string, error) {
	paperSet, err := s.GetPaperSet(ctx)
	if err != nil {
		return "", err
	}

	if _, ok := paperSet.Papers[title]; !ok {
		return "", fmt.Errorf("Paper %s not exist", title)
	}
	return paperSet.Papers[title], nil
}

func (s *SmartContract) GetPaper(ctx contractapi.TransactionContextInterface, ID string) (*Paper, error) {
	paperJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, err
	}
	var paper Paper
	err = json.Unmarshal(paperJSON, &paper)
	if err != nil {
		return nil, err
	}
	return &paper, nil
}

func (s *SmartContract) GetPaperInfo(ctx contractapi.TransactionContextInterface, title string) (*PaperInfo, error) {
	paperID, err := s.GetPaperID(ctx, title)
	if err != nil {
		return nil, err
	}
	paper, err := s.GetPaper(ctx, paperID)
	if err != nil {
		return nil, err
	}

	var authorList []string
	for _, authorID := range paper.AuthorList {
		author, err := s.ReadAuthor(ctx, authorID)
		if err != nil {
			return nil, err
		}
		authorList = append(authorList, author.Name)
	}

	return &PaperInfo{
		Title:      paper.Title,
		KeyWords:   paper.KeyWords,
		AuthorList: authorList,
		ReviewList: paper.ReviewList,
	}, nil
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
