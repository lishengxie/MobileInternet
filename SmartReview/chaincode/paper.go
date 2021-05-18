package chaincode

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	coinPerPaper = 3.0
	punishCoin   = 1.0
)

// 论文结构体
type Paper struct {
	Title          string            `json:"title"`
	ID             string            `json:"id"`
	KeyWords       []string          `json:"keywords"`
	AuthorList     []string          `json:"authorlist"`
	ReviewerList   []string          `json:"reviewerlist"`
	ReviewList     map[string]Review `json:"reviewlist"`
	AltCoin        float64           `json:"altcoin"`
	Submitter      string            `json:"submitter"`
	ReviewFinished bool              `json:"reviewfinished"`
	Successed      bool              `json:"success"`
	Passed         bool              `json:"passed"`
	StorePath      string            `json:"storepath"`
}

// PaperInfo结构体, 用于向用户展示已有的论文的信息
type PaperInfo struct {
	Title      string            `json:"title"`
	KeyWords   []string          `json:"keywords"`
	AuthorList []string          `json:"authorlist"`
	ReviewList map[string]Review `json:"reviewlist"`
	StorePath  string            `json:"storepath"`
}

func (s *SmartContract) GetPaperSet(ctx contractapi.TransactionContextInterface) (*PaperSet, error) {
	paperSetJSON, err := ctx.GetStub().GetState("paperset")
	if err != nil {
		return nil, err
	}

	var paperSet PaperSet
	err = json.Unmarshal(paperSetJSON, &paperSet)
	if err != nil {
		return nil, err
	}

	return &paperSet, nil
}

func (s *SmartContract) AddToPaperSet(ctx contractapi.TransactionContextInterface, title string, id string) error {
	paperSet, err := s.GetPaperSet(ctx)
	if err != nil {
		return err
	}

	if _, ok := paperSet.Papers[title]; ok {
		return fmt.Errorf("paper %s already exists", title)
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
		return fmt.Errorf("paper %s not exists", origin_title)
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

func (s *SmartContract) CreatePaper(ctx contractapi.TransactionContextInterface, submitter string, title string, id string, authorList string, keywords string, path string) error {
	submitterID, err := s.GetUserID(ctx, submitter)
	if err != nil {
		return err
	}

	user, err := s.GetUser(ctx, submitterID)
	if err != nil {
		return err
	}

	if user.AltCoin < coinPerPaper {
		return fmt.Errorf("altcoin of %s not enough", submitter)
	}

	authorList = strings.Trim(authorList, " ")
	authorsArr := strings.Split(authorList, "/")

	keywords = strings.Trim(keywords, " ")
	keywordsArr := strings.Split(keywords, "/")

	err = s.AddToPaperSet(ctx, title, id)
	if err != nil {
		return err
	}

	var authorIDList []string
	authors := make(map[string]string)
	for _, authorName := range authorsArr {
		authorID, err := s.GetUserID(ctx, authorName)
		if err != nil {
			return err
		}
		authorIDList = append(authorIDList, authorID)
		authors[authorID] = authorName
	}

	reviewerIDList, err := s.DistributePaper(ctx, keywordsArr, authors, 3)
	if err != nil {
		return err
	}

	paper := Paper{
		Title:          title,
		ID:             id,
		KeyWords:       keywordsArr,
		AuthorList:     authorIDList,
		ReviewerList:   reviewerIDList,
		ReviewList:     make(map[string]Review),
		AltCoin:        coinPerPaper,
		Submitter:      submitterID,
		ReviewFinished: false,
		Successed:      false,
		Passed:         true,
		StorePath:      path,
	}
	paperJSON, err := json.Marshal(paper)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, paperJSON)
	if err != nil {
		return err
	}

	for _, each := range authorIDList {
		author, err := s.GetUser(ctx, each)
		if err != nil {
			return err
		}

		if each == submitterID {
			author.AltCoin = author.AltCoin - coinPerPaper
		}

		author.CommittedPaper = append(author.CommittedPaper, id)

		newAuthorJSON, err := json.Marshal(*author)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(author.ID, newAuthorJSON)
		if err != nil {
			return err
		}
	}

	for _, each := range reviewerIDList {
		reviewer, err := s.GetUser(ctx, each)
		if err != nil {
			return err
		}

		reviewer.UNReviewedPaper = append(reviewer.UNReviewedPaper, id)

		newReviewerJSON, err := json.Marshal(*reviewer)
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
		authorID, err := s.GetUserID(ctx, authorName)
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

	paper.Title = new_title
	paper.AuthorList = append(paper.AuthorList, addedAuthorIDList...)

	newPaperJSON, err := json.Marshal(*paper)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(paper.ID, newPaperJSON)
	if err != nil {
		return err
	}

	for _, each := range addedAuthorIDList {

		author, err := s.GetUser(ctx, each)
		if err != nil {
			return err
		}

		author.CommittedPaper = append(author.CommittedPaper, paper.ID)

		newAuthorJSON, err := json.Marshal(*author)
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
		return "", fmt.Errorf("paper %s not exist", title)
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
		author, err := s.GetUser(ctx, authorID)
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
		StorePath:  paper.StorePath,
	}, nil
}

type reviewerScore struct {
	reviewerID string
	score      float64
}

func (s *SmartContract) DistributePaper(ctx contractapi.TransactionContextInterface, keywords []string, authors map[string]string, reviewedNum int) ([]string, error) {
	userSet, err := s.GetUserSet(ctx)
	if err != nil {
		return nil, err
	}

	reviewerScoreSet := make([]reviewerScore, 0)
	reviewerIDSet := make([]string, 0)

	keys := make([]string, 0, len(userSet.Users))
	for k := range userSet.Users {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		id := userSet.Users[k]
		if _, ok := authors[id]; ok {
			continue
		}

		reviewer, err := s.GetUser(ctx, id)
		if err != nil {
			return nil, err
		}

		researchTargets := reviewer.ResearchTarget
		similarity, err := s.Similarity(ctx, keywords, researchTargets)
		if err != nil {
			return nil, err
		}
		reviewerScoreSet = append(reviewerScoreSet,
			reviewerScore{
				id,
				similarity,
			})
	}

	sort.SliceStable(reviewerScoreSet, func(i, j int) bool {
		return reviewerScoreSet[i].score > reviewerScoreSet[j].score
	})

	for i := 0; i < reviewedNum; i++ {
		reviewerIDSet = append(reviewerIDSet, reviewerScoreSet[i].reviewerID)
	}

	return reviewerIDSet, nil
}

func (s *SmartContract) Similarity(ctx contractapi.TransactionContextInterface, keywords []string, researchTargets []string) (float64, error) {
	score := 0.0
	for _, key := range keywords {
		for _, area := range researchTargets {
			pair := key + "+" + area
			tmp, err := s.GetSimilarity(ctx, pair)
			if err != nil {
				return 0.0, err
			}
			score = math.Max(score, tmp)
		}
	}

	return score, nil
}
