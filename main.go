package main

import (
	"fmt"
	"log"
	"net"
	"os"
	s "strings"

	"github.com/miekg/dns"
)

var records map[string]map[string][]string
var typeToString = map[uint16]string{
	dns.TypeA:     "A",
	dns.TypeAAAA:  "AAAA",
	dns.TypeCNAME: "CNAME",
	dns.TypeTXT:   "TXT",
	dns.TypePTR:   "PTR",
}

// var stringToType = map[string]uint16{
// 	"A":     dns.TypeA,
// 	"AAAA":  dns.TypeAAAA,
// 	"CNAME": dns.TypeCNAME,
// 	"TXT":   dns.TypeTXT,
// 	"PTR":   dns.TypePTR,
// }

func loadIps(filename string) map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("reading file %s: %s", filename, err)
	}
	for _, line := range s.Split(string(data), "\n") {
		record := s.Fields(line)
		key, rType, value := record[0], record[2], record[3]
		if _, ok := result[key]; !ok {
			result[key] = make(map[string][]string)
		}
		result[key][rType] = append(result[key][rType], value)
	}
	return result
}

type handler struct{}

func makeRecord(question dns.Question, value string) dns.RR {
	var record dns.RR
	header := dns.RR_Header{
		Name:   question.Name,
		Rrtype: question.Qtype,
		Class:  dns.ClassINET,
		Ttl:    60,
	}
	switch question.Qtype {
	case dns.TypeA:
		record = &dns.A{
			Hdr: header,
			A:   net.ParseIP(value),
		}
	case dns.TypeAAAA:
		record = &dns.AAAA{
			Hdr:  header,
			AAAA: net.ParseIP(value),
		}
	case dns.TypeCNAME:
		record = &dns.CNAME{
			Hdr:    header,
			Target: value,
		}
	case dns.TypeTXT:
		record = &dns.TXT{
			Hdr: header,
			Txt: []string{value},
		}
	case dns.TypePTR:
		record = &dns.PTR{
			Hdr: header,
			Ptr: value,
		}
	}
	return record
}

func (h *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	log.Println(*m)
	msg := new(dns.Msg)
	msg.SetReply(m)
	msg.Authoritative = true

	for _, question := range m.Question {
		keyTypesToValue, ok := records[question.Name]
		if !ok {
			fmt.Println("no such key")
			continue
		}
		values, ok := keyTypesToValue[typeToString[question.Qtype]]
		if !ok {
			fmt.Println("no such type")
			continue
		}
		for _, value := range values {
			record := makeRecord(question, value)
			msg.Answer = append(msg.Answer, record)
		}
	}
	w.WriteMsg(msg)
}

func main() {
	filename := os.Args[1]
	records = loadIps(filename)
	fmt.Println(records)
	handle := new(handler)
	server := &dns.Server{
		Addr:    ":5003",
		Net:     "udp",
		Handler: handle,
	}

	fmt.Println("Starting DNS server on port 5003")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
