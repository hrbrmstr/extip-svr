package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
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
			log.Printf("A Query for %s from %s\n", q.Name, s)
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, s))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		case dns.TypeAAAA:
			log.Printf("AAAA Query for %s from %s\n", q.Name, s)
			rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, s))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		case dns.TypeTXT:
			log.Printf("TXT Query for %s from %s\n", q.Name, s)
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

	var args struct {
		Port int    `arg:"-p,--port,env:EXTIP_PORT" help:"bind port" placeholder:"PORT" default:"53"`
		TLD  string `arg:"-t,--tld,env:EXTIP_TLD" help:"TLD to handle" placeholder:"TLD" default:"is."`
		Quiet bool  `arg:"-q,--quiet" help:"Disable log messages"`
	}

	arg.MustParse(&args)

	if args.Quiet {
	  log.SetOutput(ioutil.Discard)
	}

  if !strings.HasSuffix(args.TLD, ".") {
		args.TLD = args.TLD + "."
	}

  // handle the specified domain
	dns.HandleFunc(args.TLD, handleDnsRequest)

	server := &dns.Server{Addr: ":" + strconv.Itoa(args.Port), Net: "udp"}

	log.Printf("Binding to port: %d\n", args.Port)
	log.Printf("Registering TLD: %s\n", args.TLD)

	err := server.ListenAndServe()

	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

}