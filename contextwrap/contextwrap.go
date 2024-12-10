package contextwrap

import (
	"context"

	"github.com/langgeng-jbt/langgengpkg/basicdto/response"

	"github.com/langgeng-jbt/langgengpkg/log/entity"
)

func GetBodyFromContext(ctx context.Context) []byte {
	lr := ctx.Value(BodyKey)
	if l, ok := lr.([]byte); ok {
		return l
	} else {
		return []byte("")
	}
}

func GetProcessIDFromContext(ctx context.Context) string {
	lr := ctx.Value(ProcessIDKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func GetTraceFromContext(ctx context.Context) []interface{} {
	lr := ctx.Value(TraceKey)
	if l, ok := lr.([]interface{}); ok {
		return l
	} else {
		return []interface{}{}
	}
}

func SetTraceFromContext(ctx context.Context, trace []interface{}) context.Context {
	ctx = context.WithValue(ctx, TraceKey, trace)
	return ctx
}

func GetResponseFromContext(ctx context.Context) *response.Response {
	lr := ctx.Value(RespKey)
	if lr == nil {
		return &response.Response{}
	}
	if l, ok := lr.(*response.Response); ok {
		return l
	} else {
		return &response.Response{}
	}
}

func SetResponseFromContext(ctx context.Context, resp *response.Response) context.Context {
	ctx = context.WithValue(ctx, RespKey, resp)
	return ctx
}

func GetLogResponseFromContext(ctx context.Context) *entity.Responselog {
	lr := ctx.Value(LogRespKey)
	if lr == nil {
		return &entity.Responselog{}
	}
	if l, ok := lr.(*entity.Responselog); ok {
		return l
	} else {
		return &entity.Responselog{}
	}
}
