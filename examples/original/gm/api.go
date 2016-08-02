package gm

import (
	"errors"

	"github.com/myteksi/go/commons/util/test/methodman/examples/original/bookingsdk"
)

// A highly simplified implementation for demo purpose

var (
	// ErrInvalidBookingState ...
	ErrInvalidBookingState = errors.New("ErrInvalidBookingState")
)

// BookingCharge ...
// in GM, api to charge a booking, highly simplified
func BookingCharge(bookingCode string) (err error) {
	bookingInfo, err := bookingsdk.GetBookingState(bookingCode)
	if err != nil {
		return err
	}
	if bookingInfo.State != "complete" {
		return ErrInvalidBookingState
	}
	err = doRealCharging(bookingInfo.UserID, bookingInfo.Amount)
	return
}

func doRealCharging(userID int64, amount float64) (err error) {
	// assume here handles detail of money transaction
	return nil
}
