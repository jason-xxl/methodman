package bookingsdk

var (
	defaultImpl = &sdkImpl{}
)

// GetImpl ...
func GetImpl() (o SDK) {
	o = defaultImpl
	return
}

// SDK ...
type SDK interface {
	GetBookingState(bookingCode string) (bookingInfo *BookingInfo, err error)
}

// After setting up SDK, then use
//
//     go get github.com/vektra/mockery/.../;
//     mockery -name=SDK
//
// to generate the mock type in file sdk_mock.go
//
