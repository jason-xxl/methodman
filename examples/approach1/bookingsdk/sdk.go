package bookingsdk

import "errors"

// A highly simplified implementation for demo purpose

// BookingInfo ...
type BookingInfo struct {
	Code   string
	State  string
	UserID int64
	Amount float64
}

var (
	// ErrSthWrong ...
	ErrSthWrong = errors.New("ErrSthWrong")
)

type sdkImpl struct{}

// BookingSDK, api to get booking detail, highly simplified
func (*sdkImpl) GetBookingState(bookingCode string) (bookingInfo *BookingInfo, err error) {
	// assumes it does real db query
	// and return
	bookingInfo = &BookingInfo{
		Code:   bookingCode,
		State:  "complete",
		UserID: 100,
		Amount: 10.0,
	}
	return
}
