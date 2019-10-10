package pwengine

import "context"

func (e *engine) GetStatus(context.Context, *Void) (*Status, error) {
	return &Status{
		EverythingIsOK: true,
	}, nil
}
