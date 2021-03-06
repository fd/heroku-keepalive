package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/fd/heroku-keepalive/api"
	"github.com/fd/heroku-keepalive/pinger"
	"github.com/fd/options"
)

var spec = options.MustParse(`
heroku-keepalive - Keep heroku websites alive.
Usage: heroku-keepalive --api-key=HEROKU_KEY
--
!api-key=  --api-key,HEROKU_KEY   Heroku API key.
port=      --port,PORT   Heroku API key.
--
--
--

`)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	opts := spec.MustInterpret(os.Args, os.Environ())
	if len(opts.Args) != 0 {
		spec.PrintUsageAndExit()
	}

	p := pinger.P{ApiKey: opts.Get("api-key")}

	log.Printf("--- INSERT COIN ---")
	log.Printf("> Let the zombie apocalypse begin!")

	p.Run()
	api.ListenAndServe(fmt.Sprintf(":%s", opts.Get("port")))

	terminate := make(chan os.Signal)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-terminate

	p.Stop()

	log.Printf("> Aarrggg!!!")
	log.Printf("--- GAME OVER ---")
}

/*
func ping_app(app * App, done chan bool)() {
  var err error

  start_at := time.Now()

  resp, err := http.Get(app.WebUrl)
  if err != nil {
    log.Printf("> Found a corps (%s - %s)", app.Name, err)
    done <- true
    return
  }

  defer resp.Body.Close()
  _, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Printf("> Found a corps (%s - %s)", app.Name, err)
    done <- true
    return
  }

  if time.Now().Sub(start_at) > (10 * time.Second) {
    log.Printf("> Yeah! killed another zombie (%s - %s)", app.Name, resp.Status)
  }

  done <- true
  return
}
*/
