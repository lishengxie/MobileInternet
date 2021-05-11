package controller

import (
	"MobileInternet/service"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Application struct {
	Service *service.ServiceSetup
}

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
	Content      string              `json:"content"`
	RebuttalList map[string]Rebuttal `json:"rebuttallist"`
}

type comittedPaper struct {
	Name       string            `json:"name"`
	AuthorList []string          `json:"authorlist"`
	Reviews    map[string]Review `json:"reviews"`
}

type reviewedPaper struct {
	Name         string              `json:"name"`
	Review       string              `json:"review"`
	RebuttalList map[string]Rebuttal `json:"rebuttallist"`
	StorePath    string              `json:"storepath"`
}

type unReviewedPaper struct {
	Name      string `json:"name"`
	StorePath string `json:"storepath"`
}

type AuthorInfo struct {
	Name           string   `json:"name"`
	Passwd         string   `json:"passwd"`
	Email          string   `json:"email"`
	CommittedPaper []string `json:"committedpaper"`
}

type ReviewerInfo struct {
	Name            string   `json:"name"`
	Passwd          string   `json:"passwd"`
	Email           string   `json:"email"`
	ResearchTarget  []string `json:"researchtarget"`
	ReviewedPaper   []string `json:"reviewedpaper"`
	UNReviewedPaper []string `json:"unreviewedpaper"`
}

type PaperInfo struct {
	Title      string            `json:"title"`
	KeyWords   []string          `json:"keywords"`
	AuthorList []string          `json:"authorlist"`
	ReviewList map[string]Review `json:"reviewlist"`
}

func (app *Application) HomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	passwd := r.Form.Get("passwd")
	identity := r.Form.Get("identity")
	if identity == "author" {
		arguments := []string{name, passwd}
		resp, err := app.Service.InvokeChaincode("ValidateAuthor", arguments)
		if err != nil {
			log.Fatalf("Failed to invoke chaincode %s : %s", "ValidateAuthor", err)
		}
		if string(resp.Payload) == "true" {
			app.AuthorHomeView(w, r)
		} else {
			log.Fatalf("No such author %s", name)
			//app.HomeView(w, r)
		}
	} else if identity == "reviewer" {
		arguments := []string{name, passwd}
		resp, err := app.Service.InvokeChaincode("ValidateReviewer", arguments)
		if err != nil {
			log.Fatalf("Failed to invoke chaincode %s : %s", "ValidateReviewer", err)
		}
		if string(resp.Payload) == "true" {
			app.ReviewerHomeView(w, r)
		} else {
			log.Fatalf("No such author %s", name)
			//app.HomeView(w, r)
		}
	}
}

func (app *Application) AuthorHomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	arguments := []string{name}
	resp, err := app.Service.InvokeChaincode("AuthorCommittedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AuthorCommittedPaper", err)
	}
	var paper []comittedPaper
	if resp.Payload == nil {
		paper = []comittedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &paper)
		if err != nil {
			log.Fatalf("%v", err)
		}
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
	showView(w, r, "authorCommit.html", data)
}

func (app *Application) AuthorUpdateView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	arguments := []string{name}
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("GetAuthorInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "GetAuthorInfo", err)
	}

	var authorInfo AuthorInfo
	if resp.Payload == nil {
		authorInfo = AuthorInfo{}
	} else {
		err = json.Unmarshal(resp.Payload, &authorInfo)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	data := &struct {
		Name  string   `json:"name"`
		Email string   `json:"email"`
		Paper []string `json:"committedpaper"`
	}{
		Name:  authorInfo.Name,
		Email: authorInfo.Email,
		Paper: authorInfo.CommittedPaper,
	}
	showView(w, r, "updateAuthorInfo.html", data)
}

func (app *Application) AuthorUpdateCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}

	name := r.Form.Get("name")
	old_passwd := r.Form.Get("old_passwd")
	new_passwd := r.Form.Get("new_passwd")
	email := r.Form.Get("email")

	arguments := []string{name, old_passwd, new_passwd, email}
	resp, err := app.Service.InvokeChaincode("UpdateAuthorInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "UpdateAuthorInfo", err)
	}
	fmt.Println(resp.TxValidationCode)

	app.ShowInfo(w, r, "Update Author Info Successfully.")
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
	ID, path, err := app.Upload(w, r)
	if err != nil {
		log.Fatalf("Failed to Upload Paper: %s", err)
	}
	arguments := []string{title, ID, authorlist, keywords, path}
	fmt.Println(arguments)
	resp, err := app.Service.InvokeChaincode("CreatePaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddPaper", err)
	}
	fmt.Println(resp.TxValidationCode)
	app.ShowInfo(w, r, fmt.Sprintf("Commit Paper \"%s\" successfully", title))
}

