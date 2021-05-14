package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Application struct {
	//Service *service.ServiceSetup
}

type Rebuttal struct {
	AuthorID   string `json:"authorid"`
	ReviewerID string `json:"reviewerid"`
	Question   string `json:"question"`
	Reply      string `json:"reply"`
	IsReplyed  bool   `json:"isreplyed"`
}

// 审稿内容结构体
type Review struct {
	ReviewerID   string   `json:"reviewerid"`
	Content      string   `json:"content"`
	RebuttalList Rebuttal `json:"rebuttallist"`
}

type comittedPaper struct {
	Name       string            `json:"name"`
	AuthorList []string          `json:"authorlist"`
	Reviews    map[string]Review `json:"reviews"`
}

type reviewedPaper struct {
	Name         string   `json:"name"`
	Review       string   `json:"review"`
	RebuttalList Rebuttal `json:"rebuttallist"`
}

func (app *Application) HomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	// fmt.Println(name, passwd, identity)
	name := r.Form.Get("name")
	paper := []comittedPaper{
		comittedPaper{
			Name:       "Bitcoin: Peer to peer ecoin system",
			AuthorList: []string{"hywang"},
			Reviews: map[string]Review{
				"id_reviewer": Review{
					ReviewerID: "id_reviewer",
					Content:    "Very good paper",
					RebuttalList: Rebuttal{
						
					},
				},
			},
		},
	}

	data := &struct {
		Paper []comittedPaper
		Name  string
	}{
		Paper: paper,
		Name:  name,
	}
	showView(w, r, "index.html", data)
}

func (app *Application) AuthorHomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	paper := []comittedPaper{
		comittedPaper{
			Name:       "Bitcoin: Peer to peer ecoin system",
			AuthorList: []string{"hywang"},
			Reviews: map[string]Review{
				"id_reviewer": Review{
					ReviewerID: "id_reviewer",
					Content:    "Very good paper",
					RebuttalList: Rebuttal{
						
					},
				},
			},
		},
	}

	data := &struct {
		Paper []comittedPaper
		Name  string
	}{
		Paper: paper,
		Name:  name,
	}
	showView(w, r, "authorHome.html", data)
}

func (app *Application) AuthorCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	data := &struct {
		Name string
	}{
		Name: name,
	}
	showView(w, r, "commit_paper.html", data)
}

func (app *Application) AuthorUpdateView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	data := &struct {
		Name string
		Email string
		Paper []string
	}{
		Name: name,
		Email: "1141751053@qq.com",
		Paper: []string{
			"Bitcoin: Peer to peer ecoin system",
			"Bitcoin: Peer to peer ecoin system",
			"Bitcoin: Peer to peer ecoin system",
		},
	}
	showView(w, r, "updateAuthorInfo.html", data)
}

func (app *Application) ReviewerUpdateView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	data := &struct {
		Name string
		Email string
		ResearchTarget []string
		ReviewedPaper []string
		UNReviewedPaper []string
	}{
		Name: name,
		Email: "1141751053@qq.com",
		ResearchTarget: []string{
			"Blockchain",
			"AI security",
		},
		ReviewedPaper: []string{
			"Bitcoin: Peer to peer ecoin system",
			"Bitcoin: Peer to peer ecoin system",
			"Bitcoin: Peer to peer ecoin system",
		},
		UNReviewedPaper: []string{
			"Bitcoin: Peer to peer ecoin system",
			"Bitcoin: Peer to peer ecoin system",
			"Bitcoin: Peer to peer ecoin system",
		},
	}
	showView(w, r, "updateReviewerInfo.html", data)
}

func (app *Application) RebuttalreviewerView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	name := r.Form.Get("name")
	data := &struct{
		Name string
		Title string
	}{
		Name : name,
		Title : title,
	}
	showView(w, r, "rebuttal_reviewer.html", data)
}

func (app *Application) RebuttalauthorView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	name := r.Form.Get("name")
	data := &struct{
		Name string
		Title string
	}{
		Name : name,
		Title : title,
	}
	showView(w, r, "rebuttal_reviewer.html", data)
}

func (app *Application) SeereviewView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	name := r.Form.Get("name")
	review := "Very good paper"
	data := &struct{
		Name string
		Title string
		Review string
	}{
		Name : name,
		Title : title,
		Review: review,
	}
	showView(w, r, "seereview.html", data)
}

func (app *Application) ReplyView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	reviewerName := r.Form.Get("reviewername")
	title := r.Form.Get("title")
	
	data := &struct{
		Name string
		Title string
	}{
		Name : reviewerName,
		Title : title,
	}
	showView(w, r, "reply.html", data)
}

func (app *Application) ReviewerHomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	paper := []reviewedPaper{
		{
			Name:   "Bitcoin: peer to peer ecoin system",
			Review: "Very good paper",
			RebuttalList: Rebuttal{
				ReviewerID: "id_reviewer",
				AuthorID:   "id_author",
				Question:   "How are you?",
				Reply:      "",
				IsReplyed:  false,
			},
		},
	}
	data := &struct {
		Paper []reviewedPaper
		Name  string
	}{
		Paper: paper,
		Name:  name,
	}
	showView(w, r, "reviewerHome.html", data)
}

