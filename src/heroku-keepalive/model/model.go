package model

import (
  "time"
)

var apps     map[string]*Application
var chan_get chan(chan map[string]*Application)
var chan_set chan map[string]*Application

func init () {
  chan_get = make(chan(chan map[string]*Application), 50)
  chan_set = make(chan map[string]*Application, 1)

  go manage_storage()
}

func manage_storage () {
  for {
    select {
    case resp := <- chan_get:
      resp <- apps
    case a := <- chan_set:
      for _, app := range a {
        app.DomainsArray = make([]*Domain, 0, len(app.Domains))
        for _, domain := range app.Domains {
          app.DomainsArray = append(app.DomainsArray, domain)
        }
      }
      apps = a
    }
  }
}


type Application struct {
  Name           string             `json:"name"`
  InternalDomain *Domain            `json:"domain_name,omitempty"`
  Domains        map[string]*Domain `json:"-"`
  DomainsArray   []*Domain          `json:"domains,omitempty"`
}

type Domain struct {
  Name     string           `json:"domain"`
  Status   uint             `json:"status"`
  Message  string           `json:"message"`
  LoadTime time.Duration    `json:"load_time"`
}

func Get()(map[string]*Application) {
  resp := make(chan map[string]*Application, 1)
  chan_get <- resp
  return <- resp
}

func Set(a map[string]*Application) {
  chan_set <- a
}
