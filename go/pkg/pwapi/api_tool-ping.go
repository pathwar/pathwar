package pwapi

import "context"

func (svc *service) ToolPing(context.Context, *Void) (*Void, error) {
	return &Void{}, nil
}
