package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/seancfoley/ipaddress-go/ipaddr"
)

// we have to strip the port from the incoming address
var portRegEx = regexp.MustCompile(`:[0-9]+$`)

// we need to strip the brackets from the IPv6 address
var	bracketRegEx = regexp.MustCompile(`[\[\]]`)

func parseQuery(m *dns.Msg, s string) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, s))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		case dns.TypeAAAA:
			log.Printf("Query for %s\n", q.Name)
			rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, s))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		case dns.TypeTXT:
			log.Printf("Query for %s\n", q.Name)
			rr, err := dns.NewRR(fmt.Sprintf("%s TXT %s", q.Name, s))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	// get the requestor's IP address:port
  addressOfRequester := w.RemoteAddr()
	
	// get rid of the port
	justTheAddress := portRegEx.ReplaceAllString(addressOfRequester.String(), "")

	// get rid of any brackets
	justTheAddress = bracketRegEx.ReplaceAllString(justTheAddress, "")

	// if IPv6, expand it just in case it came in abbreviated
	if strings.Contains(justTheAddress, ":") {
	  justTheAddress = ipaddr.NewIPAddressString(justTheAddress).GetAddress().ToFullString()
	}

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, justTheAddress)
	}

	w.WriteMsg(m)

}

func main() {

	// handle the Iceland domain
	dns.HandleFunc("is.", handleDnsRequest)

  // TODO make this configuratble
	port := 53

	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}

	log.Printf("Starting at %d\n", port)

	err := server.ListenAndServe()

	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

}