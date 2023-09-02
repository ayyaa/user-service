package repository

import (
	"context"
	"database/sql"
	"testing"
)

func TestRepository_Insert(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		user UserReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				DB: tt.fields.DB,
			}
			gotId, err := r.Insert(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("Repository.Insert() = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}
