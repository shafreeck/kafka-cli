# kafka-cli

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

## Consume topics

```
./kafka-cli --brokers 127.0.0.1:9092 consume default-topic unexpected-logs
```

## Produce messages

TODO
