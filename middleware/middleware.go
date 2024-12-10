package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/langgeng-jbt/langgengpkg/log"

	"github.com/langgeng-jbt/langgengpkg/basicdto/response"

	"github.com/langgeng-jbt/langgengpkg/contextwrap"

	"github.com/google/uuid"
)

type middlewareImpl struct{}

type Middleware interface {
	DumbMiddleware(next http.Handler) http.Handler
	BodyReader(next http.Handler) http.Handler
	GeneratePid(next http.Handler) http.Handler
	Setup(next http.Handler) http.Handler
	Finally(ctx context.Context, w http.ResponseWriter, resp *response.Response)
}

func New() Middleware {
	return &middlewareImpl{}
}

func (h *middlewareImpl) DumbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload interface{}
		var otherInfo interface{}
		reqbody := "ofoeif"
		header := "jewoijgoew"

		payload = reqbody
		otherInfo = header
		log.LogInbound("dumbinbound", &payload, &otherInfo)
	})
}

func (h *middlewareImpl) BodyReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		switch r.Method {
		case http.MethodGet:
			ctx := context.WithValue(r.Context(), contextwrap.ElapsedKey, start)
			next.ServeHTTP(w, r.WithContext(ctx))
		default:
			defer r.Body.Close()
			var i interface{}
			var body []byte
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&i)
			if err != nil {
				//no body or empty body maybe
				return
			} else {
				body, _ = json.Marshal(i)
				// // And now set a new body, which will simulate the same data we read:
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}

			ctx := context.WithValue(r.Context(), contextwrap.BodyKey, body)
			ctx = context.WithValue(ctx, contextwrap.ElapsedKey, start)
			next.ServeHTTP(w, r.WithContext(ctx))

		}
	})
}

func (h *middlewareImpl) GeneratePid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()
		//context.Set(req, "start", start)

		processIDObject := uuid.Must(uuid.NewRandom())
		processID := strings.Replace(fmt.Sprintf("%v", processIDObject), "-", "", -1)

		ctx := context.WithValue(r.Context(), contextwrap.ProcessIDKey, processID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *middlewareImpl) Setup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		var inboundInfo interface{}

		switch r.Method {
		case http.MethodGet:
			trxType, endpath := extractTrxType(r.URL)

			if endpath == "healthz" {
				return
			}

			inboundInfo = r.Header
			log.LogInbound(trxType, &data, &inboundInfo)

			ctx := context.WithValue(r.Context(), contextwrap.TrxTypeKey, trxType)
			ctx = context.WithValue(ctx, contextwrap.IpAddressSourceKey, r.Header.Get("IP-Address"))
			ctx = context.WithValue(ctx, contextwrap.AgentKey, r.Header.Get("User-Agent"))

			// transaction := apm.TransactionFromContext(ctx)
			// transaction.Context.SetCustom("request_data", data)

			next.ServeHTTP(w, r.WithContext(ctx))
		default:
			if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
				_ = json.NewEncoder(w).Encode("Content type not supported")
				return
			}

			body := contextwrap.GetBodyFromContext(r.Context())

			reqbody := map[string]interface{}{}

			_ = json.Unmarshal(body, &reqbody)

			reqbody["request_id"] = contextwrap.GetProcessIDFromContext(r.Context())

			data = reqbody

			trxType, endpath := extractTrxType(r.URL)

			if endpath == "healthz" {
				return
			}

			inboundInfo = r.Header
			log.LogInbound(trxType, &data, &inboundInfo)

			ctx := context.WithValue(r.Context(), contextwrap.TrxTypeKey, trxType)
			ctx = context.WithValue(ctx, contextwrap.IpAddressSourceKey, r.Header.Get("IP-Address"))
			ctx = context.WithValue(ctx, contextwrap.AgentKey, r.Header.Get("User-Agent"))

			// transaction := apm.TransactionFromContext(ctx)
			// transaction.Context.SetCustom("request_data", data)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func extractTrxType(url *url.URL) (string, string) {
	s := strings.Split(url.String(), "/")

	trxType := ""
	for i := len(s) - 1; i >= len(s)-3; i-- {
		v := strings.Title(s[i])
		if i == len(s)-3 && string(v[0]) == "V" && v != "V1" {
			trxType += "-"
		}

		if i != 1 && v != "V1" {
			trxType += v
		}
	}

	//handling underscore
	u := strings.Split(trxType, "_")

	if len(u) == 1 {
		return trxType, s[len(s)-1]
	}

	trxType = ""
	for i := 0; i < len(u); i++ {
		v := strings.Title(u[i])
		trxType += v
	}

	return trxType, s[len(s)-1]
}
