package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Rebuttal struct {
	AuthorID   string `json:"authorid"`
	ReviewerID string `json:"reviewerid"`
	RebuttalID string `json:"rebuttalid"`
	Question   string `json:"question"`
	Reply      string `json:"reply"`
	IsReplyed  bool   `json:"isreplyed"`
}

// 审稿内容结构体
type Review struct {
	ReviewerID   string              `json:"reviewerid"`
	Valid        bool                `json:"valid"`
	Content      string              `json:"content"`
	RebuttalList map[string]Rebuttal `json:"rebuttallist"` //rebuttalID => rebuttal
}

func (s *SmartContract) GetReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]Paper, error) {
	reviewerID, err := s.GetUserID(ctx, name)
	if err != nil {
		return nil, err
	}

	reviewer, err := s.GetUser(ctx, reviewerID)
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
	reviewerID, err := s.GetUserID(ctx, name)
	if err != nil {
		return nil, err
	}

	reviewer, err := s.GetUser(ctx, reviewerID)
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
	Name         string              `json:"name"`
	Review       string              `json:"review"`
	RebuttalList map[string]Rebuttal `json:"rebuttallist"`
	Keywords     []string            `json:"keywords"`
	StorePath    string              `json:"storepath"`
}

type unReviewedPaper struct {
	Name      string `json:"name"`
	StorePath string `json:"storepath"`
}

