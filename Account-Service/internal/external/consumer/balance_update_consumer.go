package consumer

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"
)

// BalanceUpdateMessage represents the decoded Avro message
type BalanceUpdateMessage struct {
	AccountID string  `json:"accountid"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
}

var balanceUpdateAvroSchema = `
{
  "type": "record",
  "name": "BalanceUpdate",
  "fields": [
    {"name": "accountid", "type": "string"},
    {"name": "amount", "type": "double"},
    {"name": "type", "type": "string"}
  ]
}
`

// AccountService is the interface required to update account balances
type AccountService interface {
	UpdateBalance(ctx context.Context, accountID uuid.UUID, amount float64, txType string) error
}

// BalanceUpdateConsumer handles consuming and processing balance updates
type BalanceUpdateConsumer struct {
	kafkaReader   *kafka.Reader
	accountSvc    AccountService
	avroCodec     *goavro.Codec
	isRunning     bool
	stopChan      chan struct{}
}

// NewBalanceUpdateConsumer creates a new balance update consumer
func NewBalanceUpdateConsumer(brokers []string, topic, groupID string, accountSvc AccountService) (*BalanceUpdateConsumer, error) {
	// Create the Avro codec
	codec, err := goavro.NewCodec(balanceUpdateAvroSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create Avro codec: %w", err)
	}

	// Create Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,          // Make sure this is CONSISTENT
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset, // <-- CHANGE THIS to read from beginning
		MaxWait:     time.Second,
	})

	// Create topic if it doesn't exist
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers[0], topic, 0)
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v", err)
	} else {
		defer conn.Close()
		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		}
		err = conn.CreateTopics(topicConfigs...)
		if err != nil {
			log.Printf("Topic creation failed (may already exist): %v", err)
		}
	}

	return &BalanceUpdateConsumer{
		kafkaReader: reader,
		accountSvc:  accountSvc,
		avroCodec:   codec,
		stopChan:    make(chan struct{}),
	}, nil
}

// Start begins consuming messages
func (c *BalanceUpdateConsumer) Start(ctx context.Context) {
	if c.isRunning {
		return
	}

	c.isRunning = true
	go c.consumeMessages(ctx)
}

// Stop stops consuming messages
func (c *BalanceUpdateConsumer) Stop() {
	if !c.isRunning {
		return
	}

	c.isRunning = false
	close(c.stopChan)
	c.kafkaReader.Close()
}

// consumeMessages is the main message processing loop
func (c *BalanceUpdateConsumer) consumeMessages(ctx context.Context) {
	log.Println("Balance update consumer started")
	defer log.Println("Balance update consumer stopped")

	// Add a backoff counter
	var consecutiveTimeouts int

	for {
		select {
		case <-c.stopChan:
			return
		default:
			// Set a timeout for reading messages
			readCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			message, err := c.kafkaReader.ReadMessage(readCtx)
			cancel()

			if err != nil {
				if err != context.DeadlineExceeded {
					// Only log errors that aren't timeouts
					log.Printf("Error reading Kafka message: %v", err)
				} else {
					// For timeouts, log at debug level or with lower frequency
					if time.Now().Second()%60 == 0 { // Log once per minute
						log.Printf("No messages received in the last minute")
					}
					consecutiveTimeouts++
					// Exponential backoff: wait longer when repeatedly finding no messages
					backoffTime := time.Duration(math.Min(float64(consecutiveTimeouts*2), 30)) * time.Second
					time.Sleep(backoffTime)
				}
				continue
			}

			// Reset the backoff counter on successful read
			consecutiveTimeouts = 0

			// Process the message
			if err := c.processMessage(ctx, message); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}
}

// processMessage handles a single Kafka message
func (c *BalanceUpdateConsumer) processMessage(ctx context.Context, message kafka.Message) error {
	// Decode Avro message
	native, _, err := c.avroCodec.NativeFromBinary(message.Value)
	if err != nil {
		return fmt.Errorf("failed to decode Avro message: %w", err)
	}

	// Convert to map
	nativeMap, ok := native.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected message format")
	}

	// Create BalanceUpdateMessage for easier handling
	update := BalanceUpdateMessage{
		AccountID: nativeMap["accountid"].(string),
		Amount:    nativeMap["amount"].(float64),
		Type:      nativeMap["type"].(string),
	}

	// Add validation for empty transaction type
	if update.Type == "" {
		log.Printf("Warning: Received message with empty transaction type, defaulting to 'adjustment'")
		update.Type = "adjustment" // Provide a default
	}

	// Log more details about the received message for debugging
	log.Printf("Message details: %+v", nativeMap)

	log.Printf("Processing balance update: Account=%s, Amount=%.2f, Type=%s",
		update.AccountID, update.Amount, update.Type)

	// Parse the account ID
	accountID, err := uuid.Parse(update.AccountID)
	if err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	// Call the account service to update the balance
	err = c.accountSvc.UpdateBalance(ctx, accountID, update.Amount, update.Type)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	log.Printf("Successfully updated balance for account %s", update.AccountID)
	return nil
}

// ResetOffsets resets the Kafka reader's offsets
func (c *BalanceUpdateConsumer) ResetOffsets() error {
	return c.kafkaReader.SetOffset(kafka.FirstOffset)
}
