package filter

type UserFilter struct {
	Base     Filter
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
}

func Filters(queries map[string]string) *UserFilter {
	f := NewFilter(queries)
	switch {
	case has(queries, "username"):
		fallthrough
	case has(queries, "mobile"):
		f.Search = true
	}

	return &UserFilter{
		Base:     *f,
		Username: queries["username"],
		Mobile:   queries["mobile"],
		Email:    queries["email"],
	}
}

func has(m map[string]string, field string) bool {
	_, ok := m[field]
	return ok
}
