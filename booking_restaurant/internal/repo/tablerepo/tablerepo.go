package tablerepo

import (
	"context"
	"fmt"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/google/uuid"
)

//Бизнес логика работы со столами
type TableRestStore interface {
	//Create создает столик принимает контекс и сущность стол которая приходит из внешнего слоя
	Create(ctx context.Context, tableRest entities.TableRest, rest_id uuid.UUID) (*uuid.UUID, error)
	//Read производит чтение по id
	Read(ctx context.Context, uid uuid.UUID) (*entities.TableRest, error)
	//Update производит обновление по входящему столику
	UpdateTable(ctx context.Context, ur entities.TableRest) error
	//Delete производит удаление по входящему id
	Delete(ctx context.Context, uid uuid.UUID) error
	//SearchRelevantRest формирует подходящие столики и отправляет их в канал
	//Принимает id ресторана, desiredDate в формате "2006.01.02" и desiredTime "15:04".
	GetAavailableTable(ctx context.Context, rid uuid.UUID, desiredData, desiredTime string) (chan entities.TableRest, error)
}

type TableRests struct {
	tblstore TableRestStore
}

func NewTableRests(tblstore TableRestStore) *TableRests {
	return &TableRests{
		tblstore: tblstore,
	}
}

func (tr *TableRests) Create(ctx context.Context, tableRest entities.TableRest, rest_id uuid.UUID) (*entities.TableRest, error) {
	tableRest.ID = uuid.New()
	id, err := tr.tblstore.Create(ctx, tableRest, rest_id)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}
	tableRest.ID = *id
	return &tableRest, nil
}

func (tr *TableRests) Read(ctx context.Context, uid uuid.UUID) (*entities.TableRest, error) {
	tblr, err := tr.tblstore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("read table error: %w", err)
	}
	return tblr, nil
}

func (tr *TableRests) UpdateTable(ctx context.Context, tab entities.TableRest) (*entities.TableRest, error) {
	table, err := tr.tblstore.Read(ctx, tab.ID)
	if err != nil {
		return nil, fmt.Errorf("search Table error: %w", err)
	}
	return table, tr.tblstore.UpdateTable(ctx, tab)
}

func (tr *TableRests) Delete(ctx context.Context, uid uuid.UUID) (*entities.TableRest, error) {
	tblr, err := tr.tblstore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("search Restaurants error: %w", err)
	}

	return tblr, tr.tblstore.Delete(ctx, uid)
}

func (tr *TableRests) GetAavailableTable(ctx context.Context, rid uuid.UUID, desiredData, desiredTime string) (chan entities.TableRest, error) {
	chintbl, err := tr.tblstore.GetAavailableTable(ctx, rid, desiredData, desiredTime)
	if err != nil {
		return nil, err
	}
	choutbl := make(chan entities.TableRest, 100)
	go func() {
		defer close(choutbl)
		for {
			select {
			case <-ctx.Done():
				return
			case table, ok := <-chintbl:
				if !ok {
					return
				}
				choutbl <- table
			}
		}
	}()
	return choutbl, nil
}
