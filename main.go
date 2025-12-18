package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const prefix = "/api/v1/image"

func worker() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "image-topic",
		GroupID: "group0",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Ошибка при получении:", err)
		}

		msgStr := string(msg.Value)

		offset := strings.Index(msgStr, "---END FILE NAME---")

		fmt.Println("Получен файл из Kafka: " + msgStr[0:offset])
	}
}

func handler() {
	http.HandleFunc(prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("file")

		if err != nil {
			http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
			return
		}
		defer file.Close()

		ctx := context.Background()

		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "image-topic",
		})
		defer writer.Close()

		fileBytes, err := io.ReadAll(file)

		bytesResult := append([]byte(handler.Filename), "---END FILE NAME---"...)
		bytesResult = append(bytesResult, fileBytes...)

		err = writer.WriteMessages(ctx, kafka.Message{
			Value: bytesResult,
		})
		if err != nil {
			log.Fatal("Ошибка при отправке:", err)
		}
	})

	log.Fatal(http.ListenAndServe(":5267", nil))
}

func main() {
	var wg sync.WaitGroup

	wg.Go(handler)
	wg.Go(worker)

	wg.Wait()
}
