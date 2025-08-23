package models

type Event struct {
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	RequestContext        struct {
		HTTP struct {
			Path string `json:"path"`
		} `json:"http"`
	} `json:"requestContext"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type ResponseWrapper struct {
	StatusCode int      `json:"statusCode"`
	Body       Response `json:"body"`
}
