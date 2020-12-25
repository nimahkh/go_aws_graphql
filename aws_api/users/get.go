package aws_api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"os"
)

type Users struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetAll(getAdmin bool) (users []Users, err error) {
	sess, err := session.NewSession()

	svc := iam.New(sess, &aws.Config{Region: aws.String("us-west-2")})

	user_filter := "User"
	input := &iam.GetAccountAuthorizationDetailsInput{Filter: []*string{&user_filter}}
	resp, err := svc.GetAccountAuthorizationDetails(input)
	if err != nil {
		fmt.Println("Got error getting account details")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	adminName := "AdministratorAccess"

	for _, user := range resp.UserDetailList {
		aws_user := Users{
			ID:   *user.Arn,
			Name: *user.UserName,
		}

		if getAdmin {
			isAdmin := IsUserAdmin(svc, user, adminName)
			if isAdmin {
				users = append(users, aws_user)
			}
		} else {
			users = append(users, aws_user)
		}
	}

	return users, nil
}
