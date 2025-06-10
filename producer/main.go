package producer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"strings"
	"time"
)

// Event - Message structure to map the recentchange stream
// Event - structure that wraps the `meta` and the actual change data
// Flatten Event struct for easier unmarshalling
type Event struct {
	Meta struct {
		URI       string `json:"uri"`
		RequestID string `json:"request_id"`
		ID        string `json:"id"`
		DT        string `json:"dt"`
		Domain    string `json:"domain"`
		Stream    string `json:"stream"`
		Topic     string `json:"topic"`
		Partition int    `json:"partition"`
		Offset    int64  `json:"offset"`
	} `json:"meta"`

	// Recent change fields
	ID               int64  `json:"id"`
	Type             string `json:"type"`
	Namespace        int    `json:"namespace"`
	Title            string `json:"title"`
	TitleURL         string `json:"title_url"`
	Comment          string `json:"comment"`
	Timestamp        int64  `json:"timestamp"`
	User             string `json:"user"`
	Bot              bool   `json:"bot"`
	NotifyURL        string `json:"notify_url"`
	ServerURL        string `json:"server_url"`
	ServerName       string `json:"server_name"`
	ServerScriptPath string `json:"server_script_path"`
	Wiki             string `json:"wiki"`
	ParsedComment    string `json:"parsedcomment"`
}

func main() {
	fmt.Println("Hello Producer")
}

// PrintTopics Exported Functions start with a capitalized Letter
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

// readStream Internally read the stream data and produce messages for the topic
func readStream() (chan Event, error) {
	// Channel to send stream data
	changesChannel := make(chan Event)

	// Start streaming from Wikimedia recent changes
	fmt.Println("Starting to read from Wikimedia stream...")

	resp, err := http.Get("https://stream.wikimedia.org/v2/stream/recentchange")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Wikimedia stream: %v", err)
	}

	fmt.Printf("Stream response status: %s\n", resp.Status)

	// Declare A function Literal
	// Close the response body once done
	go func() {
		defer resp.Body.Close()

		// Create a scanner to read the stream line by line
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()

			// Skip lines that are not JSON (e.g., `event: message`)
			if len(line) < 6 || line[:5] != "data:" {
				// We skip non-JSON lines, such as `event: message`
				continue
			}

			// Remove the `data:` prefix to get the JSON part
			jsonData := strings.TrimPrefix(line, "data:")

			// Debugging: print the raw JSON
			//log.Printf("Raw JSON: %s\n", jsonData)

			// Attempt to decode the JSON part
			var event Event
			if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
				log.Println("Failed to decode JSON data:", err)
				continue // Skip invalid JSON and move to the next line
			}

			// Send the event to the channel for further processing
			changesChannel <- event
		}

		// If there is an error reading the stream, log it
		if err := scanner.Err(); err != nil {
			log.Println("Error reading stream:", err)
		}

		// Close the channel once the stream ends
		close(changesChannel)
	}()

	// Return the channel for the producer to consume
	return changesChannel, nil
}

// ProduceMessages TO TEST: kafka-console-consumer --bootstrap-server localhost:9092 --topic wikimedia-stream-golang
func ProduceMessages(topic string) {
	//topic := "wikimedia-stream-golang"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Get the stream data channel
	changesChannel, err := readStream()
	if err != nil {
		log.Fatal("Error reading stream:", err)
	}

	// Continuously read changes from the channel and send to Kafka
	for event := range changesChannel {
		// Create Kafka message from the recent change data
		messageValue := fmt.Sprintf("ID: %d, Type: %s, Title: %s, User: %s, Comment: %s",
			event.ID, event.Type, event.Title, event.User, event.Comment)

		// Write the message to Kafka
		_, writeErr := conn.WriteMessages(
			kafka.Message{Value: []byte(messageValue)},
		)
		if writeErr != nil {
			log.Fatal("failed to write message to Kafka:", err)
			continue // Move to the next message even if one fails to write
		}

		fmt.Printf("Sent message: %s\n", messageValue)
	}

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
