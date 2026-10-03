package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/config"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	notification "github.com/abubakar-sadiq001/learning-go-shop/internal/notifications"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/providers"
)

func main() {
	fmt.Print("Starting notification service!")

	// 1. Parent context
	ctx := context.Background()

	// 2. Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 3. Initialize email notifier object
	emailConfig := notification.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}

	emailNotifier := notification.NewEmailNotifier(&emailConfig)

	// 4. Create config for AWS sqs
	awsConfig, err := providers.CreateAWSConfig(ctx, cfg.AWS.S3Endpoint, cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Failed to create AWS config: %v", err)
	}

	// 5. Create SQS subscriber
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{AWSConfig: awsConfig}, logger)
	if err != nil {
		log.Fatal("")
	}

	// 6. Subscribe to messages
	messages, err := subscriber.Subscribe(ctx, cfg.AWS.EventQueueName)

	// 7. Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-messages:
			if err := processMessages(msg, emailNotifier); err != nil {
				log.Printf("Error processing message: %d", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case <-sigChan:
			log.Panicln("Shutting down notification service...")
			subscriber.Close()
		}
	}

}

func processMessages(msg *message.Message, emailNotifier *notification.EmailNotifier) error {
	eventType := msg.Metadata.Get("event_type")
	switch eventType {
	case notification.UserLoggedIn:
		return handleUserLoggedIn(msg, emailNotifier)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func handleUserLoggedIn(msg *message.Message, emailNotifier *notification.EmailNotifier) error {
	var user models.User
	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return err
	}

	userName := user.FirstName + " " + user.LastName

	if userName == "" {
		userName = "User"
	}

	log.Printf("Sending login notification to %s", user.Email)

	return emailNotifier.SendLoginNotification(user.Email, userName)
}
