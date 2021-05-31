package downloaders

import (
	"net/url"
	"os"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/flavioribeiro/gonfig"
	"github.com/secsy/goftp"
	"github.com/gleidsonnunes/db"
)

// FTPDownload downloads the file from FTP. Job Source should be
// in format: ftp://login:password@host/path
func FTPDownload(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, jobID string) error {
	log := logger.Session("ftp-download")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	u, err := url.Parse(job.Source)
	if err != nil {
		return err
	}

	pw, isSet := u.User.Password()
	if !isSet {
		pw = ""
	}

	ftpConfig := goftp.Config{
		User:               u.User.Username(),
		Password:           pw,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, err := goftp.DialConfig(ftpConfig, u.Host+":21")
	if err != nil {
		log.Error("dial-config-failed", err)
		return err
	}

	outputFile, err := os.Create(job.LocalSource)
	if err != nil {
		log.Error("creating-local-source-failed", err)
		return err
	}

	err = client.Retrieve(u.Path, outputFile)
	if err != nil {
		log.Error("retrieving-output-failed", err)
		return err
	}

	return nil
}
