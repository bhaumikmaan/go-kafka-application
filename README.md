<h1 align="middle"> <img src="https://skillicons.dev/icons?i=go" width="35"/> Go & Kafka Application <img src="https://skillicons.dev/icons?i=kafka" width="35" /></h1>

<div align="center">

This repository houses a practical Go application demonstrating core Apache Kafka integration, pulling live data from the Wikimedia recent changes stream. It serves as a hands-on follow-up to theoretical Kafka and Go fundamentals.

[![bhaumikmaan - go-kafka-application](https://img.shields.io/static/v1?label=bhaumikmaan&message=go-kafka-application&color=00ADD8&logo=github)](https://github.com/bhaumikmaan/go-kafka-application "Go to GitHub repo")
&nbsp; [![stars - go-kafka-application](https://img.shields.io/github/stars/bhaumikmaan/go-kafka-application?style=social)](https://github.com/bhaumikmaan/go-kafka-application)
&nbsp; [![forks - go-kafka-application](https://img.shields.io/github/forks/bhaumikmaan/go-kafka-application?style=social)](https://github.com/bhaumikmaan/go-kafka-application)
&nbsp; [![License: MIT](https://img.shields.io/badge/License-MIT-00ADD8.svg)](https://opensource.org/licenses/MIT)
<br/>

<a href="https://medium.com/@bhaumikmaan/go-kafka-from-theory-to-practice-96e0757e02ca" target="_blank">
<img src="https://img.shields.io/badge/Read%20Blog%20on-Medium-black?logo=medium&style=for-the-badge" alt="Read Blog on Medium"/>
</a>
&nbsp;
<a href="https://github.com/bhaumikmaan/go-kafka-application/stargazers" target="_blank">
 <img src="https://img.shields.io/badge/Star-This%20Repo-yellow?style=for-the-badge&logo=github" alt="Star this repo"/>
</a>
</div>

---

## Project Overview

This application is a Command-Line Interface (CLI) tool built in Go that interacts with a Kafka cluster. It showcases fundamental Kafka operations:
* **Listing available topics**.
* **Producing messages** by consuming real-time data from the **Wikimedia recent changes stream**.
* **Consuming messages** from specified topics.

It utilizes the [segmentio/kafka-go](https://pkg.go.dev/github.com/segmentio/kafka-go) library for Kafka communication.

> **This serves as a hands-on follow-up to our series of Go and Kafka blogs, translating theoretical knowledge into practical application.** You can find the same below:
> <br> 
> * GoLang Series: [Medium](https://medium.com/@bhaumikmaan/list/golang-92f51617836c) | [Github](https://github.com/bhaumikmaan/Understanding-GoLang)
> * Kafka Series: [Medium](https://medium.com/@bhaumikmaan/list/kafka-549d42bfaa31) | [Github](https://github.com/bhaumikmaan/Understanding-Kafka-Basics)

---

## Project Structure & Dependencies

The application is organized into a modular Go project:

* **`main.go`**: The CLI entry point. It parses command-line arguments (e.g., `list-topics`, `start-producer <topic>`, `start-consumer <topic>`) using `os.Args` and a `switch` statement, then dispatches calls to the respective packages.
* **`producer/` package**: Contains logic for ingesting data and sending it to Kafka, and for listing topics.
* **`consumer/` package**: Holds logic for reading data from Kafka topics.

Dependencies are managed by **Go Modules** via `go.mod`:
```go
module go-kafka-application

go 1.24.3

require (
    github.com/klauspost/compress v1.15.9 // indirect
    github.com/pierrec/lz4/v4 v4.1.15 // indirect
    github.com/segmentio/kafka-go v0.4.48 // indirect
)
```
The go.mod file ensures consistent and reproducible builds by defining the module path, Go version, and managing direct (github.com/segmentio/kafka-go) and indirect dependencies.

## Core Components Explained

### **Producer (`producer` package)**
The producer handles data ingestion and sending to Kafka:
* **Wikimedia Stream Integration**: The `readStream()` function connects to the live Wikimedia EventSource stream using `net/http` and `bufio.NewScanner`. It runs in a **goroutine**, asynchronously parsing JSON events into an `Event` struct and sending them via a **Go channel** for efficient, non-blocking data flow.
* **Kafka Interaction (`ProduceMessages`)**: It connects to a specific Kafka partition leader using `kafka.DialLeader` (or `kafka.NewWriter` for higher-level use cases). It sets a **write deadline** to prevent hangs and sends messages with a `Value` (the event data) and optionally a `Key` (for consistent partitioning).
* **Topic Listing (`PrintTopics`)**: Uses `kafka.Dial` and `conn.ReadPartitions()` to discover and display all topics on the connected Kafka broker.

### **Consumer (`consumer` package)**
The consumer handles reading messages from Kafka:
* **Low-Level Consumer (`ConsumeMessages`)**: Demonstrates a direct connection to a **specific partition** using `kafka.DialLeader`. It uses `conn.SetReadDeadline` for timeouts and `conn.ReadBatch` for efficient fetching of multiple messages at once. This is suitable for basic testing but lacks scalability features.
* **Scalable Consumer (`ConsumeWithGroup`)**: This is the **recommended approach** for production. It uses `kafka.NewReader` configured with a `GroupID`. This client handles **consumer group coordination**, distributing partitions among consumer instances, enabling **parallel processing**, **automatic offset management**, and **fault tolerance** (rebalancing partitions if consumers join/leave).

---

## Getting Started

To run this application, you'll need a Kafka cluster running locally (e.g., via Docker Compose).

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/bhaumikmaan/go-kafka-application.git
    cd go-kafka-application
    ```
2.  **Initialize Go modules**:
    ```bash
    go mod tidy
    ```
3.  **Run commands**:
    * List topics: `go run main.go list-topics`
    * Start producer (e.g., to topic `wikimedia-changes`): `go run main.go start-producer wikimedia-changes`
    * Start consumer (e.g., from topic `wikimedia-changes`): `go run main.go start-consumer wikimedia-changes`
    * Start consumer group (e.g., topic `my-topic`, group `my-group`): `go run main.go start-consumer-group my-topic my-group`

---

## Testing

The project includes both unit and integration tests for reliability:

* **Unit Tests**: `TestConsumeWithFakeReader` (in `consumer_test.go`) uses a `fakeReader` to mock Kafka, allowing isolated, fast testing of consumer logic without a live Kafka instance.
* **Integration Tests**: `TestProduceMessagesIntegration` (in `producer_test.go`) and `TestKafkaIntegration` (in `consumer_test.go`) connect to a **real Kafka broker**. They verify end-to-end functionality by writing messages and then reading them back, ensuring correct communication with Kafka.

---

## Best Practices & Further Exploration

This application provides a strong foundation. For production readiness, consider:
* **Robust Error Handling**: Implement retries with backoff and comprehensive logging.
* **Configuration**: Use environment variables or dedicated config files for Kafka settings.
* **Monitoring**: Track Kafka metrics and application health.

You can further explore:
* **Schema Registry**: For message schema validation (e.g., Avro, Protobuf).
* **Exactly-Once Semantics**: For guaranteed message processing.
* **Kafka Streams API**: For complex real-time data transformations.

---
**License**

This project is licensed under the MIT License.