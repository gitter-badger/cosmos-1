package main

type ErrorJson map[string]interface{}

func NewErrorJson(code int, errmsg string) ErrorJson {
	json := make(ErrorJson)
	json["code"] = code
	json["message"] = errmsg
	return json
}
