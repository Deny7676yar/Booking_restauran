package pgstore

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/test"
	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(RestTestSuite))

type RestTestSuite struct {
	test.SuitBase
	db *sql.DB
}

func Test(t *testing.T) {
	gc.TestingT(t)
}

func (s *RestTestSuite) SetUpSuite(c *gc.C) {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		c.Skip("Missing PGDB_DSN envvar; skipping postgresdb-backed test suite")
	}

	r, err := NewRestaurants(dsn)
	c.Assert(err, gc.IsNil)
	s.SetRest(r)
	s.db = r.db
	fmt.Println(r)
}

func (s *RestTestSuite) SetUpTest(c *gc.C) {
	s.flushDB(c)
}

func (s *RestTestSuite) TearDownSuite(c *gc.C) {
	if s.db != nil {
		s.flushDB(c)
		c.Assert(s.db.Close(), gc.IsNil)
	}
}

func (s *RestTestSuite) flushDB(c *gc.C) {
	_, err := s.db.Exec("DELETE FROM links")
	c.Assert(err, gc.IsNil)
}
