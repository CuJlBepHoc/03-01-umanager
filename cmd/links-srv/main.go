package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database/links"
	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/env"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := runMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func runMain(ctx context.Context) error {
	e, err := env.Setup(ctx)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	_ = e

	create, err := e.LinksRepository.Create(
		ctx, links.CreateReq{
			ID:     primitive.NewObjectID(),
			URL:    "https://ya.ru",
			Title:  "ya main page",
			Tags:   []string{"search", "yandex"},
			Images: []string{},
			UserID: "uuid", // created user id
		},
	)
	if err != nil {
		return err
	}

	found, err := e.LinksRepository.FindByUserAndURL(ctx, "https://ya.ru", "uuid")
	if err != nil {
		return err
	}

	foundBy, err := e.LinksRepository.FindByCriteria(
		ctx, links.Criteria{
			Tags: []string{"yandex"},
		},
	)
	if err != nil {
		return err
	}

	// Улучшенный форматированный вывод
	fmt.Println("Created link:")
	printJSON(create)

	fmt.Println("Found link by URL and UserID:")
	printJSON(found)

	fmt.Println("Found links by criteria:")
	printJSON(foundBy)

	return nil
}

// printJSON выводит объект в формате JSON с отступами
func printJSON(v interface{}) {
	encoded, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	fmt.Println(string(encoded))
}
