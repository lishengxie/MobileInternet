package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 论文作者结构体
type Author struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Passwd         string   `json:"passwd"`
	Email          string   `json:"email"`
	CommittedPaper []string `json:"committedpaper"`
}

type AuthorInfo struct {
	Name           string   `json:"name"`
	Passwd         string   `json:"passwd"`
	Email          string   `json:"email"`
	CommittedPaper []string `json:"committedpaper"`
}


func (s *SmartContract) AddtoAuthorSet(ctx contractapi.TransactionContextInterface, name string, id string) error {
	authorSetJson, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return err
	}
	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJson, &authorSet)
	if err != nil {
		return err
	}

	authors := authorSet.Authors
	if _, ok := authors[name]; ok {
		return fmt.Errorf("Author %s exists", name)
	}
	authors[name] = id
	newAuthorSet := AuthorSet{
		Authors: authors,
	}
	newAuthorSetJSON, err := json.Marshal(newAuthorSet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("authorset", newAuthorSetJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) CreateAuthor(ctx contractapi.TransactionContextInterface, name string, id string, passwd string, email string) error {
	err := s.AddtoAuthorSet(ctx, name, id)
	if err != nil {
		return err
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
		return err
	}
	err = ctx.GetStub().PutState(id, authorJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) UpdateAuthorInfo(ctx contractapi.TransactionContextInterface, name string, old_passwd string, new_passwd string, email string) error {
	author, err := s.ReadAuthor(ctx,name)
	if err != nil{
		return err
	}

	if old_passwd != author.Passwd {
		return fmt.Errorf("wrong passwd")
	}

	newAuthor := Author{
		ID : author.ID,
		Name: author.Name,
		Passwd: new_passwd,
		Email: email,
		CommittedPaper: author.CommittedPaper,
	}
	authorJSON, err := json.Marshal(newAuthor)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(author.ID, authorJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) GetAuthorID(ctx contractapi.TransactionContextInterface, name string) (string, error) {
	authorSetJson, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return "", err
	}
	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJson, &authorSet)
	if err != nil {
		return "", err
	}
	if _, ok := authorSet.Authors[name]; !ok {
		return "", fmt.Errorf("Author %s doesn't exist",name)
	}
	return authorSet.Authors[name], nil
}

func (s *SmartContract) GetAuthorInfo(ctx contractapi.TransactionContextInterface, name string) (*AuthorInfo,error){
	author, err := s.ReadAuthor(ctx,name)
	if err!= nil {
		return nil,err
	}
	return &AuthorInfo{
		Name: author.Name,
		Passwd: author.Passwd,
		Email: author.Email,
		CommittedPaper: author.CommittedPaper,
	}, nil
}

func (s *SmartContract) ReadAuthor(ctx contractapi.TransactionContextInterface, name string) (*Author, error) {
	authorID, err := s.GetAuthorID(ctx, name)
	if err != nil {
		return nil, err
	}

	authorJSON, err := ctx.GetStub().GetState(authorID)
	if err != nil {
		return nil, err
	}

	var author Author
	err = json.Unmarshal(authorJSON, &author)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (s *SmartContract) GetCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	author, err := s.ReadAuthor(ctx, name)
	if err != nil {
		return nil, err
	}
	var committedPaper []Paper
	for _, each := range author.CommittedPaper {
		paper, err := s.GetPaper(ctx, each)
		if err != nil {
			return nil, err
		}
		committedPaper = append(committedPaper, *paper)
	}
	return committedPaper, nil
}

type comittedPaper struct {
	Name       string   `json:"name"`
	AuthorList []string `json:"authorlist"`
	Reviews    map[string]Review `json:"reviews"`
}

func (s *SmartContract) AuthorCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]comittedPaper, error) {
	papers, err := s.GetCommittedPaper(ctx, name)
	if err != nil {
		return nil, err
	}
	var res []comittedPaper
	for _, paper := range papers {
		tmp := comittedPaper{
			Name:       paper.Title,
			AuthorList: paper.AuthorList,
			Reviews:    paper.ReviewList,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SmartContract) ValidateAuthor(ctx contractapi.TransactionContextInterface, name string, passwd string) (bool, error) {
	author, err := s.ReadAuthor(ctx, name)
	if err != nil {
		return false, err
	}
	if author == nil {
		return false, fmt.Errorf("Author %s not exist",author)
	}
	return author.Passwd == passwd, nil
}

func (s *SmartContract) AddRebuttal(ctx contractapi.TransactionContextInterface,title string, author_name string, reviewerID string, question string)error{
	authorID, err := s.GetAuthorID(ctx,author_name)
	if err != nil {
		return err
	}
	paper, err := s.GetPaper(ctx,title)
	if err != nil {
		return err
	}

	rebuttal := Rebuttal{
		AuthorID: authorID,
		ReviewerID: reviewerID,
		Question: question,
		Reply: "",
		IsReplyed: false,
	}

	review := Review{
		ReviewerID: reviewerID,
		Content: paper.ReviewList[reviewerID].Content,
		RebuttalList: rebuttal,
	}

	newReviewList := paper.ReviewList
	newReviewList[reviewerID] = review
	newPaper := Paper{
		Title:        paper.Title,
		ID:           paper.ID,
		AuthorList:   paper.AuthorList,
		ReviewerList: paper.ReviewerList,
		ReviewList:   newReviewList,
	}

	newPaperJSON, err := json.Marshal(newPaper)
	if err != nil{
		return err
	}

	err = ctx.GetStub().PutState(paper.ID,newPaperJSON)
	if err != nil{
		return err
	}
	return nil
}