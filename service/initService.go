package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type UserInfo struct {
	Name           string
	ID             string
	Passwd         string
	Email          string
	ResearchTarget string
}

type PaperInfo struct {
	Submitter  string
	Title      string
	ID         string
	AuthorList string
	Keywords   string
	Path       string
}

func (s *ServiceSetup) InitUser(authorFile string, reviewerFile string) error {
	author := make(map[string]UserInfo)
	reviewer := make(map[string]UserInfo)

	bytes, _ := ioutil.ReadFile(authorFile)

	err := json.Unmarshal(bytes, &author)
	if err != nil {
		return err
	}

	bytes, _ = ioutil.ReadFile(reviewerFile)

	err = json.Unmarshal(bytes, &reviewer)
	if err != nil {
		return err
	}

	for key, value := range author {
		arguments := []string{value.Name, value.ID, value.Passwd, value.Email, value.ResearchTarget}

		_, err = s.InvokeChaincode("CreateUser", arguments)
		if err != nil {
			return err
		}

		fmt.Printf("add Author %s successfully\n", key)
	}

	i := 0
	for key, value := range reviewer {
		arguments := []string{value.Name, value.ID, value.Passwd, value.Email, value.ResearchTarget}

		_, err = s.InvokeChaincode("CreateUser", arguments)
		if err != nil {
			return err
		}
		i += 1
		fmt.Printf("add Reviewer %s successfully %d\n", key, i)
	}

	return nil
}

func (s *ServiceSetup) InitPaper(paperFile string) error {
	paper := make(map[string]PaperInfo)

	bytes, _ := ioutil.ReadFile(paperFile)

	err := json.Unmarshal(bytes, &paper)
	if err != nil {
		return err
	}

	for key, value := range paper {
		arguments := []string{value.Submitter, value.Title, value.ID, value.AuthorList, value.Keywords, value.Path}

		_, err = s.InvokeChaincode("CreatePaper", arguments)
		if err != nil {
			return err
		}

		fmt.Printf("add Paper %s successfully\n", key)
	}

	return nil
}

func (s *ServiceSetup) InitSimilarityPair(similarityFile string) error {
	similarity := make(map[string]float64)

	bytes, _ := ioutil.ReadFile(similarityFile)


	err := json.Unmarshal(bytes, &similarity)
	if err != nil {
		return err
	}

	i := 0

	for key, value := range similarity {
		score := strconv.FormatFloat(value, 'f', 10, 64)
		arguments := []string{key,score}
		_, err :=s.InvokeChaincode("AddSimilarityPair",arguments)
		if err != nil {
			return err
		}
		i+=1
		fmt.Printf("add %s: %s successfully %d\n", key, score, i)
	}
	return nil
}