func (app *Application) ReviewerCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	data := &struct {
		Paper []string
		Name  string
	}{
		Paper: []string{"Ethereum: Second generation blockchain"},
		Name:  name,
	}
	showView(w, r, "reviewerCommit.html", data)
}

func (app *Application) RegisterView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	identity := r.Form.Get("identity")
	if identity == "author" {
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		email := r.Form.Get("email")
		ID := randomID()
		arguments := []string{name, ID, passwd, email}

		fmt.Println(arguments)
		app.LoginView(w, r)
	} else if identity == "reviewer" {
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		email := r.Form.Get("email")
		researchTarget := r.Form.Get("researchTarget")
		ID := randomID()
		arguments := []string{name, ID, passwd, email, researchTarget}
		app.LoginView(w, r)
		fmt.Println(arguments)
	}
}

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "sign_in.html", nil)
}

func (app *Application) ReviewView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	paper := []reviewedPaper{
		{
			Name:   "Bitcoin: peer to peer ecoin system",
			Review: "Very good paper",
			RebuttalList: Rebuttal{
				ReviewerID: "id_reviewer",
				AuthorID:   "id_author",
				Question:   "How are you?",
				Reply:      "",
				IsReplyed:  false,
			},
		},
	}
	data := &struct {
		Paper []reviewedPaper
		Name  string
	}{
		Paper: paper,
		Name:  name,
	}
	showView(w, r, "review.html", data)
}

func (app *Application) ReviewedPaperView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	paper := []reviewedPaper{
		{
			Name:   "Bitcoin: peer to peer ecoin system",
			Review: "Very good paper",
			RebuttalList: Rebuttal{
				ReviewerID: "id_reviewer",
				AuthorID:   "id_author",
				Question:   "How are you?",
				Reply:      "",
				IsReplyed:  false,
			},
		},
	}
	data := &struct {
		Paper []reviewedPaper
		Name  string
	}{
		Paper: paper,
		Name:  name,
	}
	showView(w, r, "reviewed_paper.html", data)
}

func (app *Application) CommittedPaperView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	paper := []comittedPaper{
		comittedPaper{
			Name:       "Bitcoin: Peer to peer ecoin system",
			AuthorList: []string{"hywang"},
			Reviews: map[string]Review{
				"id_reviewer": Review{
					ReviewerID: "id_reviewer",
					Content:    "Very good paper",
					RebuttalList: Rebuttal{
						
					},
				},
			},
		},
	}

	data := &struct {
		Paper []comittedPaper
		Name  string
	}{
		Paper: paper,
		Name:  name,
	}
	showView(w, r, "committed_paper.html", data)
}

func (app *Application) RegisterAuthorView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "sign_up.html", nil)
}

func (app *Application) RegisterReviewerView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "registerReviewer.html", nil)
}

func (app *Application) ReviewPaperView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	title := r.Form.Get("title")
	data := &struct {
		Title string
		Name  string
	}{
		Title: title,
		Name:  name,
	}
	showView(w, r, "review_paper.html", data)
}

func (app *Application) CommitPaperView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	authorlist := r.Form.Get("authorlist")
	keywords := r.Form.Get("keywords")
	fmt.Println(title, authorlist, keywords)
	ID, err := app.Upload(w, r)
	if err != nil {
		log.Fatalf("Failed to Upload Paper: %s", err)
	}
	arguments := []string{title, ID, authorlist, keywords}
	fmt.Println(arguments)
	fmt.Fprintf(w, "Commit Paper %s successfully", title)
}

func (app *Application) Upload(w http.ResponseWriter, r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("paper")
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadFile(handler.Filename)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	randStr := make([]byte, 10)
	for i := 0; i < 10; i++ {
		b := rand.Intn(26) + 65
		randStr[i] = byte(b)
	}

	h := sha256.New()
	h.Write([]byte(string(content) + string(randStr)))
	sum := h.Sum(nil)
	s := hex.EncodeToString(sum)
	return string(s), nil
}

func (app *Application) CommitReviewView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	reviewerName := r.Form.Get("name")
	reviewContent := r.Form.Get("reviewContent")
	arguments := []string{title, reviewerName, reviewContent}
	fmt.Println(arguments)

	fmt.Fprintf(w, "Add Review to %s successfully %s", title, reviewContent)
}

func randomID() string {
	rand.Seed(time.Now().UnixNano())
	randStr := make([]byte, 10)
	for i := 0; i < 10; i++ {
		b := rand.Intn(26) + 65
		randStr[i] = byte(b)
	}
	h := sha256.New()
	h.Write([]byte(string(randStr)))
	sum := h.Sum(nil)
	id := hex.EncodeToString(sum)
	return id
}
