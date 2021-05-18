package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
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

type Info struct {
	Keywords []string
	Name     []string
	Fields   []string
}

func InitSimilarityPair(similarityFile string) error {
	similarity := make(map[string]float64)

	bytes, _ := ioutil.ReadFile(similarityFile)
	fmt.Println(string(bytes))
	err := json.Unmarshal(bytes, &similarity)
	if err != nil {
		return err
	}

	for key, value := range similarity {
		fmt.Println(key,value)
	}
	return nil
}

func InitPaper(paperFile string) error {
	paper := make(map[string]PaperInfo)

	bytes, _ := ioutil.ReadFile(paperFile)

	err := json.Unmarshal(bytes, &paper)
	if err != nil {
		return err
	}

	for key, value := range paper {
		arguments := []string{value.Submitter, value.Title, value.ID, value.AuthorList, value.Keywords, value.Path}

		fmt.Println(arguments)
		if err != nil {
			return err
		}

		fmt.Println("Add Paper %s successfully\n", key)
	}

	return nil
}

func ReadAuthor() (map[string]UserInfo, error) {
	m := make(map[string]Info)

	bytes, _ := ioutil.ReadFile("paper_info.json")

	err := json.Unmarshal(bytes, &m)
	if err != nil {
		fmt.Println(err.Error())
	}

	author := make(map[string]UserInfo)
	for key, value := range m {
		// value := x.(Info)
		name := "author" + key
		id := "id_author" + key
		passwd := "123456"
		email := "1141751053@qq.com"
		researchtargets := ""
		num := len(value.Fields)
		for i := 0; i < num-1; i++ {
			researchtargets += value.Fields[i] + "/"
		}
		researchtargets += value.Fields[num-1]

		author[name] = UserInfo{
			Name:           name,
			ID:             id,
			Passwd:         passwd,
			Email:          email,
			ResearchTarget: researchtargets,
		}
	}

	return author, nil
}

func ReadPaper() (map[string]PaperInfo, error) {
	m := make(map[string]Info)

	bytes, _ := ioutil.ReadFile("paper_info.json")

	err := json.Unmarshal(bytes, &m)
	if err != nil {
		fmt.Println(err.Error())
	}

	paper := make(map[string]PaperInfo)
	for key, value := range m {
		// value := x.(Info)
		title := value.Name[0]
		title = strings.Replace(title, ":", "", -1)
		title = strings.Replace(title, "?", "", -1)

		id := "id_paper" + key
		keywords := ""
		num := len(value.Keywords)
		for i := 0; i < num-1; i++ {
			keywords += value.Keywords[i] + "/"
		}
		keywords += value.Keywords[num-1]

		paper[value.Name[0]] = PaperInfo{
			Title:      value.Name[0],
			ID:         id,
			Submitter:  "author" + key,
			AuthorList: "author" + key,
			Keywords:   keywords,
			Path:       "./upload/" + title + ".pdf",
		}
	}

	return paper, nil
}

func ReadReviewer() (map[string]UserInfo, error) {
	m := make(map[string][]string)

	bytes, _ := ioutil.ReadFile("reviewer_field.json")

	err := json.Unmarshal(bytes, &m)
	if err != nil {
		fmt.Println(err.Error())
	}

	reviewer := make(map[string]UserInfo)
	for key, value := range m {
		// value := x.(Info)
		name := "reviewer" + key
		id := "id_reviewer" + key
		passwd := "123456"
		email := "1141751053@qq.com"
		researchtargets := ""
		num := len(value)
		for i := 0; i < num-1; i++ {
			researchtargets += value[i] + "/"
		}
		researchtargets += value[num-1]

		reviewer[name] = UserInfo{
			Name:           name,
			ID:             id,
			Passwd:         passwd,
			Email:          email,
			ResearchTarget: researchtargets,
		}
	}

	return reviewer, nil
}

func Write(data []byte, filename string) {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func Similarity(keywords []string, researchTargets []string) (float64, error) {
	similarity := make(map[string]float64)

	bytes, _ := ioutil.ReadFile("similarity.json")

	err := json.Unmarshal(bytes, &similarity)
	if err != nil {
		return 0.0, err
	}

	score := 0.0
	for _, key := range keywords {
		for _, area := range researchTargets {
			// fmt.Println(similarity[key+"+"+area])
			score = math.Max(score, similarity[key+"+"+area])
		}
	}

	return score, nil
}

func GenerateInfo() {
	author, err := ReadAuthor()
	if err != nil {
		fmt.Println(err.Error())
	}

	authorJSON, err := json.Marshal(author)
	if err != nil {
		fmt.Println(err.Error())
	}

	Write(authorJSON, "authors.json")
	fmt.Println(len(author))

	paper, err := ReadPaper()
	if err != nil {
		fmt.Println(err.Error())
	}

	paperJSON, err := json.Marshal(paper)
	if err != nil {
		fmt.Println(err.Error())
	}

	Write(paperJSON, "papers.json")

	reviewer, err := ReadReviewer()
	if err != nil {
		fmt.Println(err.Error())
	}

	reviewerJSON, err := json.Marshal(reviewer)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(len(reviewer))

	Write(reviewerJSON, "reviewers.json")

}

func main() {
	GenerateInfo()
	//InitPaper("papers.json")
	//_ =InitSimilarityPair("similarity.json")
	// for title, p := range paper {
	// 	for name, r := range reviewer {
	// 		fmt.Println(title)
	// 		fmt.Println(name)
	// 		fmt.Println(Similarity(strings.Split(p.Keywords, "/"), strings.Split(r.ResearchTarget, "/")))
	// 	}
	// }
}
