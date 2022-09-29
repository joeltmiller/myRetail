package main

type RepositoryMock struct {
	Responses RepositoryResponses
}

type RepositoryResponses struct {
	getProductModel    *PriceModel
	updateProductModel *PriceModel
}

func (rm RepositoryMock) GetProductModel(id int) (*PriceModel, error) {
	return rm.Responses.getProductModel, nil
}

func (rm RepositoryMock) UpdateProductModel(id int, price float64) (*PriceModel, error) {
	return rm.Responses.updateProductModel, nil
}
