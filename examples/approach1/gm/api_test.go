package gm

import (
	"flag"
	"os"
	"testing"

	"github.com/myteksi/go/commons/util/test/methodman/examples/approach1/bookingsdk"
	"github.com/myteksi/go/commons/util/test/methodman/examples/approach1/bookingsdk/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()

	// setup and teardown

	os.Exit(m.Run())
}

func TestBookingCharge(t *testing.T) {

	mockBookingSDK := new(mocks.SDK)
	mockBookingSDK.On("GetBookingState", "and-12345678").Return(nil, bookingsdk.ErrSthWrong)

	err := BookingCharge(mockBookingSDK, "and-12345678")
	assert.Equal(t, err, bookingsdk.ErrSthWrong)
}
