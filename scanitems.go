package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/user"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func exitWithError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

const tableName = "dev_gracenote_program"

/*type HTTPClientSettings struct {
	Connect          time.Duration
	ConnKeepAlive    time.Duration
	ExpectContinue   time.Duration
	IdleConn         time.Duration
	MaxAllIdleConns  int
	MaxHostIdleConns int
	ResponseHeader   time.Duration
	TLSHandshake     time.Duration
}

func NewHTTPClientWithTimeouts(httpSettings HTTPClientSettings) *http.Client {
	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeader,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: httpSettings.ConnKeepAlive,
			Timeout:   httpSettings.Connect,
		}).DialContext,
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConn,
		TLSHandshakeTimeout:   httpSettings.TLSHandshake,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinue,
	}
	// So client makes HTTP/2 requests
	_ = http2.ConfigureTransport(tr)

	return &http.Client{
		Transport: tr,
	}
}*/

func main() {
	cfg := Config{tableName, "us-west-2"}
	fmt.Println("Using", cfg)

	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 100
	defaultTransport.DisableCompression = false
	defaultTransport.MaxIdleConnsPerHost = 100

	// So client makes HTTP/2 requests
	_ = http2.ConfigureTransport(defaultTransportPointer)

	httpClient := &http.Client{Transport: &defaultTransport}

	// Create the config specifying the Region for the DynamoDB table.
	// If Config.Region is not set the region must come from the shared
	// config or AWS_REGION environment variable.
	/*	logLevel := aws.LogLevel(aws.LogDebugWithHTTPBody)
		LogLevel: logLevel,*/
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.HomeDir)
	awsCfg := &aws.Config{HTTPClient: httpClient, Credentials: credentials.NewSharedCredentials(usr.HomeDir+"/.aws/credentials", "default")}
	if len(cfg.Region) > 0 {
		awsCfg.WithRegion(cfg.Region)
	}
	var producedMessages, consumedMessages uint64

	// Create the session that the DynamoDB service will use.
	sess := session.Must(session.NewSession(awsCfg))

	var waitGroup sync.WaitGroup
	noOfThreads := 100

	ticker := time.NewTicker(300 * time.Second)
	done := make(chan bool)

	printPeriodically(done, ticker, &producedMessages, &consumedMessages)

	// Create the DynamoDB service client to make the query request with.
	svc := dynamodb.New(sess)
	start := time.Now()
	timestamp := makeTimestamp()
	fmt.Println(now(), timestamp)
	for i := 0; i < noOfThreads; i++ {
		waitGroup.Add(1)
		go scan(cfg, svc, int64(i), int64(noOfThreads), &producedMessages, &consumedMessages, &waitGroup, timestamp)
	}
	waitGroup.Wait()
	for {
		if producedMessages == consumedMessages {
			break
		}
		printLogs(&producedMessages, &consumedMessages)
		time.Sleep(5 * time.Second)
	}
	printLogs(&producedMessages, &consumedMessages)
	//ticker.Stop()
	elapsed := time.Since(start)
	fmt.Println(producedMessages)
	log.Printf("Processing took %s\n", elapsed)
	done <- true
	runtime.GC()
	printLogs(&producedMessages, &consumedMessages)
}

func printPeriodically(done chan bool, ticker *time.Ticker, producerCount *uint64, consumerCount *uint64) {
	go func(producerCount *uint64, consumerCount *uint64) {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				printLogs(producerCount, consumerCount)
			}
		}
	}(producerCount, consumerCount)
}

func printLogs(producerCount *uint64, consumerCount *uint64) {
	format := "%07d"
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	exitReport := fmt.Sprintf("%s | %s | goroutinesCount = %d | memAllocBytes = %s | "+"heapAllocBytes = %s | heapObjects = %d | stacksInUse = %d  |",
		fmt.Sprintf(format, *producerCount), fmt.Sprintf(format, *consumerCount), runtime.NumGoroutine(), ByteCountDecimal(stats.Alloc), ByteCountDecimal(stats.HeapAlloc), stats.HeapObjects, stats.StackInuse)
	log.Println(exitReport)
}

func now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func ByteCountDecimal(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

/**
  1576674964891
  13639
  26797

*/

func scan(cfg Config, svc *dynamodb.DynamoDB, segment int64, totalSegments int64, pr *uint64, co *uint64, waitGroup *sync.WaitGroup, timestamp int64) {
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String(cfg.Table), ProjectionExpression: aws.String("tmsId"),
		Segment: aws.Int64(segment), TotalSegments: aws.Int64(totalSegments), Limit: aws.Int64(500),
	}
	for {
		// Make the DynamoDB Query API call
		result, err := svc.Scan(params)
		if err != nil {
			exitWithError(fmt.Errorf("failed to make Query API call, %v", err))
		}

		var items []Item

		//Use CHANNEL 103 , 104 and 115
		c := make(chan string, 200)
		go consume(c, co, timestamp, svc)

		// Unmarshal the Items field in the result value to the Item Go type.
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
		if err != nil {
			exitWithError(fmt.Errorf("failed to unmarshal Query result items, %v", err))
		}

		// Print out the items returned
		for _, item := range items {
			//Use CHANNEL
			c <- item.TmsId
			//updateItem(item.TmsId, timestamp, svc)
			atomic.AddUint64(pr, uint64(1))

		}

		/*	if *co > 1 {
			os.Exit(0)
		}*/
		//fmt.Println(itemSize, pr)
		params.ExclusiveStartKey = result.LastEvaluatedKey
		if result.LastEvaluatedKey == nil {
			//close(c)
			break
		}
		result = nil
		items = nil
		close(c)
	}

	waitGroup.Done()
}

func consume(in chan string, ops *uint64, timestamp int64, svc *dynamodb.DynamoDB) {
	for s := range in {
		if len(s) > 0 {
			updateItem(s, timestamp, svc)
			atomic.AddUint64(ops, uint64(1))
		}
	}
}

func updateItem(tmsId string, timestamp int64, svc *dynamodb.DynamoDB) {
	tmsIdKey, err := dynamodbattribute.Marshal(tmsId)
	if err != nil {
		fmt.Println("Got error marshalling item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var primaryKey = map[string]*dynamodb.AttributeValue{
		"tmsId": tmsIdKey,
	}

	marshal, _ := dynamodbattribute.Marshal(timestamp)
	var upExpr = map[string]*dynamodb.AttributeValue{
		":reprocessedTime": marshal,
	}
	input := &dynamodb.UpdateItemInput{
		ConditionExpression:       aws.String("attribute_exists(tmsId)"),
		TableName:                 aws.String(tableName),
		Key:                       primaryKey,
		UpdateExpression:          aws.String("SET reprocessedTime = :reprocessedTime"),
		ExpressionAttributeValues: upExpr,
		//ReturnValues:              aws.String("UPDATED_NEW"),
	}
	_, err = svc.UpdateItem(input)
	if err != nil {
		fmt.Println("ERROR ----------------- " + tmsId)
		fmt.Println(err.Error())
		return
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type Item struct {
	TmsId string
}

type Config struct {
	Table  string // required
	Region string // optional

}
