package main

type GraphQLRequest struct {
	Query     string
	Variables map[string]interface{}
}