func (app *Application) PaperUpdateView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	name := r.Form.Get("name")
	arguments := []string{title}
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("GetPaperInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "GetPaperInfo", err)
	}

	var paperInfo PaperInfo
	if resp.Payload == nil {
		paperInfo = PaperInfo{}
	} else {
		err = json.Unmarshal(resp.Payload, &paperInfo)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	data := &struct {
		Name       string
		Title      string   `json:"title"`
		KeyWords   []string `json:"keywords"`
		AuthorList []string `json:"authorlist"`
		//ReviewList   map[string]Review `json:"reviewlist"`
	}{
		Name:       name,
		Title:      paperInfo.Title,
		KeyWords:   paperInfo.KeyWords,
		AuthorList: paperInfo.AuthorList,
	}
	showView(w, r, "updatePaperInfo.html", data)
}

func (app *Application) PaperUpdateCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}

	title := r.Form.Get("title")
	new_title := r.Form.Get("new_title")
	addedauthor := r.Form.Get("addedauthor")

	arguments := []string{title, new_title, addedauthor}
	resp, err := app.Service.InvokeChaincode("UpdatePaperInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "UpdatePaperInfo", err)
	}
	fmt.Println(resp.TxValidationCode)

	app.ShowInfo(w, r, fmt.Sprintf("Update Paper Info Successfully."))
}

func (app *Application) Upload(w http.ResponseWriter, r *http.Request) (string, string, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("paper")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	os.Mkdir("./upload", os.ModePerm)
	cur, err := os.Create("./upload/" + handler.Filename)
	defer cur.Close()
	if err != nil {
		return "", "", err
	}
	io.Copy(cur, file)

	content, err := ioutil.ReadFile(handler.Filename)
	if err != nil {
		return "", "", err
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
	return string(s), "./upload/" + handler.Filename, nil
}

func (app *Application) ReviewerHomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	arguments := []string{name}
	resp, err := app.Service.InvokeChaincode("ReviewerReviewedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "ReviewerReviewedPaper", err)
	}
	var paper []reviewedPaper
	if resp.Payload == nil {
		paper = []reviewedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &paper)
		if err != nil {
			log.Fatalf("%v", err)
		}
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
	arguments := []string{name}
	resp, err := app.Service.InvokeChaincode("ReviewerUNReviewedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "ReviewerReviewedPaper", err)
	}
	var paper []unReviewedPaper
	if resp.Payload == nil {
		paper = []unReviewedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &paper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	data := &struct {
		Paper []unReviewedPaper
		Name  string
	}{
		Paper: paper,
		Name:  name,
	}
	showView(w, r, "reviewerCommit.html", data)
}

func (app *Application) ReviewerUpdateView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	arguments := []string{name}
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("GetReviewerInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "GetReviewerInfo", err)
	}

	var reviewerInfo ReviewerInfo
	if resp.Payload == nil {
		reviewerInfo = ReviewerInfo{}
	} else {
		err = json.Unmarshal(resp.Payload, &reviewerInfo)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	data := &struct {
		Name            string   `json:"name"`
		Email           string   `json:"email"`
		ResearchTarget  []string `json:"researchtarget"`
		ReviewedPaper   []string `json:"reviewedpaper"`
		UNReviewedPaper []string `json:"unreviewedpaper"`
	}{
		Name:            reviewerInfo.Name,
		Email:           reviewerInfo.Email,
		ResearchTarget:  reviewerInfo.ResearchTarget,
		ReviewedPaper:   reviewerInfo.ReviewedPaper,
		UNReviewedPaper: reviewerInfo.UNReviewedPaper,
	}
	showView(w, r, "updateReviewerInfo.html", data)
}

func (app *Application) ReviewerUpdateCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}

	name := r.Form.Get("name")
	old_passwd := r.Form.Get("old_passwd")
	new_passwd := r.Form.Get("new_passwd")
	email := r.Form.Get("email")
	researchTarget := r.Form.Get("researchTarget")

	arguments := []string{name, old_passwd, new_passwd, email, researchTarget}
	resp, err := app.Service.InvokeChaincode("UpdateReviewerInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "UpdateReviewerInfo", err)
	}
	fmt.Println(resp.TxValidationCode)

	app.ShowInfo(w, r, "Update Reviewer Info Successfully.")
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
		resp, err := app.Service.InvokeChaincode("CreateAuthor", arguments)
		if err != nil {
			log.Fatalf("Failed to invoke chaincode %s : %s", "CreateAuthor", err)
		}
		fmt.Println(resp.TxValidationCode)
		app.LoginView(w, r)
	} else if identity == "reviewer" {
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		email := r.Form.Get("email")
		researchTarget := r.Form.Get("researchTarget")
		ID := randomID()
		arguments := []string{name, ID, passwd, email, researchTarget}
		fmt.Println(arguments)
		resp, err := app.Service.InvokeChaincode("CreateReviewer", arguments)
		if err != nil {
			log.Fatalf("Failed to invoke chaincode %s : %s", "CreateReviewer", err)
		}
		fmt.Println(resp.TxValidationCode)
		app.LoginView(w, r)

	}
}

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "login.html", nil)
}

