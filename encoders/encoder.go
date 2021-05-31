package encoders

import (
	"code.cloudfoundry.org/lager"
	"github.com/gleidsonnunes/snickers2/db"
	"github.com/gleidsonnunes/snickers2/types"
)

// EncodeFunc is a function type for the multiple
// possible ways to encode the job
type EncodeFunc func(logger lager.Logger, dbInstance db.Storage, jobID string) error

// GetEncodeFunc returns the encode function
// based on the job.
func GetEncodeFunc(job types.Job) EncodeFunc {
	return FFMPEGEncode
}
