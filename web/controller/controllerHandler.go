package controller

import (
	"net/http"
)

func IndexView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "index.html", nil)
}

func SetInfoView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "setinfo.html", nil)
}