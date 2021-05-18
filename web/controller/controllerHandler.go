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
	"strconv"
	"time"
)

const (
	INT_MAX = int(^uint32((0)) >> 1)
	INT_MIN = ^INT_MAX
)

func min(a,b int) int {
	if a<= b{
		return a
	}
	return b
}

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
	Valid 		 bool 				 `json:"valid"`
	Content      string              `json:"content"`
	RebuttalList map[string]Rebuttal `json:"rebuttallist"` //rebuttalID => rebuttal
}

type comittedPaper struct {
	Name       		string            `json:"name"`
	AuthorList 		[]string          `json:"authorlist"`
	Reviews    		map[string]Review `json:"reviews"`
	Successed  		bool			  `json:"successed"`
	Passed          bool              `json:"passed"`
	ReviewFinished	bool 		   	  `json:"reviewfinished"`
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

type UserInfo struct {
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	ResearchTarget  []string `json:"researchtarget"`
	ReviewedPaper   []string `json:"reviewedpaper"`
	UNReviewedPaper []string `json:"unreviewedpaper"`
	CommittedPaper  []string `json:"committedpaper"`
	AltCoin         float64  `json:"altcoin"`
}

type PaperInfo struct {
	Title      string            `json:"title"`
	KeyWords   []string          `json:"keywords"`
	AuthorList []string          `json:"authorlist"`
	ReviewList map[string]Review `json:"reviewlist"`
	StorePath      string            `json:"storepath"`
}

type IndexReview struct {
	ReviewerID string
	Review Review
}

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "sign_in.html", nil)
}

func (app *Application) RegisterView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "sign_up.html", nil)
}

func (app *Application) RegisterCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}

	name := r.Form.Get("name")
	passwd := r.Form.Get("passwd")
	confirmpasswd := r.Form.Get("confirmpasswd")
	email := r.Form.Get("email")
	researchTarget := r.Form.Get("researchTarget")
	if passwd != confirmpasswd {
		log.Fatalf("passwd doesn't matching")
	}

	ID := randomID()
	arguments := []string{name, ID, passwd, email, researchTarget}
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("CreateUser", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "GetPaperInfo", err)
	}
	fmt.Println(resp.TxValidationCode)
	app.LoginView(w, r)
}

