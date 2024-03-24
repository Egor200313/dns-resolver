package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

type handler struct{}

func (h *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	log.Println(*m)
	msg := new(dns.Msg)
	msg.SetReply(m)
	msg.Authoritative = true

	ans := &dns.A{
		Hdr: dns.RR_Header{
			Name:   m.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    60},
		A: net.ParseIP("1.1.1.1"),
	}
	log.Println(ans.A.String())

	msg.Answer = append(msg.Answer, ans)
	log.Println("Answer:", msg.Answer)
	w.WriteMsg(msg)
}

func main() {
	handle := new(handler)
	server := &dns.Server{
		Addr:    ":5003",
		Net:     "udp",
		Handler: handle,
	}

	fmt.Println("Starting DNS server on port 53")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
