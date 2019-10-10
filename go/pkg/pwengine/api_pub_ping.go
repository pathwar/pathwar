package pwengine

import "context"

func (e *engine) Ping(context.Context, *Void) (*Void, error) {
	return &Void{}, nil
}
