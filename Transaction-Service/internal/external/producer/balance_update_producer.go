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
	codec, err := goavro.NewCodec(balanceUpdateAvroSchema)
	if err != nil {
		log.Println("Failed to create Avro codec:", err)
		return err
	}

	native := map[string]interface{}{
		"accountid": accountID.String(),
		"amount":    amount,
		"type":      txType,
	}

	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		log.Println("Failed to encode Avro message:", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(accountID.String()),
		Value: binary,
	}

	return writer.WriteMessages(context.Background(), msg)
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