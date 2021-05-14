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
	http.HandleFunc("/register", app.RegisterAuthorView)
	http.HandleFunc("/registerReviewer", app.RegisterReviewerView)
	http.HandleFunc("/authorHome", app.AuthorHomeView)
	http.HandleFunc("/authorCommit", app.AuthorCommitView)
	http.HandleFunc("/commitPaper",app.CommitPaperView)
	http.HandleFunc("/updateAuthor",app.AuthorUpdateView)
	http.HandleFunc("/committedPaper",app.CommittedPaperView)

	http.HandleFunc("/rebuttal_reviewer",app.RebuttalreviewerView)
	http.HandleFunc("/rebuttal_author",app.RebuttalauthorView)
	http.HandleFunc("/reply",app.ReplyView)
	
	http.HandleFunc("/reviewerHome", app.ReviewerHomeView)
	http.HandleFunc("/reviewerCommit", app.ReviewerCommitView)

	//paper的review, 对应seereview.html
	http.HandleFunc("/review",app.SeereviewView)
	
	http.HandleFunc("/updateReviewer",app.ReviewerUpdateView)
	// reviewPaper review主页面，对应review.html
	http.HandleFunc("/reviewPaper",app.ReviewView)
	//review_paper review paper的具体页面 ，对应review_paper.html
	http.HandleFunc("/review_paper",app.ReviewPaperView)
	http.HandleFunc("/reviewedPaper",app.ReviewedPaperView)

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}
}
