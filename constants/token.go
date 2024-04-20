package constants

import "time"

const TokenLength int = 30
const TokenValidityDuration = 300 * time.Second

//Types

const (
	AccessTokenType int = iota
	OTPTokenType    int = iota
)
