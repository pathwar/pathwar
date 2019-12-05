package pwapi

import "context"

func (svc *service) ToolStatus(context.Context, *GetStatus_Input) (*GetStatus_Output, error) {
	return &GetStatus_Output{
		EverythingIsOK: true,
	}, nil
}
