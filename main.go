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
	ch := make(chan string)

	// Print header
	fmt.Println("DB Name, Region, CA Cert")

	for _, region := range regionsFlag {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			getRDSInstances(region, ch)
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

func getRDSInstances(region string, ch chan<- string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profileFlag,
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := rds.New(sess)
	input := &rds.DescribeDBInstancesInput{}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		fmt.Println("Error getting RDS instances:", err)
		return
	}

	for _, instance := range result.DBInstances {
		ch <- fmt.Sprintf("%s, %s, %s", region, *instance.CACertificateIdentifier, *instance.DBInstanceIdentifier)
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
