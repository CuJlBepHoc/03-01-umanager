package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/google/uuid"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database/users"
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

	// Очистка таблицы перед созданием новой записи для избежания дубликатов
	log.Println("Clearing users table...")
	if err := e.UsersRepository.ClearTable(ctx); err != nil {
		return fmt.Errorf("clear table: %w", err)
	}

	create, err := e.UsersRepository.Create(
		ctx, users.CreateUserReq{
			ID:       uuid.New(),
			Username: "random",
			Password: "password",
		},
	)
	if err != nil {
		return err
	}

	found, err := e.UsersRepository.FindByID(ctx, create.ID)
	if err != nil {
		return err
	}

	foundBy, err := e.UsersRepository.FindByUsername(ctx, "random")
	if err != nil {
		return err
	}

	// Улучшенный форматированный вывод
	fmt.Println("Created user:")
	printJSON(create)

	fmt.Println("Found user by ID:")
	printJSON(found)

	fmt.Println("Found user by username:")
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
