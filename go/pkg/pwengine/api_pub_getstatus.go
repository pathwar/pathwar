package pwengine

import "context"

func (c *client) GetStatus(context.Context, *Void) (*Status, error) {
	return &Status{
		EverythingIsOK: true,
	}, nil
}
