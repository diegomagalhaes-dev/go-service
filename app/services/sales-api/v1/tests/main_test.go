package tests

import (
	"fmt"
	"testing"

	"github.com/diegomagalhaes-dev/go-service/business/data/dbtest"
	"github.com/diegomagalhaes-dev/go-service/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer dbtest.StopDB(c)

	m.Run()
}
