package domain

import (
	"time"

	appErrors "github.com/dharmasaputraa/cinema-api/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cinema struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	City      string    `gorm:"not null"`
	Address   string    `gorm:"not null"`
	Phone     *string
	Email     *string
	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Screens []Screen `gorm:"foreignKey:CinemaID"`
}

func (c *Cinema) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}

func (c *Cinema) Deactivate() error {
	if !c.IsActive {
		return appErrors.New("CINEMA_ALREADY_INACTIVE", "cinema is already inactive", 422)
	}
	c.IsActive = false
	return nil
}

// ── Screen ───────────────────────────────────────────────────

type ScreenType string

const (
	ScreenTypeRegular ScreenType = "regular"
	ScreenTypeVIP     ScreenType = "vip"
	ScreenTypeIMAX    ScreenType = "imax"
)

type Screen struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	CinemaID   uuid.UUID  `gorm:"type:uuid;not null"`
	Name       string     `gorm:"not null"`
	Capacity   int        `gorm:"not null;default:0"`
	ScreenType ScreenType `gorm:"not null;default:'regular'"`
	HasDolby   bool       `gorm:"default:false"`
	HasIMAX    bool       `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Cinema Cinema `gorm:"foreignKey:CinemaID"`
	Seats  []Seat `gorm:"foreignKey:ScreenID"`
}

func (s *Screen) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}

// ── Seat ─────────────────────────────────────────────────────

type SeatType string

const (
	SeatTypeRegular SeatType = "regular"
	SeatTypePremium SeatType = "premium"
	SeatTypeVIP     SeatType = "vip"
)

type Seat struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ScreenID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_seat"`

	Row    string `gorm:"not null;uniqueIndex:idx_seat"`
	Number int    `gorm:"not null;uniqueIndex:idx_seat"`

	SeatType SeatType `gorm:"not null;default:'regular'"`
	IsActive bool     `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Screen Screen `gorm:"foreignKey:ScreenID"`
}

func (s *Seat) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}

func (s *Seat) Deactivate() error {
	if !s.IsActive {
		return appErrors.ErrSeatInactive
	}
	s.IsActive = false
	return nil
}
