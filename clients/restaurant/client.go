package restaurant

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/SebastianCoetzee/blog-order-service-example/models"
	"github.com/pkg/errors"
)

// Client is an interface that describes a RestaurantService client.
type Client interface {
	GetRestaurantsByIDs(ids []int) (models.Restaurants, error)
}

// NewClient creates a new Restaurant client.
func NewClient() *client {
	return &client{}
}

// client is an implementation of a RestaurantService client interface.
type client struct {
	baseURL string
}

// SetBaseURL overrides the default base URL for the restaurants service.
func (c *client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *client) getBaseURL() string {
	if c.baseURL != "" {
		return c.baseURL
	}

	c.baseURL = os.Getenv("RESTAURANT_SERVICE_BASE_URL")
	return c.baseURL
}

// GetRestaurantsByIDs retrieves the Restaurants from the RestaurantService
// using a slice of integer IDs.
func (c *client) GetRestaurantsByIDs(ids []int) (models.Restaurants, error) {
	if len(ids) == 0 {
		return []*models.Restaurant{}, nil
	}

	idStrings := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrings = append(idStrings, strconv.Itoa(id))
	}

	url := fmt.Sprintf(
		"%s/v1/restaurants?id=%s",
		c.getBaseURL(),
		strings.Join(idStrings, ","),
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("error retrieving restaurants from RestaurantService")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsedBody := models.Restaurants{}
	if err = json.Unmarshal(body, &parsedBody); err != nil {
		return nil, err
	}

	return parsedBody, nil
}
