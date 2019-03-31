package repositories

import (
	"github.com/SebastianCoetzee/blog-order-service-example/application"
	"github.com/SebastianCoetzee/blog-order-service-example/models"
	"github.com/go-pg/pg/orm"
)

// OrderRepository is the interface that an order repository should conform to.
type OrderRepository interface {
	FindAllOrdersByUserID(userID int) (models.Orders, error)
}

// NewOrderRepository returns a new implementation of an order repository.
func NewOrderRepository(db orm.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

// orderRepository is an implementation of an OrderRepository.
type orderRepository struct {
	db orm.DB
}

func (r *orderRepository) SetDB(db orm.DB) {
	r.db = db
}

func (r *orderRepository) getDB() orm.DB {
	if r.db != nil {
		return r.db
	}

	r.db = application.ResolveDB()
	return r.db
}

func (r *orderRepository) FindAllOrdersByUserID(userID int) (models.Orders, error) {
	orders := models.Orders{}
	err := r.getDB().Model(&orders).Where("user_id = ?", userID).Order("placed_at DESC").Select()
	return orders, err
}
