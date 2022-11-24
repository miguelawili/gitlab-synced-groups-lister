package services

type Query struct {
	Key   string
	Value string
}

func (q *Query) GetKey() string {
	return q.Key
}

func (q *Query) GetValue() string {
	return q.Value
}
