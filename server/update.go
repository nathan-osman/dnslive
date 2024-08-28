package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type updateParams struct {
	Name string `json:"name"`
	Ipv4 string `json:"ipv4"`
	Ipv6 string `json:"ipv6"`
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	var v updateParams
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := func() error {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		// Ensure name ends in "."
		if !strings.HasSuffix(v.Name, ".") {
			v.Name += "."
		}

		e := s.entries[v.Name]
		e.LastUpdate = time.Now()
		e.Ipv4 = v.Ipv4
		e.Ipv6 = v.Ipv6
		s.entries[v.Name] = e
		return s.save()
	}(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
