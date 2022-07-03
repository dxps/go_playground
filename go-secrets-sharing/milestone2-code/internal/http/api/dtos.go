package api

import "encoding/json"

type AddSecretInput struct {
	PlainText string `json:"plain_text"`
}

type AddSecretOutput struct {
	ID string `json:"id"`
}

func NewAddSecretOutput(id string) *AddSecretOutput {
	return &AddSecretOutput{id}
}

// TODO: We have have a generic function to unmarshal bytes to specific type.
//       This would replace all of these ...FromBytes() functions.

func NewAddSecretOutputFromBytes(data []byte) (*AddSecretOutput, error) {
	var res AddSecretOutput
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type GetSecretOutput struct {
	Data string `json:"data"`
}

func NewGetSecretOutput(data string) GetSecretOutput {
	return GetSecretOutput{data}
}

func NewGetSecretOutputFromBytes(data []byte) (*GetSecretOutput, error) {
	var res GetSecretOutput
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type ResponseError struct {
	Error string `json:"error"`
}

func NewResponseErrorFromBytes(data []byte) (*ResponseError, error) {
	var res ResponseError
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
