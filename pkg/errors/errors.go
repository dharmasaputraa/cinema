package errors

import "net/http"

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Message
}

// Constructor helpers
func New(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, HTTPStatus: status}
}

// ── Common errors ────────────────────────────────────────────

var (
	// 400
	ErrBadRequest  = New("BAD_REQUEST", "bad request", http.StatusBadRequest)
	ErrValidation  = New("VALIDATION_ERROR", "validation error", http.StatusBadRequest)
	ErrInvalidUUID = New("INVALID_UUID", "invalid uuid format", http.StatusBadRequest)

	// 401 / 403
	ErrUnauthorized = New("UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
	ErrForbidden    = New("FORBIDDEN", "forbidden", http.StatusForbidden)

	// 404
	ErrNotFound = New("NOT_FOUND", "resource not found", http.StatusNotFound)

	// 409
	ErrConflict = New("CONFLICT", "resource already exists", http.StatusConflict)

	// 422
	ErrUnprocessable = New("UNPROCESSABLE", "unprocessable entity", http.StatusUnprocessableEntity)

	// 500
	ErrInternal = New("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
)

// Domain-specific errors
var (
	// Cinema
	ErrCinemaNotFound = New("CINEMA_NOT_FOUND", "cinema not found", http.StatusNotFound)
	ErrScreenNotFound = New("SCREEN_NOT_FOUND", "screen not found", http.StatusNotFound)
	ErrSeatNotFound   = New("SEAT_NOT_FOUND", "seat not found", http.StatusNotFound)
	ErrSeatInactive   = New("SEAT_INACTIVE", "seat is not active", http.StatusUnprocessableEntity)

	// Movie
	ErrMovieNotFound = New("MOVIE_NOT_FOUND", "movie not found", http.StatusNotFound)

	// Showtime
	ErrShowtimeNotFound = New("SHOWTIME_NOT_FOUND", "showtime not found", http.StatusNotFound)
	ErrShowtimeInactive = New("SHOWTIME_INACTIVE", "showtime is not active", http.StatusUnprocessableEntity)

	// Booking
	ErrBookingNotFound          = New("BOOKING_NOT_FOUND", "booking not found", http.StatusNotFound)
	ErrBookingNotPending        = New("BOOKING_NOT_PENDING", "booking is not in pending state", http.StatusUnprocessableEntity)
	ErrBookingCannotBeCancelled = New("BOOKING_CANNOT_BE_CANCELLED", "booking cannot be cancelled", http.StatusUnprocessableEntity)
	ErrBookingExpired           = New("BOOKING_EXPIRED", "booking has expired", http.StatusUnprocessableEntity)
	ErrSeatAlreadyLocked        = New("SEAT_ALREADY_LOCKED", "seat is already locked by another user", http.StatusConflict)
	ErrSeatAlreadyBooked        = New("SEAT_ALREADY_BOOKED", "seat is already booked", http.StatusConflict)

	// Payment
	ErrPaymentNotFound  = New("PAYMENT_NOT_FOUND", "payment not found", http.StatusNotFound)
	ErrPaymentDuplicate = New("PAYMENT_DUPLICATE", "duplicate payment", http.StatusConflict)

	// Ticket
	ErrTicketNotFound = New("TICKET_NOT_FOUND", "ticket not found", http.StatusNotFound)
	ErrTicketUsed     = New("TICKET_ALREADY_USED", "ticket has already been used", http.StatusUnprocessableEntity)
)

// IsAppError checks if an error is our AppError type
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
