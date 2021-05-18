package main

import (
	"MobileInternet/service"
	"MobileInternet/web"
	"MobileInternet/web/controller"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	// "fmt"
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

	configPath := "./e2e.yaml"
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}

	serv := service.ServiceSetup{
		Sdk:sdk,
		Client:client,
	}


	// err = serv.InitUser("json/authors.json", "json/reviewers.json")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// err = serv.InitSimilarityPair("json/similarity.json")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// err = serv.InitPaper("json/papers.json")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	app := controller.Application{
		Service: &serv,
	}
	web.WebStart(&app)
}
