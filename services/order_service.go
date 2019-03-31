package services

import (
	"github.com/SebastianCoetzee/blog-order-service-example/application"
	"github.com/SebastianCoetzee/blog-order-service-example/clients/restaurant"
	"github.com/SebastianCoetzee/blog-order-service-example/models"
	"github.com/SebastianCoetzee/blog-order-service-example/repositories"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
)

// OrderService represents the business-logic layer for Orders in the system.
type OrderService interface {
	FindAllOrdersByUserID(userID int) (models.Orders, error)
}

// NewOrderService creates an order service.
func NewOrderService() *orderService {
	return &orderService{}
}

type orderService struct {
	db               orm.DB
	restaurantClient restaurant.Client
	orderRepository  repositories.OrderRepository
}

func (s *orderService) SetOrderRepository(r repositories.OrderRepository) {
	s.orderRepository = r
}

func (s *orderService) getOrderRepository() repositories.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = repositories.NewOrderRepository(application.ResolveDB())
	return s.orderRepository
}

func (s *orderService) SetRestaurantClient(c restaurant.Client) {
	s.restaurantClient = c
}

func (s *orderService) getRestaurantClient() restaurant.Client {
	if s.restaurantClient != nil {
		return s.restaurantClient
	}

	s.restaurantClient = restaurant.NewClient()
	return s.restaurantClient
}

func (s *orderService) FindAllOrdersByUserID(userID int) (models.Orders, error) {
	orders, err := s.getOrderRepository().FindAllOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return orders, nil
	}

	restaurantIDs := make([]int, 0, len(orders))
	for _, order := range orders {
		restaurantIDs = append(restaurantIDs, order.RestaurantID)
	}

	restaurants, err := s.getRestaurantClient().GetRestaurantsByIDs(restaurantIDs)
	if err != nil {
		return nil, err
	}

	restaurantsByID := make(map[int]*models.Restaurant)
	for _, restaurant := range restaurants {
		restaurantsByID[restaurant.ID] = restaurant
	}

	for _, order := range orders {
		restaurant, ok := restaurantsByID[order.RestaurantID]
		if !ok {
			return nil, errors.Errorf("restaurant with ID %d not found", order.RestaurantID)
		}

		order.Restaurant = restaurant
	}

	return orders, nil
}
