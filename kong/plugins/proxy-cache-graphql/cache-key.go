package main

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

func GenerateHashKey(requestBody string, requestHeader string) (string, error) {

}

func Az(requestBody string) ([]byte, error) {
	var graphQLReq GraphQLRequest
	if err := json.Unmarshal([]byte(requestBody), &graphQLReq); err != nil {
		return nil, fmt.Errorf("err unmarshal request body: %w", err)
	}

	source := source.NewSource(&source.Source{
		Body: []byte(graphQLReq.Query),
		Name: "",
	})
	graphQLAst, err := parser.Parse(parser.ParseParams{
		Source: source,
		Options: parser.ParseOptions{
			NoSource:   true,
			NoLocation: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("err parsing graphql req: %w", err)
	}

	graphQLAstBytes, err := json.MarshalIndent(graphQLAst, "\t", "\t")
	if err != nil {
		return nil, fmt.Errorf("err marshal indent: %w", err)
	}

	return graphQLAstBytes, nil
}
