package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
}

func (s *Service) NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateCacheKey(requestBody []byte, requestHeader []byte) (string, error) {
	graphQLAstBytes, err := GetGraphQLAst(requestBody)
	if err != nil {
		return "", err
	}
	gKong.Log.Notice(string(graphQLAstBytes))

	request := append(requestHeader, graphQLAstBytes...)
	requestHashBytes := md5.Sum(request)
	requestHash = fmt.Sprintf("%x", string(requestHashBytes[:]))

	return requestHash, nil
}

func GetGraphQLAst(requestBody []byte) ([]byte, error) {
	var graphQLReq GraphQLRequest
	if err := json.Unmarshal(requestBody, &graphQLReq); err != nil {
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
