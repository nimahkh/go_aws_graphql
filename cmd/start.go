package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	users "github.com/nimahkh/go_aws_graphql/aws_api/users"
	"net/http"
)

// Product contains information about one product

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Users",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) admin users
			   http://localhost:8085/users?query={admin{name}}
			*/
			"admin": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "get Admin users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users, _ = users.GetAll(true)
					fmt.Printf("User: %v \n", users)
					return users, nil
				},
			},
			/* Get (read) users list
			   http://localhost:8085/users?query={list{name}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "Get users list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var users, _ = users.GetAll(false)
					return users, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	fmt.Println("Server is running on port 8085")
	http.ListenAndServe(":8085", nil)
}
