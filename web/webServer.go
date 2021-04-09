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
	http.HandleFunc("/register", app.RegisterView)
	http.HandleFunc("/registerReviewer", app.RegisterReviewerView)

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}
}
