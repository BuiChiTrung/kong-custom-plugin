package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"log"
)

func main() {
	source := source.NewSource(&source.Source{
		//Body: []byte("{\n    country(code: \"VN\") {\n        native,\n        capital,\n        emoji,    \n        name,\n    }\n}"),
		//Body: []byte("query asdlasdj {\n    country(code: \"VN\") {\n        native,\n        capital,\n        emoji,    \n        name,\n    }\n}"),
		//Body: []byte("query {\n    country(code: \"VN\") {\n        native,\n        capital,\n        emoji,    \n        name,\n    }\n}"),
		//Body: []byte("query Repository($name: String!, $owner: String!, $followRenames: Boolean) {\n  repository(name: $name, owner: $owner, followRenames: $followRenames) {\n    allowUpdateBranch\n    autoMergeAllowed\n    createdAt\n    id\n    isPrivate\n    owner {\n      avatarUrl\n      id\n      login\n      url\n    }\n  }\n}\n"),
		Body: []byte("mutation AddReactionToIssue {\n  addReaction(input:{subjectId:\"MDU6SXNzdWUyMzEzOTE1NTE=\",content:HOORAY}) {\n    reaction {\n      content\n    }\n    subject {\n      id\n    }\n  }\n}"),
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

	//normalizeNode(reflect.ValueOf(a).Elem())
	fmt.Println(string(PrintJSON(a)))
}

func PrintJSON(obj interface{}) []byte {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	return bytes
}
