package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

var (
	profileFlag string
	regionsFlag stringSlice
)

func init() {
	flag.StringVar(&profileFlag, "profile", "default", "AWS profile to use")
	flag.Var(&regionsFlag, "region", "AWS region(s) to check. Can be specified multiple times.")
}

func main() {
	flag.Parse()

	if len(regionsFlag) == 0 {
		regionsFlag = append(regionsFlag, "us-west-1") // default region if none is specified
	}

	var wg sync.WaitGroup

	// Use a buffered channel
	ch := make(chan string, len(regionsFlag))

	// Print header
	fmt.Println("DB Name, Region, CA Cert")

	// Create a single session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           profileFlag,
		SharedConfigState: session.SharedConfigEnable,
	}))

	for _, region := range regionsFlag {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			getRDSInstances(sess, region, ch)
		}(region)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}
}

func getRDSInstances(sess *session.Session, region string, ch chan<- string) {
	// Use the session and just change the region for the service client
	svc := rds.New(sess, &aws.Config{Region: aws.String(region)})

	input := &rds.DescribeDBInstancesInput{}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		// Consider using log.Printf instead of fmt for better error handling
		fmt.Println("Error getting RDS instances for region", region, ":", err)
		return
	}

	for _, instance := range result.DBInstances {
		ch <- fmt.Sprintf("%s, %s, %s, %v", region, *instance.CACertificateIdentifier, *instance.DBInstanceIdentifier)
	}
}

// stringSlice allows to specify multiple flags with the same name.
type stringSlice []string

func (i *stringSlice) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}
