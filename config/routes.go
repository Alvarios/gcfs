package config

// import nefts_config "github.com/Alvarios/nefts-go/config"

type Routes struct {
	Ping         string `json:"ping"`
	PingDatabase string `json:"ping_database"`
	Insert       string `json:"insert"`
	Delete       string `json:"delete"`
	Get          string `json:"get"`
	Update       string `json:"update"`
	Search       string `json:"search"`
	// SearchOptions nefts_config.Options `json:"search_options"`
}

func (r *Routes) Provided() bool {
	return r.Insert != "" || r.Delete != "" || r.Get != "" || r.PingDatabase != "" || r.Ping != "" || r.Search != ""
}
