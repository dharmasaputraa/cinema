package helper

import (
	"github.com/dharmasaputraa/cinema-api/pkg/response"
	"github.com/dharmasaputraa/cinema-api/pkg/validator"
	"github.com/gin-gonic/gin"
)

var validate = validator.New()

// BindAndValidate:
// - bind JSON
// - validate struct
// - auto return error response kalau gagal
func BindAndValidate(c *gin.Context, req interface{}) bool {
	// bind request
	if err := c.ShouldBindJSON(req); err != nil {
		response.ValidationError(c, map[string][]string{
			"_": {"invalid request body"},
		})
		return false
	}

	// validate struct
	if errs := validate.Validate(req); errs != nil {
		response.ValidationError(c, errs)
		return false
	}

	return true
}
