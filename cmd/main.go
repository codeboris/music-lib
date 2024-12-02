package main

import (
	"net/http"

	"github.com/codeboris/music-lib/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.Infoln("App logrus view.")

	srv := new(server.Server)
	srv.Run("8000", http.HandlerFunc(startString))
}

func startString(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Добро пожаловать на апи!"))
}
