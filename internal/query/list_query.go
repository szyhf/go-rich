package query

type ListQuery struct {
	*Query
}

func NewListQuery(q *Query) *ListQuery {
	return &ListQuery{
		Query: q,
	}
}
