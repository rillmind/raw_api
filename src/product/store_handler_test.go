package product

import (
	"context"
)

type fakeStore struct {
	CreateFn func(ctx context.Context, product *Product) error
	GetByIDFn func(ctx context.Context, id int64) (*Product, error)
	GetAllFn func(ctx context.Context) ([]Product, error)
	UpdateFn func(ctx context.Context, product *Product) error
	DeleteFn func(ctx context.Context, id int64) error
}

func (fs *fakeStore) Create(ctx context.Context, product *Product) error {
	return fs.CreateFn(ctx, product)
}

func (fs *fakeStore) GetByID(ctx context.Context, id int64) (*Product, error) {
	return fs.GetByIDFn(ctx, id)
}

func (fs *fakeStore) GetAll(ctx context.Context) ([]Product, error) {
	return fs.GetAllFn(ctx)
}

func (fs *fakeStore) Update(ctx context.Context, product *Product) error {
	return fs.UpdateFn(ctx, product)
}

func (fs *fakeStore) Delete(ctx context.Context, id int64) error {
	return fs.DeleteFn(ctx, id)
}

