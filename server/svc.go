package server

import "context"

type svc struct {
}

func (s *svc) Ping(_ context.Context, _ *Void) (*Void, error) {
	return &Void{}, nil
}
