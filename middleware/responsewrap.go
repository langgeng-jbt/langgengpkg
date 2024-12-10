package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/langgeng-jbt/langgengpkg/basicdto/response"

	"github.com/langgeng-jbt/langgengpkg/contextwrap"

	"github.com/langgeng-jbt/langgengpkg/log"

	"golang.org/x/exp/rand"
)

func (h *middlewareImpl) Finally(ctx context.Context, w http.ResponseWriter, resp *response.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	_ = json.NewEncoder(w).Encode(resp)

	logResp := contextwrap.GetLogResponseFromContext(ctx)

	body := contextwrap.GetBodyFromContext(ctx)

	traceCtx := contextwrap.GetTraceFromContext(ctx)

	trace := ""
	traceRaw, err := json.Marshal(traceCtx)
	if err != nil {
		fmt.Println("fail to marshal tracectx : ", err)
	}

	trace = string(traceRaw)

	reqbody := map[string]interface{}{}

	_ = json.Unmarshal(body, &reqbody)

	//init the loc
	// loc, _ := time.LoadLocation("Asia/Jakarta")
	//set timezone,
	// dtnow := time.Now().In(loc)

	// logResp.ThirdParty = ""
	logResp.Trace = trace
	logResp.ResponseHeader = w.Header()
	logResp.ResponseBody = resp
	logResp.ResponseCode = resp.Code

	// accountDebet := contextwrap.GetAccountDebetFromContext(ctx)

	log.LogRespBasic(logResp)
}

func randomRefnum() string {
	s := ""
	for i := 0; i < 12; i++ {
		s += (string)(rand.Intn(10) + 48)
	}

	return s
}
