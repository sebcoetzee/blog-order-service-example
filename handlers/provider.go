package handlers

import (
	"github.com/SebastianCoetzee/blog-order-service-example/services"
)

// Provider is the endpoint provider that holds the dependencies for the
// endpoints.
type Provider struct {
	orderService services.OrderService
}

// SetOrderService sets the OrderService dependency on the Provider.
func (p *Provider) SetOrderService(s services.OrderService) {
	p.orderService = s
}

func (p *Provider) getOrderService() services.OrderService {
	if p.orderService != nil {
		return p.orderService
	}

	p.orderService = services.NewOrderService()
	return p.orderService
}
