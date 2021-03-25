package main

import (
	"log"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"MobileInternet/web"
	"MobileInternet/web/controller"
)

const (
	orgName  = "Org1"
	orgAdmin = "Admin"
)

func main() {
	configPath := "./e2e.yaml"
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	defer sdk.Close()
	queryChannel(sdk)

	app := controller.Application{
	}

	web.WebStart(&app)
}

func queryChannel(sdk *fabsdk.FabricSDK) {
	configBackend, err := sdk.Config()
	if err != nil {
		log.Fatalf("Failed to get config backend from SDK: %s", err)
	}
	targets, err := orgTargetPeers([]string{orgName}, configBackend)
	if err != nil {
		log.Fatalf("creating peers failed: %s", err)
	}

	clientContext := sdk.Context(fabsdk.WithUser("User1"), fabsdk.WithOrg("Org1"))
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("failed to query channel management client:%s", err)
	}
	channelQueryResponse, err := resMgmtClient.QueryChannels(
		resmgmt.WithTargetEndpoints(targets[0]), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("QueryChannels return error: %s", err)
	}
	for _, channel := range channelQueryResponse.Channels {
		log.Printf("***  Channel :%s\n", channel.ChannelId)
	}
}

func orgTargetPeers(orgs []string, configBackend ...core.ConfigBackend) ([]string, error) {
	networkConfig := fab.NetworkConfig{}
	err := lookup.New(configBackend...).UnmarshalKey("organizations", &networkConfig.Organizations)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get organizations from config ")
	}

	var peers []string
	for _, org := range orgs {
		orgConfig, ok := networkConfig.Organizations[strings.ToLower(org)]
		if !ok {
			continue
		}
		peers = append(peers, orgConfig.Peers...)
	}
	return peers, nil
}