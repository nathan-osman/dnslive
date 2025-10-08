package server

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

func (s *Server) respond(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.Authoritative = true
	m.SetReply(r)
	func() {
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		// Loop over each question
		for _, q := range r.Question {

			fqdn := q.Name

			// Handle SOA requests to the zone
			if fqdn == s.zone {
				if q.Qtype == dns.TypeSOA {
					m.Answer = append(m.Answer, &dns.SOA{
						Hdr:    dns.RR_Header{Name: fqdn, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 3600},
						Ns:     s.nameserver,
						Mbox:   s.mailbox,
						Serial: 1,
					})
				}
				continue
			}

			// If the name is within the zone, remove the zone
			before, found := strings.CutSuffix("."+fqdn, s.zone)
			if !found {
				continue
			}

			// Now extract the name (the last part of the value)
			var (
				parts = strings.Split(before, ".")
				name  = parts[len(parts)-1]
			)

			// Check to see if the name is valid
			e, ok := s.entries[name]
			if !ok {
				return
			}

			var rrA, rrAAAA dns.RR

			// Create a record for the IPv4 address
			if ipv4 := net.ParseIP(e.Ipv4); ipv4 != nil {
				rrA = &dns.A{
					Hdr: dns.RR_Header{Name: fqdn, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
					A:   ipv4,
				}
			}

			// Create a record for the IPv6 address
			if ipv6 := net.ParseIP(e.Ipv6); ipv6 != nil {
				rrAAAA = &dns.AAAA{
					Hdr:  dns.RR_Header{Name: fqdn, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600},
					AAAA: ipv6,
				}
			}

			// Return the desired response
			switch q.Qtype {
			case dns.TypeA:
				if rrA != nil {
					m.Answer = append(m.Answer, rrA)
				}
			case dns.TypeAAAA:
				if rrAAAA != nil {
					m.Answer = append(m.Answer, rrAAAA)
				}
			}
		}
	}()
	if err := w.WriteMsg(m); err != nil {
		s.logger.Error().Msg(err.Error())
	}
}
