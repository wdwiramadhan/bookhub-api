package usecase

import (
	"context"
	"time"

	"github.com/wdwiramadhan/bookhub-api/domain"
)

// ProductUseCase represent the product use case struct
type ProductUseCase struct {
	productRepo    domain.ProductRepository
	contextTimeout time.Duration
}

// NewProductUsecase will create new an productUsecase object representation of domain.ProductUsecase interface
func NewProductUsecase(p domain.ProductRepository, timeout time.Duration) domain.ProductUseCase {
	return &ProductUseCase{
		productRepo:    p,
		contextTimeout: timeout,
	}
}

func (p *ProductUseCase) Fetch(c context.Context) (res []domain.Product, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	res, err = p.productRepo.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return
}

func (p *ProductUseCase) Store(c context.Context, m *domain.Product) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	err = p.productRepo.Store(ctx, m)
	return
}

func (p *ProductUseCase) GetByID(c context.Context, id int) (res domain.Product, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	res, err = p.productRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (p *ProductUseCase) Update(c context.Context, m *domain.Product, id int) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	err = p.productRepo.Update(ctx, m, id)
	if err != nil {
		return
	}
	return
}

func (p *ProductUseCase) Delete(c context.Context, id int) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	err = p.productRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}
