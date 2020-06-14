
package methods

import (
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
)

func Delete(fileId string) *responses.Error {
	_, err := database.Bucket.DefaultCollection().Remove(fileId, nil)
	if err != nil {
		return &responses.Error{
			Code: 500,
			Message: err.Error(),
		}
	}

	return nil
}
