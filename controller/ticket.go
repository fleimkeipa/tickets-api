package controller

import (
	"net/http"

	"github.com/fleimkeipa/tickets-api/models"
	"github.com/fleimkeipa/tickets-api/uc"

	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	ticketUC *uc.TicketUC
}

func NewTicketHandler(ticketUC *uc.TicketUC) *TicketHandler {
	return &TicketHandler{
		ticketUC: ticketUC,
	}
}

// CreateTicket godoc
//
//	@Summary		CreateTicket creates a new ticket
//	@Description	This endpoint creates a new ticket by providing name, description, and allocation.
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			{object}	models.CreateRequest	true	"Ticket creation input"
//	@Success		201				{object}	models.TicketResponse			"Created ticket details"
//	@Failure		400				{object}	models.FailureResponse	"Error message including details on failure"
//	@Router			/tickets [post]
func (rc *TicketHandler) CreateTicket(c echo.Context) error {
	var request models.CreateRequest

	if err := c.Bind(&request); err != nil {
		return HandleEchoError(c, err)
	}

	ticket, err := rc.ticketUC.Create(c.Request().Context(), &request)
	if err != nil {
		return HandleEchoError(c, err)
	}

	response := fillTicketResponse(ticket)

	return c.JSON(http.StatusCreated, response)
}

// PurchaseTicket godoc
//
//	@Summary		PurchaseTicket purchases a new ticket
//	@Description	This endpoint purchases a new ticket by providing id and quantity.
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			{object}	models.PurchaseRequest	true	"Ticket purchase input"
//	@Success		204				"Purchase successful, no content"
//	@Failure		400				{object}	models.FailureResponse	"Error message including details on failure"
//	@Router			/tickets/{id}/purchase [post]
func (rc *TicketHandler) PurchaseTicket(c echo.Context) error {
	id := c.Param("id")

	var request models.PurchaseRequest
	if err := c.Bind(&request); err != nil {
		return HandleEchoError(c, err)
	}

	_, err := rc.ticketUC.Purchase(c.Request().Context(), id, &request)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetByID godoc
// GetByID godoc
//
//	@Summary		Get a ticket by ID
//	@Description	Retrieves a ticket from the database by its ID.
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			id				path		string					true	"ID of the ticket"
//	@Success		200				{object}	models.Ticket			"Details of the requested ticket"
//	@Failure		400				{object}	models.FailureResponse	"Error message including details on failure"
//	@Router			/tickets/{id} [get]
func (rc *TicketHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	ticket, err := rc.ticketUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return HandleEchoError(c, err)
	}

	response := fillTicketResponse(ticket)

	return c.JSON(http.StatusOK, response)
}

func fillTicketResponse(ticket *models.Ticket) *models.TicketResponse {
	if ticket == nil {
		return &models.TicketResponse{}
	}

	return &models.TicketResponse{
		ID:          ticket.ID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Allocation:  ticket.Allocation,
	}
}
