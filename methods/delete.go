
package methods

import (
	"github.com/Alvarios/gcfs/config/errors"
	"github.com/Alvarios/gcfs/database"
)

func Delete(fileId string) *errors.Error {
	_, err := database.Bucket.DefaultCollection().Remove(fileId, nil)
	if err != nil {
		return &errors.Error{
			Code: 500,
			Message: err.Error(),
		}
	}

	return nil
}
