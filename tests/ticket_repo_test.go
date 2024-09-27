package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/tickets-api/models"
	"github.com/fleimkeipa/tickets-api/pkg"
	"github.com/fleimkeipa/tickets-api/repositories"

	"github.com/go-pg/pg"
)

func TestTicketRepository_Create(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx    context.Context
		ticket *models.Ticket
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
				db: test_db,
			},
			args: args{
				ctx: context.TODO(),
				ticket: &models.Ticket{
					ID:          1,
					Name:        "batman",
					Description: "batman returns",
					Allocation:  100,
				},
			},
			want: &models.Ticket{
				ID:          1,
				Name:        "batman",
				Description: "batman returns",
				Allocation:  100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := repositories.NewTicketRepository(tt.fields.db)
			got, err := rc.Create(tt.args.ctx, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketRepository.Create() = %v, want %v", got, tt.want)
			}
			if err := clearTable(); err != nil {
				t.Errorf("TicketRepository.Create() clearTable error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTicketRepository_Update(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx    context.Context
		ticket *models.Ticket
	}
	type tempData struct {
		ticket *models.Ticket
	}
	tests := []struct {
		name     string
		fields   fields
		tempData tempData
		args     args
		want     *models.Ticket
		wantErr  bool
	}{
		{
			name: "success",
			fields: fields{
				db: test_db,
			},
			tempData: tempData{
				ticket: &models.Ticket{
					ID:          1,
					Name:        "joker",
					Description: "joker down",
					Allocation:  100,
				},
			},
			args: args{
				ctx: context.TODO(),
				ticket: &models.Ticket{
					ID:          1,
					Name:        "joker",
					Description: "joker up",
					Allocation:  99,
				},
			},
			want: &models.Ticket{
				ID:          1,
				Name:        "joker",
				Description: "joker up",
				Allocation:  99,
			},
			wantErr: false,
		},
		{
			name: "error - updating a non-existent ticket",
			fields: fields{
				db: test_db,
			},
			tempData: tempData{
				ticket: &models.Ticket{
					ID:          1,
					Name:        "joker",
					Description: "joker up",
					Allocation:  99,
				},
			},
			args: args{
				ctx: context.TODO(),
				ticket: &models.Ticket{
					ID:          3,
					Name:        "joker",
					Description: "joker up",
					Allocation:  99,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addTempData(tt.tempData.ticket); err != nil {
				t.Errorf("TicketRepository.Update() addTempData error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			rc := repositories.NewTicketRepository(tt.fields.db)
			got, err := rc.Update(tt.args.ctx, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketRepository.Update() = %v, want %v", got, tt.want)
			}
			if err := clearTable(); err != nil {
				t.Errorf("TicketRepository.Update() clearTable error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTicketRepository_GetByID(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx context.Context
		id  string
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
			name: "success",
			fields: fields{
				db: test_db,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "devil",
						Description: "devil may cry",
						Allocation:  100,
					},
					{
						ID:          2,
						Name:        "wanted",
						Description: "wanted follows you",
						Allocation:  23,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "1",
			},
			want: &models.Ticket{
				ID:          1,
				Name:        "devil",
				Description: "devil may cry",
				Allocation:  100,
			},
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				db: test_db,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "devil",
						Description: "devil may cry",
						Allocation:  100,
					},
					{
						ID:          2,
						Name:        "wanted",
						Description: "wanted follows you",
						Allocation:  23,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "2",
			},
			want: &models.Ticket{
				ID:          2,
				Name:        "wanted",
				Description: "wanted follows you",
				Allocation:  23,
			},
			wantErr: false,
		},
		{
			name: "error - finding a non-existent ticket",
			fields: fields{
				db: test_db,
			},
			tempDatas: tempDatas{
				ticket: []models.Ticket{
					{
						ID:          1,
						Name:        "devil",
						Description: "devil may cry",
						Allocation:  100,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  "2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.tempDatas.ticket {
				if err := addTempData(&v); err != nil {
					t.Errorf("TicketRepository.GetByID() addTempData error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			rc := repositories.NewTicketRepository(tt.fields.db)
			got, err := rc.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketRepository.GetByID() = %v, want %v", got, tt.want)
			}
			if err := clearTable(); err != nil {
				t.Errorf("TicketRepository.GetByID() clearTable error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
