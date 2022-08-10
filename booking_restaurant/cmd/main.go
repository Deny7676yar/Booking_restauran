package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/infrastructure/api/handler"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/infrastructure/api/routergin"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/infrastructure/db/pgstore"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/infrastructure/server"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/bookingrepo"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/restaurantrepo"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/tablerepo"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	rst, err := pgstore.NewRestaurants(os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	tst, err := pgstore.NewTableRests(os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	bst, err := pgstore.NewBookings(os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	stable := tablerepo.NewTableRests(tst)
	hst := handler.NewHandlerTables(stable)

	srest := restaurantrepo.NewRestaurant(rst)
	hsr := handler.NewHandlerRests(srest)

	sbooking := bookingrepo.NewBookings(bst, tst)
	hbooking := handler.NewHandlerBookings(sbooking)

	h := routergin.NewRouterGinRest(hsr, hst, hbooking)
	srv := server.NewServer(":"+os.Getenv("PORT"), h)

	srv.Start(srest, stable, sbooking)
	log.WithFields(log.Fields{
		"Start": time.Now(),
	}).Info()

	// прослушивание системных вызовов
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	for {
		select {
		case <-ctx.Done():
			return
		case <-sigCh:
			log.WithFields(log.Fields{
				"SIGINT": <-sigCh,
			}).Info("cencel context")
			srv.Stop()
			cancel() //Если пришёл сигнал SigInt - завершаем контекст
			rst.Close()
			tst.Close()
			bst.Close()
		}
	}
}
