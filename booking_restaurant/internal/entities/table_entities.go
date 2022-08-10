package entities

import "github.com/google/uuid"

// Table - столик в ресторане.
type TableRest struct {
	ID uuid.UUID
	//RestaurantID - id ресторана к которому относится столик
	RestaurantID uuid.UUID
	//CapacityT - вместимость столика
	CapacityT int
}
