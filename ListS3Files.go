// snippet-comment:[These are tags for the AWS doc team's sample catalog. Do not remove.]
// snippet-sourceauthor:[Doug-AWS]
// snippet-sourcedescription:[Lists the items in an S3 bucket.]
// snippet-keyword:[Amazon Simple Storage Service]
// snippet-keyword:[Amazon S3]
// snippet-keyword:[ListObjectsV2 function]
// snippet-keyword:[Go]
// snippet-sourcesyntax:[go]
// snippet-service:[s3]
// snippet-keyword:[Code Sample]
// snippet-sourcetype:[full-example]
// snippet-sourcedate:[2019-04-06]
/*
   Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/net/http2"
	"net/http"
	"os"
)

// Lists the items in the specified S3 Bucket
//
// Usage:
//    go run s3_list_objects.go BUCKET_NAME
func main() {
	/*	if len(os.Args) != 2 {
		exitErrorf("Bucket name required\nUsage: %s bucket_name",
			os.Args[0])
	}*/

	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 250
	defaultTransport.DisableCompression = false
	defaultTransport.MaxIdleConnsPerHost = 250

	// So client makes HTTP/2 requests
	_ = http2.ConfigureTransport(defaultTransportPointer)

	httpClient := &http.Client{Transport: &defaultTransport}
	bucket := "aeg-gracenote-program-metadata-sandbox"
	logLevel := aws.LogLevel(aws.LogDebugWithHTTPBody)
	awsCfg := &aws.Config{HTTPClient: httpClient, LogLevel: logLevel, Region: aws.String("us-west-2")}

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(awsCfg)

	// Create S3 service client
	svc := s3.New(sess)

	_ = &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String("aeg-gracenote-program-metadata-sandbox/EP000000060001.xml"),
		Key:        aws.String("EP000000060001.xml"),
	}

	//copyObjectOutput, err := svc.CopyObject(copyInput)

	//if err != nil {
	//	exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	//}
	//
	//fmt.Println(copyObjectOutput)

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	fmt.Println("")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
