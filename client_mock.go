package main

type ClientMock struct {
	Responses ClientResponses
}

type ClientResponses struct {
	getProductName  string
	getProductError error
}

func (cm ClientMock) GetProductName(id string) (string, error) {
	return cm.Responses.getProductName, cm.Responses.getProductError
}
