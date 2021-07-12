package b9s

import (
	"net/url"
	"strconv"
	"strings"
)

// Request is a B9s GET request
type Request struct {
	ProjectID string
	APIKey    string
	Table     string

	params   url.Values
	pageSize uint
	page     uint
}

// String builds request URL
func (r *Request) String() string {
	r.setPagination()

	parts := []string{
		host,
		r.ProjectID,
		r.APIKey,
		"data",
		r.Table,
	}

	u, _ := url.Parse(strings.Join(parts, "/"))
	u.RawQuery = r.params.Encode()

	return u.String()
}

// SetWhere param
func (r *Request) SetWhere(where string) {
	r.params.Set("where", where)
}

// SetPage for pagination
func (r *Request) SetPage(page uint) {
	r.page = page
}

// SetPageSize for pagination
func (r *Request) SetPageSize(pageSize uint) {
	r.pageSize = pageSize
}

// setPagination params
func (r *Request) setPagination() {
	if r.pageSize == 0 {
		r.pageSize = 10
	}
	if r.page == 0 {
		r.page = 1
	}

	offset := r.pageSize * (r.page - 1)

	r.params.Set("pageSize", strconv.FormatUint(uint64(r.pageSize), 10))
	r.params.Set("offset", strconv.FormatUint(uint64(offset), 10))
}

// SetRelations to load
func (r *Request) SetRelations(relations ...string) {
	rs := strings.Join(relations, ",")
	r.params.Set("loadRelations", rs)
}

// SetRelationsSizes page size & depth
func (r *Request) SetRelationsSizes(pageSize uint, depth uint) {
	r.params.Set("relationsPageSize", strconv.FormatUint(uint64(pageSize), 10))
	r.params.Set("relationsDepth", strconv.FormatUint(uint64(depth), 10))
}

// SetOrder sorting
func (r *Request) SetOrder(field string, asc bool) {
	s := "`" + field + "` "
	if asc {
		s += " asc"
	} else {
		s += " desc"
	}
	r.params.Set("sortBy", s)
}
