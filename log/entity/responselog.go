package entity

import (
	"net/http"
)

type Responselog struct {
	ResponseHeader http.Header
	ResponseBody   interface{}
	ResponseCode   string
	Trace          string
	Timestamp      string
	Elapsed        string
}
