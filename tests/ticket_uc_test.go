package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/tickets-api/models"
	"github.com/fleimkeipa/tickets-api/pkg"
	"github.com/fleimkeipa/tickets-api/repositories"
	"github.com/fleimkeipa/tickets-api/repositories/interfaces"
	"github.com/fleimkeipa/tickets-api/uc"
)

var testTicketValidator *pkg.CustomValidator

func init() {
	testTicketValidator = pkg.NewValidator()
}

func TestTicketUC_Create(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	testTicketRepo := repositories.NewTicketRepository(test_db)
	type fields struct {
		ticketRepo interfaces.TicketInterfaces
		validator  *pkg.CustomValidator
	}
	type args struct {
		ctx     context.Context
		request *models.CreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Ticket
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			args: args{
				ctx: context.TODO(),
				request: &models.CreateRequest{
					Name:        "spiderman",
					Description: "spiderman homecoming",
					Allocation:  23,
				},
			},
			want: &models.Ticket{
				ID:          1,
				Name:        "spiderman",
				Description: "spiderman homecoming",
				Allocation:  23,
			},
			wantErr: false,
		},
		{
			name: "error - invalid allocation value",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			args: args{
				ctx: context.TODO(),
				request: &models.CreateRequest{
					Name:        "superman",
					Description: "superman returns",
					Allocation:  -10,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - missing required fields",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			args: args{
				ctx: context.TODO(),
				request: &models.CreateRequest{
					Description: "nameless ticket",
					Allocation:  23,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := uc.NewTicketUC(tt.fields.ticketRepo, tt.fields.validator)
			got, err := rc.Create(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketUC.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketUC.Create() = %v, want %v", got, tt.want)
			}
			if err := clearTable(); err != nil {
				t.Errorf("TicketUC.Create() clearTable error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTicketUC_Purchase(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	testTicketRepo := repositories.NewTicketRepository(test_db)
	type fields struct {
		ticketRepo interfaces.TicketInterfaces
		validator  *pkg.CustomValidator
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
			name: "success - corret quantity",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "vangogh",
						Description: "vangogh ear",
						Allocation:  73,
					},
					{
						ID:          2,
						Name:        "pearl",
						Description: "girl with a pearl earring",
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
			want: &models.Ticket{
				ID:          1,
				Name:        "vangogh",
				Description: "vangogh ear",
				Allocation:  3,
			},
			wantErr: false,
		},
		{
			name: "error - quantity exceeds allocation",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "vangogh",
						Description: "vangogh ear",
						Allocation:  73,
					},
					{
						ID:          2,
						Name:        "pearl",
						Description: "girl with a pearl earring",
						Allocation:  68,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "2",
				ticket: &models.PurchaseRequest{
					UserID:   "344b6d2d-599a-4b23-b358-8f26512079a9",
					Quantity: 70,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - buying on zero allocation ticket",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "arrival",
						Description: "never arrival",
						Allocation:  0,
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
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - updating a non-existent ticket",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "erik",
						Description: "erik satie gnossienne no:1",
						Allocation:  10,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "2",
				ticket: &models.PurchaseRequest{
					UserID:   "344b6d2d-599a-4b23-b358-8f26512079a9",
					Quantity: 1,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - failed validation on user id",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "starrynight",
						Description: "starry night sky",
						Allocation:  50,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "1",
				ticket: &models.PurchaseRequest{
					UserID:   "", // Invalid UserID
					Quantity: 10,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - failed validation on negative quantity",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "starrynight",
						Description: "starry night sky",
						Allocation:  50,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "1",
				ticket: &models.PurchaseRequest{
					UserID:   "344b6d2d-599a-4b23-b358-8f26512079a9", // Invalid UserID
					Quantity: -10,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - failed validation on negative quantity",
			fields: fields{
				ticketRepo: testTicketRepo,
				validator:  testTicketValidator,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "starrynight",
						Description: "starry night sky",
						Allocation:  50,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "1",
				ticket: &models.PurchaseRequest{
					UserID:   "344b6d2d-599a-4b23-b358-8f26512079a9", // Invalid UserID
					Quantity: -10,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.tempDatas.ticket {
				if err := addTempData(&v); err != nil {
					t.Errorf("TicketUC.Purchase() addTempData error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			rc := uc.NewTicketUC(tt.fields.ticketRepo, tt.fields.validator)
			got, err := rc.Purchase(tt.args.ctx, tt.args.id, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketUC.Purchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketUC.Purchase() = %v, want %v", got, tt.want)
			}
			if err := clearTable(); err != nil {
				t.Errorf("TicketUC.Purchase() clearTable error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
