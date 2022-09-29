package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	GetProductModel(id int) (*PriceModel, error)
	UpdateProductModel(id int, price float64) (*PriceModel, error)
}

type productRepository struct {
	client *mongo.Client
}

func (pr productRepository) GetProductModel(id int) (*PriceModel, error) {
	coll := pr.client.Database("myRetail").Collection("pricing")

	result := bson.M{}
	err := coll.FindOne(context.Background(), bson.D{{"product_id", id}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the id %v\n", id)
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return nil, err
	}

	priceModel := PriceModel{}
	unmarshalErr := json.Unmarshal(jsonData, &priceModel)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &priceModel, nil
}

func (pr productRepository) UpdateProductModel(id int, price float64) (*PriceModel, error) {
	coll := pr.client.Database("myRetail").Collection("pricing")

	filter := bson.D{{"product_id", id}}
	update := bson.D{{"$set", bson.D{{"price", price}}}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return nil, err
	}

	updateModel := mongo.UpdateResult{}
	unmarshalErr := json.Unmarshal(jsonData, &updateModel)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	productModel, err := pr.GetProductModel(id)
	if err != nil {
		return nil, err
	}

	return productModel, nil
}

type PriceModel struct {
	ID           string  `json:"_id"`
	Price        float64 `json:"price"`
	ProductID    int     `json:"product_id"`
	CurrencyCode string  `json:"currency_code"`
}
