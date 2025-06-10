package producer

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

// TestProduceMessagesIntegration runs ProduceMessages for a short time and verifies
// it produces messages to Kafka.
// Requires Kafka running at localhost:9092 and topic "test-stream-golang" created.
func TestProduceMessagesIntegration(t *testing.T) {
	// Run ProduceMessages in a goroutine, it runs indefinitely reading the Wikimedia stream
	done := make(chan struct{})
	go func() {
		ProduceMessages("test-stream-golang")
		close(done)
	}()

	// Allow the producer to generate some messages
	select {
	case <-done:
		t.Log("ProduceMessages finished before timeout")
	case <-time.After(5 * time.Second):
		t.Log("Timeout reached, stopping ProduceMessages test")
	}

	// Now check if Kafka received messages
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "test-stream-golang",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		t.Fatalf("failed to read message from kafka: %v", err)
	}

	t.Logf("Received message from Kafka topic: %s", string(msg.Value))

	_ = reader.Close()
}
