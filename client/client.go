package kclient

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/nick-uk/ksql/ksql"
)

var defaultclient *ksql.Client
var clientWithTimeout *ksql.Client
var clientWithcancel *ksql.Client

// NewClientWithCancel creates a KSQL client and returns a cancel function
func NewClientWithCancel(host string) (*ksql.Client, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	client := ksql.NewClientContext(ctx, host)
	return client, cancel
}

// NewClientWithTimeout creates a KSQL client with configurable timeout
func NewClientWithTimeout(host string, timeout time.Duration) (*ksql.Client, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	client := ksql.NewClientContext(ctx, host)
	return client, cancel
}

// NewClient creates the default KSQL client
func NewClient(host string) *ksql.Client {
	return ksql.NewClient(host)
}

func printOutput1(QueryResponse []*ksql.QueryResponse) {
	for row := range QueryResponse {
		//log.Printf("%d - %+v\n", row, QueryResponse[row].Row.Columns)
		if QueryResponse[row].Row != nil {
			for col := range QueryResponse[row].Row.Columns {
				ErrTXT := ""
				if QueryResponse[row].Row.ErrorMessage != "" || QueryResponse[row].Row.FinalMessage != "" {
					ErrTXT = QueryResponse[row].Row.ErrorMessage + " " + QueryResponse[row].Row.FinalMessage
				}
				// Field can be any type
				log.Printf("Value: %v : %s (%s)\n", QueryResponse[row].Row.Columns[col], ErrTXT, reflect.TypeOf(QueryResponse[row].Row.Columns[col]).String())
			}
		} else {
			log.Printf("MSG: %s %s\n", QueryResponse[row].FinalMessage, QueryResponse[row].ErrorMessage)
		}
	}
}

func printOutput2(Kres []ksql.Response) {
	if len(Kres) < 1 {
		log.Println("Didn't get enough responses")
	}
	for row := range Kres {

		log.Println("Response:", row)

		if Kres[row].Topics != nil {
			topics := Kres[row].Topics
			for i := range topics {
				log.Println("Topic", i, ":", topics[i].Name, ", Registered:", topics[i].Registered, ", Consumers:", topics[i].Consumers, ", GroupConsumers:", topics[i].GroupConsumers)
			}
		}

		if Kres[row].Streams != nil {
			streams := Kres[row].Streams
			for i := range streams {
				log.Println("Stream", i, ":", streams[i].Name, ", Topic:", streams[i].Topic, ", Format:", streams[i].Format)
			}
		}

		if Kres[row].Tables != nil {
			tables := Kres[row].Tables
			for i := range tables {
				log.Println("Table", i, ":", tables[i].Name, ", Topic:", tables[i].Topic, ", Format:", tables[i].Format, ", Window:", tables[i].Windowed)
			}
		}

		if Kres[row].Message != "" {
			log.Println("=> MSG:", Kres[row].Message)
		}
	}
}

// Ksql is a more generic & simple KSQL client with output handling,
// set timeout = 0 for streams to wait forever
func Ksql(req ksql.Request, ksqlserver string, timeout time.Duration) (context.CancelFunc, error) {

	var cancel context.CancelFunc
	// Get timings
	start := time.Now()
	defer func() {
		elapsed := time.Since(start).Seconds()
		comments := ""
		if elapsed > 0.5 {
			comments = "(That was definitely slow...)"
		}
		log.Printf("== Took " + strconv.FormatFloat(elapsed, 'f', 4, 64) + "secs " + comments + "\n\n")
	}()

	// /query
	if strings.Contains(strings.ToLower(req.KSQL[0:7]), "select") {

		if timeout == 0 {
			if clientWithcancel == nil || clientWithcancel.Err() != nil {
				// create new client
				clientWithcancel, cancel = NewClientWithCancel(ksqlserver)
			}
			ch := make(chan *ksql.QueryResponse)

			go clientWithcancel.Query(req, ch)
			for {
				select {
				case msg := <-ch:
					res := append([]*ksql.QueryResponse{}, msg)
					printOutput1(res)
				}
			}
		} else {
			if clientWithTimeout == nil || clientWithTimeout.Err() != nil {
				// create new client
				clientWithTimeout, cancel = NewClientWithTimeout(ksqlserver, timeout)
			}
			res, qerr := clientWithTimeout.LimitQuery(req)
			if qerr != nil {
				return nil, errors.New("LimitQuery ERROR: " + qerr.Error())
			}
			printOutput1(res)
		}
	} else { // /ksql

		if defaultclient == nil || defaultclient.Err() != nil {
			// create new client
			defaultclient = NewClient(ksqlserver)
		}
		KsqlResponse, qerr := defaultclient.Do(req)
		if qerr != nil {
			return nil, qerr
		}
		printOutput2(KsqlResponse)

	}

	return cancel, nil
}
