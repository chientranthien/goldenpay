package common

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/chientranthien/goldenpay/internal/proto"
)

const (
	defaultLimit = 11
	maxLimit     = 501
)

type Cond struct {
	IsOr           bool
	CompositeConds []CompositeCond
}

func (c Cond) ToWhere() string {
	join := " AND "
	if c.IsOr {
		join = " OR "
	}

	whereConds := []string{}
	for _, cond := range c.CompositeConds {
		whereConds = append(whereConds, cond.ToWhere())
	}

	return strings.Join(whereConds, join)
}

type CompositeCond struct {
	IsOr        bool
	SingleConds []SingleCond
}

func (c CompositeCond) ToWhere() string {
	join := " AND "
	if c.IsOr {
		join = " OR "
	}

	whereConds := []string{}
	for _, cond := range c.SingleConds {
		whereConds = append(whereConds, cond.ToWhere())
	}

	return strings.Join(whereConds, join)
}

type SingleCond struct {
	IsOr bool
	Col  string
	GTE  any
	LTE  any
	EQ   any
}

func (c SingleCond) ToWhere() string {
	join := " AND "
	if c.IsOr {
		join = " OR "
	}
	whereConds := []string{}
	if c.GTE != 0 {
		whereConds = append(whereConds, fmt.Sprintf("%s >= %v", c.Col, c.GTE))
	}
	if c.LTE != 0 {
		whereConds = append(whereConds, fmt.Sprintf("%s <= %v", c.Col, c.LTE))
	}
	if c.EQ != 0 {
		whereConds = append(whereConds, fmt.Sprintf("%s = %v", c.Col, c.EQ))
	}
	s := strings.Join(whereConds, join)

	return fmt.Sprintf("(%s)", s)
}

func OrderDesc(col string) string {
	return col + " DESC"
}

func DefaultLimit(l uint32) *int {
	if l <= 0 {
		return Int(defaultLimit)
	}

	if l > 500 {
		return Int(maxLimit)
	}

	return Int(int(l + 1))
}

func NexPagination(prev *proto.Pagination, getVal func(interface{}) int64, elements interface{}) *proto.Pagination {
	hasMore := false
	val := int64(0)

	elementValues := reflect.ValueOf(elements)
	if n := elementValues.Len(); n >= int(prev.Limit) {
		hasMore = true
		e := elementValues.Index(n-1).Interface()
		val = getVal(e)
	}

	return &proto.Pagination{
		Val:     val,
		HasMore: hasMore,
		Limit:   prev.Limit - 1,
	}
}

func EnsurePagination(p *proto.Pagination) *proto.Pagination {
	if p == nil {
		return &proto.Pagination{
			Limit: uint32(*DefaultLimit(0)),
		}
	}

	n := &proto.Pagination{
		Val:     p.Val,
		HasMore: p.HasMore,
		Limit:   uint32(*DefaultLimit(p.Limit)),
	}
	return n
}
