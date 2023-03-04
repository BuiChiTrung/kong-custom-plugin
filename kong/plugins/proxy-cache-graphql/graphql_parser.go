package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"log"
	"reflect"
)

func main() {
	requestBody := "{\"query\":\"query Query{country(code: \\\"TL\\\"){native,capital,name,emoji}}\"}"
	//requestBody = "{\"query\":\"# query Query{country(code: \\\"CN\\\") {\\n#     name,\\n#     native,\\n#     capital,    \\n#         emoji  #aaasd\\n# }}\\nmutation {\\n  createPerson(name: \\\"Bob\\\", age: 36) {\\n    name\\n    age\\n  }\\n}\",\"variables\":{}}"
	var graphQLReq GraphQLRequest

	if err := json.Unmarshal([]byte(requestBody), &graphQLReq); err != nil {
		panic(err)
	}

	source := source.NewSource(&source.Source{
		Body: []byte(graphQLReq.Query),
		Name: "",
	})

	a, err := parser.Parse(parser.ParseParams{
		Source: source,
		Options: parser.ParseOptions{
			NoSource:   true,
			NoLocation: true,
		},
	})

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	normalizeNode(reflect.ValueOf(a).Elem())
	fmt.Println(string(convertInterfaceToBytes(a)))
}

func convertInterfaceToBytes(obj interface{}) []byte {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	return bytes
}

func normalizeNode(nodeVal reflect.Value) {
	if nodeVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < nodeVal.NumField(); i++ {
		fmt.Println(nodeVal.Field(i).Type(), nodeVal.Field(i).Type().Kind())

		switch nodeVal.Field(i).Kind() {
		case reflect.Ptr:
			normalizeNode(nodeVal.Field(i).Elem())
		case reflect.Struct:
			normalizeNode(nodeVal.Field(i))
		case reflect.Slice:
			fmt.Println(nodeVal.Field(i).Len())
			for j := 0; j < nodeVal.Field(i).Len(); j++ {
				fmt.Println(reflect.ValueOf(nodeVal.Field(i).Index(j).Interface()).Elem().Kind(), "aaaaaaaaaaaa")
				normalizeNode(reflect.ValueOf(nodeVal.Field(i).Index(j).Interface()).Elem())
			}

			for j := 0; j < nodeVal.Field(i).Len(); j++ {
				hashNodeJBytes := md5.Sum(convertInterfaceToBytes(nodeVal.Field(i).Index(j).Interface()))
				hashNodeJ := fmt.Sprintf("%x", string(hashNodeJBytes[:]))

				for l := j + 1; l < nodeVal.Field(i).Len(); l++ {
					hashNodeLBytes := md5.Sum(convertInterfaceToBytes(nodeVal.Field(i).Index(l).Interface()))
					hashNodeL := fmt.Sprintf("%x", string(hashNodeLBytes[:]))

					fmt.Println(j, hashNodeJ, l, hashNodeL)

					if hashNodeJ > hashNodeL {
						tmp := reflect.ValueOf(nodeVal.Field(i).Index(j).Interface())
						nodeVal.Field(i).Index(j).Set(nodeVal.Field(i).Index(l))
						nodeVal.Field(i).Index(l).Set(tmp)
					}
				}
			}
		}
		//fmt.Println(nodeVal.Field(i), nodeVal.Field(i).Kind())
	}
}
