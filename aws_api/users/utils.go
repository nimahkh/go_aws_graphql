package aws_api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"os"
)

func UserPolicyHasAdmin(user *iam.UserDetail, admin string) bool {
	for _, policy := range user.UserPolicyList {
		if *policy.PolicyName == admin {
			return true
		}
	}
	return false
}

func AttachedUserPolicyHasAdmin(user *iam.UserDetail, admin string) bool {
	for _, policy := range user.AttachedManagedPolicies {
		if *policy.PolicyName == admin {
			return true
		}
	}

	return false
}

func GroupPolicyHasAdmin(svc *iam.IAM, group *iam.Group, admin string) bool {
	input := &iam.ListGroupPoliciesInput{
		GroupName: group.GroupName,
	}

	result, err := svc.ListGroupPolicies(input)
	if err != nil {
		fmt.Println("Got error calling ListGroupPolicies for group", group.GroupName)
	}

	// Wade through policies
	for _, policyName := range result.PolicyNames {
		if
		*policyName == admin {
			return true
		}
	}

	return false
}

func AttachedGroupPolicyHasAdmin(svc *iam.IAM, group *iam.Group, admin string) bool {
	input := &iam.ListAttachedGroupPoliciesInput{GroupName: group.GroupName}
	result, err := svc.ListAttachedGroupPolicies(input)
	if err != nil {
		fmt.Println("Got error getting attached group policies:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, policy := range result.AttachedPolicies {
		if *policy.PolicyName == admin {
			return true
		}
	}

	return false
}

func IsUserAdmin(svc *iam.IAM, user *iam.UserDetail, admin string) bool {
	// Check policy, attached policy, and groups (policy and attached policy)
	policyHasAdmin := UserPolicyHasAdmin(user, admin)
	if policyHasAdmin {
		return true
	}

	attachedPolicyHasAdmin := AttachedUserPolicyHasAdmin(user, admin)
	if attachedPolicyHasAdmin {
		return true
	}

	userGroupsHaveAdmin := UsersGroupsHaveAdmin(svc, user, admin)
	if userGroupsHaveAdmin {
		return true
	}

	return false
}

func UsersGroupsHaveAdmin(svc *iam.IAM, user *iam.UserDetail, admin string) bool {
	input := &iam.ListGroupsForUserInput{UserName: user.UserName}
	result, err := svc.ListGroupsForUser(input)
	if err != nil {
		fmt.Println("Got error getting groups for user:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, group := range result.Groups {
		groupPolicyHasAdmin := GroupPolicyHasAdmin(svc, group, admin)

		if groupPolicyHasAdmin {
			return true
		}

		attachedGroupPolicyHasAdmin := AttachedGroupPolicyHasAdmin(svc, group, admin)

		if attachedGroupPolicyHasAdmin {
			return true
		}
	}

	return false
}
