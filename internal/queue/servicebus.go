package queue

import (
	"context"
	"log"
	"time"

	azservicebus "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// ServiceBusPublisher manages the connection and publishing to Azure Service Bus.
type ServiceBusPublisher struct {
	sender *azservicebus.Sender
	queueName string
}

// NewServiceBusPublisher initializes the Service Bus client and sender.
func NewServiceBusPublisher(connString, queueName string) (*ServiceBusPublisher, error) {
	log.Printf("Connecting to Azure Service Bus queue: %s", queueName)
	
	// Create a client
	client, err := azservicebus.NewClientFromConnectionString(connString, nil)
	if err != nil {
		return nil, err
	}

	// Create a sender
	sender, err := client.NewSender(queueName, nil)
	if err != nil {
		return nil, err
	}

	return &ServiceBusPublisher{
		sender: sender,
		queueName: queueName,
	}, nil
}

// Publish sends the raw payload body to the configured Service Bus queue.
func (p *ServiceBusPublisher) Publish(ctx context.Context, payloadBody []byte) error {
	
	// The message body is the raw GitHub JSON payload
	message := azservicebus.NewMessage(payloadBody)

	// Add a short timeout to the context to prevent blocking forever
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.sender.SendMessage(ctx, message, nil); err != nil {
		return err
	}
	return nil
}

// Close gracefully closes the Service Bus sender connection.
func (p *ServiceBusPublisher) Close(ctx context.Context) {
	p.sender.Close(ctx)
}
