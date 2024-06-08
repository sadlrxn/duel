package utils

import (
	"bytes"
	"image"
	"image/png"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/*
* @Internal
* Create a new session with AWS credentials.
 */
func getAWSSession() *session.Session {
	// 1. Get aws configurations.
	config := config.Get()
	awsRegion := config.AWSRegion
	awsAccessID := config.AWSAccessID
	awsSecretKey := config.AWSSecretKey

	// 2. Create a new session with aws credentials.
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(awsRegion),
				Credentials: credentials.NewStaticCredentials(
					awsAccessID,
					awsSecretKey,
					"",
				),
			},
		),
	)
	return sess
}

/*
* @Internal
* Get AWS S3 bucket uploader.
 */
func getAWSS3Uploader() *s3manager.Uploader {
	// 1. Get session.
	sess := getAWSSession()

	// 2. Create & return s3 bucket uploader.
	s3ManagerClient := s3manager.NewUploader(
		sess,
		func(u *s3manager.Uploader) {
			u.PartSize = int64(0)
			u.Concurrency = 0
		})
	return s3ManagerClient
}

/*
* @Internal
* Get AWS S3 bucket service client.
 */
func getAWSS3ServiceClient() *s3.S3 {
	// 1. Get session.
	sess := getAWSSession()

	// 2. Return sevice client.
	return s3.New(sess)
}

/*
* @External
* Upload image to the AWS S3 bucket with key.
 */
func UploadImage(
	image image.Image,
	fileType string,
	key string,
) (string, error) {
	// 1. Get AWS bucket name.
	bucket := config.Get().S3BucketName

	// 2. Encode image into PNG.
	buf := new(bytes.Buffer)
	png.Encode(buf, image)

	// 3. Create new byte reader.
	reader := bytes.NewReader(buf.Bytes())

	// 4. Get S3 bucket uploader.
	s3ManagerClient := getAWSS3Uploader()

	// 5. Upload image and return new url.
	s3UploadOutput, err := s3ManagerClient.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        reader,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/png"),
	})
	return s3UploadOutput.Location, err
}

/*
* @External
* Delete image from AWS S3 bucket.
 */
func DeleteImage(key string) error {
	// 1. Get aws s3 service client object.
	bucket := config.Get().S3BucketName
	svc := getAWSS3ServiceClient()

	// 2. Delete object from bucket.
	if _, err := svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	); err != nil {
		return MakeError(
			"utils_aws_s3_manager",
			"DeleteImage",
			"failed to delete object",
			err,
		)
	}

	// 3. Wait until the object not exists.
	if err := svc.WaitUntilObjectNotExists(
		&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	); err != nil {
		return MakeError(
			"utils_aws_s3_manager",
			"DeleteImage",
			"failed to wait until object not exists",
			err,
		)
	}
	return nil
}
