package pagination

import (
	"errors"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

var ErrInvalidPagination = errors.New("invalid_pagination")

type Params struct {
	Limit  int
	Page   int
	Search string
	Offset int
}

type Pagination struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Pages  int    `json:"pages"`
	Search string `json:"search"`
	Total  int    `json:"total"`
}

type Response struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func Parse(c *gin.Context) (Params, error) {
	page := 1
	limit := defaultLimit
	search := c.DefaultQuery("search", "")

	if raw := c.Query("page"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return Params{}, ErrInvalidPagination
		}
		page = parsed
	}

	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return Params{}, ErrInvalidPagination
		}
		limit = parsed
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	return Params{
		Limit:  limit,
		Page:   page,
		Search: search,
		Offset: offset,
	}, nil
}

func NewResponse(data any, total int, params Params) Response {
	pages := int(math.Ceil(float64(total) / float64(params.Limit)))
	if pages == 0 {
		pages = 1
	}

	return Response{
		Data: data,
		Pagination: Pagination{
			Limit:  params.Limit,
			Page:   params.Page,
			Pages:  pages,
			Search: params.Search,
			Total:  total,
		},
	}
}
