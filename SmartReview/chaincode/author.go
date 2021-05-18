package chaincode

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type comittedPaper struct {
	Name           string            `json:"name"`
	AuthorList     []string          `json:"authorlist"`
	Reviews        map[string]Review `json:"reviews"`
	Successed      bool              `json:"successed"`
	Passed         bool              `json:"passed"`
	ReviewFinished bool              `json:"reviewfinished"`
}

func (s *SmartContract) GetCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	authorID, err := s.GetUserID(ctx, name)
	if err != nil {
		return nil, err
	}

	author, err := s.GetUser(ctx, authorID)
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

func (s *SmartContract) AuthorCommittedPaper(ctx contractapi.TransactionContextInterface, name string) ([]comittedPaper, error) {
	papers, err := s.GetCommittedPaper(ctx, name)
	if err != nil {
		return nil, err
	}
	var res []comittedPaper

	for _, paper := range papers {
		var authorList []string
		for _, authorID := range paper.AuthorList {
			author, err := s.GetUser(ctx, authorID)
			if err != nil {
				return nil, err
			}
			authorList = append(authorList, author.Name)
		}
		tmp := comittedPaper{
			Name:           paper.Title,
			AuthorList:     authorList,
			Reviews:        paper.ReviewList,
			Successed:      paper.Successed,
			Passed:         paper.Passed,
			ReviewFinished: paper.ReviewFinished,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SmartContract) AddRebuttal(ctx contractapi.TransactionContextInterface, title string, author_name string, reviewerID string, question string) error {
	authorID, err := s.GetUserID(ctx, author_name)
	if err != nil {
		return err
	}

	paperID, err := s.GetPaperID(ctx, title)
	if err != nil {
		return err
	}
	paper, err := s.GetPaper(ctx, paperID)
	if err != nil {
		return err
	}

	rebuttalID := len(paper.ReviewList[reviewerID].RebuttalList)

	rebuttal := Rebuttal{
		AuthorID:   authorID,
		ReviewerID: reviewerID,
		RebuttalID: strconv.Itoa(rebuttalID),
		Question:   question,
		Reply:      "",
		IsReplyed:  false,
	}

	rebuttalList := paper.ReviewList[reviewerID].RebuttalList
	rebuttalList[strconv.Itoa(rebuttalID)] = rebuttal

	review := Review{
		ReviewerID:   reviewerID,
		Content:      paper.ReviewList[reviewerID].Content,
		RebuttalList: rebuttalList,
	}

	newReviewList := paper.ReviewList
	newReviewList[reviewerID] = review
	paper.ReviewList = newReviewList

	newPaperJSON, err := json.Marshal(*paper)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(paper.ID, newPaperJSON)
	if err != nil {
		return err
	}
	return nil
}
