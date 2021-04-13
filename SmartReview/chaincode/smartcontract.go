package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// 审稿人结构体
type Reviewer struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Passwd          string   `json:"passwd"`
	Email           string   `json:"email"`
	ResearchTarget  []string `json:"researchtarget"`
	ReviewedPaper   []string `json:"reviewedpaper"`
	UNReviewedPaper []string `json:"unreviewedpaper"`
}

// 论文作者结构体
type Author struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Passwd         string   `json:"passwd"`
	Email          string   `json:"email"`
	CommittedPaper []string `json:"paper"`
}

// 审稿内容结构体
type Review struct {
	ReviewerName string `json:"reviewername"`
	Content      string `json:"content"`
}

// 论文结构体
type Paper struct {
	Title        string   `json:"title"`
	ID           string   `json:"id"`
	AuthorList   []string `json:"authorlist"`
	ReviewerList []string `json:"reviewerlist"`
	ReviewList   []Review `json:"reviewlist"`
}

// 审稿人集合
type ReviewerSet struct {
	ReviewerNameSet []string `json:"reviewers"`
}

// 作者集合
type AuthorSet struct {
	AuthorNameSet []string `json:"authors"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	reviewerSet := ReviewerSet{
		ReviewerNameSet: []string{},
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
		AuthorNameSet: []string{},
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

func (s *SmartContract) CreateReviewer(ctx contractapi.TransactionContextInterface, name string, id string, passwd string, email string, researchTarget string) error {
	researchTarget = strings.Trim(researchTarget, " ")
	researchTargetArr := strings.Split(researchTarget, "/")

	exists, err := s.AssetExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the reviewer %s already exist", name)
	}
	reviewer := Reviewer{
		ID:              id,
		Name:            name,
		Passwd:          passwd,
		Email:           email,
		ResearchTarget:  researchTargetArr,
		ReviewedPaper:   []string{},
		UNReviewedPaper: []string{},
	}

	reviewerJSON, err := json.Marshal(reviewer)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, []byte(name))
	if err != nil {
		return nil
	}

	err = ctx.GetStub().PutState(name, reviewerJSON)
	if err != nil {
		return nil
	}

	reviewerSetJson, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return err
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJson, &reviewerSet)
	if err != nil {
		return err
	}
	newReviewerSet := ReviewerSet{
		ReviewerNameSet: append(reviewerSet.ReviewerNameSet, name),
	}
	newReviewerSetJSON, err := json.Marshal(newReviewerSet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("reviewerset", newReviewerSetJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) CreateAuthor(ctx contractapi.TransactionContextInterface, name string, id string, passwd string, email string) error {
	exists, err := s.AssetExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the author %s already exist", name)
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
	err = ctx.GetStub().PutState(id, []byte(name))
	if err != nil {
		return nil
	}
	err = ctx.GetStub().PutState(name, authorJSON)
	if err != nil {
		return nil
	}

	authorSetJson, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return err
	}
	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJson, &authorSet)
	if err != nil {
		return err
	}
	newAuthorSet := AuthorSet{
		AuthorNameSet: append(authorSet.AuthorNameSet, name),
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
		return fmt.Errorf("the paper %s already exist", title)
	}
	reviewerList, err := s.distributePaper(ctx, keywordsArr)
	if err != nil {
		return err
	}
	paper := Paper{
		Title:        	title,
		ID:           	id,
		AuthorList:   	authorListArr,
		ReviewerList: 	reviewerList,
		ReviewList: 	[]Review{},
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
	var reviewer Reviewer
	for _, each := range authorListArr {
		authorJSON, err := ctx.GetStub().GetState(each)
		if err != nil {
			return err
		}
		if authorJSON == nil {
			return fmt.Errorf("the author doesn't register yet:%v", each)
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
		err = ctx.GetStub().PutState(author.Name, newAuthorJSON)
		if err != nil {
			return err
		}
	}
	for _, each := range reviewerList {
		reviewerJSON, err := ctx.GetStub().GetState(each)
		if err != nil {
			return err
		}
		if reviewerJSON == nil {
			return fmt.Errorf("the reviwer doesn't register yet:%v", each)
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
		err = ctx.GetStub().PutState(reviewer.Name, newReviewerJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SmartContract) AddReview(ctx contractapi.TransactionContextInterface, name string, reviewerName string, content string) error {
	paperJSON, err := ctx.GetStub().GetState(name)
	if err != nil {
		return fmt.Errorf("failed to read from world state:%v", err)
	}
	if paperJSON == nil {
		return fmt.Errorf("the paper %s does not exist", name)
	}

	var paper Paper
	err = json.Unmarshal(paperJSON, &paper)
	if err != nil {
		return err
	}
	review := Review{
		ReviewerName: reviewerName,
		Content:      content,
	}
	newPaper := Paper{
		Title:        paper.Title,
		ID:           paper.ID,
		AuthorList:   paper.AuthorList,
		ReviewerList: paper.ReviewerList,
		ReviewList:   append(paper.ReviewList, review),
	}
	newPaperJSON, err := json.Marshal(newPaper)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(paper.Title, newPaperJSON)
	if err != nil {
		return err
	}

	var reviewer Reviewer
	reviewerJSON, err := ctx.GetStub().GetState(reviewerName)
	if err != nil {
		return err
	}
	if reviewerJSON == nil {
		return fmt.Errorf("the reviwer doesn't register yet:%v", reviewerName)
	}
	err = json.Unmarshal(reviewerJSON, &reviewer)
	if err != nil {
		return err
	}
	var index int
	for i,each := range reviewer.ReviewedPaper{
		if each == name{
			index = i
			break
		}
	}
	unReviewedPaper := append(reviewer.UNReviewedPaper[0:index],reviewer.UNReviewedPaper[index+1:]...)
	newReviewer := Reviewer{
		ID:              reviewer.ID,
		Name:            reviewer.Name,
		Passwd:          reviewer.Passwd,
		Email:           reviewer.Email,
		ResearchTarget:  reviewer.ResearchTarget,
		ReviewedPaper:   append(reviewer.ReviewedPaper, name),
		UNReviewedPaper: unReviewedPaper,
	}
	newReviewerJSON, err := json.Marshal(newReviewer)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(reviewer.Name, newReviewerJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) ReadReviewer(ctx contractapi.TransactionContextInterface, name string) (*Reviewer, error) {
	reviewerSetJson, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return nil,err
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJson, &reviewerSet)
	if err != nil {
		return nil,err
	}
	index := len(reviewerSet.ReviewerNameSet)
	for i:=0;i< len(reviewerSet.ReviewerNameSet);i++{
		if reviewerSet.ReviewerNameSet[i]==name {
			index = i
			break;
		}
	}
	if index == len(reviewerSet.ReviewerNameSet) {
		return nil,fmt.Errorf("reviewer %s not registered",name)
	}

	reviewerJSON, err := ctx.GetStub().GetState(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state:%v", err)
	}
	if reviewerJSON == nil {
		return nil, fmt.Errorf("the reviewer %s does not exist", name)
	}

	var reviewer Reviewer
	err = json.Unmarshal(reviewerJSON, &reviewer)
	if err != nil {
		return nil, err
	}
	return &reviewer, nil
}

func (s *SmartContract) ReadAuthor(ctx contractapi.TransactionContextInterface, name string) (*Author, error) {
	authorSetJson, err := ctx.GetStub().GetState("authorset")
	if err != nil {
		return nil,err
	}
	var authorSet AuthorSet
	err = json.Unmarshal(authorSetJson, &authorSet)
	if err != nil {
		return nil,err
	}
	index := len(authorSet.AuthorNameSet)
	for i:=0;i< len(authorSet.AuthorNameSet);i++{
		if authorSet.AuthorNameSet[i]==name {
			index = i
			break;
		}
	}
	if index == len(authorSet.AuthorNameSet) {
		return nil,fmt.Errorf("author %s not registered",name)
	}

	authorJSON, err := ctx.GetStub().GetState(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state:%v", err)
	}
	if authorJSON == nil {
		return nil, fmt.Errorf("the author %s does not exist", name)
	}

	var author Author
	err = json.Unmarshal(authorJSON, &author)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (s *SmartContract) GetPaper(ctx contractapi.TransactionContextInterface, title string)(*Paper,error){
	paperJSON,err := ctx.GetStub().GetState(title)
	if err != nil{
		return nil, err
	}
	if paperJSON==nil{
		return nil, fmt.Errorf("paper %s doesn't exist",title)
	}
	var paper Paper
	err = json.Unmarshal(paperJSON,&paper)
	if err!=nil {
		return nil,err
	}
	return &paper,nil
}

func (s *SmartContract) GetCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	author,err := s.ReadAuthor(ctx,name)
	if err != nil {
		return nil, err
	}

	var committedPaper []Paper
	for _, each := range author.CommittedPaper {
		paper,err := s.GetPaper(ctx,each)
		if err != nil {
			return nil, err
		}
		committedPaper = append(committedPaper, *paper)
	}
	return committedPaper, nil
}

func (s *SmartContract) GetReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	reviewer,err := s.ReadReviewer(ctx,name)
	if err != nil {
		return nil, err
	}

	var reviewedPaper []Paper
	for _, each := range reviewer.ReviewedPaper {
		paper,err := s.GetPaper(ctx,each)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		reviewedPaper = append(reviewedPaper, *paper)
	}
	return reviewedPaper, nil
}

func (s *SmartContract) GetUNReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	reviewer,err := s.ReadReviewer(ctx,name)
	if err != nil {
		return nil, err
	}

	var unreviewedPaper []Paper
	for _, each := range reviewer.UNReviewedPaper {
		paper,err := s.GetPaper(ctx,each)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		unreviewedPaper = append(unreviewedPaper, *paper)
	}
	return unreviewedPaper, nil
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
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
	return reviewerSet.ReviewerNameSet, nil
}

type comittedPaper struct {
	Name       string   `json:"name"`
	AuthorList []string `json:"authorlist"`
	Reviews    []string `json:"reviews"`
}

type reviewedPaper struct {
	Name   string `json:"name"`
	Review string `json:"review"`
}

func (s *SmartContract) ValidateAuthor(ctx contractapi.TransactionContextInterface,name string, passwd string)(bool,error){
	author,err := s.ReadAuthor(ctx,name)
	if err != nil {
		return false, err
	}
	if author == nil {
		return false, err
	}
	return author.Passwd == passwd, nil
}

func (s *SmartContract) ValidateReviewer(ctx contractapi.TransactionContextInterface,name string, passwd string)(bool,error){
	reviewer,err := s.ReadReviewer(ctx,name)
	if err != nil {
		return false, err
	}
	if reviewer == nil {
		return false, err
	}
	return reviewer.Passwd == passwd, nil
}

func (s *SmartContract) AuthorCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]comittedPaper, error) {
	papers, err := s.GetCommittedPaper(ctx, name)
	if err != nil {
		return nil, err
	}
	var res []comittedPaper
	for _, paper := range papers {
		reviews := []string{}
		if len(paper.ReviewList) > 0{
			for _, review := range paper.ReviewList {
				reviews = append(reviews, review.Content)
			}
		}
		tmp := comittedPaper{
			Name:       paper.Title,
			AuthorList: paper.AuthorList,
			Reviews:    reviews,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SmartContract) ReviewerReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]reviewedPaper, error) {
	papers, err := s.GetReviewedPaper(ctx, name)
	if err != nil {
		return nil, err
	}
	var res []reviewedPaper
	for _, paper := range papers {
		var content string
		for _, review := range paper.ReviewList {
			if review.ReviewerName == name {
				content = review.Content
				break
			}
		}
		tmp := reviewedPaper{
			Name:   paper.Title,
			Review: content,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SmartContract) ReviewerUNReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]string, error) {
	papers, err := s.GetUNReviewedPaper(ctx, name)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, paper := range papers {
		res = append(res, paper.Title)
	}
	return res, nil
}
