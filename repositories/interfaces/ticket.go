package interfaces

import (
	"context"

	"github.com/fleimkeipa/tickets-api/models"
)

type TicketInterfaces interface {
	Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
	Update(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
	GetByID(ctx context.Context, ticketID string) (*models.Ticket, error)
}
