package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql/language/ast"
	"reflect"
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
	graphQLAstBytes, err := s.GetAndNormalizeGraphQLAst(requestBody)
	if err != nil {
		return "", err
	}
	gKong.Log.Notice(string(graphQLAstBytes))

	request := append(requestHeader, graphQLAstBytes...)
	requestHashBytes := md5.Sum(request)
	requestHash := fmt.Sprintf("%s/%x", requestPath, string(requestHashBytes[:]))

	return requestHash, nil
}

func (s *Service) GetAndNormalizeGraphQLAst(requestBody []byte) ([]byte, error) {
	graphQLAST, err := s.GetGraphQLAst(requestBody)
	if err != nil {
		return nil, err
	}

	s.NormalizeGraphQLAST(reflect.ValueOf(graphQLAST).Elem())

	graphQLAstBytes, err := json.Marshal(graphQLAST)
	if err != nil {
		return nil, fmt.Errorf("err marshal indent: %w", err)
	}

	return graphQLAstBytes, nil
}

func (s *Service) GetGraphQLAst(requestBody []byte) (*ast.Document, error) {
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

	return graphQLAst, nil
}

func (s *Service) NormalizeGraphQLAST(nodeVal reflect.Value) {
	if nodeVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < nodeVal.NumField(); i++ {
		fmt.Println(nodeVal.Field(i).Type(), nodeVal.Field(i).Kind())
		subNodeVal := nodeVal.Field(i)

		switch subNodeVal.Kind() {
		case reflect.Interface:
			s.NormalizeGraphQLAST(subNodeVal.Elem().Elem())
			//s.NormalizeGraphQLAST(reflect.ValueOf(subNodeVal.Interface()).Elem())
		case reflect.Ptr:
			s.NormalizeGraphQLAST(subNodeVal.Elem())
		case reflect.Struct:
			s.NormalizeGraphQLAST(subNodeVal)
		case reflect.Slice:
			for j := 0; j < subNodeVal.Len(); j++ {
				s.NormalizeGraphQLAST(reflect.ValueOf(subNodeVal.Index(j).Interface()).Elem())
			}

			s.sortSliceNode(subNodeVal)
		}
	}
}

func (s *Service) sortSliceNode(nodeVal reflect.Value) {
	for i := 0; i < nodeVal.Len(); i++ {
		for j := i + 1; j < nodeVal.Len(); j++ {
			hashNodeI := s.hashNodeVal(nodeVal.Index(i))
			hashNodeJ := s.hashNodeVal(nodeVal.Index(j))

			//fmt.Println(i, hashNodeI, j, hashNodeJ)

			if hashNodeI > hashNodeJ {
				tmp := reflect.ValueOf(nodeVal.Index(i).Interface())
				nodeVal.Index(i).Set(nodeVal.Index(j))
				nodeVal.Index(j).Set(tmp)
			}
		}
	}
}

func (s *Service) hashNodeVal(nodeVal reflect.Value) string {
	hashNodeBytes := md5.Sum(getObjBytes(nodeVal.Interface()))
	hashNode := fmt.Sprintf("%x", string(hashNodeBytes[:]))

	return hashNode
}
