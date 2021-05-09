package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

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

func (s *SmartContract) GetAuthorSet(ctx contractapi.TransactionContextInterface) (*AuthorSet, error) {
	authorSetJson, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return nil,err
	}
	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJson, &authorSet)
	if err != nil {
		return nil,err
	}
	return &authorSet, nil
}

func (s *SmartContract) AddtoAuthorSet(ctx contractapi.TransactionContextInterface, name string, id string) error {
	authorSet, err := s.GetAuthorSet(ctx)
	if err != nil {
		return err
	}

	if _, ok := authorSet.Authors[name]; ok {
		return fmt.Errorf("Author %s exists", name)
	}
	authorSet.Authors[name] = id
	newAuthorSet := AuthorSet{
		Authors: authorSet.Authors,
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
	authorID, err := s.GetAuthorID(ctx,name)
	if err != nil {
		return err
	}

	author, err := s.ReadAuthor(ctx,authorID)
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
	authorSet, err := s.GetAuthorSet(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := authorSet.Authors[name]; !ok {
		return "", fmt.Errorf("Author %s doesn't exist",name)
	}
	return authorSet.Authors[name], nil
}

func (s *SmartContract) GetAuthorInfo(ctx contractapi.TransactionContextInterface, name string) (*AuthorInfo,error){
	authorID, err := s.GetAuthorID(ctx, name)
	if err != nil {
		return nil, err
	}

	author, err := s.ReadAuthor(ctx,authorID)
	if err!= nil {
		return nil,err
	}

	var committedPaper []string

	for _, each := range author.CommittedPaper {
		paper, err := s.GetPaper(ctx,each)
		if err != nil {
			return nil,err
		}
		committedPaper = append(committedPaper,paper.Title)
	}

	return &AuthorInfo{
		Name: author.Name,
		Passwd: author.Passwd,
		Email: author.Email,
		CommittedPaper: committedPaper,
	}, nil
}

func (s *SmartContract) ReadAuthor(ctx contractapi.TransactionContextInterface, ID string) (*Author, error) {
	authorJSON, err := ctx.GetStub().GetState(ID)
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
	authorID, err := s.GetAuthorID(ctx, name)
	if err != nil {
		return nil, err
	}

	author, err := s.ReadAuthor(ctx, authorID)
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
		var authorList []string
		for _, authorID := range paper.AuthorList{
			author,err := s.ReadAuthor(ctx, authorID)
			if err != nil {
				return nil, err
			}
			authorList = append(authorList,author.Name)
		}
		tmp := comittedPaper{
			Name:       paper.Title,
			AuthorList: authorList,
			Reviews:    paper.ReviewList,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SmartContract) ValidateAuthor(ctx contractapi.TransactionContextInterface, name string, passwd string) (bool, error) {
	authorID, err := s.GetAuthorID(ctx, name)
	if err != nil {
		return false, err
	}

	author, err := s.ReadAuthor(ctx, authorID)
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

	paperID, err := s.GetPaperID(ctx,title)
	if err!= nil{
		return err
	}
	paper, err := s.GetPaper(ctx,paperID)
	if err != nil {
		return err
	}

	rebuttalID := len(paper.ReviewList[reviewerID].RebuttalList)

	rebuttal := Rebuttal{
		AuthorID: authorID,
		ReviewerID: reviewerID,
		RebuttalID: strconv.Itoa(rebuttalID),
		Question: question,
		Reply: "",
		IsReplyed: false,
	}

	rebuttalList := paper.ReviewList[reviewerID].RebuttalList
	rebuttalList[strconv.Itoa(rebuttalID)] = rebuttal

	review := Review{
		ReviewerID: reviewerID,
		Content: paper.ReviewList[reviewerID].Content,
		RebuttalList: rebuttalList,
	}

	newReviewList := paper.ReviewList
	newReviewList[reviewerID] = review
	newPaper := Paper{
		Title:        paper.Title,
		ID:           paper.ID,
		KeyWords: 	  paper.KeyWords,
		AuthorList:   paper.AuthorList,
		ReviewerList: paper.ReviewerList,
		ReviewList:   newReviewList,
		StorePath: paper.StorePath,
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