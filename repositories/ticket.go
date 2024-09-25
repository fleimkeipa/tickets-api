package repositories

import (
	"context"
	"fmt"

	"github.com/fleimkeipa/tickets-api/models"

	"github.com/go-pg/pg"
)

type TicketRepository struct {
	db *pg.DB
}

func NewTicketRepository(db *pg.DB) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (rc *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	_, err := rc.db.Model(ticket).Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

func (rc *TicketRepository) Update(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	_, err := rc.db.Model(ticket).WherePK().Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	return ticket, nil
}

func (rc *TicketRepository) GetByID(ctx context.Context, id string) (*models.Ticket, error) {
	var ticket = new(models.Ticket)

	err := rc.db.
		Model(ticket).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket [%s] id, error: %w", id, err)
	}

	return ticket, nil
}
