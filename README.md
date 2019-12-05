# GO-for-KSQL
Checking KSQL benchmarks and sql capabilities for Kafka streams


### A few words about KSQL
KSQL is a SQL engine for Kafka. It allows you to write SQL queries to analyze a stream of data in real time.

Since a stream is an unbounded data set, a query with KSQL will keep generating results until you stop it.

KSQL is built on top of Kafka Streams. When you submit a query, this query will be parsed and a Kafka Streams topology will be built and executed.

This means that KSQL offers similar concepts as to what Kafka Streams offers, but all with a SQL language: streams (KStreams), tables (KTables), joins, windowing functions, etc.

To use KSQL, you need to start one or more instances of KSQL Server. You can add more capacity easily by just starting more instances.

All the instances will work together by exchanging information through a private topic (_confluent-ksql-default__command_topic).

##### Users can interact with KSQL through a REST API or, more conveniently, through KSQL CLI, a command line interface

### What is the repository for?
- I'm using this repository to do KSQL experiments by using a more simplified KSQL Go client version
- Includes basic tests, ready to use docker-compose local stack and sample data generation.

```Go
client, qerr = kclient.Ksql(ksql.Request{KSQL: SQL}, ksqlserver, timeout)
if qerr != nil {
  log.Println("Query error:", qerr)
}
```

### How to use this repo (test & build)
```bash
$ ./test.sh
=== RUN   TestNewClient
--- PASS: TestNewClient (0.00s)
=== RUN   TestNewClientWithTimeout
=== RUN   TestNewClientWithTimeout/Test_ksql_client_with_timeout
--- PASS: TestNewClientWithTimeout (0.00s)
    --- PASS: TestNewClientWithTimeout/Test_ksql_client_with_timeout (0.00s)
=== RUN   TestNewClientWithCancel
=== RUN   TestNewClientWithCancel/Test_ksql_client_with_cancel
--- PASS: TestNewClientWithCancel (0.00s)
    --- PASS: TestNewClientWithCancel/Test_ksql_client_with_cancel (0.00s)
PASS
ok      client/kclient/client    (cached)
Creating network "go-ksql_default" with the default driver
Creating go-ksql_zookeeper_1 ... done
Creating go-ksql_kafka_1     ... done
Creating go-ksql_schema-registry_1 ... done
Creating go-ksql_datagen_1         ... done
Creating go-ksql_ksql-server_1     ... done
* Wait sample data to be generated...
* Run Go ksql client...
2019/12/05 09:22:44 =>>> Topics
2019/12/05 09:22:45 Response: 0
2019/12/05 09:22:45 Topic 0 : _confluent-metrics , Registered: false , Consumers: 0 , GroupConsumers: 0
2019/12/05 09:22:45 Topic 1 : _schemas , Registered: false , Consumers: 0 , GroupConsumers: 0
2019/12/05 09:22:45 Topic 2 : pageviews , Registered: false , Consumers: 0 , GroupConsumers: 0
2019/12/05 09:22:45 Topic 3 : users , Registered: false , Consumers: 0 , GroupConsumers: 0
2019/12/05 09:22:45 == Took 0.7921secs (That was definitely slow...)

2019/12/05 09:22:45 =>>> Create some streams & tables...
2019/12/05 09:22:51 Response: 0
2019/12/05 09:22:51 Response: 1
2019/12/05 09:22:51 Response: 2
2019/12/05 09:22:51 Response: 3
2019/12/05 09:22:51 Response: 4
2019/12/05 09:22:51 Response: 5
2019/12/05 09:22:51 == Took 6.2813secs (That was definitely slow...)

2019/12/05 09:22:51 =>>> Streams
2019/12/05 09:22:51 Response: 0
2019/12/05 09:22:51 Stream 0 : PAGEVIEWS_ENRICHED , Topic: PAGEVIEWS_ENRICHED , Format: DELIMITED
2019/12/05 09:22:51 Stream 1 : PAGEVIEWS_ORIGINAL , Topic: pageviews , Format: DELIMITED
2019/12/05 09:22:51 Stream 2 : PAGEVIEWS_FEMALE_LIKE_89 , Topic: pageviews_enriched_r8_r9 , Format: DELIMITED
2019/12/05 09:22:51 Stream 3 : PAGEVIEWS_FEMALE , Topic: PAGEVIEWS_FEMALE , Format: DELIMITED
2019/12/05 09:22:51 == Took 0.2010secs

2019/12/05 09:22:51 =>>> Tables
2019/12/05 09:22:52 Table 0: PAGEVIEWS_REGIONS
2019/12/05 09:23:07 == Took 15.0008secs (That was definitely slow...)

2019/12/05 09:23:07 Query error: LimitQuery ERROR: context deadline exceeded
2019/12/05 09:23:07 Table 1: USERS_ORIGINAL
2019/12/05 09:23:10 Value: 1.57553055866e+12 :  (float64)
2019/12/05 09:23:10 Value: User_3 :  (string)
2019/12/05 09:23:10 Value: 1.50484662532e+12 :  (float64)
2019/12/05 09:23:10 Value: MALE :  (string)
...
2019/12/05 09:23:10 MSG: Limit Reached
2019/12/05 09:23:10 == Took 3.0792secs (That was definitely slow...)

2019/12/05 09:23:10 Execution complete
Stopping go-ksql_ksql-server_1     ... done
Stopping go-ksql_schema-registry_1 ... done
Stopping go-ksql_kafka_1           ... done
Stopping go-ksql_zookeeper_1       ... done
Removing go-ksql_ksql-server_1     ... done
Removing go-ksql_datagen_1         ... done
Removing go-ksql_schema-registry_1 ... done
Removing go-ksql_kafka_1           ... done
Removing go-ksql_zookeeper_1       ... done
Removing network go-ksql_default
```

### Comments
KSQL is a very powerful way to combine streams and do processing by using SQL commands, perfect for real-time applications. By default KSQL reads the topics for streams and tables from the latest offset, it is a behavior that you can change but this doesn't means you can use Kafka as a SQL database.

KSQL server is much more practical to use Kafka with the ability to attach it easily with Load balancers and use clients from different sources avoiding the complexity of ADVERTISED addresses
