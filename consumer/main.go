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
