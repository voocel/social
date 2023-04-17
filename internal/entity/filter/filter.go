package filter

import (
	"strconv"
	"strings"
)

const (
	paginationDefaultPage = 1
	paginationDefaultSize = 30

	queryParamPage          = "page"
	queryParamLimit         = "limit"
	queryParamOffset        = "offset"
	queryParamSort          = "sort"
	queryParamDisablePaging = "disable_paging"
)

type Filter struct {
	Page          int
	Offset        int
	Limit         int
	DisablePaging bool

	Sort   map[string]string
	Search bool
}

func NewFilter(queries map[string]string) *Filter {
	var page, limit, offset int
	page, err := strconv.Atoi(queries[queryParamPage])
	if err != nil {
		page = paginationDefaultPage
	}
	limit, err = strconv.Atoi(queries[queryParamLimit])
	if err != nil {
		limit = paginationDefaultSize
	}

	offset, err = strconv.Atoi(queries[queryParamOffset])
	if err != nil {
		offset = limit * (page - 1) // calculates offset
	}

	disablePaging, _ := strconv.ParseBool(queries[queryParamDisablePaging])

	sortKey := make(map[string]string)
	_, ok := queries[queryParamSort]
	if ok {
		s := queries[queryParamSort]
		//key, order, found := strings.Cut(s, ",")
		eq := strings.Index(s, ",")
		if eq == -1 {
			sortKey[s[:eq]] = "asc"
		} else {
			sortKey[s[:eq]] = strings.ToUpper(s[eq+1:])
		}
	}

	return &Filter{
		Page:          page,
		Offset:        offset,
		Limit:         limit,
		DisablePaging: disablePaging,
		Sort:          sortKey,
	}
}
