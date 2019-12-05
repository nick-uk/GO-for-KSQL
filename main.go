package main

import (
	kclient "client/kclient/client"
	"log"
	"time"

	"github.com/nick-uk/ksql/ksql"
)

// A few words about KSQL
// KSQL is a SQL engine for Kafka. It allows you to write SQL queries to analyze a stream of data in real time.
// Since a stream is an unbounded data set, a query with KSQL will keep generating results until you stop it.
// KSQL is built on top of Kafka Streams. When you submit a query, this query will be parsed and a Kafka Streams topology will be built and executed.
// This means that KSQL offers similar concepts as to what Kafka Streams offers, but all with a SQL language: streams (KStreams), tables (KTables), joins, windowing functions, etc.
// To use KSQL, you need to start one or more instances of KSQL Server. You can add more capacity easily by just starting more instances.
// All the instances will work together by exchanging information through a private topic (_confluent-ksql-default__command_topic).

// Users can interact with KSQL through a REST API or, more conveniently, through KSQL CLI, a command line interface.

func main() {

	defer func() {
		log.Println("Execution complete")
	}()

	var qerr error
	ksqlserver := "http://localhost:8088"
	timeout := time.Second * 15
	SQLCmd := ""
	Defaultclient := kclient.NewClient(ksqlserver)

	// List existing topics
	log.Println("=>>> Topics")
	_, qerr = kclient.Ksql(ksql.Request{KSQL: "SHOW topics;"}, ksqlserver, 0)
	if qerr != nil {
		log.Panic(qerr)
	}

	// Create some streams & tables
	log.Println("=>>> Create some streams & tables...")
	SQL := `CREATE STREAM pageviews_original (viewtime bigint, userid varchar, pageid varchar) WITH (kafka_topic='pageviews', value_format='DELIMITED');
    CREATE TABLE users_original (registertime BIGINT, gender VARCHAR, regionid VARCHAR, userid VARCHAR) WITH (kafka_topic='users', value_format='JSON', key = 'userid');
    -- Create a persistent query by using the CREATE STREAM keywords to precede the SELECT statement
    CREATE STREAM pageviews_enriched AS SELECT users_original.userid AS userid, pageid, regionid, gender FROM pageviews_original LEFT JOIN users_original ON pageviews_original.userid = users_original.userid;
    -- Create a new persistent query where a condition limits the streams content, using WHERE
    CREATE STREAM pageviews_female AS SELECT * FROM pageviews_enriched WHERE gender = 'FEMALE';
    -- Create a new persistent query where another condition is met, using LIKE
    CREATE STREAM pageviews_female_like_89 WITH (kafka_topic='pageviews_enriched_r8_r9') AS SELECT * FROM pageviews_female WHERE regionid LIKE '%_8' OR regionid LIKE '%_9';
    -- Create a new persistent query that counts the pageviews for each region and gender combination in a tumbling window of 30 seconds when the count is greater than one
	CREATE TABLE pageviews_regions WITH (VALUE_FORMAT='avro') AS SELECT gender, regionid , COUNT(*) AS numusers FROM pageviews_enriched WINDOW TUMBLING (size 30 second) GROUP BY gender, regionid HAVING COUNT(*) > 1;`

	_, qerr = kclient.Ksql(ksql.Request{KSQL: SQL}, ksqlserver, 0)
	if qerr != nil {
		log.Println("Query error:", qerr)
	}

	// Check streams
	log.Println("=>>> Streams")
	_, qerr = kclient.Ksql(ksql.Request{KSQL: "SHOW streams;"}, ksqlserver, 0)
	if qerr != nil {
		log.Panic(qerr)
	}

	// Check tables
	log.Println("=>>> Tables")
	tables, err := Defaultclient.ListTables()
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range tables {
		log.Printf("Table %d: %s", i, v.Name)

		// Check select speed for earliest records
		// By default KSQL reads the topics for streams and tables from the latest offset, so let's change this to see if makes sense to use Kafka to "Search" in a table
		SQLCmd = "SELECT * from " + v.Name + " limit 5;"
		streamPro := make(map[string]string)
		streamPro["ksql.streams.auto.offset.reset"] = "earliest"

		_, qerr = kclient.Ksql(ksql.Request{
			KSQL:              SQLCmd,
			StreamsProperties: streamPro,
		}, ksqlserver, timeout)

		if qerr != nil {
			log.Println("Query error:", qerr)
		}
	}

	// TODO: write real-time anomaly detection
	// ksqlDB is a good fit for identifying patterns or anomalies on real-time data. By processing the stream as data arrives you can identify
	// and properly surface out of the ordinary events with millisecond latency.

	SQLCmd = `CREATE TABLE possible_fraud AS
	SELECT card_number, count(*)
	FROM authorization_attempts
	WINDOW TUMBLING (SIZE 5 SECONDS)
	GROUP BY card_number
	HAVING count(*) > 3 EMIT CHANGES;`
	// ch := make(chan *ksql.QueryResponse)
	// go Defaultclient.Query(ksql.Request{
	// 	KSQL:              SQLCmd,
	// }, ch)
	// for {
	// 	select {
	// 	case msg := <-ch:
	// 		log.Println(msg.Row)
	// 	}
	// }

}
