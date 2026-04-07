package http

import (
	_ "github.com/dharmasaputraa/cinema-api/internal/cinema/domain"
	"github.com/dharmasaputraa/cinema-api/internal/cinema/usecase"
	appErrors "github.com/dharmasaputraa/cinema-api/pkg/errors"
	"github.com/dharmasaputraa/cinema-api/pkg/helper"
	"github.com/dharmasaputraa/cinema-api/pkg/pagination"
	"github.com/dharmasaputraa/cinema-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CinemaHandler struct {
	uc usecase.CinemaUsecase
}

func NewCinemaHandler(uc usecase.CinemaUsecase) *CinemaHandler {
	return &CinemaHandler{uc}
}

// @Description Saving new cinema data to database
// @Tags cinemas
// @Accept json
// @Produce json
// @Param request body usecase.CreateCinemaInput true "Data Cinema"
// @Success 201 {object} domain.Cinema
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /cinemas [post]
func (h *CinemaHandler) Create(c *gin.Context) {
	var input usecase.CreateCinemaInput

	if !helper.BindAndValidate(c, &input) {
		return
	}

	cinema, err := h.uc.CreateCinema(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, cinema)
}

func (h *CinemaHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(appErrors.New("INVALID_ID", "invalid id format", 400))
		return
	}

	cinema, err := h.uc.GetCinema(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, cinema)
}

// List godoc
// @Summary Get Cinema list
// @Description Displays a list of cinemas with pagination and search features by city.
// @Tags cinemas
// @Accept json
// @Produce json
// @Param city query string false "Filter by city"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Amount of data per page" default(10)
// @Success 200 {object} response.SuccessResponse
// @Router /cinemas [get]
func (h *CinemaHandler) List(c *gin.Context) {
	p := pagination.Parse(c)
	city := c.Query("city")

	cinemas, total, err := h.uc.ListCinemas(c.Request.Context(), city, p.Page, p.PerPage)
	if err != nil {
		c.Error(err)
		return
	}

	response.OKWithMeta(c, cinemas, &response.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, p.PerPage),
	})
}

func (h *CinemaHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(appErrors.New("INVALID_ID", "invalid id format", 400))
		return
	}

	var input usecase.UpdateCinemaInput
	if !helper.BindAndValidate(c, &input) {
		return
	}

	cinema, err := h.uc.UpdateCinema(c.Request.Context(), id, input)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, cinema)
}

func (h *CinemaHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(appErrors.New("INVALID_ID", "invalid id format", 400))
		return
	}

	if err := h.uc.DeleteCinema(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}

func (h *CinemaHandler) AddScreen(c *gin.Context) {
	cinemaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(appErrors.New("INVALID_ID", "invalid id format", 400))
		return
	}

	var input usecase.AddScreenInput
	if !helper.BindAndValidate(c, &input) {
		return
	}

	screen, err := h.uc.AddScreen(c.Request.Context(), cinemaID, input)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, screen)
}

func (h *CinemaHandler) GetScreens(c *gin.Context) {
	cinemaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(appErrors.New("INVALID_ID", "invalid id format", 400))
		return
	}

	screens, err := h.uc.GetScreens(c.Request.Context(), cinemaID)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, screens)
}

func (h *CinemaHandler) AddSeats(c *gin.Context) {
	screenID, err := uuid.Parse(c.Param("screen_id"))
	if err != nil {
		c.Error(appErrors.New("INVALID_ID", "invalid id format", 400))
		return
	}

	var input usecase.AddSeatsInput
	if !helper.BindAndValidate(c, &input) {
		return
	}

	seats, err := h.uc.AddSeats(c.Request.Context(), screenID, input)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, seats)
}
