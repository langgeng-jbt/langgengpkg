package trace

type TraceHttp struct {
	Url      string      `json:"url"`
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
	Elapsed  string      `json:"elapsed"`
}

type TraceDatabase struct {
	Host     string      `json:"host"`
	Query    string      `json:"query"`
	Response interface{} `json:"result"`
	Elapsed  string      `json:"elapsed"`
}

type TraceMinio struct {
	Host       string `json:"host"`
	ObjectName string `json:"object_name"`
	BucketName string `json:"bucket_name"`
	Elapsed    string `json:"elapsed"`
}
