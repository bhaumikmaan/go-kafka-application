package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello Consumer")
}

// ConsumeMessages To Test - kafka-console-producer --bootstrap-server localhost:9092 --topic first_topic
// Input and see your own data being consumed
func ConsumeMessages(topic string) {
	fmt.Println("Starting to read from topic: ", topic)
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer conn.Close()

	for {
		// Set a read deadline to avoid blocking forever in case there are no messages. Here, 60 Seconds
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		batch := conn.ReadBatch(10, 1e6) // fetch 10 bytes min, 1MB max

		b := make([]byte, 10e3) // 10KB max per message
		for {
			n, err := batch.Read(b)
			if err != nil {
				// If the error is EOF, i.e. batch is finished, so break inner loop to get a new batch
				break
			}
			fmt.Println(string(b[:n]))
		}

		if err := batch.Close(); err != nil {
			log.Fatal("failed to close batch:", err)
		}
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

// ConsumeWithGroup consumes messages using Kafka consumer groups
// To Test: kafka-console-consumer --bootstrap-server localhost:9092 --topic first_topic --group my-first-application
func ConsumeWithGroup(topic, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  groupID,
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
	})

	fmt.Printf("Starting consumer group reader: topic=%s, group ID=%s\n", topic, groupID)
	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		fmt.Printf("Message at %s/%d/%d: %s = %s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
