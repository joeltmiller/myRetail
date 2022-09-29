package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	GetProductName(id string) (string, error)
}

type client struct {
	BaseUrl string
	ApiKey  string
}

func NewClient(baseUrl string, apiKey string) Client {
	return client{BaseUrl: baseUrl, ApiKey: apiKey}
}

func (c client) GetProductName(id string) (string, error) {
	url := fmt.Sprintf("%s%s&tcin=%s", c.BaseUrl, c.ApiKey, id)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching product details from internal resource")
		return "", err
	}

	if resp.StatusCode == 404 {
		fmt.Printf("%v with id: %v\n", NotFound, id)
		return "", errors.New(NotFound)
	}

	body := ProductDetailsResponse{}
	bytes, _ := io.ReadAll(resp.Body)
	unmarshalError := json.Unmarshal(bytes, &body)

	if unmarshalError != nil {
		return "", unmarshalError
	}

	return body.Data.Product.Item.ProductDescription.Title, nil
}

type ProductDetailsResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Product Product `json:"product"`
}

type Product struct {
	ItemNumber string `json:"tcin"`
	Item       Item   `json:"item"`
}

type Item struct {
	ProductDescription ProductDescription `json:"product_description"`
}

type ProductDescription struct {
	Title string `json:"title"`
}
