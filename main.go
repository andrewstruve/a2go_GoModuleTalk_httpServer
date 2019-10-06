package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/olebedev/config"
	logger "github.com/sirupsen/logrus"
)

func main() {
	var environment string
	var logLevelString string
	var loglevel logger.Level
	var err error
	var file []byte
	environment = os.Getenv("ENVIRONMENT")
	fmt.Println("environment=", environment)
	// Verify that the environment variable is set to an expected value
	if environment != "development" &&
		environment != "testing" &&
		environment != "production" {
		fmt.Println("ENVIRONMENT Variable not set to a valid Environment")
		panic(err)
	}
	file, err = ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)

	// log levels defined in logrus.go
	if environment == "development" {
		logLevelString, err = cfg.String("development.loglevel")
	} else if environment == "production" {
		logLevelString, err = cfg.String("production.loglevel")
	} else if environment == "testing" {
		logLevelString, err = cfg.String("testing.loglevel")
	}

	loglevel, _ = logger.ParseLevel(logLevelString)

	logger.SetLevel(logger.Level(loglevel))
	logger.Info("Starting WebServer")
	http.Handle("/", loggingMiddleware(http.HandlerFunc(handler)))
	http.ListenAndServe(":8000", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	logger.Trace("func Entrance handler")
	var randomNum int
	if r.RequestURI == "/" {
		fmt.Fprintf(w, "Root")
	} else if r.RequestURI == "/rng" {
		fmt.Fprintf(w, "rng")
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		randomNum = r1.Intn(1000)
		logger.Debug("randomNumber ", randomNum)
		w.Write([]byte(strconv.Itoa(randomNum)))
	} else {
		fmt.Fprintf(w, "Request Not Supported")
	}

}
func loggingMiddleware(next http.Handler) http.Handler {
	logger.Trace("func Entrance loggingMiddleWare")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("uri: ", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
