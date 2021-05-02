package chaincode

import (
	"encoding/json"
	"fmt"
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

type Rebuttal struct{
	AuthorID		string 	`json:"authorid"`
	ReviewerID		string 	`json:"reviewerid"`
	Question 		string	`json:"question"`
	Reply			string	`json:"reply"`
	IsReplyed		bool	`json:"isreplyed"`
}

// 审稿内容结构体
type Review struct {
	ReviewerID 		string `json:"reviewerid"`
	Content    		string `json:"content"`
	RebuttalList	Rebuttal `json:"rebuttallist"`
}

func (s *SmartContract) AddtoReviewerSet(ctx contractapi.TransactionContextInterface, name string, id string) error {
	reviewerSetJson, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return err
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJson, &reviewerSet)
	if err != nil {
		return err
	}

	if _, ok := reviewerSet.Reviewers[name]; ok {
		return fmt.Errorf("Reviewer %s exists", name)
	}

	reviewerSet.Reviewers[name] = id
	newReviewerSet := ReviewerSet{
		Reviewers: reviewerSet.Reviewers,
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

func (s *SmartContract) CreateReviewer(ctx contractapi.TransactionContextInterface, name string, id string, passwd string, email string, researchTarget string) error {
	researchTarget = strings.Trim(researchTarget, " ")
	researchTargetArr := strings.Split(researchTarget, "/")

	err := s.AddtoReviewerSet(ctx, name, id)
	if err != nil{
		return err
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

	err = ctx.GetStub().PutState(id, reviewerJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) UpdateReviewerInfo(ctx contractapi.TransactionContextInterface, name string, passwd string, email string, researchTarget string) error {
	researchTarget = strings.Trim(researchTarget, " ")
	researchTargetArr := strings.Split(researchTarget, "/")

	reviewer, err := s.ReadReviewer(ctx,name)
	if err != nil {
		return err
	}
	newReviewer := Reviewer {
		ID: reviewer.ID,
		Name: reviewer.Name,
		Passwd: passwd,
		Email: email,
		ResearchTarget: researchTargetArr,
		ReviewedPaper: reviewer.ReviewedPaper,
		UNReviewedPaper: reviewer.UNReviewedPaper,
	}

	reviewerJSON,err := json.Marshal(newReviewer)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(reviewer.ID, reviewerJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) GetReviewerID(ctx contractapi.TransactionContextInterface, name string) (string, error) {
	reviewerSetJSON, err := ctx.GetStub().GetState("reviewerset")
	if err != nil {
		return "", err
	}
	var reviewerSet ReviewerSet
	err = json.Unmarshal(reviewerSetJSON, &reviewerSet)
	if err != nil {
		return "", err
	}
	if _, ok := reviewerSet.Reviewers[name]; !ok {
		return "", fmt.Errorf("Reviewer %s not exist", name)
	}
	return reviewerSet.Reviewers[name], nil
}

func (s *SmartContract) AddReview(ctx contractapi.TransactionContextInterface, title string, reviewerName string, content string) error {
	paper, err := s.GetPaper(ctx,title)
	if err != nil{
		return err
	}

	reviewerID, err := s.GetReviewerID(ctx, reviewerName)
	if err != nil {
		return err
	}

	review := Review{
		ReviewerID: reviewerID,
		Content:    content,
		RebuttalList: Rebuttal{},
	}
	if _,ok := paper.ReviewList[reviewerID]; ok{
		return fmt.Errorf("Review has been added by %s to %s.",reviewerID,paper.Title)
	}
	paper.ReviewList[reviewerID] = review;
	newPaper := Paper{
		Title:        paper.Title,
		ID:           paper.ID,
		AuthorList:   paper.AuthorList,
		ReviewerList: paper.ReviewerList,
		ReviewList:   paper.ReviewList,
	}
	newPaperJSON, err := json.Marshal(newPaper)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(paper.ID, newPaperJSON)
	if err != nil {
		return err
	}

	var reviewer Reviewer
	reviewerJSON, err := ctx.GetStub().GetState(reviewerID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(reviewerJSON, &reviewer)
	if err != nil {
		return err
	}
	var index int
	for i, each := range reviewer.ReviewedPaper {
		if each == title {
			index = i
			break
		}
	}
	unReviewedPaper := append(reviewer.UNReviewedPaper[0:index], reviewer.UNReviewedPaper[index+1:]...)
	newReviewer := Reviewer{
		ID:              reviewer.ID,
		Name:            reviewer.Name,
		Passwd:          reviewer.Passwd,
		Email:           reviewer.Email,
		ResearchTarget:  reviewer.ResearchTarget,
		ReviewedPaper:   append(reviewer.ReviewedPaper, title),
		UNReviewedPaper: unReviewedPaper,
	}
	newReviewerJSON, err := json.Marshal(newReviewer)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(reviewer.ID, newReviewerJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) ReadReviewer(ctx contractapi.TransactionContextInterface, name string) (*Reviewer, error) {
	reviewerID, err := s.GetReviewerID(ctx, name)
	if err != nil {
		return nil, err
	}
	reviewerJSON, err := ctx.GetStub().GetState(reviewerID)
	if err != nil {
		return nil, err
	}

	var reviewer Reviewer
	err = json.Unmarshal(reviewerJSON, &reviewer)
	if err != nil {
		return nil, err
	}
	return &reviewer, nil
}

func (s *SmartContract) GetReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	reviewer, err := s.ReadReviewer(ctx, name)
	if err != nil {
		return nil, err
	}

	var reviewedPaper []Paper
	for _, each := range reviewer.ReviewedPaper {
		paper, err := s.GetPaper(ctx, each)
		if err != nil {
			return nil, err
		}
		reviewedPaper = append(reviewedPaper, *paper)
	}
	return reviewedPaper, nil
}

func (s *SmartContract) GetUNReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	reviewer, err := s.ReadReviewer(ctx, name)
	if err != nil {
		return nil, err
	}

	var unreviewedPaper []Paper
	for _, each := range reviewer.UNReviewedPaper {
		paper, err := s.GetPaper(ctx, each)
		if err != nil {
			return nil, err
		}
		unreviewedPaper = append(unreviewedPaper, *paper)
	}
	return unreviewedPaper, nil
}

type reviewedPaper struct {
	Name   string `json:"name"`
	Review string `json:"review"`
	RebuttalList Rebuttal `json:"rebuttallist"`
}

func (s *SmartContract) ValidateReviewer(ctx contractapi.TransactionContextInterface, name string, passwd string) (bool, error) {
	reviewer, err := s.ReadReviewer(ctx, name)
	if err != nil {
		return false, err
	}
	if reviewer == nil {
		return false, err
	}
	return reviewer.Passwd == passwd, nil
}

func (s *SmartContract) ReviewerReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]reviewedPaper, error) {
	papers, err := s.GetReviewedPaper(ctx, name)
	if err != nil {
		return nil, err
	}

	reviewerID, err := s.GetReviewerID(ctx, name)
	if err != nil {
		return nil, err
	}

	var res []reviewedPaper

	for _, paper := range papers {
		review := paper.ReviewList[reviewerID]

		tmp := reviewedPaper{
			Name:   paper.Title,
			Review: review.Content,
			RebuttalList: review.RebuttalList,
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

func (s *SmartContract) AddReply(ctx contractapi.TransactionContextInterface, title string, reviewer_name string, reply string) error {
	reviewerID, err := s.GetReviewerID(ctx, reviewer_name)
	if err != nil{
		return err
	}

	paper, err := s.GetPaper(ctx,title)
	if err != nil {
		return err
	}

	rebuttal := paper.ReviewList[reviewerID].RebuttalList
	newRebuttal := Rebuttal{
		AuthorID: rebuttal.AuthorID,
		ReviewerID: reviewerID,
		Question: rebuttal.Question,
		Reply: reply,
		IsReplyed: true,
	}

	review := Review{
		ReviewerID: reviewerID,
		Content: paper.ReviewList[reviewerID].Content,
		RebuttalList: newRebuttal,
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
