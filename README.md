# maschera

---

`maschera` is a toy Golang application to demonstrate a simple Kafka
consumer and producer.

The application reads JSON messages from an input topic, and writes new messages
to an output topic with the specified top-level JSON fields hashed using SHA256.
The imaginary use case is a scenario where an organization wants to use, share,
or analyze data containing PII, without exposing the PII data itself. SHA256 is
used as the hashing algorithm because it is fast, and because the output is
deterministic. This allows downstream consumers of the masked data to query or
partition data by the masked PII fields, even though the real underlying values
are unknown to the consumer.

Possible improvements:
- Improve test coverage.
- Add an LRU cache in front of value hashing to avoid re-hashing the same
  values.
- Measure performance and tune the number of goroutines used for hashing and
  writing to Kafka.
- Instrument the application with Prometheus and/or statsd/dogstatsd metrics.
- Instead of simply hashing the values, map each hash to a unique identifier,
  persist the mapping in a database, and replace the PII values in the output
  stream with the identifier for the hash output. This would allow the operators
  to insulate consumers from changes like switching to a different hashing
  algorithm, or rolling the keys for the hash algorithm, while still protecting
  the PII data.
