package service

import (
	"context"
	"testing"

	"restapirian/internal/domain"
)

type productRepositoryStub struct {
	query   domain.ListQuery
	created *domain.Product
}

func (s *productRepositoryStub) List(_ context.Context, query domain.ListQuery) ([]domain.Product, int64, error) {
	s.query = query
	return []domain.Product{{ID: 1, Name: "Keyboard"}}, 21, nil
}
func (s *productRepositoryStub) FindByID(context.Context, int32) (*domain.Product, error) {
	if s.created == nil {
		return nil, domain.ErrNotFound
	}
	return s.created, nil
}
func (*productRepositoryStub) CodeExists(context.Context, string, int32) (bool, error) {
	return false, nil
}
func (s *productRepositoryStub) Create(_ context.Context, product *domain.Product) error {
	product.ID = 1
	s.created = product
	return nil
}
func (*productRepositoryStub) Update(context.Context, *domain.Product) error { return nil }
func (*productRepositoryStub) Delete(context.Context, int32) error           { return nil }

type categoryRepositoryStub struct{}

func (*categoryRepositoryStub) List(context.Context, domain.ListQuery) ([]domain.Category, int64, error) {
	return nil, 0, nil
}
func (*categoryRepositoryStub) FindByID(_ context.Context, id int32) (*domain.Category, error) {
	return &domain.Category{ID: id, Name: "Accessories"}, nil
}
func (*categoryRepositoryStub) CodeExists(context.Context, string, int32) (bool, error) {
	return false, nil
}
func (*categoryRepositoryStub) Create(context.Context, *domain.Category) error { return nil }
func (*categoryRepositoryStub) Update(context.Context, *domain.Category) error { return nil }
func (*categoryRepositoryStub) Delete(context.Context, int32) error            { return nil }

func TestProductServiceListNormalizesPagination(t *testing.T) {
	products := &productRepositoryStub{}
	service := NewProductService(products, &categoryRepositoryStub{})

	page, err := service.List(context.Background(), domain.ListQuery{Page: -1, PerPage: 500, Order: "DESC"})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if products.query.Page != 1 || products.query.PerPage != 100 || products.query.Order != "desc" {
		t.Fatalf("query was not normalized: %+v", products.query)
	}
	if page.Meta.Total != 21 || page.Meta.TotalPages != 1 {
		t.Fatalf("unexpected pagination metadata: %+v", page.Meta)
	}
}

func TestProductServiceCreateRequiresExistingCategory(t *testing.T) {
	products := &productRepositoryStub{}
	service := NewProductService(products, &categoryRepositoryStub{})

	product, err := service.Create(context.Background(), ProductInput{
		Name: "Keyboard", Code: "key-01", Price: 250000, CategoryID: 2,
	}, 9)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if product == nil || product.CategoryID == nil || *product.CategoryID != 2 || product.CreatedBy != 9 || product.Code != "KEY-01" {
		t.Fatalf("unexpected product: %+v", product)
	}
}
