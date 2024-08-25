package server

import (
	"encoding/json"
	"os"
)

func (s *Server) load() error {
	f, err := os.Open(s.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.entries)
}

func (s *Server) save() error {
	f, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(&s.entries)
}
