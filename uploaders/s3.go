package uploaders

import (
	"os"

	"code.cloudfoundry.org/lager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gleidsonnunes/snickers2/db"
	"github.com/gleidsonnunes/snickers2/helpers"
	"github.com/gleidsonnunes/snickers2/types"
)

// S3Upload sends the file to S3 bucket. Job Destination should be
// in format: http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT
func S3Upload(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("s3-upload")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	file, err := os.Open(job.LocalSource)
	if err != nil {
		return err
	}

	err = helpers.SetAWSCredentials(job.Destination)
	if err != nil {
		return err
	}

	bucket, err := helpers.GetAWSBucket(job.Destination)
	if err != nil {
		return err
	}

	key, err := helpers.GetAWSKey(job.Destination)
	if err != nil {
		return err
	}

	job.Status = types.JobUploading
	job.Progress = "0%"
	dbInstance.UpdateJob(job.ID, job)

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	job.Progress = "100%"
	dbInstance.UpdateJob(job.ID, job)

	return nil
}
