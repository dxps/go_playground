package api

type AddSecretInput struct {
	PlainText string `json:"plain_text"`
}

type AddSecretOutput struct {
	ID string `json:"id"`
}

func NewAddSecretOutput(id string) *AddSecretOutput {
	return &AddSecretOutput{id}
}

type GetSecretOutput struct {
	Data string `json:"data"`
}

func NewGetSecretOutput(data string) *GetSecretOutput {
	return &GetSecretOutput{data}
}

type ResponseError struct {
	Error string `json:"error"`
}