func (app *Application) RegisterAuthorView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "registerAuthor.html", nil)
}

func (app *Application) RegisterReviewerView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "registerReviewer.html", nil)
}

func (app *Application) ReviewView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	title := r.Form.Get("title")
	path := r.Form.Get("path")
	data := &struct {
		Title     string
		Name      string
		StorePath string
	}{
		Title:     title,
		Name:      name,
		StorePath: path,
	}
	showView(w, r, "paperReview.html", data)
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
	resp, err := app.Service.InvokeChaincode("AddReview", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddReview", err)
	}
	fmt.Println(resp.TxValidationCode)

	app.ShowInfo(w, r, fmt.Sprintf("Add Review to \"%s\" successfully: \"%s\"", title, reviewContent))
}

func (app *Application) RebuttalView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	authorName := r.Form.Get("authorname")
	title := r.Form.Get("title")
	reviewerID := r.Form.Get("reviewerid")
	data := &struct {
		Name       string
		Title      string
		ReviewerID string
	}{
		Name:       authorName,
		Title:      title,
		ReviewerID: reviewerID,
	}
	showView(w, r, "rebuttal.html", data)
}

func (app *Application) ReplyView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	reviewerName := r.Form.Get("reviewername")
	title := r.Form.Get("title")
	rebuttalID := r.Form.Get("rebuttalid")

	data := &struct {
		Name       string
		Title      string
		RebuttalID string
	}{
		Name:       reviewerName,
		Title:      title,
		RebuttalID: rebuttalID,
	}
	showView(w, r, "reply.html", data)
}

func (app *Application) CommitRebuttalView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	authorName := r.Form.Get("name")
	title := r.Form.Get("title")
	reviewerID := r.Form.Get("reviewerid")
	question := r.Form.Get("rebuttalcontent")

	arguments := []string{title, authorName, reviewerID, question}

	resp, err := app.Service.InvokeChaincode("AddRebuttal", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddRebuttal", err)
	}
	fmt.Println(resp.TxValidationCode)

	app.ShowInfo(w, r, fmt.Sprintf("Add Rebuttal to \"%s\" successfully: \"%s\"", title, question))
}

func (app *Application) CommitReplyView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	reviewerName := r.Form.Get("name")
	title := r.Form.Get("title")
	reply := r.Form.Get("replycontent")
	rebuttalID := r.Form.Get("rebuttalid")

	arguments := []string{title, reviewerName, reply, rebuttalID}
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("AddReply", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddReply", err)
	}
	fmt.Println(resp.TxValidationCode)
	app.ShowInfo(w, r, fmt.Sprintf("Add Reply to \"%s\" successfully:\"%s\"", title, reply))
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

func (app *Application) ShowInfo(w http.ResponseWriter, r *http.Request, info string) {
	data := &struct {
		Content string
	}{
		Content: info,
	}
	showView(w, r, "blank.html", data)
}
