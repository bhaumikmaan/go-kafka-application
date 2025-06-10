package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"os"
	"testing"
	"time"
)

// UNIT TEST: Using a Fake Reader (No Kafka Required)
// fakeReader implements kafka.Reader interface partially for testing
type fakeReader struct {
	messages []kafka.Message
	index    int
}

// ReadMessage returns the next message in the slice or simulates EOF.
func (r *fakeReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if r.index >= len(r.messages) {
		return kafka.Message{}, context.Canceled // Simulate end of messages
	}
	msg := r.messages[r.index]
	r.index++
	return msg, nil
}

// Close: Stubbed method for closing
func (r *fakeReader) Close() error {
	return nil
}

// TestConsumeWithFakeReader tests that we can read and log messages from a mock Kafka reader.
func TestConsumeWithFakeReader(t *testing.T) {
	testMsgs := []kafka.Message{
		{Topic: "test-topic", Partition: 0, Offset: 1, Key: []byte("key1"), Value: []byte("value1")},
		{Topic: "test-topic", Partition: 0, Offset: 2, Key: []byte("key2"), Value: []byte("value2")},
	}

	reader := &fakeReader{messages: testMsgs}

	// simulate reading from the reader
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		t.Logf("Test received message: %s = %s", string(msg.Key), string(msg.Value))
	}
}

// INTEGRATION TEST: Using a Real Kafka Broker
// TestKafkaIntegration writes a message to Kafka and reads it back.
// Prerequisites: Kafka broker is running on localhost:9092 and a topic named test topic is present
func TestKafkaIntegration(t *testing.T) {
	// kafka-topics --bootstrap-server localhost:9092 --topic test-topic --create

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092" // default for local
	}

	// Create a writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "test-topic",
	})

	// Write messages to the topic
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test"),
		Value: []byte("integration test"),
	})
	if err != nil {
		t.Fatalf("failed to write message: %v", err)
	}
	time.Sleep(10 * time.Second)

	// Create a Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  "test-group",
		Topic:    "test-topic",
		MaxBytes: 10e6,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	t.Logf("Received message from topic: %s", string(msg.Value))

	_ = reader.Close()
	_ = writer.Close()
}
