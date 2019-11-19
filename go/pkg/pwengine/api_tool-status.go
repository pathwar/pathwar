package pwengine

import "context"

func (e *engine) ToolStatus(context.Context, *GetStatus_Input) (*GetStatus_Output, error) {
	return &GetStatus_Output{
		EverythingIsOK: true,
	}, nil
}
