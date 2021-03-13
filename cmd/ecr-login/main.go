package main

import (
	"context"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	ecrinternal "github.com/jacoelho/ecr-login/internal/ecr"
)

type Config struct {
	Region string
}

func main() {
	var c Config
	flag.StringVar(&c.Region, "region", "eu-west-1", "aws region")

	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ecr.New(sess, aws.NewConfig().WithRegion(c.Region))

	tokens, err := ecrinternal.GetAuthorizationTokens(context.Background(), svc)
	if err != nil {
		log.Fatal(err.Error())
	}

	creds, err := ecrinternal.DecodeCredentials(tokens)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, c := range creds {
		if err := ecrinternal.IsCredentialValid(c); err != nil {
			log.Fatal(err.Error())
		}
	}

	for _, c := range creds {
		if err := ecrinternal.DockerLogin(c); err != nil {
			log.Fatal(err.Error())
		}
	}
}
