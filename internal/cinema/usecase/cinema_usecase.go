package usecase

import (
	"context"

	"github.com/dharmasaputraa/cinema-api/internal/cinema/domain"
	appErrors "github.com/dharmasaputraa/cinema-api/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CinemaUsecase interface {
	CreateCinema(ctx context.Context, input CreateCinemaInput) (*domain.Cinema, error)
	GetCinema(ctx context.Context, id uuid.UUID) (*domain.Cinema, error)
	ListCinemas(ctx context.Context, city string, page, limit int) ([]domain.Cinema, int64, error)
	UpdateCinema(ctx context.Context, id uuid.UUID, input UpdateCinemaInput) (*domain.Cinema, error)
	DeleteCinema(ctx context.Context, id uuid.UUID) error

	AddScreen(ctx context.Context, cinemaID uuid.UUID, input AddScreenInput) (*domain.Screen, error)
	GetScreens(ctx context.Context, cinemaID uuid.UUID) ([]domain.Screen, error)

	AddSeats(ctx context.Context, screenID uuid.UUID, input AddSeatsInput) ([]domain.Seat, error)
	GetSeats(ctx context.Context, screenID uuid.UUID) ([]domain.Seat, error)
}

// ── Input DTOs (tinggal di usecase layer, bukan domain) ──────

type CreateCinemaInput struct {
	Name    string `validate:"required,min=2,max=255"`
	City    string `validate:"required"`
	Address string `validate:"required"`
	Phone   *string
	Email   *string `validate:"omitempty,email"`
}

type UpdateCinemaInput struct {
	Name    *string
	City    *string
	Address *string
	Phone   *string
	Email   *string
}

type AddScreenInput struct {
	Name       string            `validate:"required"`
	ScreenType domain.ScreenType `validate:"required,oneof=regular vip imax"`
	HasDolby   bool
	HasIMAX    bool
}

type AddSeatsInput struct {
	Rows     []string        `validate:"required,min=1"` // ["A","B","C"]
	PerRow   int             `validate:"required,min=1"`
	SeatType domain.SeatType `validate:"required,oneof=regular premium vip"`
}

// ── Implementation ───────────────────────────────────────────

type cinemaUsecase struct {
	cinemaRepo domain.CinemaRepository
	screenRepo domain.ScreenRepository
	seatRepo   domain.SeatRepository
	db         *gorm.DB
}

func NewCinemaUsecase(
	cinemaRepo domain.CinemaRepository,
	screenRepo domain.ScreenRepository,
	seatRepo domain.SeatRepository,
	db *gorm.DB,
) CinemaUsecase {
	return &cinemaUsecase{cinemaRepo, screenRepo, seatRepo, db}
}

func (uc *cinemaUsecase) CreateCinema(ctx context.Context, input CreateCinemaInput) (*domain.Cinema, error) {
	cinema := &domain.Cinema{
		Name:    input.Name,
		City:    input.City,
		Address: input.Address,
		Phone:   input.Phone,
		Email:   input.Email,
	}
	if err := uc.cinemaRepo.Create(ctx, cinema); err != nil {
		return nil, err
	}
	return cinema, nil
}

func (uc *cinemaUsecase) GetCinema(ctx context.Context, id uuid.UUID) (*domain.Cinema, error) {
	cinema, err := uc.cinemaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrCinemaNotFound
	}
	return cinema, nil
}

func (uc *cinemaUsecase) ListCinemas(ctx context.Context, city string, page, limit int) ([]domain.Cinema, int64, error) {
	return uc.cinemaRepo.FindAll(ctx, city, page, limit)
}

func (uc *cinemaUsecase) UpdateCinema(ctx context.Context, id uuid.UUID, input UpdateCinemaInput) (*domain.Cinema, error) {
	cinema, err := uc.cinemaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrCinemaNotFound
	}
	if input.Name != nil {
		cinema.Name = *input.Name
	}
	if input.City != nil {
		cinema.City = *input.City
	}
	if input.Address != nil {
		cinema.Address = *input.Address
	}
	if input.Phone != nil {
		cinema.Phone = input.Phone
	}
	if input.Email != nil {
		cinema.Email = input.Email
	}

	if err := uc.cinemaRepo.Update(ctx, cinema); err != nil {
		return nil, err
	}
	return cinema, nil
}

func (uc *cinemaUsecase) DeleteCinema(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.cinemaRepo.FindByID(ctx, id); err != nil {
		return appErrors.ErrCinemaNotFound
	}
	return uc.cinemaRepo.SoftDelete(ctx, id)
}

func (uc *cinemaUsecase) AddScreen(ctx context.Context, cinemaID uuid.UUID, input AddScreenInput) (*domain.Screen, error) {
	if _, err := uc.cinemaRepo.FindByID(ctx, cinemaID); err != nil {
		return nil, appErrors.ErrCinemaNotFound
	}
	screen := &domain.Screen{
		CinemaID:   cinemaID,
		Name:       input.Name,
		ScreenType: input.ScreenType,
		HasDolby:   input.HasDolby,
		HasIMAX:    input.HasIMAX,
	}
	if err := uc.screenRepo.Create(ctx, screen); err != nil {
		return nil, err
	}
	return screen, nil
}

func (uc *cinemaUsecase) GetScreens(ctx context.Context, cinemaID uuid.UUID) ([]domain.Screen, error) {
	if _, err := uc.cinemaRepo.FindByID(ctx, cinemaID); err != nil {
		return nil, appErrors.ErrCinemaNotFound
	}
	return uc.screenRepo.FindByCinemaID(ctx, cinemaID)
}

func (uc *cinemaUsecase) AddSeats(ctx context.Context, screenID uuid.UUID, input AddSeatsInput) ([]domain.Seat, error) {
	screen, err := uc.screenRepo.FindByID(ctx, screenID)
	if err != nil {
		return nil, appErrors.ErrScreenNotFound
	}

	var seats []domain.Seat
	for _, row := range input.Rows {
		for num := 1; num <= input.PerRow; num++ {
			seats = append(seats, domain.Seat{
				ScreenID: screenID,
				Row:      row,
				Number:   num,
				SeatType: input.SeatType,
			})
		}
	}

	tx := uc.db.WithContext(ctx).Begin()

	if err := tx.Create(&seats).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	screen.Capacity += len(seats)
	if err := tx.Save(screen).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return seats, nil
}

func (uc *cinemaUsecase) GetSeats(ctx context.Context, screenID uuid.UUID) ([]domain.Seat, error) {
	if _, err := uc.screenRepo.FindByID(ctx, screenID); err != nil {
		return nil, appErrors.ErrScreenNotFound
	}
	return uc.seatRepo.FindByScreenID(ctx, screenID)
}
