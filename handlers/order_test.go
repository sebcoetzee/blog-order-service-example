package handlers_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/SebastianCoetzee/blog-order-service-example/handlers"
	"github.com/SebastianCoetzee/blog-order-service-example/mock_handlers"
	"github.com/SebastianCoetzee/blog-order-service-example/mock_services"
	"github.com/SebastianCoetzee/blog-order-service-example/models"
	"github.com/SebastianCoetzee/blog-order-service-example/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = Describe("FindOrdersForUser", func() {
	var (
		c            handlers.Context
		p            *handlers.Provider
		orderService services.OrderService
		ctrl         *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	JustBeforeEach(func() {
		p = &handlers.Provider{}
		p.SetOrderService(orderService)
	})

	Describe("with an invalid ID", func() {
		BeforeEach(func() {
			mockContext := mock_handlers.NewMockContext(ctrl)
			mockContext.EXPECT().Param(gomock.Eq("id")).Return("invalid_id")
			mockContext.EXPECT().Status(gomock.Eq(400))
			c = mockContext
		})

		It("should return a 400", func() {
			p.FindOrdersForUser(c)
		})
	})

	Describe("with a valid ID", func() {
		Describe("when an error is returned from the OrderService", func() {
			BeforeEach(func() {
				mockContext := mock_handlers.NewMockContext(ctrl)
				mockContext.EXPECT().Param(gomock.Eq("id")).Return("5")
				mockContext.EXPECT().Status(gomock.Eq(500))
				c = mockContext

				mockOrderService := mock_services.NewMockOrderService(ctrl)
				mockOrderService.EXPECT().FindAllOrdersByUserID(gomock.Eq(5)).Return(nil, errors.New("some error"))
				orderService = mockOrderService
			})

			It("should return a 500", func() {
				p.FindOrdersForUser(c)
			})
		})

		Describe("when the OrderService returns an order", func() {
			BeforeEach(func() {
				orders := models.Orders{}
				orders = append(orders, &models.Order{
					ID: 5,
					Restaurant: &models.Restaurant{
						ID:   9,
						Name: "Nando's",
					},
				})

				mockContext := mock_handlers.NewMockContext(ctrl)
				mockContext.EXPECT().Param(gomock.Eq("id")).Return("5")
				mockContext.EXPECT().JSON(gomock.Eq(200), gomock.Eq(orders))
				c = mockContext

				mockOrderService := mock_services.NewMockOrderService(ctrl)
				mockOrderService.EXPECT().FindAllOrdersByUserID(gomock.Eq(5)).Return(orders, error(nil))
				orderService = mockOrderService
			})

			It("should return a 200 with the JSON response", func() {
				p.FindOrdersForUser(c)
			})
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
