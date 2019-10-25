package pwengine

import "context"

func (e *engine) ToolStatus(context.Context, *Void) (*Status, error) {
	return &Status{
		EverythingIsOK: true,
	}, nil
}
