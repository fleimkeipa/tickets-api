package interfaces

import (
	"context"

	"github.com/fleimkeipa/tickets-api/models"
)

type TicketInterfaces interface {
	Create(context.Context, *models.Ticket) (*models.Ticket, error)
	Update(context.Context, *models.Ticket) (*models.Ticket, error)
	GetByID(context.Context, string) (*models.Ticket, error)
}
