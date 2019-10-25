package pwengine

import "context"

func (e *engine) ToolPing(context.Context, *Void) (*Void, error) {
	return &Void{}, nil
}
