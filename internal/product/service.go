package product

import (
	"log"

	"github.com/ncostamagna/gocourse_product/internal/domain"
)

type (
	Filters struct {
		Name string
	}

	Service interface {
		Create(name string, price float64) (*domain.Product, error)
		Get(id string) (*domain.Product, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Product, error)
		Delete(id string) error
		Update(id string, name *string, price *float64) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(name string, price float64) (*domain.Product, error) {

	product := &domain.Product{
		Name:  name,
		Price: price,
	}

	if err := s.repo.Create(product); err != nil {
		s.log.Println(err)
		return nil, err
	}

	return product, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Product, error) {

	products, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return products, nil
}

func (s service) Get(id string) (*domain.Product, error) {
	product, err := s.repo.Get(id)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return product, nil
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Update(id string, name *string, price *float64) error {

	return s.repo.Update(id, name, price)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
