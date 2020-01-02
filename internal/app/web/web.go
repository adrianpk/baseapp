package web

import (
	//"encoding/gob"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/app/svc"
)

type (
	Endpoint struct {
		*kbs.WebEndpoint
		Service *svc.Service
	}
)

const (
	// Generic
	CannotProcErrID  = "cannot_proc_err_msg"
	InputValuesErrID = "input_values_err_msg"
	// Fields
	RequiredErrID   = "required_err_msg"
	MinLengthErrID  = "min_length_err_msg"
	MaxLengthErrID  = "max_length_err_msg"
	NotAllowedErrID = "not_allowed_err_msg"
	NotEmailErrID   = "not_email_err_msg"
	ConfMatchErrID  = "conf_match_err_msg"
)

func NewEndpoint(cfg *kbs.Config, log *kbs.Logger, name string) (*Endpoint, error) {
	//registerGobTypes()

	wep, err := kbs.MakeWebEndpoint(cfg, log, pathFxs)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		WebEndpoint: wep,
	}, nil
}

// registerGobTypes
func registerGobTypes() {
	// gob.Register(CustomType1{})
	// gob.Register(CustomType2{})
	// gob.Register(CustomType3{})
}
