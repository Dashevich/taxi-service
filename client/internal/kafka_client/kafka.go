package kafka_client

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func SendToTopic(conn *kafka.Conn, message []byte) error {
	_, err := conn.WriteMessages(
		kafka.Message{Value: message},
	)
	return err
}

func ReadFromTopic(conn *kafka.Conn) ([]byte, error) {
	b := make([]byte, 10e3)
	n, err := conn.Read(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}

func ConnectKafka(ctx context.Context, address string, topic string, partition int) (*kafka.Conn, error) {
	return kafka.DialLeader(ctx, "tcp", address, topic, partition)
}
