/**
 * helper functions related to maps.
 *
**/

package common

import (
	"lp_customer_portal/models"
	"testing"
)

func TestGetZone(t *testing.T) {
	type args struct {
		lat   float64
		long  float64
		zones []models.Zone
	}
	tests := []struct {
		name    string
		args    args
		want    models.Zone
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test case for lat long inside zone",
			args: args{
				lat:  2,
				long: 2,
				zones: []models.Zone{
					{
						Name: "zone1",
						ZoneBoundaries: []models.ZoneBoundary{
							{Lat: 1, Long: 1},
							{Lat: 1, Long: 3},
							{Lat: 3, Long: 3},
							{Lat: 3, Long: 1},
						},
					},
				},
			},
			want: models.Zone{
				Name: "zone1",
			},
			wantErr: false,
		},
		{
			name: "test case for lat long when outside zone",
			args: args{
				lat:  5,
				long: 2,
				zones: []models.Zone{
					{
						Name: "zone1",
						ZoneBoundaries: []models.ZoneBoundary{
							{Lat: 1, Long: 1},
							{Lat: 1, Long: 3},
							{Lat: 3, Long: 3},
							{Lat: 3, Long: 1},
						},
					},
				},
			},
			want:    models.Zone{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetZone(tt.args.lat, tt.args.long, tt.args.zones)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want.Name != got.Name {
				t.Errorf("GetZone() = %v, want %v", got, tt.want)
			}
		})
	}
}
