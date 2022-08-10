package bookingrepo

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/tablerepo"
	"github.com/google/uuid"
)

type BookingStore interface {
	Create(ctx context.Context, bking entities.Booking, tableIDs ...uuid.UUID) (*uuid.UUID, error)
}

type Bookings struct {
	bookingstore BookingStore
	tablestore   tablerepo.TableRestStore
}

func NewBookings(bookingstore BookingStore, tablestore tablerepo.TableRestStore) *Bookings {
	return &Bookings{
		bookingstore: bookingstore,
		tablestore:   tablestore,
	}
}

func (b *Bookings) Create(ctx context.Context, bking entities.Booking) (*entities.Booking, error) {
	desiredDataTime := strings.Split(bking.DataOfVisit, " ")

	// получаем доступные для брони столики в выбранном ресторане
	tables, err := b.tablestore.GetAavailableTable(ctx, bking.RestaurantID, desiredDataTime[0], desiredDataTime[1])
	if err != nil {
		return nil, err
	}

	peopleNum := bking.NumberOfPeople

	// подсчёт общего количество доступных мест в ресторане
	avalaibleSeatsNum := 0
	for {
		table, ok := <-tables
		if ok == false {
			return nil, fmt.Errorf("read chan %w", err)
		} else {
			avalaibleSeatsNum += table.CapacityT
		}
	}

	// если суммарное количество доступных мест меньше, чем хочет прийти людей
	if avalaibleSeatsNum < bking.NumberOfPeople {
		return nil, fmt.Errorf("the number of seats is less: %w", err)
	}

	//сортируем столики по возрастанию колличества мест
	var tableSlice []entities.TableRest
	sort.SliceStable(tables, func(i, j int) bool {
		table := <-tables
		tableSlice = append(tableSlice, table)
		return tableSlice[i].CapacityT < tableSlice[j].CapacityT
	})

	// столики, которые будут забронированы после создания брони (изначально пусто)
	bookedTables := make([]uuid.UUID, 0)
	bookedCap := 0
	for i := 0; i < len(tables) && bookedCap < peopleNum; i++ {
		bookedTables = append(bookedTables, tableSlice[i].ID)
		bookedCap += tableSlice[i].CapacityT
	}

	dateTime, err := time.Parse("2006.01.02 15:04", bking.DataOfVisit)
	if err != nil {
		return nil, err
	}

	bookingstore := entities.Booking{
		ClientName:       bking.ClientName,
		ClientPhone:      bking.ClientPhone,
		DataBooking:      dateTime,
		BookingTimeFirst: dateTime,
	}

	bking.BookingID = uuid.New()
	id, err := b.bookingstore.Create(ctx, bookingstore, bookedTables...)
	if err != nil {
		return nil, fmt.Errorf("create Booking error: %w", err)
	}
	bking.BookingID = *id

	return &bking, nil
}
