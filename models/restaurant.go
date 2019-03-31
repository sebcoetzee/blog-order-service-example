package models

// Restaurant is the model representation of a restaurant. Restaurants are
// stored in the RestaurantService.
type Restaurant struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

// Restaurants is a slice of Restaurant pointers.
type Restaurants []*Restaurant
