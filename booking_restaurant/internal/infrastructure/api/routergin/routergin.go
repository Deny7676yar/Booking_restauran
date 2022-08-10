package routergin

import (
	"fmt"
	"net/http"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/infrastructure/api/auth"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/infrastructure/api/handler"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/google/uuid"
)

type RouterGinRest struct {
	*gin.Engine
	hr *handler.HandlerRests
	ht *handler.HandlerTables
	hb *handler.HandlerBookings
}

func NewRouterGinRest(hr *handler.HandlerRests, ht *handler.HandlerTables, hb *handler.HandlerBookings) *RouterGinRest {
	r := gin.Default()
	rgr := &RouterGinRest{
		hr: hr,
		ht: ht,
		hb: hb,
	}

	r.Use(auth.GinAuthMW)

	r.GET("/", rgr.CreateBooking)
	r.POST("/createrest", rgr.CreateRestaurant)
	r.GET("/readrest/:id", rgr.ReadRestaurant)
	r.PATCH("/updaterest/:id", rgr.UpdateRestaurant)
	r.DELETE("/deleterest/:id", rgr.DeleteRestaurant)
	r.GET("/searchrelevantrest/:data/:time/:nump", rgr.SearchRelevantRest)
	r.POST("/tables/create", rgr.CreateTable)
	r.GET("/tables/read/:id", rgr.ReadTable)
	r.PATCH("/tables/update/:id", rgr.UpdateTable)
	r.DELETE("/tables/delete/:id", rgr.DeleteTable)
	r.GET("/tables/getavailabletable/:rid/:", rgr.GetAvailableTable)

	rgr.Engine = r
	return rgr
}

type Restaurant handler.Restaurant
type Table handler.TableRest
type Booking handler.Booking

func (r *RouterGinRest) CreateBooking(c *gin.Context) {
	booking := Booking{}
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bookingcreate, err := r.hb.Create(c.Request.Context(), handler.Booking(booking))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookingcreate)
}

func (r *RouterGinRest) CreateRestaurant(c *gin.Context) {
	rest := Restaurant{}
	if err := c.ShouldBindJSON(&rest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restcreate, err := r.hr.CreateRest(c.Request.Context(), handler.Restaurant(rest))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restcreate)
}

func (r *RouterGinRest) CreateTable(c *gin.Context) {
	rest := c.Value("restaurant").(*Restaurant)
	table := Table{}
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tablecreate, err := r.ht.CreateTable(c.Request.Context(), handler.TableRest(table), rest.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tablecreate)
}

func (r *RouterGinRest) UpdateRestaurant(c *gin.Context) {
	rest := Restaurant{}
	if err := c.ShouldBindJSON(&rest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restupdate, err := r.hr.UpdateRest(c.Request.Context(), handler.Restaurant(rest))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, restupdate)
}

func (r *RouterGinRest) UpdateTable(c *gin.Context) {
	table := Table{}
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableupdate, err := r.ht.UpdateTable(c.Request.Context(), handler.TableRest(table))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tableupdate)
}

func (r *RouterGinRest) ReadRestaurant(c *gin.Context) {
	sid := c.Param("id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restread, err := r.hr.ReadRest(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restread)
}

func (r *RouterGinRest) ReadTable(c *gin.Context) {
	sid := c.Param("id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableread, err := r.ht.ReadTable(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tableread)
}

func (r *RouterGinRest) DeleteRestaurant(c *gin.Context) {
	sid := c.Param("id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restdel, err := r.hr.DeleteRest(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restdel)
}

func (r *RouterGinRest) DeleteTable(c *gin.Context) {
	sid := c.Param("id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tabledel, err := r.ht.DeleteTable(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tabledel)
}

func (r *RouterGinRest) SearchRelevantRest(c *gin.Context) {
	desiredData := c.Param("data")
	desiredTime := c.Param("time")
	numbp := c.Param("number_people")
	w := c.Writer
	fmt.Fprintln(w, "[")
	comma := false
	err := r.hr.SearchRelevantRest(c.Request.Context(), desiredData, desiredTime, numbp, func(resth handler.Restaurant) error {
		if comma {
			fmt.Fprintln(w, ",")
		} else {
			comma = true
		}
		(render.JSON{Data: resth}).Render(w)
		w.Flush()
		return nil
	})
	if err != nil {
		if comma {
			fmt.Fprintln(w, ",")
		}
		(render.JSON{Data: err}).Render(w)
	}
	fmt.Fprintln(w, "]")
}

func (r *RouterGinRest) GetAvailableTable(c *gin.Context) {
	desiredData := c.Param("data")
	desiredTime := c.Param("time")
	rid := c.Param("restaurant_id")

	uid, err := uuid.Parse(rid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w := c.Writer
	fmt.Fprintln(w, "[")
	comma := false
	err = r.ht.GetAavailableTable(c.Request.Context(), uid, desiredData, desiredTime, func(tableh handler.TableRest) error {
		if comma {
			fmt.Fprintln(w, ",")
		} else {
			comma = true
		}
		(render.JSON{Data: tableh}).Render(w)
		w.Flush()
		return nil
	})
	if err != nil {
		if comma {
			fmt.Fprintln(w, ",")
		}
		(render.JSON{Data: err}).Render(w)
	}
	fmt.Fprintln(w, "]")
}
