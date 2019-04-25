package server

import "context"

func (s *svc) Ping(_ context.Context, _ *Void) (*Void, error) {
	return &Void{}, nil
}
