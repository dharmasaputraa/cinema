package response

import (
	"net/http"
	"sort"

	"github.com/dharmasaputraa/cinema-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
	Total      int64 `json:"total,omitempty"`
}

type successResponse struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta,omitempty"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Type    string `json:"type,omitempty"`
}

type errorResponse struct {
	Errors []errorDetail `json:"errors"`
}

// ── Success ────────────────────────────────────────────────

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, successResponse{Data: data})
}

func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, successResponse{Data: data, Meta: meta})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, successResponse{Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ── Error ─────────────────────────────────────────────────

func Error(c *gin.Context, err error) {
	if appErr, ok := errors.IsAppError(err); ok {
		c.JSON(appErr.HTTPStatus, errorResponse{
			Errors: []errorDetail{
				{
					Code:    appErr.Code,
					Message: appErr.Message,
					Type:    "application_error",
				},
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, errorResponse{
		Errors: []errorDetail{
			{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
				Type:    "internal_error",
			},
		},
	})
}

// multi error + sorted + 422
func ValidationError(c *gin.Context, fields map[string][]string) {
	details := make([]errorDetail, 0)

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, field := range keys {
		for _, msg := range fields[field] {
			details = append(details, errorDetail{
				Code:    "VALIDATION_ERROR",
				Message: msg,
				Field:   field,
				Type:    "validation_error",
			})
		}
	}

	c.JSON(http.StatusUnprocessableEntity, errorResponse{
		Errors: details,
	})
}
