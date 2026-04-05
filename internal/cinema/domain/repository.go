package domain

import (
	"context"

	"github.com/google/uuid"
)

type CinemaRepository interface {
	Create(ctx context.Context, cinema *Cinema) error
	FindByID(ctx context.Context, id uuid.UUID) (*Cinema, error)
	FindAll(ctx context.Context, city string, page, limit int) ([]Cinema, int64, error)
	Update(ctx context.Context, cinema *Cinema) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type ScreenRepository interface {
	Create(ctx context.Context, screen *Screen) error
	FindByID(ctx context.Context, id uuid.UUID) (*Screen, error)
	FindByCinemaID(ctx context.Context, cinemaID uuid.UUID) ([]Screen, error)
	Update(ctx context.Context, screen *Screen) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type SeatRepository interface {
	BulkCreate(ctx context.Context, seats []Seat) error
	FindByScreenID(ctx context.Context, screenID uuid.UUID) ([]Seat, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Seat, error)
	Update(ctx context.Context, seat *Seat) error
}
