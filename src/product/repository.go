package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("Produto nao encontrado")
)

type StoreRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (s *StoreRepository) Create(ctx context.Context, p *Product) error {
	query := `
		insert into products (name, price, stock)
		values ($1, $2, $3)
		returning id, created_at
	`

	return s.db.QueryRowContext(ctx, query, p.Name, p.Price, p.Stock).Scan(&p.ID, &p.CreatedAt)
}

func (s *StoreRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	query := `
		select id, name, price, stock, created_at
		from products
		where id = $1
	`

	var p Product

	err := s.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("Buscando produto: %v", err)
	}

	return &p, nil
}

func (s *StoreRepository) GetAll(ctx context.Context) ([]Product, error) {
	query := `
		select id, name, price, stock, created_at
		from products
		order by id
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Listando produtos: %v", err)
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("Lendo linha: %v", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterando resultados: %v", err)
	}

	return products, nil
}

func (s *StoreRepository) Update(ctx context.Context, p *Product) error {
	query := `
		update products
		set name = $1, price = $2, stock = $3
		where id = $4
	`

	result, err := s.db.ExecContext(ctx, query, p.Name, p.Price, p.Stock, p.ID)
	if err != nil {
		return fmt.Errorf("Atualizando produto: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Checando linhas afetadas: %v", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *StoreRepository) Delete(ctx context.Context, id int64) error {
	query := `
		delete
		from products
		where id = $1
	`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Deletando produto: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Checando linhas afetadas: %v", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
