package router

type StubsMapper interface {
	MapToStubs() []*Stub
}

type Stub struct {
	HttpMethod   string
	Path         string
	Body         string
	Status       int
	RequestIdKey string
}
