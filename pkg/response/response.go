package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/HoangQuan74/goodie-api/pkg/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page      int   `json:"page,omitempty"`
	PerPage   int   `json:"per_page,omitempty"`
	Total     int64 `json:"total,omitempty"`
	TotalPage int   `json:"total_page,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		c.JSON(appErr.Code, Response{
			Success: false,
			Error: &ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorBody{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		},
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
