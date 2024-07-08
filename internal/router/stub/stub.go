package stub

const (
	CollectTypeJsonPath   = "JsonPath"
	CollectTypeQueryParam = "QueryParam"
	CollectTypePathParam  = "PathParam"
)

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
	Type  string
	Value string
}
