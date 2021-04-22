package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

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

// 审稿内容结构体
type Review struct {
	ReviewerID string `json:"reviewerid"`
	Content      string `json:"content"`
}

// 审稿人集合
type ReviewerSet struct {
	Reviewers map[string]string `json:"reviewers"`
}

func (s *SmartContract) AddtoReviewerSet(ctx contractapi.TransactionContextInterface,name string, id string) int {
	reviewerSetJson, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return getStateError
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJson, &reviewerSet)
	if err != nil {
		return jsonUnMarshalError
	}
	cur_reviewers := reviewerSet.Reviewers

	if _,ok := cur_reviewers[name]; ok{
		return existsError
	}

	cur_reviewers[name] = id
	newReviewerSet := ReviewerSet{
		Reviewers: cur_reviewers,
	}
	newReviewerSetJSON, err := json.Marshal(newReviewerSet)
	if err != nil {
		return jsonMarshalError
	}
	err = ctx.GetStub().PutState("reviewerset", newReviewerSetJSON)
	if err != nil {
		return jsonUnMarshalError
	}
	return 0
}

func (s *SmartContract) CreateReviewer(ctx contractapi.TransactionContextInterface, name string, id string, passwd string, email string, researchTarget string) int {
	researchTarget = strings.Trim(researchTarget, " ")
	researchTargetArr := strings.Split(researchTarget, "/")

	err2 := s.AddtoReviewerSet(ctx,name,id)
	if err2!=0{
		return err2
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
		return jsonMarshalError
	}

	err = ctx.GetStub().PutState(id, reviewerJSON)
	if err != nil {
		return putStateError
	}

	return 0
}

func (s *SmartContract) GetReviewerID(ctx contractapi.TransactionContextInterface, name string)(string, int){
	reviewerSetJSON, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return "",getStateError
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJSON, &reviewerSet)
	if err != nil {
		return "",jsonUnMarshalError
	}
	if _, ok := reviewerSet.Reviewers[name];!ok{
		return "",notExistsError
	}
	return reviewerSet.Reviewers[name],0
}

func (s *SmartContract) AddReview(ctx contractapi.TransactionContextInterface, name string, reviewerName string, content string) int {
	paperJSON, err := ctx.GetStub().GetState(name)
	if err != nil {
		return getStateError
	}
	if paperJSON == nil {
		return notExistsError
	}
	var paper Paper
	err = json.Unmarshal(paperJSON, &paper)
	if err != nil {
		return jsonUnMarshalError
	}

	reviewerID,err2 := s.GetReviewerID(ctx,reviewerName)
	if err2!=0{
		return err2
	}

	review := Review{
		ReviewerID: reviewerID,
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
		return jsonMarshalError
	}
	err = ctx.GetStub().PutState(paper.Title, newPaperJSON)
	if err != nil {
		return putStateError
	}

	var reviewer Reviewer
	reviewerJSON, err := ctx.GetStub().GetState(reviewerID)
	if err != nil {
		return getStateError
	}
	err = json.Unmarshal(reviewerJSON, &reviewer)
	if err != nil {
		return jsonUnMarshalError
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
		return jsonMarshalError
	}
	err = ctx.GetStub().PutState(reviewer.ID, newReviewerJSON)
	if err != nil {
		return putStateError
	}
	return 0
}

func (s *SmartContract) ReadReviewer(ctx contractapi.TransactionContextInterface, name string) (*Reviewer, int) {
	reviewerID,err2 := s.GetReviewerID(ctx,name)
	if err2!=0{
		return nil, err2
	}
	reviewerJSON, err := ctx.GetStub().GetState(reviewerID)
	if err != nil {
		return nil, getStateError
	}

	var reviewer Reviewer
	err = json.Unmarshal(reviewerJSON, &reviewer)
	if err != nil {
		return nil, jsonUnMarshalError
	}
	return &reviewer, 0
}

func (s *SmartContract) GetReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, int) {
	reviewer,err := s.ReadReviewer(ctx,name)
	if err != 0 {
		return nil, err
	}

	var reviewedPaper []Paper
	for _, each := range reviewer.ReviewedPaper {
		paper,err := s.GetPaper(ctx,each)
		if err != 0 {
			return nil, err
		}
		if err != 0 {
			return nil, err
		}
		reviewedPaper = append(reviewedPaper, *paper)
	}
	return reviewedPaper, 0
}

func (s *SmartContract) GetUNReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, int) {
	reviewer,err := s.ReadReviewer(ctx,name)
	if err != 0 {
		return nil, err
	}

	var unreviewedPaper []Paper
	for _, each := range reviewer.UNReviewedPaper {
		paper,err := s.GetPaper(ctx,each)
		if err != 0 {
			return nil, err
		}
		if err != 0 {
			return nil, err
		}
		unreviewedPaper = append(unreviewedPaper, *paper)
	}
	return unreviewedPaper, 0
}

type reviewedPaper struct {
	Name   string `json:"name"`
	Review string `json:"review"`
}

func (s *SmartContract) ValidateReviewer(ctx contractapi.TransactionContextInterface,name string, passwd string)(bool,int){
	reviewer,err := s.ReadReviewer(ctx,name)
	if err != 0 {
		return false, err
	}
	if reviewer == nil {
		return false, err
	}
	return reviewer.Passwd == passwd, 0
}

func (s *SmartContract) ReviewerReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]reviewedPaper, int) {
	papers, err := s.GetReviewedPaper(ctx, name)
	if err != 0 {
		return nil, err
	}

	reviewerID,err2 := s.GetReviewerID(ctx,name)
	if err2!=0 {
		return nil,err2
	}

	var res []reviewedPaper

	for _, paper := range papers {
		var content string
		for _, review := range paper.ReviewList {
			if review.ReviewerID == reviewerID {
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
	return res, 0
}

func (s *SmartContract) ReviewerUNReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]string, int) {
	papers, err := s.GetUNReviewedPaper(ctx, name)
	if err != 0 {
		return nil, err
	}
	var res []string
	for _, paper := range papers {
		res = append(res, paper.Title)
	}
	return res, 0
}
