package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/langgeng-jbt/langgengpkg/basicdto/trace"
	"github.com/langgeng-jbt/langgengpkg/contextwrap"
	"github.com/langgeng-jbt/langgengpkg/log"
)

type HttpMicro interface {
	Call(ctx context.Context, requestBody map[string]interface{}, header http.Header, path string, method string) (context.Context, []byte, http.Header, error)
	// GenerateHeaderLivvik(body, clientID, clientKey string) http.Header
	// GenerateHeaderGenMicro(device *entity.Device, ipSource, agent string) http.Header
}

type HttpMicroImpl struct {
	httpc   *http.Client
	baseUrl string
}

func New(timeout int, baseUrl string) HttpMicro {
	hclientWConfig := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// client := apmhttp.WrapClient(hclientWConfig)

	custom := &HttpMicroImpl{
		httpc:   hclientWConfig,
		baseUrl: baseUrl,
	}

	return custom
}

func (c *HttpMicroImpl) Call(ctx context.Context, requestBody map[string]interface{}, header http.Header, path string, method string) (context.Context, []byte, http.Header, error) {
	start := time.Now()
	jsonRequest, _ := json.Marshal(requestBody)
	payload := bytes.NewReader(jsonRequest)

	currentTrace := contextwrap.GetTraceFromContext(ctx)

	request, err := http.NewRequest(method, c.baseUrl+path, payload)
	if err != nil {
		return ctx, nil, nil, err
	}

	request.Header = header

	response, err := c.httpc.Do(request.WithContext(ctx))
	if err != nil {
		return ctx, nil, nil, err
	}

	defer response.Body.Close()

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//log.LogWarn(err.Error(), "read esb response")
		return ctx, nil, nil, err
	}

	elapsed := time.Since(start).String()

	tr := &trace.TraceHttp{
		Url:     c.baseUrl + path,
		Request: log.Minify(requestBody),
		Elapsed: elapsed,
	}

	currentTrace = append(currentTrace, tr)

	fmt.Println("current trace ", currentTrace)

	ctx = context.WithValue(ctx, contextwrap.TraceKey, currentTrace)

	// check for valid json response
	var js map[string]interface{}
	err = json.Unmarshal(responseByte, &js)
	if err != nil {
		// log.LogWarn("invalid json", "invalid json")
		// return ctx, nil, nil, err
		fmt.Println("response body not in json format")
	}

	tr.Response = log.Minify(js)
	response.Header.Add("statusCode", response.Status)

	return ctx, responseByte, response.Header, nil
}
