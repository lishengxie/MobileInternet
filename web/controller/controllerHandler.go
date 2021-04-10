package controller

import (
	"MobileInternet/service"
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
	Service *service.ServiceSetup
}

func (app *Application) HomeView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil{
		log.Fatalf("Failed to parse Form: %s", err)
	}
	name := r.Form.Get("name")
	passwd := r.Form.Get("passwd")
	identity := r.Form.Get("identity")
	if identity == "author" {
		if name == "lishengxie" && passwd == "123456" {
			app.AuthorHomeView(w, r)
		}
	} else if identity == "reviewer" {
		if name == "lishengxie" && passwd == "123456" {
			app.ReviewerHomeView(w, r)
		}
	}
}

func (app *Application) AuthorHomeView(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name []string
	}{
		[]string{"1", "2"},
	}
	showView(w, r, "authorHome.html", data)
}

func (app *Application) AuthorCommitView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "authorCommit.html", nil)
}

func (app *Application) ReviewerHomeView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "reviewerHome.html", nil)
}

func (app *Application) ReviewerCommitView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "reviewerCommit.html", nil)
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
		fmt.Println(name+" "+passwd+" "+email)
	}else if identity == "reviewer"{
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		email := r.Form.Get("email")
		researchTarget := r.Form.Get("researchTarget")
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
	showView(w, r, "paperReview.html", nil)
}

func (app *Application) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	content, err := ioutil.ReadFile(handler.Filename)
	if err != nil {
		fmt.Println(err)
		return
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
	fmt.Fprintln(w, string(s))
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
