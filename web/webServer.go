package web

import (
	"fmt"
	"net/http"
	"MobileInternet/web/controller"
)

func WebStart(app *controller.Application) {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.IndexView)
	http.HandleFunc("/index", app.IndexView)
	http.HandleFunc("/setInfo", app.SetInfoView)

	fmt.Println("启动Web服务, 监听端口号: 8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}
}