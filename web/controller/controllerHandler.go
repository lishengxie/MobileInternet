package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Application struct {
}

func (app *Application)HomeView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	passwd := r.Form.Get("passwd")
	identity := r.Form.Get("identity")
	if identity == "author" {
		if name == "lishengxie" && passwd == "123456" {
			app.AuthorHomeView(w,r)
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
		[]string{"1","2"},
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

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "login.html", nil)
}

func (app *Application) RegisterView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "register.html", nil)
}

func (app *Application) RegisterReviewerView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "registerReviewer.html", nil)
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
