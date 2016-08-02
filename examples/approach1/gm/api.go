package gm

import (
	"errors"

	"github.com/myteksi/go/commons/util/test/methodman/examples/approach1/bookingsdk"
)

// A highly simplified implementation for demo purpose

var (
	// ErrInvalidBookingState ...
	ErrInvalidBookingState = errors.New("ErrInvalidBookingState")
)

// BookingCharge ...
//
// in GM, api to charge a booking, highly simplified
//
// So in real world you will change your original way calling the API to this way,
//
//     err := BookingCharge(bookingSDK.GetImpl(), bookingCode)
//
func BookingCharge(bookingSDK bookingsdk.SDK, bookingCode string) (err error) {
	bookingInfo, err := bookingSDK.GetBookingState(bookingCode)
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
