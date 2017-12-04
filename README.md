# kafka-cli

## Install

```
go get github.com/shafreeck/kafka-cli
```

## Usage

```
kafka-cli is a console util tool to access kafka cluster

Usage:
  kafka-cli [command]

Available Commands:
  consume     consume topic from kafka
  monitor     display kafka cluster metrics, TODO
  produce     TODO
  topics      list all topics

Flags:
      --brokers string                    broker list, delimited by comma (default "127.0.0.1:9092")
      --buffersize int                    internal channel buffer size (default 256)
      --clientid string                   a user-provided string sent with every request to the brokers for logging debugging, and auditing purposes (default "kafka-cli")
      --config string                     config file (default is $HOME/.kafka-cli.yaml)
      --metadata.refresh duration         metadata refresh frequency (default 10m0s)
      --metadata.retry.backoff duration   backoff between retrying (default 250ms)
      --metadata.retry.max int            total number to request metadata when the cluster has a leader election (default 3)
      --net.dialtimeout duration          timeout of dialing to brokers (default 30s)
      --net.keepalive duration            keepalive period, 0 means disabled
      --net.maxopenrequests int           how many outstanding requests a connection is allowed to have before sending on it blocks (default 5)
      --net.readtimeout duration          timeout of reading messages (default 30s)
      --net.writetimeout duration         timeout of writing messages (default 30s)
  -v, --verbose                           print log messages
      --zookeepers string                 zookeeper server list, delimited by comma, only use when operate topic (default "127.0.0.1:9093")

Use "kafka-cli [command] --help" for more information about a command.
```

## List all topics

```
./kafka-cli --brokers 127.0.0.1:9092 topics
topic           partition[replicaid...]:offset ...
--------------------------------------------------
unexpected-logs 0[0]:12285581
output-pod2     0[0]:9725580
default-topic   0[0]:13971093
test            0[0]:5
filebeats       0[0]:12488893   1[0]:11301355   2[0]:11301353   3[0]:11301197   4[0]:11301197   5[0]:11301197   6[0]:11301200   7[0]:11301356
```

## Create topics

```
./kafka-cli --brokers 127.0.0.1:9092 --zookeepers 127.0.0.1:2181 topics create  t1 t2 t3 --partitions 10 --replicas 3
```

## Delete topics

```
./kafka-cli --brokers 127.0.0.1:9092 --zookeepers 127.0.0.1:2181 topics delete  t1 t2 t3
```

## Consume topics

```
./kafka-cli --brokers 127.0.0.1:9092 consume default-topic unexpected-logs
```

## Produce messages

TODO
