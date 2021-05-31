package pipeline

import (
	"io"
	"os"
	"reflect"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gleidsonnunes/gleidsonnunes/db"
	"github.com/gleidsonnunes/gleidsonnunes/downloaders"
	"github.com/gleidsonnunes/gleidsonnunes/types"
)

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

var _ = Describe("Pipeline", func() {
	var (
		cfg        gonfig.Gonfig
		dbInstance db.Storage
	)

	BeforeEach(func() {
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
	})

	Context("SetupJob function", func() {
		It("Should set the local source and local destination on Job", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://flv.io/source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}

			dbInstance.StoreJob(exampleJob)
			SetupJob(exampleJob.ID, dbInstance, cfg)
			changedJob, _ := dbInstance.RetrieveJob("123")

			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")

			sourceExpected := swapDir + "123/src/source_here.mp4"
			Expect(changedJob.LocalSource).To(Equal(sourceExpected))

			destinationExpected := swapDir + "123/dst/source_here_240p.mp4"
			Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
		})
	})

	Context("Pipeline", func() {
		It("Should get the HTTPDownload function if source is HTTP", func() {
			jobSource := "http://flv.io/KailuaBeach.mp4"
			downloadFunc := downloaders.GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(downloaders.HTTPDownload).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})

		It("Should get the S3Download function if source is S3", func() {
			jobSource := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
			downloadFunc := downloaders.GetDownloadFunc(jobSource)
			funcPointer := reflect.ValueOf(downloadFunc).Pointer()
			expected := reflect.ValueOf(downloaders.S3Download).Pointer()
			Expect(funcPointer).To(BeIdenticalTo(expected))
		})
	})

	Context("when calling Swap Cleaner", func() {
		It("should remove local source and local destination", func() {
			currentDir, _ := os.Getwd()

			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      "/tmp/123/src/KailuaBeach.mp4",
				LocalDestination: "/tmp/123/dst/KailuaBeach.webm",
			}

			dbInstance.StoreJob(exampleJob)

			os.MkdirAll("/tmp/123/src/", 0777)
			os.MkdirAll("/tmp/123/dst/", 0777)

			cp(exampleJob.LocalSource, currentDir+"/../fixtures/videos/nyt.mp4")
			cp(exampleJob.LocalDestination, currentDir+"/../fixtures/videos/nyt.mp4")

			Expect(exampleJob.LocalSource).To(BeAnExistingFile())
			Expect(exampleJob.LocalDestination).To(BeAnExistingFile())

			CleanSwap(dbInstance, exampleJob.ID)

			Expect(exampleJob.LocalSource).To(Not(BeAnExistingFile()))
			Expect(exampleJob.LocalDestination).To(Not(BeAnExistingFile()))
		})
	})

	Context("StartJob function", func() {
		It("should set error message to Details if errors occur", func() {
			exampleJob := types.Job{
				ID:      "123",
				Source:  "http://source.here.mp4",
				Details: "",
			}

			dbInstance.StoreJob(exampleJob)
			logger := lagertest.NewTestLogger("StartJob")
			StartJob(logger, cfg, dbInstance, exampleJob)

			changedJob, _ := dbInstance.RetrieveJob("123")
			Expect(changedJob.Details).To(ContainSubstring("no such host"))
		})
	})
})
