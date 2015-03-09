package main

type ResJson map[string]interface{}

func NewResJson() ResJson {
	return make(ResJson)
}
