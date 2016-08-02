package gm

import (
	"flag"
	"os"
	"testing"

	"github.com/myteksi/go/commons/util/test/methodman"
	"github.com/myteksi/go/commons/util/test/methodman/examples/approach3/bookingsdk"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()

	methodman.EnableMock(&bookingsdk.GetBookingState, "bookingsdk.GetBookingState")

	// setup and teardown

	os.Exit(m.Run())
}

func TestBookingCharge(t *testing.T) {

	defer methodman.RestoreMock()

	methodman.Expect(&bookingsdk.GetBookingState, nil, bookingsdk.ErrSthWrong)

	err := BookingCharge("and-12345678")
	assert.Equal(t, err, bookingsdk.ErrSthWrong)
}
