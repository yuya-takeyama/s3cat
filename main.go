package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	flags "github.com/jessevdk/go-flags"
)

const appName = "s3cat"

type options struct {
}

var opts options

func main() {
	parser := flags.NewParser(&opts, flags.Default^flags.PrintErrors)
	parser.Name = appName
	parser.Usage = "s3://..."

	args, err := parser.Parse()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sess := session.Must(session.NewSession())
	client := s3.New(sess)

	s3cat(client, args, os.Stdout, os.Stderr)
}

func s3cat(client *s3.S3, files []string, stdout io.Writer, stderr io.Writer) {
	for _, file := range files {
		fmt.Fprintln(stderr, "FILE: "+file)

		s3URL, urlErr := url.Parse(file)

		if urlErr != nil {
			fmt.Fprintln(stderr, urlErr)
			continue
		}

		obj, err := client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3URL.Host),
			Key:    aws.String(s3URL.Path),
		})

		if err != nil {
			fmt.Fprintln(stderr, err)
			continue
		}

		draw(obj, stdout, stderr)
	}
}

func draw(obj *s3.GetObjectOutput, stdout io.Writer, stderr io.Writer) {
	defer func() {
		obj.Body.Close()
	}()

	io.Copy(stdout, obj.Body)
}
