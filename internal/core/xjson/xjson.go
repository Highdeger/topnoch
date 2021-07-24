package xjson

import "encoding/json"

func StructToJson(parsingStruct interface{}) string {
	byts, _ := json.Marshal(parsingStruct)
	return string(byts)
}

func JsonToStruct(parsingJson string, out interface{}) {
	_ = json.Unmarshal([]byte(parsingJson), out)
}
