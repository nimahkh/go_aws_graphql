package aws_api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"os"
)

type Users struct {
	ID   string  `json:"id"`
	Name string `json:"name"`
}

func GetAll(getAdmin bool) (users []Users, err error) {
	sess, err := session.NewSession()

	svc := iam.New(sess, &aws.Config{Region: aws.String("us-west-2")})

	//user_filter := "User"
	//input := &iam.GetAccountAuthorizationDetailsInput{Filter: []*string{&user_filter}}
	//resp, err := svc.GetAccountAuthorizationDetails(input)
	if err != nil {
		fmt.Println("Got error getting account details")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//adminName := "AdministratorAccess"
	//numUsers:=0
	//numAdmins:=0
	// Wade through resulting users
	//for _, user := range resp.UserDetailList {
	//	numUsers += 1
	//	if getAdmin {
	//	//	isAdmin := IsUserAdmin(svc, user, adminName)
	//
	//		//if isAdmin {
	//		//	fmt.Println(*user.UserName)
	//		//	numAdmins += 1
	//		//}
	//	}
	//}

	result, err := svc.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(10),
	})

	if err != nil {
		fmt.Println("Error", err)
		return users , err
	}

	for _, user := range result.Users {
		if user == nil {
			continue
		}

		//if getAdmin {
		//	isAdmin := IsUserAdmin(svc, user, adminName)
		//}
		aws_user := Users{
			ID: *user.Arn,
			Name: *user.UserName,
		}
		users = append(users, aws_user)
		fmt.Printf("User: %v \n", user)

		//fmt.Printf("%d user %s created %v\n", i, *user.UserName, user.CreateDate)
	}
	fmt.Printf("User: %v \n", users)

	return users, nil
}
