package web

import (
	"MobileInternet/web/controller"
	"fmt"
	"net/http"
)

func WebStart(app *controller.Application) {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	ds := http.FileServer(http.Dir("upload"))
	http.Handle("/upload/", http.StripPrefix("/upload/", ds))

	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.LoginView)
	http.HandleFunc("/home", app.HomeView)
	http.HandleFunc("/register", app.RegisterView)
	http.HandleFunc("/registerCommit", app.RegisterCommitView)

	http.HandleFunc("/authorCommit", app.AuthorCommitView)
	http.HandleFunc("/commitPaper",app.CommitPaperView)

	http.HandleFunc("/committedPaper",app.CommittedPaperView)

	http.HandleFunc("/updateUser",app.UpdateUserView)
	http.HandleFunc("/updateUserCommit",app.UpdateUserCommitView)

	http.HandleFunc("/updatePaper", app.UpdatePaperView)
	http.HandleFunc("/commitUpdatePaper", app.UpdatePaperCommitView)

	http.HandleFunc("/rebuttal_reviewer",app.RebuttalreviewerView)
	http.HandleFunc("/commitReply",app.CommitReplyView)

	http.HandleFunc("/rebuttal_author",app.RebuttalauthorView)
	http.HandleFunc("/commitRebuttal",app.CommitRebuttalView)


	//paper的review, 对应seereview.html author视角
	http.HandleFunc("/review",app.SeereviewView)
	http.HandleFunc("/reviewPaper",app.ReviewView) // reviewPaper review主页面，对应review.html
	http.HandleFunc("/review_paper",app.ReviewPaperView) //review_paper review paper的具体页面 ，对应review_paper.html
	http.HandleFunc("/reviewedPaper",app.ReviewedPaperView)

	http.HandleFunc("/commitReview",app.ReviewCommitView)

	http.HandleFunc("/reject",app.RejectView)

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}
}
