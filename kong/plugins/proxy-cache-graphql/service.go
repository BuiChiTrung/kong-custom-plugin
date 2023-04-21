package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/BuiChiTrung/kong-custom-plugin/kong/logger"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdbMaster       *redis.Client
	rdbReplicas     *redis.Client
	rdbWrite        *redis.Client
	rdbRead         *redis.Client
	rdbCtx          context.Context
	lastHealthCheck time.Time
}

func NewService() *Service {
	rdbCtx := context.Background()

	rdbMaster := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv(EnvRedisMasterHost), os.Getenv(EnvRedisMasterPort)),
	})
	rdbMaster.Do(context.Background(), "SLAVEOF", "NO", "ONE").Result()

	rdbReplicas := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv(EnvRedisReplicasHost), os.Getenv(EnvRedisReplicasPort)),
	})
	rdbReplicas.Do(context.Background(), "REPLICAOF", os.Getenv(EnvRedisMasterHost), os.Getenv(EnvRedisMasterPort)).Result()

	return &Service{
		rdbMaster:   rdbMaster,
		rdbReplicas: rdbReplicas,
		rdbWrite:    rdbMaster,
		rdbRead:     rdbReplicas,
		rdbCtx:      rdbCtx,
	}
}

func (s *Service) GetCacheKey(cacheKey string) (string, error) {
	val, err := s.rdbRead.Get(s.rdbCtx, cacheKey).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (s *Service) InsertCacheKey(cacheKey string, value string, expireNanoSec int64) error {
	_, err := s.rdbWrite.Get(s.rdbCtx, cacheKey).Result()
	if err == redis.Nil {
		if err := s.rdbWrite.Set(s.rdbCtx, cacheKey, value, time.Duration(expireNanoSec)).Err(); err != nil {
			return fmt.Errorf("err InsertCacheKey: %v", err)
		}
		return nil
	}

	return err
}

func (s *Service) GenerateCacheKey(requestBody string, requestHeader string, requestPath string) (cacheKey string, shouldCached bool, err error) {
	defer func() {
		message := recover()
		if message != nil {
			fmt.Println(message)
		}
	}()

	var graphQLReq GraphQLRequest
	if err := json.Unmarshal([]byte(requestBody), &graphQLReq); err != nil {
		return "", false, fmt.Errorf("err GenerateCacheKey unmarshal request body: %w", err)
	}

	graphQLAST, err := s.GetGraphQLAst(graphQLReq.Query)
	if err != nil {
		return "", false, err
	}

	if shouldCached = s.reqOperationIsQuery(graphQLAST); !shouldCached {
		return "", shouldCached, err
	}

	if !gConf.DisableNormalize {
		s.NormalizeOperationName(graphQLAST)
		s.NormalizeGraphQLAST(reflect.ValueOf(graphQLAST).Elem())
	}

	graphQLAstBytes, err := json.Marshal(graphQLAST)
	if err != nil {
		return "", false, fmt.Errorf("err GenerateCacheKey marshal graphQLAst: %w", err)
	}

	fmt.Printf("%v", graphQLAST.Definitions)

	request := fmt.Sprintf("%s%v%v", requestHeader, string(graphQLAstBytes), graphQLReq.Variables)
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

func (s *Service) HealthCheckRedis() {
	logger.Infof("[RdbRead] %s", s.rdbRead.String())
	logger.Infof("[RdbWrite] %s", s.rdbWrite.String())
	logger.Infof("[RdbReplicas] %s", s.rdbReplicas.String())
	logger.Infof("[RdbMaster] %s", s.rdbMaster.String())

	host, port := GetRdbHostPort(s.rdbMaster)
	_, errReplicas := s.rdbReplicas.Do(context.Background(), "REPLICAOF", host, port).Result()
	_, errMaster := s.rdbMaster.Do(context.Background(), "SLAVEOF", "NO", "ONE").Result()

	// Both instances are alive
	if errReplicas == nil && errMaster == nil {
		s.rdbRead = s.rdbReplicas
		return
	}

	// Both instances are death
	if errReplicas != nil && errMaster != nil {
		logger.Error("Both redis instances are death.")
		return
	}

	// Replicas instance is death
	if errReplicas != nil {
		logger.Errorf("Replicas instance is death: %s", s.rdbReplicas.String())
		s.rdbRead = s.rdbMaster
		return
	}

	// Master instance is death: promote replicas to master & update rdbRead, rdbWrite
	logger.Errorf("Master instance is death: %s", s.rdbMaster.String())
	host, port = GetRdbHostPort(s.rdbReplicas)

	_, err := s.rdbReplicas.Do(context.Background(), "SLAVEOF", "NO", "ONE").Result()
	if err != nil {
		return
	}

	tmp := s.rdbMaster
	s.rdbMaster = s.rdbReplicas
	s.rdbReplicas = tmp

	s.rdbRead = s.rdbMaster
	s.rdbWrite = s.rdbMaster
}

func GetRdbHostPort(rdb *redis.Client) (string, string) {
	options := rdb.Options()
	addr := strings.Split(options.Addr, ":")
	return addr[0], addr[1]
}
