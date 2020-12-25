package main

import (
	"encoding/json"
	"fmt"
	users "github.com/nimahkh/go_aws_graphql/aws_api/users"
	"net/http"
	"github.com/graphql-go/graphql"
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
			/* Get (read) single product by id
			   http://localhost:8080/product?query={product(id:1){name,info,price}}
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
			/* Get (read) product list
			   http://localhost:8080/product?query={list{id,name,info,price}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "Get users list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var users, _ =users.GetAll(false)
					return users, nil
				},
			},
		},
	})

//var mutationType = graphql.NewObject(graphql.ObjectConfig{
//	Name: "Mutation",
//	Fields: graphql.Fields{
//		/* Create new product item
//		http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
//		*/
//		"create": &graphql.Field{
//			Type:        productType,
//			Description: "Create new product",
//			Args: graphql.FieldConfigArgument{
//				"name": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.String),
//				},
//				"info": &graphql.ArgumentConfig{
//					Type: graphql.String,
//				},
//				"price": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.Float),
//				},
//			},
//			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//				rand.Seed(time.Now().UnixNano())
//				product := Product{
//					ID:    int64(rand.Intn(100000)), // generate random ID
//					Name:  params.Args["name"].(string),
//					Info:  params.Args["info"].(string),
//					Price: params.Args["price"].(float64),
//				}
//				products = append(products, product)
//				return product, nil
//			},
//		},
//
//		/* Update product by id
//		   http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
//		*/
//		"update": &graphql.Field{
//			Type:        productType,
//			Description: "Update product by id",
//			Args: graphql.FieldConfigArgument{
//				"id": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.Int),
//				},
//				"name": &graphql.ArgumentConfig{
//					Type: graphql.String,
//				},
//				"info": &graphql.ArgumentConfig{
//					Type: graphql.String,
//				},
//				"price": &graphql.ArgumentConfig{
//					Type: graphql.Float,
//				},
//			},
//			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//				id, _ := params.Args["id"].(int)
//				name, nameOk := params.Args["name"].(string)
//				info, infoOk := params.Args["info"].(string)
//				price, priceOk := params.Args["price"].(float64)
//				product := Product{}
//				for i, p := range products {
//					if int64(id) == p.ID {
//						if nameOk {
//							products[i].Name = name
//						}
//						if infoOk {
//							products[i].Info = info
//						}
//						if priceOk {
//							products[i].Price = price
//						}
//						product = products[i]
//						break
//					}
//				}
//				return product, nil
//			},
//		},
//
//		/* Delete product by id
//		   http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
//		*/
//		"delete": &graphql.Field{
//			Type:        productType,
//			Description: "Delete product by id",
//			Args: graphql.FieldConfigArgument{
//				"id": &graphql.ArgumentConfig{
//					Type: graphql.NewNonNull(graphql.Int),
//				},
//			},
//			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//				id, _ := params.Args["id"].(int)
//				product := Product{}
//				for i, p := range products {
//					if int64(id) == p.ID {
//						product = products[i]
//						// Remove from product list
//						products = append(products[:i], products[i+1:]...)
//					}
//				}
//
//				return product, nil
//			},
//		},
//	},
//})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
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
	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	fmt.Println("Server is running on port 8085")
	http.ListenAndServe(":8085", nil)
}
