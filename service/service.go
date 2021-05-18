package service

import (
	"log"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

type ServiceSetup struct {
	Sdk    *fabsdk.FabricSDK
	Client *channel.Client
}

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
	peer1          = "peer0.org1.example.com"
	userName       = "User1"
	chaincode      = "smartreview"
)

func (s *ServiceSetup) InvokeChaincode(function string, arguments []string) (*channel.Response, error) {
	args := packArgs(arguments)
	req := channel.Request{
		ChaincodeID: chaincode,
		Fcn:         function,
		Args:        args,
	}
	reqPeers := channel.WithTargetEndpoints("peer0.org1.example.com", "peer0.org2.example.com")
	resp, err := s.Client.Execute(req, reqPeers)
	if err != nil {
		return nil, err
	}

	log.Printf("query chaincode tx: %s", resp.TransactionID)
	log.Printf("result: %v", string(resp.Payload))
	log.Printf("Status: %v", string(resp.ChaincodeStatus))

	return &resp, nil
}

func (s *ServiceSetup) queryChannel() {
	configBackend, err := s.Sdk.Config()
	if err != nil {
		log.Fatalf("Failed to get config backend from SDK: %s", err)
	}
	targets, err := s.orgTargetPeers([]string{orgName}, configBackend)
	if err != nil {
		log.Fatalf("creating peers failed: %s", err)
	}

	clientContext := s.Sdk.Context(fabsdk.WithUser("User1"), fabsdk.WithOrg("Org1"))
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

func (s *ServiceSetup) orgTargetPeers(orgs []string, configBackend ...core.ConfigBackend) ([]string, error) {
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

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}