func (app *Application) HomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	// fmt.Println(name, passwd, identity)
	name := r.Form.Get("name")
	passwd := r.Form.Get("passwd")

	if passwd != ""{
		arguments := []string{name, passwd}
		resp, err := app.Service.InvokeChaincode("ValidateUser", arguments)
		if resp == nil {
			log.Fatalf("Null Response")
		}
		if string(resp.Payload) == "false" {
			log.Fatalf("No such user %s", name)
		}
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	

	arguments := []string{name}
	resp, err := app.Service.InvokeChaincode("AuthorCommittedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AuthorCommittedPaper", err)
	}
	var cPaper []comittedPaper
	if resp.Payload == nil {
		cPaper = []comittedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &cPaper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	arguments = []string{name}
	resp, err = app.Service.InvokeChaincode("ReviewerReviewedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "ReviewerReviewedPaper", err)
	}
	var rPaper []reviewedPaper
	if resp.Payload == nil {
		rPaper = []reviewedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &rPaper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	arguments = []string{name}
	resp, err = app.Service.InvokeChaincode("ReviewerUNReviewedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "ReviewerUNReviewedPaper", err)
	}
	var urPaper []unReviewedPaper
	if resp.Payload == nil {
		urPaper = []unReviewedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &urPaper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	arguments = []string{name}
	resp, err = app.Service.InvokeChaincode("GetUserInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "GetUserInfo", err)
	}
	var userInfo UserInfo
	if resp.Payload == nil {
		userInfo = UserInfo{}
	} else {
		err = json.Unmarshal(resp.Payload, &userInfo)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	data := &struct {
		CommittedPaper []comittedPaper
		ReviewedPaper  []reviewedPaper
		UNReviewedPaper []unReviewedPaper
		Name           string
		AltCoin			string
	}{
		CommittedPaper: cPaper,
		ReviewedPaper:  rPaper,
		UNReviewedPaper: urPaper,
		Name:           name,
		AltCoin:		strconv.FormatFloat(userInfo.AltCoin,'f',2,64),
	}
	showView(w, r, "index.html", data)
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

func (app *Application) CommitPaperView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	title := r.Form.Get("title")
	authorList := r.Form.Get("authorlist")
	keywords := r.Form.Get("keywords")
	ID, path, err := app.Upload(w, r)

	arguments := []string{name, title, ID, authorList, keywords, path}
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("CreatePaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AuthorCommittedPaper", err)
	}

	fmt.Println(resp.TxValidationCode)

	app.CommittedPaperView(w,r)
}

func (app *Application) UpdateUserView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	arguments := []string{name}
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("GetUserInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "GetUserInfo", err)
	}

	var userInfo UserInfo
	if resp.Payload == nil {
		userInfo = UserInfo{}
	} else {
		err = json.Unmarshal(resp.Payload, &userInfo)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	data := &struct {
		Name            string   `json:"name"`
		Email           string   `json:"email"`
		ResearchTarget  []string `json:"researchtarget"`
		AltCoin			string
	}{
		Name:            userInfo.Name,
		Email:           userInfo.Email,
		ResearchTarget:  userInfo.ResearchTarget,
		AltCoin:		 strconv.FormatFloat(userInfo.AltCoin,'f',2,64),
	}
	showView(w, r, "updateUser.html", data)
}

func (app *Application) UpdateUserCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}

	name := r.Form.Get("name")
	old_passwd := r.Form.Get("old_passwd")
	new_passwd := r.Form.Get("new_passwd")
	new_email := r.Form.Get("new_email")
	researchTarget := r.Form.Get("researchTarget")

	arguments := []string{name, old_passwd, new_passwd, new_email, researchTarget}
	resp, err := app.Service.InvokeChaincode("UpdateUserInfo", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "UpdateUserInfo", err)
	}
	fmt.Println(resp.TxValidationCode)
	app.UpdateUserView(w,r)
}

func (app *Application) UpdatePaperView(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) UpdatePaperCommitView(w http.ResponseWriter, r *http.Request) {
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
	app.CommittedPaperView(w,r)
}

func (app *Application) CommittedPaperView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	// fmt.Println(name, passwd, identity)
	name := r.Form.Get("name")

	arguments := []string{name}
	resp, err := app.Service.InvokeChaincode("AuthorCommittedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AuthorCommittedPaper", err)
	}
	var cPaper []comittedPaper
	if resp.Payload == nil {
		cPaper = []comittedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &cPaper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	data := &struct {
		Paper []comittedPaper
		Name  string
	}{
		cPaper,
		name,
	}
	showView(w, r, "committed_paper.html", data)
}

func (app *Application) RebuttalreviewerView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	reviewerName := r.Form.Get("name")
	title := r.Form.Get("title")

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
	fmt.Println(paperInfo)

	arguments = []string{reviewerName}
	fmt.Println(arguments)
	resp, err = app.Service.InvokeChaincode("GetUserID", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "GetUserID", err)
	}

	var reviewerID string
	if resp.Payload == nil {
		reviewerID = ""
	} else {
		reviewerID = string(resp.Payload)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	var rebuttalid int
	rebuttalid = INT_MAX
	for _, rebuttal := range paperInfo.ReviewList[reviewerID].RebuttalList{
		if rebuttal.IsReplyed == false{
			tmp,err := strconv.Atoi(rebuttal.RebuttalID)
			fmt.Println(tmp)
			rebuttalid = min(tmp,rebuttalid)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}

	fmt.Println(rebuttalid)

	data := &struct {
		Name       string
		Title      string
		Review	 	Review
		RebuttalID  string
	}{
		Name:       reviewerName,
		Title:      title,
		Review:  	paperInfo.ReviewList[reviewerID],
		RebuttalID: strconv.Itoa(rebuttalid),
	}
	showView(w, r, "rebuttal_reviewer.html", data)
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
	app.RebuttalreviewerView(w,r)
}

func (app *Application) RebuttalauthorView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	authorName := r.Form.Get("name")
	title := r.Form.Get("title")
	reviewerID := r.Form.Get("reviewerid")

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
		Title      string
		Review	 	Review
	}{
		Name:       authorName,
		Title:      title,
		Review:  	paperInfo.ReviewList[reviewerID],
	}
	showView(w, r, "rebuttal_author.html", data)
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
	fmt.Println(arguments)

	resp, err := app.Service.InvokeChaincode("AddRebuttal", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddRebuttal", err)
	}
	fmt.Println(resp.TxValidationCode)
	app.RebuttalauthorView(w,r)
}



