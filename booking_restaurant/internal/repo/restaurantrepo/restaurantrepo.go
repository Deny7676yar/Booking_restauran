package restaurantrepo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/google/uuid"
)

//Бизнес логика ресторанов
type RestaurantStore interface {
	//Create создает ресторан принимает контекс и сущность ресторана которая приходит из внешнего слоя
	Create(ctx context.Context, res entities.Restaurant) (*uuid.UUID, error)
	//Read производит чтение по id
	Read(ctx context.Context, uid uuid.UUID) (*entities.Restaurant, error)
	//Update производит обновление по входящей сущности Restaurant
	Update(ctx context.Context, res entities.Restaurant) error
	//Delete производит удаление по входящему id
	Delete(ctx context.Context, uid uuid.UUID) error
	//SearchRelevantRest формирует подходящие рестораны и отправляет их в канал
	SearchRelevantRest(ctx context.Context, desiredData, desiredTime, numbp string) (chan entities.Restaurant, error)
}

type Restaurants struct {
	restore RestaurantStore
}

func NewRestaurant(restore RestaurantStore) *Restaurants {
	return &Restaurants{
		restore: restore,
	}
}

func (r *Restaurants) Create(ctx context.Context, res entities.Restaurant) (*entities.Restaurant, error) {
	res.ID = uuid.New()
	id, err := r.restore.Create(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("create Restaurants error: %w", err)
	}
	res.ID = *id
	return &res, nil
}

func (r *Restaurants) Read(ctx context.Context, uid uuid.UUID) (*entities.Restaurant, error) {
	res, err := r.restore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("read restaurant error: %w", err)
	}
	return res, nil
}

func (r *Restaurants) Update(ctx context.Context, rest entities.Restaurant) (*entities.Restaurant, error) {
	res, err := r.restore.Read(ctx, rest.ID)
	if err != nil {
		return nil, fmt.Errorf("search Restaurants error: %w", err)
	}
	return res, r.restore.Update(ctx, rest)
}

func (r *Restaurants) Delete(ctx context.Context, uid uuid.UUID) (*entities.Restaurant, error) {
	res, err := r.restore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("search Restaurants error: %w", err)
	}

	return res, r.restore.Delete(ctx, uid)
}

func (r *Restaurants) SearchRelevantRest(ctx context.Context, desiredData, desiredTime, numbp string) (chan entities.Restaurant, error) {
	pn, err := strconv.Atoi(numbp)
	if err != nil {
		return nil, fmt.Errorf("invalid convert error: %w", err)
	}
	if pn < 1 {
		return nil, fmt.Errorf("number people must be greater than zero error: %w", numbp)
	}

	dTime, err := time.Parse("15:04", desiredTime)
	if err != nil {
		return nil, fmt.Errorf("invalid format error: %w", err)
	}

	if time.Now().After(dTime) {
		return nil, fmt.Errorf("time is out of date: %w", err)
	}
	timeH := dTime.Hour()
	timeM := dTime.Minute()
	if timeH < 9 || timeH >= 22 || timeH == 21 && timeM > 0 {
		return nil, fmt.Errorf("invalid input time error: %w", err)
	}

	chin, err := r.restore.SearchRelevantRest(ctx, desiredData, desiredTime, numbp)
	if err != nil {
		return nil, err
	}

	chout := make(chan entities.Restaurant, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case rest, ok := <-chin:
				if !ok {
					return
				}
				chout <- rest
			}
		}
	}()

	return chout, nil
}
