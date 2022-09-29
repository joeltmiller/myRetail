package main

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestGetProductModel(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("test get product model", func(mt *mtest.T) {
		expectedModel := PriceModel{Price: 12.29, ID: "qw2eqwe21e3", CurrencyCode: "USD", ProductID: 123}

		getResponse := mtest.CreateCursorResponse(1, "myRetail.pricing", mtest.FirstBatch, bson.D{
			{"_id", expectedModel.ID},
			{"price", expectedModel.Price},
			{"currency_code", expectedModel.CurrencyCode},
			{"product_id", expectedModel.ProductID},
		})
		mt.AddMockResponses(getResponse)
		repository := productRepository{client: mt.Client}
		response, _ := repository.GetProductModel(123)
		assert.Equal(t, &expectedModel, response)
	})
}

func TestUpdateProductModel(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("test update product model", func(mt *mtest.T) {
		expectedModel := PriceModel{Price: 13.19, ID: "e8wr88w8rw", CurrencyCode: "USD", ProductID: 456}

		updatedResponse := mtest.CreateCursorResponse(1, "myRetail.pricing", mtest.FirstBatch, bson.D{
			{"matched_count", 1},
			{"modified_count", 1},
			{"upserted_count", 0},
			{"upserted_id", nil},
		})
		getResponse := mtest.CreateCursorResponse(1, "myRetail.pricing", mtest.FirstBatch, bson.D{
			{"_id", expectedModel.ID},
			{"price", expectedModel.Price},
			{"currency_code", expectedModel.CurrencyCode},
			{"product_id", expectedModel.ProductID},
		})
		mt.AddMockResponses(updatedResponse, getResponse)
		repository := productRepository{client: mt.Client}
		response, _ := repository.UpdateProductModel(123, 14.46)
		assert.Equal(t, &expectedModel, response)
	})
}
