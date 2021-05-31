package encoders

import (
	"code.cloudfoundry.org/lager"
	"github.com/gleidsonnunes/gleidsonnunes/db"
	"github.com/gleidsonnunes/gleidsonnunes/types"
)

// EncodeFunc is a function type for the multiple
// possible ways to encode the job
type EncodeFunc func(logger lager.Logger, dbInstance db.Storage, jobID string) error

// GetEncodeFunc returns the encode function
// based on the job.
func GetEncodeFunc(job types.Job) EncodeFunc {
	if job.Preset.Container == "m3u8" {
		return HLSEncode
	}
	return FFMPEGEncode
}