func (s *SmartContract) ReviewerReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]reviewedPaper, error) {
	papers, err := s.GetReviewedPaper(ctx, name)
	if err != nil {
		return nil, err
	}

	reviewerID, err := s.GetUserID(ctx, name)
	if err != nil {
		return nil, err
	}

	var res []reviewedPaper

	for _, paper := range papers {
		review := paper.ReviewList[reviewerID]

		tmp := reviewedPaper{
			Name:         paper.Title,
			Review:       review.Content,
			RebuttalList: review.RebuttalList,
			Keywords:     paper.KeyWords,
			StorePath:    paper.StorePath,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SmartContract) ReviewerUNReviewedPaper(ctx contractapi.TransactionContextInterface, name string) ([]unReviewedPaper, error) {
	papers, err := s.GetUNReviewedPaper(ctx, name)
	if err != nil {
		return nil, err
	}
	var res []unReviewedPaper
	for _, paper := range papers {
		tmp := unReviewedPaper{
			Name:      paper.Title,
			StorePath: paper.StorePath,
		}
		res = append(res, tmp)
	}
	return res, nil
}

type DecideSignal struct {
	IsFinished     bool
	IsValid        bool
	ValidReviewNum int
}

func (s *SmartContract) AddReview(ctx contractapi.TransactionContextInterface, title string, reviewerName string, content string, valid bool) error {
	paperID, err := s.GetPaperID(ctx, title)
	paper, err := s.GetPaper(ctx, paperID)
	if err != nil {
		return err
	}

	reviewerID, err := s.GetUserID(ctx, reviewerName)
	if err != nil {
		return err
	}

	review := Review{
		ReviewerID:   reviewerID,
		Valid:        valid,
		Content:      content,
		RebuttalList: make(map[string]Rebuttal),
	}
	if _, ok := paper.ReviewList[reviewerID]; ok {
		return fmt.Errorf("review has been added by %s to %s", reviewerName, paper.Title)
	}
	paper.ReviewList[reviewerID] = review
	newPaperJSON, err := json.Marshal(*paper)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(paper.ID, newPaperJSON)
	if err != nil {
		return err
	}

	reviewer, err := s.GetUser(ctx, reviewerID)
	if err != nil {
		return err
	}
	var index int
	for i, each := range reviewer.UNReviewedPaper {
		if each == paperID {
			index = i
			break
		}
	}
	unReviewedPaper := append(reviewer.UNReviewedPaper[0:index], reviewer.UNReviewedPaper[index+1:]...)

	reviewer.UNReviewedPaper = unReviewedPaper
	reviewer.ReviewedPaper = append(reviewer.ReviewedPaper, paperID)
	newReviewerJSON, err := json.Marshal(*reviewer)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(reviewer.ID, newReviewerJSON)
	if err != nil {
		return err
	}

	signal, _ := s.DecideFinality(paper)
	isFinished := signal.IsFinished
	isValid := signal.IsValid
	validReviewNum := signal.ValidReviewNum

	if isFinished {
		if isValid {
			err = s.RewardReviewer(ctx, paper, validReviewNum, reviewer)
			if err != nil {
				return err
			}
			paper.Successed = true
		}
		if !isValid {
			err = s.PunishSubmitter(ctx, paper, punishCoin, reviewer)
			if err != nil {
				return err
			}
			paper.Successed = false
		}

		paper.AltCoin = 0
		paper.ReviewFinished = true
		newPaperJSON, err = json.Marshal(*paper)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(paper.ID, newPaperJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SmartContract) DecideFinality(paper *Paper) (*DecideSignal, error) {
	var validReviewNum int
	var invalidReviewNum int
	totalReviewNum := len(paper.ReviewerList)
	for _, reviewerID := range paper.ReviewerList {
		if review, ok := paper.ReviewList[reviewerID]; ok {
			if review.Valid {
				validReviewNum++
			} else {
				invalidReviewNum++
			}
		}
	}

	if validReviewNum+invalidReviewNum == totalReviewNum {
		if float64(validReviewNum) > float64(totalReviewNum)/2.0 {
			return &DecideSignal{
				IsFinished:     true,
				IsValid:        true,
				ValidReviewNum: validReviewNum,
			}, nil
		}
		if float64(invalidReviewNum) > float64(totalReviewNum)/2.0 {
			return &DecideSignal{
				IsFinished:     true,
				IsValid:        false,
				ValidReviewNum: validReviewNum,
			}, nil
		}
	} else {
		return &DecideSignal{
			IsFinished:     false,
			IsValid:        false,
			ValidReviewNum: validReviewNum,
		}, nil
	}

	return &DecideSignal{
		IsFinished:     false,
		IsValid:        false,
		ValidReviewNum: validReviewNum,
	}, nil
}

func (s *SmartContract) RewardReviewer(ctx contractapi.TransactionContextInterface, paper *Paper, validReviewNum int, finalUser *User) error {
	num := paper.AltCoin / float64(validReviewNum)

	for _, reviewerID := range paper.ReviewerList {
		if _, ok := paper.ReviewList[reviewerID]; !ok {
			return fmt.Errorf("review not exist")
		}

		if !paper.ReviewList[reviewerID].Valid {
			continue
		}

		if reviewerID == finalUser.ID {
			reviewer := finalUser
			reviewer.AltCoin += num
			reviewerJSON, err := json.Marshal(*reviewer)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(reviewerID, reviewerJSON)
			if err != nil {
				return err
			}
		} else {
			reviewer, err := s.GetUser(ctx, reviewerID)
			if err != nil {
				return err
			}
			reviewer.AltCoin += num
			reviewerJSON, err := json.Marshal(*reviewer)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(reviewerID, reviewerJSON)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SmartContract) PunishSubmitter(ctx contractapi.TransactionContextInterface, paper *Paper, punishNum float64, finalUser *User) error {
	num := punishNum / float64(len(paper.ReviewerList))

	for _, reviewerID := range paper.ReviewerList {
		if _, ok := paper.ReviewList[reviewerID]; !ok {
			return fmt.Errorf("review not exist")
		}

		if reviewerID == finalUser.ID {
			reviewer := finalUser
			reviewer.AltCoin += num
			reviewerJSON, err := json.Marshal(*reviewer)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(reviewerID, reviewerJSON)
			if err != nil {
				return err
			}
		} else {
			reviewer, err := s.GetUser(ctx, reviewerID)
			if err != nil {
				return err
			}
			reviewer.AltCoin += num
			reviewerJSON, err := json.Marshal(*reviewer)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(reviewerID, reviewerJSON)
			if err != nil {
				return err
			}
		}
	}

	submitter, err := s.GetUser(ctx, paper.Submitter)
	if err != nil {
		return err
	}

	submitter.AltCoin += paper.AltCoin - punishNum
	submitterJSON, err := json.Marshal(submitter)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(paper.Submitter, submitterJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) AddReply(ctx contractapi.TransactionContextInterface, title string, reviewer_name string, reply string, rebuttalID string) error {
	reviewerID, err := s.GetUserID(ctx, reviewer_name)
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

	rebuttalList := paper.ReviewList[reviewerID].RebuttalList
	rebuttal := rebuttalList[rebuttalID]

	rebuttalList[rebuttalID] = Rebuttal{
		AuthorID:   rebuttal.AuthorID,
		ReviewerID: reviewerID,
		RebuttalID: rebuttalID,
		Question:   rebuttal.Question,
		Reply:      reply,
		IsReplyed:  true,
	}

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

func (s *SmartContract) Reject(ctx contractapi.TransactionContextInterface, title string, rejected bool) error {
	paperID, err := s.GetPaperID(ctx, title)
	if err != nil {
		return err
	}

	paper, err := s.GetPaper(ctx, paperID)
	if err != nil {
		return err
	}

	if !paper.Passed {
		return nil
	}

	if rejected {
		paper.Passed = false
	} else {
		paper.Passed = true
	}

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
