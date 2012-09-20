package model

import (
	"encoding/json"
	"testing"
)

func TestParseHttpDomain(t *testing.T) {
	const doc1 = `
    {
      "created_at": "2012/08/06 07:34:28 -0700",
      "updated_at": "2012/08/06 07:34:28 -0700",
      "default": null,
      "domain": "b-token.de",
      "id": 5719146,
      "app_id": 6491370,
      "base_domain": "b-token.de"
    }
  `

	var d HttpDomain
	err := json.Unmarshal([]byte(doc1), &d)
	if err != nil {
		t.Error(err)
	}

	if d.Name != "b-token.de" {
		t.Errorf("expected d.Name == b-token.de but was %+v", d.Name)
	}
}

func TestParseDnsDomain(t *testing.T) {
	const doc1 = `
    {
      "created_at": "2012/08/06 07:34:28 -0700",
      "updated_at": "2012/08/06 07:34:28 -0700",
      "default": null,
      "domain": "b-token.de",
      "id": 5719146,
      "app_id": 6491370,
      "base_domain": "b-token.de"
    }
  `

	var d DnsDomain
	err := json.Unmarshal([]byte(doc1), &d)
	if err != nil {
		t.Error(err)
	}

	if d.Name != "b-token.de" {
		t.Errorf("expected d.Name == b-token.de but was %+v", d.Name)
	}
}
