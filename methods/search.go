package methods

import (
	"gcfs"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/nefts-go"
	nefts_config "github.com/Alvarios/nefts-go/config"
)

func Search(start, end int64, params nefts_config.Options) (*nefts_config.QueryResults, *nefts_config.Error) {
	params.Config.Cluster = gcfs.Cluster
	params.Config.Bucket = config.Main.Database.BucketName

	return nefts.Thread(start, end, params)
}