package db

import (
	"fmt"
	"unicode"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
	"github.com/go-pg/urlstruct"
)

const (
	// common statuses
	StatusEnabled = iota + 1
	StatusDisabled
	StatusDeleted
)

var (
	StatusFilter        = Filter{Field: "statusId", Value: []int{StatusEnabled, StatusDisabled}, SearchType: SearchTypeArray}
	StatusEnabledFilter = Filter{Field: "statusId", Value: []int{StatusEnabled}, SearchType: SearchTypeArray}
)

type SortDirection string

const (
	SortAsc            SortDirection = "asc"
	SortAscNullsFirst  SortDirection = "asc nulls first"
	SortAscNullsLast   SortDirection = "asc nulls last"
	SortDesc           SortDirection = "desc"
	SortDescNullsFirst SortDirection = "desc nulls first"
	SortDescNullsLast  SortDirection = "desc nulls last"
)

type SortField struct {
	Column    string
	Direction SortDirection
}

func NewSortField(column string, sortDesc bool) SortField {
	d := SortAsc
	if sortDesc {
		d = SortDesc
	}
	return SortField{Column: column, Direction: d}
}

// OpFunc is a function that applies different options to query.
type OpFunc func(query *orm.Query)

// WithSort is a function that adds uses specific columns to query.
func WithSort(fields ...SortField) OpFunc {
	return func(query *orm.Query) {
		for _, f := range fields {
			query.OrderExpr("? ?", types.Ident(f.Column), types.Safe(f.Direction))
		}
	}
}

// WithColumns is a function that adds user specific columns to query.
func WithColumns(cols ...string) OpFunc {
	return func(query *orm.Query) {
		for _, col := range cols {
			for _, r := range col {
				if unicode.IsLetter(r) && unicode.IsUpper(r) {
					query.Relation(col)
					break
				} else {
					query.Column(col)
					break
				}
			}
		}
	}
}

// WithoutColumns is a function that excludes user specific columns from a query.
func WithoutColumns(cols ...string) OpFunc {
	return func(query *orm.Query) {
		for _, col := range cols {
			for _, r := range col {
				if !unicode.IsLetter(r) || !unicode.IsUpper(r) {
					query.ExcludeColumn(col)
					break
				}
			}
		}
	}
}

// WithRelations is a function that adds user specific relations to query.
func WithRelations(rels ...string) OpFunc {
	return func(query *orm.Query) {
		for _, rel := range rels {
			query.Relation(rel)
		}
	}
}

// EnabledOnly is a function that adds "statusId"=1 filter to query.
func EnabledOnly() OpFunc {
	return func(query *orm.Query) {
		Filter{Field: "statusId", Value: StatusEnabled}.Apply(query)
	}
}

// applyOps applies operations to current orm query.
func applyOps(q *orm.Query, ops ...OpFunc) {
	for _, op := range ops {
		op(q)
	}
}

const (
	defaultMaxLimit = 25
	defaultNoLimit  = 999999
)

var (
	PagerDefault = Pager{PageSize: defaultMaxLimit}
	PagerNoLimit = Pager{PageSize: defaultNoLimit}
	PagerOne     = Pager{PageSize: 1}
	PagerTwo     = Pager{PageSize: 2}
)

type Pager struct {
	Page     int
	PageSize int
}

// NewPager create new Pager. If page and pageSize is zero return PagerDefault
func NewPager(page, pageSize int) Pager {
	if page == 0 && pageSize == 0 {
		return PagerDefault
	}
	return Pager{
		Page:     page,
		PageSize: pageSize,
	}
}

// Pager gets orm.Pages for go-pg
func (p Pager) Pager() *urlstruct.Pager {
	maxLimit := p.PageSize
	if maxLimit > defaultNoLimit {
		maxLimit = defaultNoLimit
	} else if maxLimit == 0 {
		maxLimit = defaultMaxLimit
	}
	pager := &urlstruct.Pager{
		Limit:    p.PageSize,
		MaxLimit: maxLimit,
	}

	pager.SetPage(p.Page)

	return pager
}

// String gets sql string from options
func (p Pager) String() (opts string) {
	pager := p.Pager()
	limit := pager.GetLimit()
	offset := pager.GetOffset()

	if limit != 0 {
		opts = fmt.Sprintf("LIMIT %d ", limit)
	}

	if offset != 0 {
		opts += fmt.Sprintf("OFFSET %d ", offset)
	}

	return
}

// Apply applies options to go-pg orm
func (p Pager) Apply(query *orm.Query) *orm.Query {
	pager := p.Pager()
	limit := pager.GetLimit()
	offset := pager.GetOffset()

	if limit != 0 {
		query = query.Limit(limit)
	}
	if offset != 0 {
		query = query.Offset(offset)
	}
	return query
}
