package controller

import (
	"fmt"
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
//	@Param			body			body		models.CreateRequest	true	"Ticket creation input"
//	@Success		201				{object}	models.Ticket			"Created ticket details"
//	@Failure		400				{object}	models.FailureResponse	"Error message including details on failure"
//	@Router			/tickets [post]
func (rc *TicketHandler) CreateTicket(c echo.Context) error {
	var request models.CreateRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.FailureResponse{
			Error:   fmt.Sprintf("Failed to validate: %v", err),
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	ticket, err := rc.ticketUC.Create(c.Request().Context(), &request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.FailureResponse{
			Error:   fmt.Sprintf("Failed to create ticket: %v", err),
			Message: "Ticket creation failed. Please check the provided details and try again.",
		})
	}

	return c.JSON(http.StatusCreated, ticket)
}

// PurchaseTicket godoc
//
//	@Summary		PurchaseTicket purchases a new ticket
//	@Description	This endpoint purchases a new ticket by providing id and quantity.
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body	models.PurchaseRequest	true	"Ticket purchase input"
//	@Success		204				"Purchase successful, no content"
//	@Failure		400				{object}	models.FailureResponse	"Error message including details on failure"
//	@Router			/tickets/{id}/purchase [post]
func (rc *TicketHandler) PurchaseTicket(c echo.Context) error {
	var id = c.Param("id")

	var request models.PurchaseRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.FailureResponse{
			Error:   fmt.Sprintf("Failed to validate: %v", err),
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	_, err := rc.ticketUC.Purchase(c.Request().Context(), id, &request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.FailureResponse{
			Error:   fmt.Sprintf("Failed to purchase ticket: %v", err),
			Message: "Ticket purchae failed. Please check the provided details and try again.",
		})
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
	var id = c.Param("id")

	ticket, err := rc.ticketUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve ticket: %v", err),
			Message: "Error fetching the ticket details. Please verify the ticket name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, ticket)
}
