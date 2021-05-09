package web

import (
	"MobileInternet/web/controller"
	"fmt"
	"net/http"
)

func WebStart(app *controller.Application) {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.LoginView)
	http.HandleFunc("/home", app.HomeView)
	http.HandleFunc("/registerAuthor", app.RegisterAuthorView)
	http.HandleFunc("/registerReviewer", app.RegisterReviewerView)
	http.HandleFunc("/authorHome", app.AuthorHomeView)
	http.HandleFunc("/authorCommit", app.AuthorCommitView)
	http.HandleFunc("/commitPaper",app.CommitPaperView)
	http.HandleFunc("/updateAuthor",app.AuthorUpdateView)
	http.HandleFunc("/commitUpdateAuthor",app.AuthorUpdateCommitView)

	http.HandleFunc("/updatePaper", app.PaperUpdateView)
	http.HandleFunc("/commitUpdatePaper",app.PaperUpdateCommitView)

	http.HandleFunc("/rebuttal",app.RebuttalView)
	http.HandleFunc("/reply",app.ReplyView)
	http.HandleFunc("/commitRebuttal",app.CommitRebuttalView)
	http.HandleFunc("/commitReply",app.CommitReplyView)

	http.HandleFunc("/previewPaper",app.PreviewView)

	http.HandleFunc("/reviewerHome", app.ReviewerHomeView)
	http.HandleFunc("/reviewerCommit", app.ReviewerCommitView)
	http.HandleFunc("/commitReview",app.CommitReviewView)
	http.HandleFunc("/register",app.RegisterView)
	http.HandleFunc("/review",app.ReviewView)
	http.HandleFunc("/updateReviewer",app.ReviewerUpdateView)
	http.HandleFunc("/commitUpdateReviewer",app.ReviewerUpdateCommitView)

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}
}
