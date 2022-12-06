package product

import (
	"fmt"
	"log"
	"strings"

	"github.com/ncostamagna/gocourse_product/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(product *domain.Product) error
		GetAll(filters Filters, offset, limit int) ([]domain.Product, error)
		Get(id string) (*domain.Product, error)
		Delete(id string) error
		Update(id string, name *string, price *float64) error
		Count(filters Filters) (int, error)
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

//NewRepo is a repositories handler
func NewRepo(l *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: l,
		db:  db,
	}
}

func (r *repo) Create(product *domain.Product) error {

	if err := r.db.Create(product).Error; err != nil {
		r.log.Printf("error: %v", err)
		return err
	}

	r.log.Println("product created with id: ", product.ID)
	return nil
}

func (r *repo) GetAll(filters Filters, offset, limit int) ([]domain.Product, error) {
	var c []domain.Product

	tx := r.db.Model(&c)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&c)

	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}

func (r *repo) Get(id string) (*domain.Product, error) {
	product := domain.Product{ID: id}
	result := r.db.First(&product)

	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *repo) Delete(id string) error {
	product := domain.Product{ID: id}
	result := r.db.Delete(&product)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repo) Update(id string, name *string, price *float64) error {

	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if price != nil {
		values["price"] = *price
	}

	if err := r.db.Model(&domain.Product{}).Where("id = ?", id).Updates(values); err.Error != nil {
		return err.Error
	}

	return nil
}

func (r *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := r.db.Model(domain.Product{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}
