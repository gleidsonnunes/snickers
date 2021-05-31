package downloaders

import (
	"os"
	"reflect"
	"runtime"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/flavioribeiro/gonfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gleidsonnunes/snickers2/db"
	"github.com/gleidsonnunes/snickers2/types"
)

var _ = Describe("Downloaders", func() {
	var (
		logger     *lagertest.TestLogger
		dbInstance db.Storage
		downloader DownloadFunc
		exampleJob types.Job
		cfg        gonfig.Gonfig
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("http-download")
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
	})

	Context("GetDownloadFunc", func() {
		It("should return S3Download if source has amazonaws", func() {
			jobSource := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/source_here.mp4"
			downloadFunc := GetDownloadFunc(jobSource)
			funcName := runtime.FuncForPC(reflect.ValueOf(downloadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/gleidsonnunes/snickers2/downloaders.S3Download"))
		})

		It("should return FTPDownload if source starts with ftp://", func() {
			jobSource := "ftp://login:password@host/source_here.mp4"
			downloadFunc := GetDownloadFunc(jobSource)
			funcName := runtime.FuncForPC(reflect.ValueOf(downloadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/gleidsonnunes/snickers2/downloaders.FTPDownload"))
		})

		It("should return HTTPDownload if source starts with http://", func() {
			jobSource := "http://source_here.mp4"
			downloadFunc := GetDownloadFunc(jobSource)
			funcName := runtime.FuncForPC(reflect.ValueOf(downloadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/gleidsonnunes/snickers2/downloaders.HTTPDownload"))
		})
	})

	runDownloadersSuite := func() {
		It("should return an error if source couldn't be fetched", func() {
			dbInstance.StoreJob(exampleJob)
			err := downloader(logger, cfg, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(SatisfyAny(
				ContainSubstring("no such host"),
				ContainSubstring("No filename could be determined"),
				ContainSubstring("The AWS Access Key Id you provided does not exist in our records")))
		})
	}

	Context("HTTP Downloader", func() {
		BeforeEach(func() {
			downloader = HTTPDownload
			exampleJob = types.Job{
				ID:          "123",
				Source:      "http://source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
		})

		runDownloadersSuite()
	})

	Context("FTP Downloader", func() {
		BeforeEach(func() {
			downloader = FTPDownload
			exampleJob = types.Job{
				ID:          "123",
				Source:      "ftp://login:password@host/source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
		})

		runDownloadersSuite()
	})

	Context("S3 Downloader", func() {
		BeforeEach(func() {
			downloader = S3Download
			exampleJob = types.Job{
				ID:               "123",
				Source:           "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/source_here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "240p", Container: "mp4"},
				Status:           types.JobCreated,
				Details:          "",
				LocalDestination: "/tmp/output_here.mp4",
			}
		})

		runDownloadersSuite()
	})
})
