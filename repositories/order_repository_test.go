package repositories_test

import (
	"testing"
	"time"

	"github.com/SebastianCoetzee/blog-order-service-example/application"
	"github.com/SebastianCoetzee/blog-order-service-example/models"
	"github.com/SebastianCoetzee/blog-order-service-example/repositories"
	"github.com/go-pg/pg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOrderRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Repository Suite")
}

var _ = Describe("OrderRespository", func() {
	var (
		tx        *pg.Tx
		orderRepo repositories.OrderRepository
		orders    models.Orders
		err       error

		userID = 5
	)

	BeforeEach(func() {
		tx, err = application.ResolveDB().Begin()
		Expect(err).To(BeNil())
		orderRepo = repositories.NewOrderRepository(tx)
	})

	Describe("FindAllOrdersByUserID", func() {
		Describe("with no records in the database", func() {
			It("returns an empty slice of orders", func() {
				orders, err = orderRepo.FindAllOrdersByUserID(userID)
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
				err = tx.Insert(order1)
				Expect(err).To(BeNil())

				order2 := &models.Order{
					Total:        2500,
					CurrencyCode: "GBP",
					UserID:       userID,
					RestaurantID: 9,
					PlacedAt:     time.Now().Add(-36 * time.Hour),
				}
				err = tx.Insert(order2)
				Expect(err).To(BeNil())

				order3 := &models.Order{
					Total:        600,
					CurrencyCode: "GBP",
					UserID:       7,
					RestaurantID: 8,
					PlacedAt:     time.Now().Add(-24 * time.Hour),
				}
				err = tx.Insert(order3)
				Expect(err).To(BeNil())
			})

			It("returns only the records belonging to the user, in order from latest palced_at first", func() {
				orders, err = orderRepo.FindAllOrdersByUserID(userID)
				Expect(err).To(BeNil())
				Expect(len(orders)).To(Equal(2))
				Expect(orders[0].RestaurantID).To(Equal(9))
				Expect(orders[1].RestaurantID).To(Equal(8))
			})
		})
	})

	AfterEach(func() {
		err = tx.Rollback()
		Expect(err).To(BeNil())
	})
})
