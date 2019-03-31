package services_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/SebastianCoetzee/blog-order-service-example/clients/mock_restaurant"
	"github.com/SebastianCoetzee/blog-order-service-example/clients/restaurant"
	"github.com/SebastianCoetzee/blog-order-service-example/mock_repositories"
	"github.com/SebastianCoetzee/blog-order-service-example/models"
	"github.com/SebastianCoetzee/blog-order-service-example/repositories"
	"github.com/SebastianCoetzee/blog-order-service-example/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOrderService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Service Suite")
}

var _ = Describe("OrderService", func() {
	var (
		restaurantClient restaurant.Client
		orderRepo        repositories.OrderRepository
		orderService     services.OrderService
		orders           models.Orders
		ctrl             *gomock.Controller
		err              error

		userID = 5
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	JustBeforeEach(func() {
		orderServiceImpl := services.NewOrderService()
		orderServiceImpl.SetOrderRepository(orderRepo)
		orderServiceImpl.SetRestaurantClient(restaurantClient)
		orderService = orderServiceImpl
	})

	Describe("FindAllOrdersByUserID", func() {
		Describe("with no records in the database", func() {
			BeforeEach(func() {
				orderRepoMock := mock_repositories.NewMockOrderRepository(ctrl)
				orderRepoMock.EXPECT().FindAllOrdersByUserID(gomock.Eq(userID))
				orderRepo = orderRepoMock
			})

			It("returns an empty slice of orders", func() {
				orders, err = orderService.FindAllOrdersByUserID(userID)
				Expect(err).To(BeNil())
				Expect(len(orders)).To(Equal(0))
			})
		})

		Describe("when a few records exist", func() {
			BeforeEach(func() {
				order1 := &models.Order{
					Total:        1000,
					CurrencyCode: "GBP",
					UserID:       userID,
					RestaurantID: 8,
					PlacedAt:     time.Now().Add(-72 * time.Hour),
				}
				order2 := &models.Order{
					Total:        2500,
					CurrencyCode: "GBP",
					UserID:       userID,
					RestaurantID: 9,
					PlacedAt:     time.Now().Add(-36 * time.Hour),
				}

				orderRepoMock := mock_repositories.NewMockOrderRepository(ctrl)
				orderRepoMock.EXPECT().
					FindAllOrdersByUserID(gomock.Eq(userID)).
					Return(models.Orders{order2, order1}, error(nil))
				orderRepo = orderRepoMock
			})

			Describe("when not all Restaurants can be found", func() {
				BeforeEach(func() {
					restaurantClientMock := mock_restaurant.NewMockClient(ctrl)
					restaurantClientMock.EXPECT().
						GetRestaurantsByIDs(gomock.Eq([]int{9, 8})).
						Return(models.Restaurants{}, error(nil))
					restaurantClient = restaurantClientMock
				})

				It("returns only the records belonging to the user, in order from latest palced_at first", func() {
					orders, err = orderService.FindAllOrdersByUserID(userID)
					Expect(err).To(MatchError("restaurant with ID 9 not found"))
				})
			})

			Describe("when all Restaurants are found", func() {
				BeforeEach(func() {
					restaurant1 := &models.Restaurant{
						ID:   9,
						Name: "Nando's",
					}

					restaurant2 := &models.Restaurant{
						ID:   8,
						Name: "KFC",
					}

					restaurantClientMock := mock_restaurant.NewMockClient(ctrl)
					restaurantClientMock.EXPECT().
						GetRestaurantsByIDs(gomock.Eq([]int{9, 8})).
						Return(models.Restaurants{restaurant1, restaurant2}, error(nil))
					restaurantClient = restaurantClientMock
				})

				It("returns only the records belonging to the user, in order from latest palced_at first", func() {
					orders, err = orderService.FindAllOrdersByUserID(userID)
					Expect(err).To(BeNil())
					Expect(len(orders)).To(Equal(2))
					Expect(orders[0].Restaurant.Name).To(Equal("Nando's"))
					Expect(orders[0].Total).To(Equal(2500))
					Expect(orders[1].Restaurant.Name).To(Equal("KFC"))
					Expect(orders[1].Total).To(Equal(1000))
				})
			})
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
