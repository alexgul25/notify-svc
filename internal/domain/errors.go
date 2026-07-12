package domain

import "errors"

var (
	ErrMsgDoubleSend = errors.New("msg double send")
)
