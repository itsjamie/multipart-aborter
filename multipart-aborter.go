package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Svc *s3.S3

func main() {
	spOp := flag.String("op", "removeMulti", "Enter operation to perform.")
	spBucket := flag.String("bucket", "", "Enter bucket to perform operation on.")
	spRegion := flag.String("region", "us-east-1", "The AWS Region to use for the bucket.")
	flag.Parse()

	config := aws.NewConfig().WithRegion(*spRegion)
	s3Svc = s3.New(config)

	switch *spOp {
	case "removeMulti":
		resp := mustListMultipartUploads(spBucket)

		for _, upload := range resp.Uploads {
			mustAbort(spBucket, upload)
			confirmDeletion(spBucket, upload)
		}
	case "removeObjects":
		params := &s3.ListObjectsInput{
			Bucket:    spBucket, // Required
			Delimiter: aws.String(""),
			Prefix:    aws.String(""),
		}
		resp, err := s3Svc.ListObjects(params)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Generic AWS Error with Code, Message, and original error (if any)
				fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
				if reqErr, ok := err.(awserr.RequestFailure); ok {
					// A service error occurred
					fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
				}
			}
		}

		delparams := &s3.DeleteObjectsInput{
			Bucket: spBucket, // Required
			Delete: &s3.Delete{ // Required
				Objects: make([]*s3.ObjectIdentifier, len(resp.Contents)),
				Quiet:   aws.Bool(true),
			},
		}

		for idx, obj := range resp.Contents {
			delparams.Delete.Objects[idx] = &s3.ObjectIdentifier{
				Key: obj.Key,
			}
		}

		_, err = s3Svc.DeleteObjects(delparams)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("No operation supported.")
	}

}

func mustListMultipartUploads(spBucket *string) *s3.ListMultipartUploadsOutput {
	params := &s3.ListMultipartUploadsInput{
		Bucket:     spBucket,
		Delimiter:  aws.String(""),
		MaxUploads: aws.Int64(1000),
		Prefix:     aws.String(""),
	}

	resp, err := s3Svc.ListMultipartUploads(params)
	if err != nil {
		panic(err)
	}

	return resp
}

func mustAbort(spBucket *string, upload *s3.MultipartUpload) {
	params := &s3.AbortMultipartUploadInput{
		Bucket:   spBucket,
		Key:      upload.Key,
		UploadId: upload.UploadId,
	}

	_, err := s3Svc.AbortMultipartUpload(params)
	if err != nil {
		panic(err)
	}
}

func confirmDeletion(spBucket *string, upload *s3.MultipartUpload) {
	params := &s3.ListPartsInput{
		Bucket:   spBucket,
		Key:      upload.Key,
		UploadId: upload.UploadId,
		MaxParts: aws.Int64(1),
	}

	resp, err := s3Svc.ListParts(params)
	if err != nil {
		fmt.Println(err)
	}

	if len(resp.Parts) > 0 {
		fmt.Println("This upload still has parts in it.")
	}
}
