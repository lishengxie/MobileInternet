package web

import (
	"net/http"
	"fmt"
	"MobileInternet/web/controller"
)

func WebStart(app *controller.Application) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.IndexView)
	http.HandleFunc("/index", app.IndexView)
	http.HandleFunc("/setInfo", app.SetInfoView)

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}
}