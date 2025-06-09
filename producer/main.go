package producer

import (
	"fmt"
	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("Hello Producer")
}

// Exported Functions start with a capitalized Letter
func PrintTopics() {
	// Connect to localhost:9092
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// Read Partitions
	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	// Map the partitions with the topic as the key
	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}

	// Print the topics
	for k := range m {
		fmt.Println(k)
	}
}
