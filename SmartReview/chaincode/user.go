package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

const (
	userInitCoin = 100.0
)

type User struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Passwd          string   `json:"passwd"`
	Email           string   `json:"email"`
	AltCoin			float64  `json:"altcoin"`
	ResearchTarget  []string `json:"researchtarget"`
	ReviewedPaper   []string `json:"reviewedpaper"`
	UNReviewedPaper []string `json:"unreviewedpaper"`
	CommittedPaper  []string `json:"committedpaper"`
}

type UserInfo struct {
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	ResearchTarget  []string `json:"researchtarget"`
	ReviewedPaper   []string `json:"reviewedpaper"`
	UNReviewedPaper []string `json:"unreviewedpaper"`
	CommittedPaper  []string `json:"committedpaper"`
	AltCoin			float64  `json:"altcoin"`
}

func (s *SmartContract) GetUserSet(ctx contractapi.TransactionContextInterface) (*UserSet, error) {
	userSetJson, err := ctx.GetStub().GetState("userset")
	if err != nil {
		return nil, err
	}

	var userSet UserSet
	err = json.Unmarshal(userSetJson, &userSet)
	if err != nil {
		return nil, err
	}

	return &userSet, nil
}

func (s *SmartContract) AddToUserSet(ctx contractapi.TransactionContextInterface, name string, id string) error {
	userSet, err := s.GetUserSet(ctx)
	if err != nil {
		return err
	}

	if _, ok := userSet.Users[name]; ok {
		return fmt.Errorf("user %s already exists", name)
	}

	userSet.Users[name] = id
	newUserSet := UserSet{
		Users: userSet.Users,
	}
	newUserSetJSON, err := json.Marshal(newUserSet)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("userset", newUserSetJSON)
	if err != nil {
		return err
	}
	return nil
}


func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface,
	name string, id string, passwd string, email string, researchTargets string) error {

	researchTargets = strings.Trim(researchTargets, " ")
	researchTargetsArr := strings.Split(researchTargets,"/")

	err := s.AddToUserSet(ctx,name,id)
	if err != nil {
		return err
	}

	user := User{
		ID: 				id,
		Name:				name,
		Passwd: 			passwd,
		Email: 				email,
		AltCoin:  			userInitCoin,
		ResearchTarget: 	researchTargetsArr,
		ReviewedPaper: 		[]string{},
		UNReviewedPaper: 	[]string{},
		CommittedPaper: 	[]string{},
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, userJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) UpdateUserInfo(ctx contractapi.TransactionContextInterface,
	name string, old_passwd string, new_passwd string, email string, researchTarget string) error {
	userID, err := s.GetUserID(ctx, name)
	if err != nil {
		return err
	}
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if old_passwd != user.Passwd {
		return fmt.Errorf("wrong passwd")
	}

	researchTarget = strings.Trim(researchTarget, " ")
	researchTargetArr := strings.Split(researchTarget, "/")
	
	user.Passwd = new_passwd
	user.Email = email
	user.ResearchTarget = researchTargetArr
	
	userJSON, err := json.Marshal(*user)
	if err != nil {
		return err
	}
	
	err = ctx.GetStub().PutState(userID, userJSON)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *SmartContract) TransferCoin(ctx contractapi.TransactionContextInterface, fromUserName string, toUserName string, transferNum float64) error {
	fromUserID, err := s.GetUserID(ctx, fromUserName)
	if err != nil {
		return err
	}
	fromUser, err := s.GetUser(ctx, fromUserID)
	if err != nil {
		return err
	}

	if fromUser.AltCoin < transferNum {
		return fmt.Errorf("alt coin num is %f, less than %f", fromUser.AltCoin, transferNum)
	}

	fromUser.AltCoin -= transferNum
	fromUserJSON, err := json.Marshal(*fromUser)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(fromUserID, fromUserJSON)


	toUserID, err := s.GetUserID(ctx, toUserName)
	if err != nil {
		return err
	}
	toUser, err := s.GetUser(ctx, toUserID)
	if err != nil {
		return err
	}

	toUser.AltCoin += transferNum
	toUserJSON, err := json.Marshal(*toUser)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(toUserID, toUserJSON)

	return nil
}


func (s *SmartContract) GetUserID(ctx contractapi.TransactionContextInterface, name string) (string, error) {
	userSet, err := s.GetUserSet(ctx)
	if err != nil {
		return "", err
	}

	if _, ok := userSet.Users[name]; !ok {
		return "", fmt.Errorf("user %s not exist", name)
	}
	return userSet.Users[name], nil
}

func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil,err
	}

	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil,err
	}

	return &user, nil
}

func (s *SmartContract) GetUserInfo(ctx contractapi.TransactionContextInterface, name string) (*UserInfo, error) {
	userID, err := s.GetUserID(ctx,name)
	if err != nil {
		return nil,err
	}

	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	reviewedPaper := make([]string, 0)
	unreviewedPaper := make([]string, 0)
	committedPaper := make([]string,0)

	for _, each := range user.ReviewedPaper {
		paper, err := s.GetPaper(ctx, each)
		if err != nil {
			return nil, err
		}
		reviewedPaper = append(reviewedPaper, paper.Title)
	}

	for _, each := range user.UNReviewedPaper {
		paper, err := s.GetPaper(ctx, each)
		if err != nil {
			return nil, err
		}
		unreviewedPaper = append(unreviewedPaper, paper.Title)
	}

	for _, each := range user.CommittedPaper {
		paper, err := s.GetPaper(ctx, each)
		if err != nil {
			return nil, err
		}
		committedPaper = append(committedPaper, paper.Title)
	}

	return &UserInfo{
		Name:user.Name,
		Email: user.Email,
		ResearchTarget: user.ResearchTarget,
		ReviewedPaper: reviewedPaper,
		UNReviewedPaper: unreviewedPaper,
		CommittedPaper: committedPaper,
		AltCoin:	user.AltCoin,
	},nil
}

func (s *SmartContract) ValidateUser(ctx contractapi.TransactionContextInterface, name string, passwd string) (bool, error){
	userID, err := s.GetUserID(ctx,name)
	if err != nil{
		return false, err
	}

	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, fmt.Errorf("user %s not exist", name)
	}

	return user.Passwd == passwd, nil
}
