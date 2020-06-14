package config

type Server struct {
	Port string `json:"port"`
}

func (s *Server) Provided() bool {
	return s.Port != ""
}