func (app *Application) SeereviewView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	name := r.Form.Get("name")

	arguments := []string{title}
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

	review := make([]IndexReview,0)

	for key, value := range paperInfo.ReviewList {
		review = append(review,IndexReview{
			ReviewerID:key,
			Review:value,
		})
	}

	data := &struct {
		Name   string
		Title  string
		Review []IndexReview
		AuthorList []string
		KeyWords []string
	}{
		Name:   name,
		Title:  title,
		Review: review,
		AuthorList: paperInfo.AuthorList,
		KeyWords: paperInfo.KeyWords,
	}
	showView(w, r, "seereview.html", data)
}

func (app *Application) ReviewedPaperView(w http.ResponseWriter, r *http.Request) {
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
	var rPaper []reviewedPaper
	if resp.Payload == nil {
		rPaper = []reviewedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &rPaper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	data := &struct {
		Paper []reviewedPaper
		Name  string
	}{
		Paper: rPaper,
		Name:  name,
	}
	showView(w, r, "reviewed_paper.html", data)
}

func (app *Application) ReviewView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}

	name := r.Form.Get("name")
	arguments := []string{name}
	resp, err := app.Service.InvokeChaincode("ReviewerUNReviewedPaper", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "ReviewerUNReviewedPaper", err)
	}
	var uRPaper []unReviewedPaper
	if resp.Payload == nil {
		uRPaper = []unReviewedPaper{}
	} else {
		err = json.Unmarshal(resp.Payload, &uRPaper)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	data := &struct {
		Paper []unReviewedPaper
		Name  string
	}{
		Paper: uRPaper,
		Name:  name,
	}
	showView(w, r, "review.html", data)
}

func (app *Application) ReviewPaperView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	title := r.Form.Get("title")

	arguments := []string{title}
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
		Title string
		Name  string
		StorePath string
		AuthorList []string
		KeyWords   []string
	}{
		Title: title,
		Name:  name,
		StorePath: paperInfo.StorePath,
		AuthorList: paperInfo.AuthorList,
		KeyWords: paperInfo.KeyWords,
	}
	showView(w, r, "review_paper.html", data)
}

func (app *Application) ReviewCommitView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	reviewerName := r.Form.Get("name")
	reviewContent := r.Form.Get("reviewcontent")
	valid := r.Form.Get("valid")
	arguments := []string{title, reviewerName, reviewContent, valid}

	resp, err := app.Service.InvokeChaincode("AddReview", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "AddReview", err)
	}
	fmt.Println(resp.TxValidationCode)

	app.ReviewedPaperView(w,r)
}

func (app *Application) RejectView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to parse Form: %s", err)
	}
	title := r.Form.Get("title")
	rejected := r.Form.Get("rejected")

	arguments := []string{title, rejected}

	resp, err := app.Service.InvokeChaincode("Reject", arguments)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode %s : %s", "Reject", err)
	}
	fmt.Println(resp.TxValidationCode)
	app.ReviewedPaperView(w,r)
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
