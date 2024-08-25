package server

import (
	"net"

	"github.com/miekg/dns"
)

func (s *Server) respond(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	func() {
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		fqdn := r.Question[0].Name

		// Look up the entry
		e, ok := s.entries[fqdn]
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
		switch r.Question[0].Qtype {
		case dns.TypeA:
			if rrA != nil {
				m.Answer = append(m.Answer, rrA)
			}
		case dns.TypeAAAA:
			if rrAAAA != nil {
				m.Answer = append(m.Answer, rrAAAA)
			}
		}
	}()
	if err := w.WriteMsg(m); err != nil {
		s.logger.Error().Msg(err.Error())
	}
}
