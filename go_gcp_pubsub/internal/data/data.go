package data

import "encoding/json"

type MyData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewMyData() *MyData {
	return &MyData{
		Name: "some testing",
	}
}

func NewMyDataFromBytes(data []byte) (*MyData, error) {
	var obj MyData
	err := json.Unmarshal(data, &obj)
	return &obj, err
}
