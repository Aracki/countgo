package db

import (
	"fmt"
	"testing"
)

func TestDatabase_GetMostFrequentVisitors(t *testing.T) {
	type fields struct {
		dbconfig Conf
	}
	tests := []struct {
		name    string
		fields  fields
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "ballab",
			fields:  struct{ dbconfig Conf }{dbconfig: db.dbconfig},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := Database{
				dbconfig: tt.fields.dbconfig,
			}
			got, err := db.GetMostFrequentVisitors()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetMostFrequentVisitors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
