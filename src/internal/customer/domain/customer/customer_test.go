package customer

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"time"
)

func TestRegisterNewCustomer(t *testing.T) {
	t.Parallel()

	type args struct {
		id           string
		nik          string
		fullname     string
		legalName    string
		placeOfBirth string
		dateOfBirth  time.Time
		wage         float64
		photoURL     string
		kTPURL       string
		processDate  time.Time
	}

	fc := FactoryConfig{
		WageLimit:  3500000,
		MinimumAge: 20,
	}

	tests := []struct {
		name    string
		config  FactoryConfig
		args    args
		want    *Customer
		wantErr error
	}{
		{
			name:   "age requirement not meet",
			config: fc,
			args: args{
				id:           "customerid-1",
				nik:          "666-999-111",
				fullname:     "john doe",
				legalName:    "john doe",
				placeOfBirth: "tangerang",
				dateOfBirth:  time.Date(2020, time.April, 1, 1, 0, 0, 0, time.UTC),
				wage:         4000000,
				photoURL:     "https://imgurl.com/photo",
				kTPURL:       "https://imgurl.com/ktp",
				processDate:  time.Date(2025, time.April, 1, 1, 0, 0, 0, time.UTC),
			},
			want:    nil,
			wantErr: AgeRequirementNotMetError{RequiredMinimumAge: 20, ApplicantAge: 5},
		},
		{
			name:   "wage requirement not meet",
			config: fc,
			args: args{
				id:           "customerid-1",
				nik:          "666-999-111",
				fullname:     "john doe",
				legalName:    "john doe",
				placeOfBirth: "tangerang",
				dateOfBirth:  time.Date(1998, time.April, 1, 1, 0, 0, 0, time.UTC),
				wage:         2000000,
				photoURL:     "https://imgurl.com/photo",
				kTPURL:       "https://imgurl.com/ktp",
				processDate:  time.Date(2025, time.April, 1, 1, 0, 0, 0, time.UTC),
			},
			want:    nil,
			wantErr: WageRequirementNotMetError{RequiredMinimumWage: 3500000, ApplicantWage: 2000000},
		},
		{
			name:   "valid",
			config: fc,
			args: args{
				id:           "customerid-1",
				nik:          "666-999-111",
				fullname:     "john doe",
				legalName:    "john doe",
				placeOfBirth: "tangerang",
				dateOfBirth:  time.Date(1998, time.April, 1, 1, 0, 0, 0, time.UTC),
				wage:         3500000,
				photoURL:     "https://imgurl.com/photo",
				kTPURL:       "https://imgurl.com/ktp",
				processDate:  time.Date(2025, time.April, 1, 1, 0, 0, 0, time.UTC),
			},
			want: &Customer{
				ID:        "customerid-1",
				NIK:       "666-999-111",
				Fullname:  "john doe",
				LegalName: "john doe",
				PlaceAndDateOfBirth: PlaceAndDateOfBirth{
					Place: "tangerang",
					Date:  time.Date(1998, time.April, 1, 1, 0, 0, 0, time.UTC),
				},
				Wage:     3500000,
				PhotoURL: "https://imgurl.com/photo",
				KTPURL:   "https://imgurl.com/ktp",
				CreateAt: time.Date(2025, time.April, 1, 1, 0, 0, 0, time.UTC),
				UpdateAt: time.Date(2025, time.April, 1, 1, 0, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := NewFactory(fc)

			got, err := f.RegisterNewCustomer(
				tt.args.id,
				tt.args.nik,
				tt.args.fullname,
				tt.args.legalName,
				tt.args.placeOfBirth,
				tt.args.dateOfBirth,
				tt.args.wage,
				tt.args.photoURL,
				tt.args.kTPURL,
				tt.args.processDate,
			)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)

		})

	}

}
