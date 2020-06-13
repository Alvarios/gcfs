package config

type Logs struct {
	Default string `json:"default"`
	Error   string `json:"error"`
}

type Server struct {
	Port string `json:"port"`
	Logs Logs   `json:"logs"`
}

func (s *Server) Provided() bool {
	return s.Port != ""
}