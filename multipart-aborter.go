package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/s3"
)

var s3Svc *s3.S3

func main() {
	spBucket := flag.String("bucket", "", "Enter bucket to clear in-progress multipart uploads.")
	flag.Parse()

	s3Svc = s3.New(aws.DefaultConfig)

	resp := mustListMultipartUploads(spBucket)

	for _, upload := range resp.Uploads {
		mustAbort(spBucket, upload)
		confirmDeletion(spBucket, upload)
	}
}

func mustListMultipartUploads(spBucket *string) *s3.ListMultipartUploadsOutput {
	params := &s3.ListMultipartUploadsInput{
		Bucket:     spBucket,
		Delimiter:  aws.String("/"),
		MaxUploads: aws.Long(1000),
		Prefix:     aws.String(""),
	}

	resp, err := s3Svc.ListMultipartUploads(params)
	if awserr := aws.Error(err); awserr != nil {
		log.Fatal("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		panic(err)
	}

	fmt.Println(awsutil.StringValue(resp))

	return resp
}

func mustAbort(spBucket *string, upload *s3.MultipartUpload) {
	params := &s3.AbortMultipartUploadInput{
		Bucket:   spBucket,
		Key:      upload.Key,
		UploadID: upload.UploadID,
	}

	resp, err := s3Svc.AbortMultipartUpload(params)
	if awserr := aws.Error(err); awserr != nil {
		log.Fatal("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		panic(err)
	}

	fmt.Println(awsutil.StringValue(resp))
}

func confirmDeletion(spBucket *string, upload *s3.MultipartUpload) {
	params := &s3.ListPartsInput{
		Bucket:   spBucket,
		Key:      upload.Key,
		UploadID: upload.UploadID,
		MaxParts: aws.Long(1),
	}

	resp, err := s3Svc.ListParts(params)
	if awserr := aws.Error(err); awserr != nil {
		fmt.Println("Error:", awserr.Code, awserr.Message)
		fmt.Println("The upload below still has Parts listed.")
		fmt.Println(awsutil.StringValue(upload))
	} else if err != nil {
		panic(err)
	}

	fmt.Println(awsutil.StringValue(resp))
}
