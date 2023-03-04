package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb      *redis.Client
	redisCtx context.Context
}

func NewService() *Service {
	return &Service{
		redis.NewClient(&redis.Options{
			Addr:     "kong-redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		context.Background(),
	}
}

func (s *Service) GetCacheKey(cacheKey string) (string, error) {
	val, err := s.rdb.Get(s.redisCtx, cacheKey).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (s *Service) InsertCacheKey(cacheKey string, value string, expireNanoSec int64) error {
	_, err := s.rdb.Get(s.redisCtx, cacheKey).Result()
	if err == redis.Nil {
		if err := s.rdb.Set(s.redisCtx, cacheKey, value, time.Duration(expireNanoSec)).Err(); err != nil {
			return err
		}
	}

	return err
}

func (s *Service) GenerateCacheKey(requestBody []byte, requestHeader []byte, requestPath string) (string, error) {
	graphQLAstBytes, err := GetGraphQLAst(requestBody)
	if err != nil {
		return "", err
	}
	gKong.Log.Notice(string(graphQLAstBytes))

	request := append(requestHeader, graphQLAstBytes...)
	requestHashBytes := md5.Sum(request)
	requestHash := fmt.Sprintf("%s/%x", requestPath, string(requestHashBytes[:]))

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
