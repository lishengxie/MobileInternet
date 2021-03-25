package controller

import (
	"net/http"
)

type Application struct {

}

func (app *Application) IndexView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "index.html", nil)
}

func (app *Application) SetInfoView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "setinfo.html", nil)
}