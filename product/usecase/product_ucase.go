package usecase

import (
	"context"
	"time"

	"github.com/wdwiramadhan/bookhub-api/domain"
)

type ProductUseCase struct{
	productRepo domain.ProductRepository
	contextTimeout time.Duration
}

func NewProductUsecase(p domain.ProductRepository,timeout time.Duration) domain.ProductUseCase {
	return &ProductUseCase{
		productRepo:    p,
		contextTimeout: timeout,
	}
}

func (p *ProductUseCase) Fetch(c context.Context)(res []domain.Product, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	res, err = p.productRepo.Fetch(ctx)
	if err != nil {
		return nil,  err
	}
	return
}

func (p *ProductUseCase) Store(c context.Context, m *domain.Product) (err error){
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	err = p.productRepo.Store(ctx, m)
	return
}

func (p *ProductUseCase) GetById(c context.Context, id string)(res domain.Product,err error){
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	res, err = p.productRepo.GetById(ctx, id)
	if err != nil {
		return
	}
	return
}