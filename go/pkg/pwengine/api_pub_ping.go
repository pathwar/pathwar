package pwengine

import (
	context "context"
)

func (c *client) Ping(context.Context, *Void) (*Void, error) {
	return &Void{}, nil
}
