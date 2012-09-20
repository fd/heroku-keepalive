package model

import (
	"crypto/tls"
	"github.com/miekg/dns"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var DNS_servers = map[string]string{
	// Google
	"ns1.google.com": "8.8.8.8:53",
	"ns2.google.com": "8.8.4.4:53",

	// Telenet
	"ns1.telenet.be": "195.130.131.139:53",
	"ns2.telenet.be": "195.130.130.139:53",
	"ns3.telenet.be": "195.130.131.11:53",
	"ns4.telenet.be": "195.130.130.11:53",
}

var dns_client = new(dns.Client)

func init() {
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func (d *Domain) Test() bool {
	return true
}

func (d *DnsDomain) Test() bool {
	if !(*Domain)(d).Test() {
		return false
	}

	fqdn := d.Name
	if strings.HasPrefix(fqdn, "*.") {
		fqdn = "a" + fqdn[1:]
	}
	if !strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn + "."
	}

	any_ok := false

	d.DNS = make([]*DnsRecords, 0, len(DNS_servers))
	for name, addr := range DNS_servers {
		records := new(DnsRecords)
		records.Server = name
		records.NS = addr
		d.DNS = append(d.DNS, records)

		req := new(dns.Msg)
		req.Id = dns.Id()
		req.RecursionDesired = true
		req.Question = []dns.Question{
			dns.Question{fqdn, dns.TypeA, dns.ClassINET},
		}

		resp, err := dns_client.Exchange(req, addr)
		if err != nil {
			records.Status = 900
			records.Message = err.Error()
			continue
		}

		records.IPs = make([]string, 0, len(resp.Answer))
		for _, rr := range resp.Answer {
			switch a := rr.(type) {
			case *dns.RR_A:
				records.IPs = append(records.IPs, a.A.String())
			}
		}

		if len(records.IPs) > 0 {
			any_ok = true
		} else {
			records.Status = 900
			records.Message = "No records"
		}
	}

	return any_ok
}

func (d *HttpDomain) Test() bool {
	if !(*DnsDomain)(d).Test() {
		return false
	}

	if strings.HasPrefix(d.Name, "*.") {
		d.Status = 900
		d.Message = "Unable to check wildcard domains."
		d.LoadTime = 0
		return false
	}

	start_at := time.Now()

	resp, err := http.Get("http://" + d.Name + "/")
	if err != nil {
		d.Status = 900
		d.Message = err.Error()
		d.LoadTime = time.Now().Sub(start_at) / time.Millisecond
		return false
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		d.Status = 900
		d.Message = err.Error()
		d.LoadTime = time.Now().Sub(start_at) / time.Millisecond
		return false
	}

	d.Status = uint(resp.StatusCode)
	d.Message = resp.Status
	d.LoadTime = time.Now().Sub(start_at) / time.Millisecond
	return true
}
