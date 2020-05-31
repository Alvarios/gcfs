
package methods

import (
	"gcfs"
	"gcfs/config"
)

func Delete(fileId string) error {
	_, err := gcfs.Cluster.Bucket(config.Main.Database.BucketName).DefaultCollection().Remove(fileId, nil)
	return err
}
