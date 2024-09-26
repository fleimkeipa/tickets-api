package uc

import (
	"context"
	"errors"
	"fmt"

	"github.com/fleimkeipa/tickets-api/models"
	"github.com/fleimkeipa/tickets-api/pkg"
	"github.com/fleimkeipa/tickets-api/repositories/interfaces"
)

type TicketUC struct {
	ticketRepo interfaces.TicketInterfaces
	validator  *pkg.CustomValidator
}

func NewTicketUC(ticketRepo interfaces.TicketInterfaces, validator *pkg.CustomValidator) *TicketUC {
	return &TicketUC{
		ticketRepo: ticketRepo,
		validator:  validator,
	}
}

func (rc *TicketUC) Create(ctx context.Context, request *models.CreateRequest) (*models.Ticket, error) {
	if err := rc.validator.Validate(request); err != nil {
		return nil, fmt.Errorf("failed to validate create request: %w", err)
	}

	ticket := models.Ticket{
		Name:        request.Name,
		Description: request.Description,
		Allocation:  request.Allocation,
	}

	return rc.ticketRepo.Create(ctx, &ticket)
}

func (rc *TicketUC) Purchase(ctx context.Context, id string, request *models.PurchaseRequest) (*models.Ticket, error) {
	if err := rc.validator.Validate(request); err != nil {
		return nil, fmt.Errorf("failed to validate purchase request: %w", err)
	}

	// ticket exist control
	existTicket, err := rc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existTicket.Allocation == 0 {
		return nil, errors.New("there is no available ticket now")
	}

	if existTicket.Allocation < request.Quantity {
		return nil, errors.New("cannot afford this quantity")
	}

	existTicket.Allocation -= request.Quantity

	return rc.ticketRepo.Update(ctx, existTicket)
}

func (rc *TicketUC) GetByID(ctx context.Context, id string) (*models.Ticket, error) {
	return rc.ticketRepo.GetByID(ctx, id)
}
