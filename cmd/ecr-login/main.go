package main

import (
	"context"
	"flag"
	"log"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tokens, err := ecrinternal.GetAuthorizationTokens(ctx, svc)
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
		if err := ecrinternal.DockerLogin(ctx, c); err != nil {
			log.Fatal(err.Error())
		}
	}
}
