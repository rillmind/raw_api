package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rillmind/raw_api/models"
	"github.com/rillmind/raw_api/store"
)

type StoreHandler struct {
	store *store.Store
}

func NewStoreHandler(s *store.Store) *StoreHandler {
	return &StoreHandler{store: s}
}

func (sh *StoreHandler) RegisterRoutes(mux *http.ServeMux) {
	// TODO resto das funcoes do handler
	mux.HandleFunc("POST /products", sh.Create)
	mux.HandleFunc("GET /products", sh.List)
}

func (sh *StoreHandler) Create(w http.ResponseWriter, req *http.Request) {
	var p models.Product

	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Corpo da requisicao invalido")
		return
	}

	if p.Name == "" || p.Price <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "Campo obrigatorio: name. Price nao pode ser 0 ou negativo")
		return
	}

	ctx, cancel := withTimeout(req)
	defer cancel()

	if err := sh.store.Create(ctx, &p); err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao criar produto")
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

func (sh *StoreHandler) List(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := withTimeout(req)
	defer cancel()

	products, err := sh.store.GetAll(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao listar produtos")
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// funcoes helper

func parseID(req *http.Request) (int64, error) {
	return strconv.ParseInt(req.PathValue("id"), 10, 64)
}

func withTimeout(req *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(req.Context(), time.Second*3)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
