package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb      *redis.Client
	redisCtx context.Context
}

func NewService() Service {
	return Service{
		redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", "kong-redis", "6379"),
			//Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
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

func (s *Service) GenerateCacheKey(requestBody string, requestHeader string, requestPath string) (cacheKey string, shouldCached bool, err error) {
	var graphQLReq GraphQLRequest
	if err := json.Unmarshal([]byte(requestBody), &graphQLReq); err != nil {
		return "", false, fmt.Errorf("err GenerateCacheKey unmarshal request body: %w", err)
	}

	graphQLAST, err := s.GetAndNormalizeGraphQLAst(graphQLReq.Query)
	if err != nil {
		return "", false, err
	}

	// TODO: trung.bc - TD should check before normalizing the ast
	shouldCached = s.reqOperationIsQuery(graphQLAST)
	if !shouldCached {
		return "", shouldCached, err
	}

	graphQLAstBytes, err := json.Marshal(graphQLAST)
	if err != nil {
		return "", false, fmt.Errorf("err GenerateCacheKey marshal graphQLAst: %w", err)
	}

	graphQlVariableStr := s.NormalizeGraphQLVariable(graphQLReq.Variables)

	request := fmt.Sprintf("%s%s%s", requestHeader, string(graphQLAstBytes), graphQlVariableStr)
	requestHashBytes := md5.Sum([]byte(request))
	requestHash := fmt.Sprintf("%s/%x", requestPath, string(requestHashBytes[:]))
	requestHash = strings.ReplaceAll(requestHash, "/", "-")

	return requestHash, true, nil
}

func (s *Service) reqOperationIsQuery(graphQLAST *ast.Document) bool {
	for _, definition := range graphQLAST.Definitions {
		operationDef, ok := definition.(*ast.OperationDefinition)
		if !ok {
			continue
		}

		if operationDef.Operation != string(Query) {
			return false
		}
	}

	return true
}

func (s *Service) GetAndNormalizeGraphQLAst(graphQLQuery string) (*ast.Document, error) {
	defer func() {
		message := recover()
		if message != nil {
			fmt.Println(message)
		}
	}()

	graphQLAST, err := s.GetGraphQLAst(graphQLQuery)
	if err != nil {
		return nil, err
	}

	if gConf.DisableNormalize {
		return graphQLAST, err
	}

	s.NormalizeOperationName(graphQLAST)
	s.NormalizeGraphQLAST(reflect.ValueOf(graphQLAST).Elem())

	return graphQLAST, err
}

func (s *Service) GetGraphQLAst(graphQLQuery string) (*ast.Document, error) {
	source := source.NewSource(&source.Source{
		Body: []byte(graphQLQuery),
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
		return nil, fmt.Errorf("err GetGraphQLAst parsing graphql req: %w", err)
	}

	return graphQLAst, nil
}

func (s *Service) NormalizeOperationName(graphQLAST *ast.Document) {
	for _, definition := range graphQLAST.Definitions {
		operationDef, ok := definition.(*ast.OperationDefinition)
		if !ok {
			continue
		}

		operationDef.Name = nil
		if operationDef.VariableDefinitions == nil {
			operationDef.VariableDefinitions = make([]*ast.VariableDefinition, 0)
		}
	}
}

func (s *Service) NormalizeGraphQLVariable(variableMp map[string]interface{}) (variableStr string) {
	keys := make([]string, 0, len(variableMp))

	for key := range variableMp {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		variableStr += fmt.Sprintf("%s:%v,", key, variableMp[key])
	}

	return variableStr
}

func (s *Service) NormalizeGraphQLAST(nodeVal reflect.Value) {
	if nodeVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < nodeVal.NumField(); i++ {
		//fmt.Println(nodeVal.Field(i).Type(), nodeVal.Field(i).Kind())
		subNodeVal := nodeVal.Field(i)

		switch subNodeVal.Kind() {
		case reflect.Interface:
			if subNodeVal.Elem().Kind() != reflect.Invalid {
				s.NormalizeGraphQLAST(subNodeVal.Elem().Elem())
			}
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
