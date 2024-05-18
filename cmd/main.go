package main

import (
	"log"
	"meatsrv/authenticate"
	"meatsrv/inithandlers"
	"meatsrv/meet"
	"meatsrv/meetops"
	"meatsrv/session"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	sessionManager := session.NewSessionManager()
	meetManager := meetops.NewMeetManager()
	authenticateBL := authenticate.NewAuthenticateBL(sessionManager)
	inithandlers.InitAuthenticate(router, authenticateBL)
	meetBL := meet.NewMeetBL(sessionManager, meetManager)
	inithandlers.InitMeet(router, meetBL)
	serviceHostPost := "localhost:8080"
	log.Println("Listening on", serviceHostPost)
	go func() {
		for {
			<-time.After(1 * time.Second)
			log.Printf("Num goroutines: %d\n", runtime.NumGoroutine())
		}
	}()
	http.ListenAndServe(serviceHostPost, router)
}
