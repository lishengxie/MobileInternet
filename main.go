package main

import (
	"MobileInternet/web"
	"MobileInternet/web/controller"
)

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
	peer1          = "peer0.org1.example.com"
	userName       = "User1"
)

func main() {

	app := controller.Application{}
	web.WebStart(&app)
}
