package entities

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	BookingID        uuid.UUID
	RestaurantID     uuid.UUID
	ClientName       string
	ClientPhone      string
	DataBooking      time.Time
	BookingTimeFirst time.Time
	BookingTimeTo    time.Time
	DataOfVisit      string
	NumberOfPeople   int
}
