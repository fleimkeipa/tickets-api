package uc

import (
	"context"
	"net/http"

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

// Create creates a new ticket based on the provided create request data.
func (rc *TicketUC) Create(ctx context.Context, request *models.CreateRequest) (*models.Ticket, error) {
	if err := rc.validator.Validate(request); err != nil {
		return nil, pkg.NewError(err, "failed to validate create request", http.StatusBadRequest)
	}

	ticket := models.Ticket{
		Name:        request.Name,
		Description: request.Description,
		Allocation:  request.Allocation,
	}

	t, err := rc.ticketRepo.Create(ctx, &ticket)
	if err != nil {
		return nil, pkg.NewError(err, "failed to create ticket", http.StatusInternalServerError)
	}

	return t, nil
}

// Purchase handles the purchasing of a ticket by the provided ticket ID.
func (rc *TicketUC) Purchase(ctx context.Context, ticketID string, request *models.PurchaseRequest) (*models.Ticket, error) {
	if err := rc.validator.Validate(request); err != nil {
		return nil, pkg.NewError(err, "failed to validate purchase request", http.StatusBadRequest)
	}

	// ticket exist control
	existTicket, err := rc.GetByID(ctx, ticketID)
	if err != nil {
		return nil, pkg.NewError(err, "failed to find ticket", http.StatusNotFound)
	}

	if existTicket.Allocation == 0 {
		return nil, pkg.NewError(err, "there is no available ticket now", http.StatusBadRequest)
	}

	if existTicket.Allocation < request.Quantity {
		return nil, pkg.NewError(err, "cannot afford this quantity", http.StatusBadRequest)
	}

	existTicket.Allocation -= request.Quantity

	t, err := rc.ticketRepo.Update(ctx, existTicket)
	if err != nil {
		return nil, pkg.NewError(err, "failed to update ticket", http.StatusInternalServerError)
	}

	return t, nil
}

// GetByID retrieves a ticket by the provided ticket ID.
func (rc *TicketUC) GetByID(ctx context.Context, ticketID string) (*models.Ticket, error) {
	t, err := rc.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, pkg.NewError(err, "failed to find ticket", http.StatusNotFound)
	}

	return t, nil
}
