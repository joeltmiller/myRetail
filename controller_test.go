package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetProductHandler(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		repository       ProductRepository
		expectedResponse ProductModel
		errorResponse    ErrorResponse
		client           Client
		expectError      bool
	}{
		{
			name:             "Success",
			id:               "13860428",
			repository:       RepositoryMock{Responses: RepositoryResponses{getProductModel: &PriceModel{Price: 12.22, CurrencyCode: "USD", ProductID: 13860428}}},
			client:           ClientMock{Responses: ClientResponses{getProductName: "The Big Lebowski (Blu-ray)"}},
			expectedResponse: ProductModel{ID: 13860428, Name: "The Big Lebowski (Blu-ray)", CurrentPrice: Price{Value: 12.22, CurrencyCode: "USD"}},
		},
		{
			name:          "Not found",
			id:            "23860428",
			client:        ClientMock{Responses: ClientResponses{getProductError: errors.New("no product found with id 23860428")}},
			errorResponse: ErrorResponse{ErrorMessages: []ErrorMessage{{Message: "no product found with id 23860428"}}},
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "", nil)

			r = mux.SetURLVars(r, map[string]string{"id": tt.id})

			if err != nil {
				t.Logf("error")
			}

			controller := NewProductController(tt.repository, tt.client)

			controller.GetProductHandler(rr, r)

			result := rr.Result()
			if tt.expectError {
				errorResponse := ErrorResponse{}
				unmarshalErr := json.NewDecoder(result.Body).Decode(&errorResponse)
				if unmarshalErr != nil {
					t.Logf("Error unmarshalling")
				}
				expectedResponse := tt.errorResponse
				assert.Equal(t, expectedResponse, errorResponse)
			} else {
				productResponse := ProductModel{}
				unmarshalErr := json.NewDecoder(result.Body).Decode(&productResponse)
				if unmarshalErr != nil {
					t.Logf("Error unmarshalling")
				}
				expectedResponse := tt.expectedResponse
				assert.Equal(t, expectedResponse, productResponse)
			}
		})
	}
}

func TestUpdateProductHandler(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		repository       ProductRepository
		expectedResponse ProductModel
		errorResponse    ErrorResponse
		client           Client
		expectError      bool
	}{
		{
			name:             "Success",
			id:               "13860428",
			repository:       RepositoryMock{Responses: RepositoryResponses{updateProductModel: &PriceModel{Price: 12.22, CurrencyCode: "USD", ProductID: 13860428}}},
			client:           ClientMock{Responses: ClientResponses{getProductName: "The Big Lebowski (Blu-ray)"}},
			expectedResponse: ProductModel{ID: 13860428, Name: "The Big Lebowski (Blu-ray)", CurrentPrice: Price{Value: 12.22, CurrencyCode: "USD"}},
		},
		{
			name:          "Not found",
			id:            "13860428",
			client:        ClientMock{Responses: ClientResponses{getProductError: errors.New("no product found with id 23860428")}},
			errorResponse: ErrorResponse{ErrorMessages: []ErrorMessage{{Message: "no product found with id 23860428"}}},
			expectError:   true,
		},
		{
			name:          "IDs do not match",
			id:            "23860428",
			client:        ClientMock{Responses: ClientResponses{getProductError: errors.New("ID in path does not match request body")}},
			errorResponse: ErrorResponse{ErrorMessages: []ErrorMessage{{Message: "ID in path does not match request body"}}},
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			requestBody := ProductModel{
				ID:   13860428,
				Name: "Foo Bar",
				CurrentPrice: Price{
					Value:        12.29,
					CurrencyCode: "USD",
				},
			}

			marshalledBody, _ := ffjson.Marshal(requestBody)

			r, err := http.NewRequest(http.MethodPut, "/product/13860428", bytes.NewBuffer(marshalledBody))

			if err != nil {
				t.Logf("error")
			}

			r = mux.SetURLVars(r, map[string]string{"id": tt.id})

			controller := NewProductController(tt.repository, tt.client)

			controller.UpdateProductHandler(rr, r)

			result := rr.Result()
			if tt.expectError {
				errorResponse := ErrorResponse{}
				unmarshalErr := json.NewDecoder(result.Body).Decode(&errorResponse)
				if unmarshalErr != nil {
					t.Logf("Error unmarshalling")
				}
				expectedResponse := tt.errorResponse
				assert.Equal(t, expectedResponse, errorResponse)
			} else {
				productResponse := ProductModel{}
				unmarshalErr := json.NewDecoder(result.Body).Decode(&productResponse)
				if unmarshalErr != nil {
					t.Logf("Error unmarshalling")
				}
				expectedResponse := tt.expectedResponse
				assert.Equal(t, expectedResponse, productResponse)
			}
		})
	}
}
