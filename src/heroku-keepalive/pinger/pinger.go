package pinger

import (
  "log"
  "time"
  "strings"
  "net/http"
  "crypto/tls"
  "encoding/json"
  "io/ioutil"

  m "heroku-keepalive/model"
)


type P struct {
  ApiKey       string
  Interval     time.Duration
  Concurrency  uint

  stop         chan bool
  done         chan bool
  ticker       *time.Ticker
}

func init() {
  http.DefaultClient.Transport = &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }
}

func (this *P) Run() {
  if this.Interval == 0 {
    this.Interval = 5 * time.Minute
  }

  if this.Concurrency == 0 {
    this.Concurrency = 100
  }

  this.stop    = make(chan bool)
  this.done    = make(chan bool)
  this.ticker  = time.NewTicker(this.Interval)

  go this.loop()
}

func (this *P) Stop() {
  this.ticker.Stop()
  this.stop <- true
  <- this.done
}

func (this *P) loop() {
  this.tick()

  for {
    select {
    case <- this.ticker.C:
      this.tick()
    case <- this.stop:
      this.done <- true
      return
    }
  }
}


func (this *P) tick() () {
  var pipe <- chan *m.Application
  pipe = this.fetch_applications()
  pipe = this.fetch_domains(pipe)
  pipe = this.load_pages(pipe)
  this.terminator(pipe)
}

func (this *P) fetch_applications() (<- chan *m.Application) {
  pipe := make(chan *m.Application)

  go func(){
    defer handle_err(pipe)
    var err error

    resp, err := http.Get("https://:"+this.ApiKey+"@api.heroku.com/apps")
    if err != nil { panic(err) }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { panic(err) }

    apps := make([]*m.Application, 0, 100)
    err  = json.Unmarshal(body, &apps)
    if err != nil { panic(err) }

    for _, application := range apps {
      pipe <- application
    }
  }()

  return pipe
}

func (this *P) fetch_domains(in <- chan *m.Application) (<- chan *m.Application) {
  return this.run_concurently(in, func(app *m.Application, pipe chan<-*m.Application){
    var err error

    resp, err := http.Get("https://:"+this.ApiKey+"@api.heroku.com/apps/"+app.Name+"/domains")
    if err != nil { panic(err) }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { panic(err) }

    domains := make([]*m.Domain, 0, 500)
    err  = json.Unmarshal(body, &domains)
    if err != nil { panic(err) }

    domains_map := make(map[string]*m.Domain, 10)
    for _, domain := range domains {
      domains_map[domain.Name] = domain
    }

    app.Domains = domains_map

    pipe <- app
  })
}

func (this *P) load_pages(in <- chan *m.Application) (<- chan *m.Application) {
  return this.run_concurently(in, func(app *m.Application, out chan<-*m.Application){
    load_domain(app.InternalDomain)

    for _, domain := range app.Domains {
      time.Sleep(time.Second)
      load_domain(domain)
    }

    out <- app
  })
}

func (this *P) terminator(in <- chan *m.Application) () {
  a_count := 0
  d_count := 0

  apps := make(map[string]*m.Application, 100)

  for app := range in {
    a_count += 1
    d_count += 1
    d_count += len(app.Domains)

    apps[app.Name] = app
  }

  m.Set(apps)

  log.Printf("validated: %d apps, %d domains", a_count, d_count)
}

func load_domain (d * m.Domain) {
  if strings.HasPrefix(d.Name, "*.") {
    d.Status   = 800
    d.Message  = "Unable to check wildcard domains."
    d.LoadTime = 0
    return
  }

  start_at := time.Now()

  resp, err := http.Get("http://"+d.Name+"/")
  if err != nil {
    d.Status   = 800
    d.Message  = err.Error()
    d.LoadTime = time.Now().Sub(start_at) / time.Millisecond
    return
  }

  defer resp.Body.Close()
  _, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    d.Status  = 800
    d.Message = err.Error()
    d.LoadTime = time.Now().Sub(start_at) / time.Millisecond
    return
  }

  d.Status  = uint(resp.StatusCode)
  d.Message = resp.Status
  d.LoadTime = time.Now().Sub(start_at) / time.Millisecond
}

func (this *P) run_concurently(in <- chan *m.Application, handler func(*m.Application, chan<-*m.Application)) (<- chan *m.Application) {
  pipe := make(chan *m.Application)

  wrapper := func(app *m.Application){
    defer handle_err(nil)
    handler(app, pipe)
  }

  d := make(chan bool, this.Concurrency)

  for i := uint(0); i < this.Concurrency; i++ {
    go func(){
      defer handle_err(nil)

      for application := range in {
        wrapper(application)
      }

      d <- true
    }()
  }

  go func(){
    for i := uint(0); i < this.Concurrency; i++ {
      <- d
    }

    close(pipe)
  }()

  return pipe
}

func handle_err(pipe chan<-*m.Application){
  if err := recover(); err != nil {
    log.Printf("[E]: %s", err)
  }
  if pipe != nil {
    close(pipe)
  }
}
