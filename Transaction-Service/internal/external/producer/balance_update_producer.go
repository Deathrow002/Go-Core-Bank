package external

import (
	"context"
	"log"
	"net/http"
	"transaction-service/internal/transaction/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"
)

type BalanceUpdateRequest struct {
	AccountID	uuid.UUID				`json:"accountid"`
	Amount		float64					`json:"amount"`
	Type		models.TransactionType	`json:"type"`
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

func SendBalanceUpdate(writer *kafka.Writer, accountID uuid.UUID, amount float64, txType string) error {
	// Log input parameters
	log.Printf("SendBalanceUpdate called with accountID=%s, amount=%.2f, txType='%s'", 
		accountID.String(), amount, txType)
	
	codec, err := goavro.NewCodec(balanceUpdateAvroSchema)
	if err != nil {
		log.Printf("Failed to create Avro codec: %v", err)
		return err
	}
	log.Printf("Avro codec created successfully with schema: %s", balanceUpdateAvroSchema)

	// Validate transaction type
	if txType == "" {
		log.Printf("WARNING: Empty transaction type received, using 'adjustment' as default")
		txType = "adjustment"
	}

	native := map[string]interface{}{
		"accountid": accountID.String(),
		"amount":    amount,
		"type":      txType,
	}
	log.Printf("Created native map for Avro: %+v", native)

	// Validate all fields are present
	for field, value := range native {
		log.Printf("Field '%s' = %v (type: %T)", field, value, value)
	}

	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		log.Printf("Failed to encode Avro message: %v", err)
		return err
	}
	log.Printf("Successfully encoded Avro message, binary size: %d bytes", len(binary))

	// Decode back to verify encoding worked correctly
	decoded, _, err := codec.NativeFromBinary(binary)
	if err != nil {
		log.Printf("Warning: Unable to decode the message we just encoded: %v", err)
	} else {
		log.Printf("Verification - decoded message: %+v", decoded)
	}

	msg := kafka.Message{
		Key:   []byte(accountID.String()),
		Value: binary,
	}
	log.Printf("Created Kafka message with key: %s", accountID.String())

	// Log Kafka writer configuration
	log.Printf("Kafka writer configured for brokers: %v, topic: %s", 
		writer.Addr, writer.Topic)

	// Write the message
	log.Printf("Attempting to write message to Kafka...")
	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}
	
	log.Printf("Successfully sent balance update to Kafka for account %s", accountID.String())
	return nil
}

func SendBalanceUpdateToKafka(writer *kafka.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BalanceUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := SendBalanceUpdate(writer, req.AccountID, req.Amount, string(req.Type)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send update", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "balance update sent"})
	}
}