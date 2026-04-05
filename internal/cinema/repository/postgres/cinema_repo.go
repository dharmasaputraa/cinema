package postgres

import (
	"context"

	"github.com/dharmasaputraa/cinema-api/internal/cinema/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cinemaRepo struct{ db *gorm.DB }

func NewCinemaRepository(db *gorm.DB) domain.CinemaRepository {
	return &cinemaRepo{db}
}

func (r *cinemaRepo) Create(ctx context.Context, cinema *domain.Cinema) error {
	return r.db.WithContext(ctx).Create(cinema).Error
}

func (r *cinemaRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Cinema, error) {
	var cinema domain.Cinema
	err := r.db.WithContext(ctx).
		Preload("Screens").
		First(&cinema, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cinema, nil
}

func (r *cinemaRepo) FindAll(ctx context.Context, city string, page, limit int) ([]domain.Cinema, int64, error) {
	var cinemas []domain.Cinema
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Cinema{})

	if city != "" {
		q = q.Where("city ILIKE ?", "%"+city+"%")
	}

	// 🔥 FIX
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Offset((page - 1) * limit).Limit(limit).Find(&cinemas).Error; err != nil {
		return nil, 0, err
	}

	return cinemas, total, nil
}

func (r *cinemaRepo) Update(ctx context.Context, cinema *domain.Cinema) error {
	return r.db.WithContext(ctx).Save(cinema).Error
}

func (r *cinemaRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Cinema{}, "id = ?", id).Error
}

// ── Screen repo ───────────────────────────────────────────────

type screenRepo struct{ db *gorm.DB }

func NewScreenRepository(db *gorm.DB) domain.ScreenRepository {
	return &screenRepo{db}
}

func (r *screenRepo) Create(ctx context.Context, screen *domain.Screen) error {
	return r.db.WithContext(ctx).Create(screen).Error
}

func (r *screenRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Screen, error) {
	var screen domain.Screen
	err := r.db.WithContext(ctx).First(&screen, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &screen, nil
}

func (r *screenRepo) FindByCinemaID(ctx context.Context, cinemaID uuid.UUID) ([]domain.Screen, error) {
	var screens []domain.Screen
	err := r.db.WithContext(ctx).
		Where("cinema_id = ?", cinemaID).
		Find(&screens).Error
	return screens, err
}

func (r *screenRepo) Update(ctx context.Context, screen *domain.Screen) error {
	return r.db.WithContext(ctx).Save(screen).Error
}

func (r *screenRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Screen{}, "id = ?", id).Error
}

// ── Seat repo ─────────────────────────────────────────────────

type seatRepo struct{ db *gorm.DB }

func NewSeatRepository(db *gorm.DB) domain.SeatRepository {
	return &seatRepo{db}
}

func (r *seatRepo) BulkCreate(ctx context.Context, seats []domain.Seat) error {
	return r.db.WithContext(ctx).CreateInBatches(seats, 100).Error
}

func (r *seatRepo) FindByScreenID(ctx context.Context, screenID uuid.UUID) ([]domain.Seat, error) {
	var seats []domain.Seat
	err := r.db.WithContext(ctx).
		Where("screen_id = ? AND is_active = true", screenID).
		Order("row ASC, number ASC").
		Find(&seats).Error
	return seats, err
}

func (r *seatRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Seat, error) {
	var seat domain.Seat
	err := r.db.WithContext(ctx).First(&seat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *seatRepo) Update(ctx context.Context, seat *domain.Seat) error {
	return r.db.WithContext(ctx).Save(seat).Error
}
