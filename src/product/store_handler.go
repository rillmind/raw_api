package product

import (
	"context"
	"errors"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Store interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int64) error
}

type StoreHandler struct {
	store StoreRepository
}

func NewStoreHandler(s StoreRepository) *StoreHandler {
	return &StoreHandler{store: s}
}

func (sh *StoreHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", sh.Create)
	mux.HandleFunc("GET /products", sh.List)
	mux.HandleFunc("GET /products/{id}", sh.GetByID)
	mux.HandleFunc("PUT /products/{id}", sh.Update)
	mux.HandleFunc("DELETE /products/{id}", sh.Delete)
}

func (sh *StoreHandler) Create(w http.ResponseWriter, req *http.Request) {
	var p Product

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

func (sh *StoreHandler) GetByID(wr http.ResponseWriter, req *http.Request) {
	id, err := parseID(req)
	if err != nil {
		writeError(wr, http.StatusBadRequest, "ID inválido!")
		return
	}

	ctx, cancel := withTimeout(req)
	defer cancel()

	product, err := sh.store.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		writeError(wr, http.StatusNotFound, "Produto não encontrado!")
		return
	}
	if err != nil {
		writeError(wr, http.StatusInternalServerError, "Erro ao buscar produto!")
		return
	}

	writeJSON(wr, http.StatusOK, product)
}

func (sh *StoreHandler) Update(wr http.ResponseWriter, req *http.Request) {
	id, err := parseID(req)
	if err != nil {
		writeError(wr, http.StatusBadRequest, "ID inválido!")
		return
	}

	var product Product
	if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
		writeError(wr, http.StatusBadRequest, "Corpo da requisição inválido!")
		return
	}
	product.ID = id

	ctx, cancel := withTimeout(req)
	defer cancel()

	err = sh.store.Update(ctx, &product)
	if errors.Is(err, ErrNotFound) {
		writeError(wr, http.StatusNotFound, "Produto não encontrado!")
		return
	}
	if err != nil {
		writeError(wr, http.StatusInternalServerError, "Erro ao atualizar produto!")
		return
	}

	writeJSON(wr, http.StatusOK, product)
}

func (sh *StoreHandler) Delete(wr http.ResponseWriter, req *http.Request) {
	id, err := parseID(req)
	if err != nil {
		writeError(wr, http.StatusBadRequest, "ID inválido!")
		return
	}

	ctx, cancel := withTimeout(req)
	defer cancel()

	err = sh.store.Delete(ctx, id)
	if errors.Is(err, ErrNotFound) {
		writeError(wr, http.StatusNotFound, "Produto não encontrado!")
		return
	}
	if err != nil {
		writeError(wr, http.StatusInternalServerError, "Erro ao deletar produto!")
		return
	}

	wr.WriteHeader(http.StatusNoContent)
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
