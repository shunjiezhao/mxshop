package handler

import (
	Potesting "server/shared/postgres/testing"
	"testing"
)

func TestInventory(t *testing.T) {

}
func TestMain(m *testing.M) {
	Potesting.RunWithMongoInDocker(m)
}
