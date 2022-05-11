package hashicorp

type HTTPError struct {
	APIName string
	Status  string
	URL     string
}

func (e *HTTPError) Error() string {
	return e.APIName + ": response " + e.Status + " url" + e.URL
}
