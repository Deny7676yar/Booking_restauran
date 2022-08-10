package entities

import "github.com/google/uuid"

type Restaurant struct {
	ID    uuid.UUID
	Title string
	// AverageWaitingTime - среднее время ожидания заказа в минутах.
	AverageWaitingTime int
	// AverageCheck - средний чек на блюдо в ресторане.
	AverageCheck float64
	// AvailableSeatsNumber - актуальное количество свободных мест.
	AvailableSeats int
}

//UpdateRestData - содержит информацию о ресторане и используется для обновления записи о нём в БД.
type UpdateRestData struct {
	Title              *string
	AverageWaitingTime *int
	AverageCheck       *float64
}
