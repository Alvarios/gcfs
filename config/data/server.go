package data

type Logs struct {
	Default string `json:"default"`
	Error   string `json:"error"`
}

type Server struct {
	Port      string `json:"port" default:"8080"`
	Logs      Logs   `json:"logs"`
}

func (s *Server) Provided() bool {
	return s.Port != ""
}