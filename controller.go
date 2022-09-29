package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pquerna/ffjson/ffjson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"strconv"
)

type ProductController struct {
	repository ProductRepository
	client     Client
}

func NewProductController(repository ProductRepository, client Client) ProductController {
	return ProductController{repository: repository, client: client}
}

func (pc ProductController) SetRoutes(router *mux.Router) {
	router.HandleFunc("/products/{id}", pc.GetProductHandler).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", pc.UpdateProductHandler).Methods(http.MethodPut)
}

func (pc ProductController) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	numericID, _ := strconv.Atoi(id)

	name, err := pc.client.GetProductName(id)

	if err != nil {
		switch err.Error() {
		case NotFound:
			respondWithError(NotFound, w, 404)
		default:
			respondWithError(err.Error(), w, 500)
		}
		return
	}

	priceModel, err := pc.repository.GetProductModel(numericID)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			respondWithError(NotFound, w, 404)
		default:
			respondWithError(err.Error(), w, 500)
		}

		return
	}

	response := ProductModel{
		ID:   numericID,
		Name: name,
		CurrentPrice: Price{
			Value:        priceModel.Price,
			CurrencyCode: priceModel.CurrencyCode,
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (pc ProductController) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	numericID, _ := strconv.Atoi(id)

	rawBody, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(err.Error(), w, http.StatusInternalServerError)
		return
	}

	body := ProductModel{}
	unmarshalError := json.Unmarshal(rawBody, &body)

	if unmarshalError != nil {
		respondWithError(unmarshalError.Error(), w, http.StatusInternalServerError)
		return
	}

	if body.ID != numericID {
		respondWithError(IdDoesNotMatch, w, http.StatusBadRequest)
		return
	}

	// Fetch name and verify exists
	name, err := pc.client.GetProductName(id)

	if err != nil {
		switch err.Error() {
		case NotFound:
			respondWithError(NotFound, w, 404)
		default:
			respondWithError(err.Error(), w, 500)
		}
		return
	}

	priceModel, err := pc.repository.UpdateProductModel(numericID, body.CurrentPrice.Value)

	if err != nil {
		respondWithError(err.Error(), w, http.StatusInternalServerError)
	}

	response := ProductModel{
		ID:   numericID,
		Name: name,
		CurrentPrice: Price{
			Value:        priceModel.Price,
			CurrencyCode: priceModel.CurrencyCode,
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(message string, w http.ResponseWriter, status int) {
	errorResponse := ErrorResponse{
		ErrorMessages: []ErrorMessage{{Message: message}},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	marshalledBody, _ := ffjson.Marshal(errorResponse)
	w.Write(marshalledBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(err.Error(), w, 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type ProductModel struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CurrentPrice Price  `json:"current_price"`
}

type Price struct {
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currency_code"`
}

type ErrorResponse struct {
	ErrorMessages []ErrorMessage `json:"errors"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
