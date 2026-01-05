package mongodb

import (
	"fmt"
	"testing"
)

func TestDatabase_GetMostFrequentVisitors(t *testing.T) {
	tests := []struct {
		name    string
		db      *Database
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "ballab",
			db:      db,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.GetMostFrequentVisitors()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetMostFrequentVisitors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
