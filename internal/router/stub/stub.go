package stub

type StubsMapper interface {
	MapToStubs() []*Stub
}

type Stub struct {
	HttpMethod    string
	Path          string
	Body          string
	Status        int
	CollectParams *CollectParams
}

type CollectParams struct {
	Value string
}
