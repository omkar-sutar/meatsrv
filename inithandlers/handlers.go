package inithandlers

import (
	"meatsrv/authenticate"
	"meatsrv/meet"

	"github.com/gorilla/mux"
)

func InitAuthenticate(router *mux.Router, bl *authenticate.BL) {
	endpoints := authenticate.AddEndpoints(bl)
	authenticate.AddHandlers(router, endpoints)
}

func InitMeet(router *mux.Router,bl *meet.BL){
	endpoints:=meet.AddEndpoints(bl)
	meet.AddHandlers(router,endpoints)
}
