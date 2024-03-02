package rpc

import "context"

type Request[Req any, Res any] struct {
	ctx     context.Context
	name    string
	payload *Req
}

func NewRequest[Req any, Res any](
	ctx context.Context,
	name string,
	payload *Req,
) *Request[Req, Res] {
	return &Request[Req, Res]{
		ctx:     ctx,
		name:    name,
		payload: payload,
	}
}
