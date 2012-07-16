package main


import (
  "os"
  "syscall"
  "os/signal"
  "time"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "haraway/common/options"
)


type App struct {
  Name   string `json:"name"`
  WebUrl string `json:"web_url"`
}


var api_key string


const desc = `
heroku-keepalive - Keep heroku websites alive.
Usage: heroku-keepalive --api-key=HEROKU_KEY
--
!api-key=  --api-key,HEROKU_KEY   Heroku API key.
--
--
--

`


func main() {
  spec, err := options.New(desc)
  if err != nil { panic(err) }

  opts, err := spec.Parse(os.Args, os.Environ())

  if err != nil {
    spec.PrintUsageWithError(err)
  }

  if len(opts.Args) != 0 {
    spec.PrintUsageAndExit()
  }

  api_key = opts.Get("api-key")

  log.Printf("--- INSERT COIN ---")
  log.Printf("> Let the zombie apocalypse begin!")
  ping_loop()
  log.Printf("> Aarrggg!!!")
  log.Printf("--- GAME OVER ---")
}


func ping_loop()(){
  terminate := make(chan os.Signal)
  signal.Notify(terminate, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

  ticker := time.Tick(15 * time.Minute)

  ping()

  for {
    select {
    case <- ticker:
      ping()
    case <- terminate:
      return
    }
  }
}

func ping()() {
  log.Printf("> Patrolling the neighborhood for zombies...")

  apps, err := load_apps()
  if err != nil {
    log.Printf("> Ouch, I can't see! (%s)", err)
    return
  }

  done := make(chan bool, len(apps))

  for _, app := range apps {
    go ping_app(app, done)
  }

  for i := 0; i < len(apps); i++ {
    <- done
  }
}


func load_apps()(apps []*App, err error) {
  resp, err := http.Get("https://:"+api_key+"@api.heroku.com/apps")
  if err != nil { return }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil { return }

  apps = make([]*App, 0, 500)
  err = json.Unmarshal(body, &apps)
  if err != nil { return }

  return
}

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
