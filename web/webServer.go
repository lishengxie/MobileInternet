package web

import (
	"net/http"
	"fmt"
	"web/controller"
)

func WebStart() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", IndexView)
	http.HandleFunc("/index", IndexView)
	http.HandleFunc("/setInfo", SetInfoView)

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}
}