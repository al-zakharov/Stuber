package router

type Stub struct {
	HttpMethod string
	Path       string
	Body       string
	Status     int
}
