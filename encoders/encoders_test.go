package encoders

import (
	"reflect"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gleidsonnunes/snickers2/types"
)

var _ = Describe("Encoders", func() {
	Context("GetEncodeFunc", func() {
		It("should return FFMPEGEncode if source is not m3u8", func() {
			job := types.Job{
				ID:          "123",
				Source:      "ftp://login:password@host/source_here.mov",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
			}
			encodeFunc := GetEncodeFunc(job)
			funcName := runtime.FuncForPC(reflect.ValueOf(encodeFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/gleidsonnunes/snickers2/encoders.FFMPEGEncode"))
		})
	})
})
