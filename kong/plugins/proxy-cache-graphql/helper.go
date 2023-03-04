package main

import "encoding/json"

func getObjBytes(obj interface{}) []byte {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	return bytes
}

func getObjJSONString(obj interface{}) string {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	return string(bytes)
}
