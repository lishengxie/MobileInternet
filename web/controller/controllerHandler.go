package controller

import (
	"MobileInternet/service"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Application struct {
	Service *service.ServiceSetup
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
		}else{
			log.Fatalf("No such author %s", name)
			//app.HomeView(w, r)
		}
	} else if identity == "reviewer" {
		arguments := []string{name,passwd}
		resp,err := app.Service.InvokeChaincode("ValidateReviewer",arguments)
		if err!=nil {
			log.Fatalf("Failed to invoke chaincode %s : %s", "ValidateReviewer", err)
		}
		if string(resp.Payload) == "true" {
			app.ReviewerHomeView(w, r)
		}else{
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
	resp,err := app.Service.InvokeChaincode("AuthorCommittedPaper",arguments)
	if err!=nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AuthorCommittedPaper", err)
	}
	var paper []comittedPaper
	if resp.Payload == nil{
		paper = []comittedPaper{}
	}else {
		err = json.Unmarshal(resp.Payload, &paper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	data := &struct{
		Paper []comittedPaper
		Name  string
	}{
		Paper: paper,
		Name: name,
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
		Name : name,
	}
	showView(w, r, "authorCommit.html", data)
}

func (app *Application) ReviewerHomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	arguments := []string{name}
	resp,err := app.Service.InvokeChaincode("ReviewerReviewedPaper",arguments)
	if err!=nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "ReviewerReviewedPaper", err)
	}
	var paper []reviewedPaper
	if resp.Payload == nil {
		paper = []reviewedPaper{}
	}else {
		err = json.Unmarshal(resp.Payload, &paper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	data := &struct{
		Paper []reviewedPaper
		Name  string
	}{
		Paper: paper,
		Name: name,
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
	resp,err := app.Service.InvokeChaincode("ReviewerUNReviewedPaper",arguments)
	if err!=nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "ReviewerReviewedPaper", err)
	}
	var paper []string
	if resp.Payload == nil {
		paper = []string{}
	}else {
		err = json.Unmarshal(resp.Payload, &paper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	data := &struct{
		Paper []string
		Name  string
	}{
		Paper: paper,
		Name: name,
	}
	showView(w, r, "reviewerCommit.html", data)
}

func (app *Application) RegisterView(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil{
		log.Fatalf("Failed to parse Form: %s", err)
	}
	identity := r.Form.Get("identity")
	if identity == "author"{
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		email := r.Form.Get("email")
		ID := randomID()
		arguments := []string{name,ID,passwd,email}
		resp,err := app.Service.InvokeChaincode("CreateAuthor",arguments)
		if err!=nil {
			log.Fatalf("Failed to invoke chaincode %s : %s", "CreateAuthor", err)
		}
		fmt.Println(resp.TxValidationCode)
		app.LoginView(w, r)
	}else if identity == "reviewer"{
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		email := r.Form.Get("email")
		researchTarget := r.Form.Get("researchTarget")
		ID := randomID()
		arguments := []string{name,ID,passwd,email,researchTarget}
		resp,err := app.Service.InvokeChaincode("CreateReviewer",arguments)
		if err!=nil {
			log.Fatalf("Failed to invoke chaincode %s : %s", "CreateReviewer", err)
		}
		fmt.Println(resp.TxValidationCode)
		app.LoginView(w, r)
		fmt.Println(name+" "+passwd+" "+email+" "+researchTarget)
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

func (app *Application) ReviewView(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	title := r.Form.Get("title")
	data := &struct{
		Title string
		Name  string
	}{
		Title: title,
		Name: name,
	}
	showView(w, r, "paperReview.html", data)
}

func (app *Application) CommitPaperView(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	authorlist := r.Form.Get("authorlist")
	keywords := r.Form.Get("keywords")
	ID,err := app.Upload(w,r)
	if err != nil{
		log.Fatalf("Failed to Upload Paper: %s", err)
	}
	arguments := []string{title,ID,authorlist,keywords}
	resp,err := app.Service.InvokeChaincode("AddPaper",arguments)
	if err!=nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddPaper", err)
	}
	fmt.Println(resp.TxValidationCode)
	fmt.Fprintf(w,"Commit Paper %s successfully",title)
}

func (app *Application) Upload(w http.ResponseWriter, r *http.Request)(string,error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("paper")
	if err != nil {
		return "",err
	}
	defer file.Close()

	content, err := ioutil.ReadFile(handler.Filename)
	if err != nil {
		return "",err
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
	return string(s),nil
}

func (app *Application) CommitReviewView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	reviewerName := r.Form.Get("name")
	reviewContent := r.Form.Get("reviewContent")
	arguments := []string{title,reviewerName,reviewContent}
	resp,err := app.Service.InvokeChaincode("AddReview",arguments)
	if err!=nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddReview", err)
	}
	fmt.Println(resp.TxValidationCode)
	fmt.Fprintf(w,"Add Review to %s successfully %s",title,reviewContent)
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