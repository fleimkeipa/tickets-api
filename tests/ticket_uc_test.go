package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/tickets-api/models"
	"github.com/fleimkeipa/tickets-api/repositories"
	"github.com/fleimkeipa/tickets-api/repositories/interfaces"
	"github.com/fleimkeipa/tickets-api/uc"
)

var testTicketRepo interfaces.TicketInterfaces

func init() {
	testTicketRepo = repositories.NewTicketRepository(test_db)
}

func TestTicketUC_Purchase(t *testing.T) {
	type fields struct {
		ticketRepo interfaces.TicketInterfaces
	}
	type args struct {
		ctx    context.Context
		id     string
		ticket *models.PurchaseRequest
	}
	type tempDatas struct {
		ticket []models.Ticket
	}
	tests := []struct {
		name      string
		tempDatas tempDatas
		fields    fields
		args      args
		want      *models.Ticket
		wantErr   bool
	}{
		{
			name: "",
			fields: fields{
				ticketRepo: testTicketRepo,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						Name:        "vangogh",
						Description: "vangogh ear",
						Allocation:  73,
					},
					{
						Name:        "pearl",
						Description: "girl with a pearl earing",
						Allocation:  68,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "1",
				ticket: &models.PurchaseRequest{
					UserID:   "344b6d2d-599a-4b23-b358-8f26512079a9",
					Quantity: 70,
				},
			},
			want:    &models.Ticket{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.tempDatas.ticket {
				if err := addTempData(&v); (err != nil) != tt.wantErr {
					t.Errorf("TicketUC.Purchase() addTempData error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			rc := uc.NewTicketUC(tt.fields.ticketRepo)
			got, err := rc.Purchase(tt.args.ctx, tt.args.id, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketUC.Purchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketUC.Purchase() = %v, want %v", got, tt.want)
			}
		})
	}
}
