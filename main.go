package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	_ "github.com/spf13/viper"
)

var wg sync.WaitGroup

type Config struct {
	Broker string
}

func GetConfig() *Config {
	return &Config{Broker: "localhost:"}
}

func main() {
	connStr := "postgresql://postgres:postgres@localhost:5432?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(db)
}

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func ProduceMessage() {
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!!!")},
		kafka.Message{Value: []byte("two!!!")},
		kafka.Message{Value: []byte("three!!!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func ConsumeMessage() {

	topic := "my-topic"

	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message

	for {
		_, err := batch.Read(b)
		if err != nil {
			fmt.Println(err)
			break
		}

		teste := Teste{}

		_ = json.Unmarshal(b, &teste)
		fmt.Println(teste)
	}

	if err := batch.Close(); err != nil {
		fmt.Println("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		fmt.Println("failed to close connection:", err)
	}
}

type Teste struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func Write() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "my-topic",
		Balancer: &kafka.LeastBytes{},
	}

	msg, _ := json.Marshal(Teste{
		Nome:  "Erick",
		Email: "erick@gmail.com",
	})

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(msg),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func Reader() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "my-topic",
		GroupID: "my-group",
	})

	fmt.Println("Iniciando Kafka reader")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		teste := Teste{}

		err = json.Unmarshal(m.Value, &teste)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), teste.Nome)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
