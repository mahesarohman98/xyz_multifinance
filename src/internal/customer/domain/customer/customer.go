package customer

import (
	"fmt"
	"time"
)

type PlaceAndDateOfBirth struct {
	Place string
	Date  time.Time
}

type Customer struct {
	ID                  string
	NIK                 string
	Fullname            string
	LegalName           string
	PlaceAndDateOfBirth PlaceAndDateOfBirth
	Wage                float64
	PhotoURL            string
	KTPURL              string
	CreateAt            time.Time
	UpdateAt            time.Time
}

type Factory struct {
	fc FactoryConfig
}

type FactoryConfig struct {
	WageLimit  float64
	MinimumAge int
}

func NewFactory(fc FactoryConfig) Factory {
	return Factory{fc: fc}
}

type WageRequirementNotMetError struct {
	RequiredMinimumWage float64
	ApplicantWage       float64
}

func (e WageRequirementNotMetError) Error() string {
	return fmt.Sprintf(
		"applicant does not meet the minimum wage requirement: required %.2f, provided %.2f",
		e.RequiredMinimumWage, e.ApplicantWage)
}

func (f *Factory) validateMinimumWage(wage float64) error {
	if f.fc.WageLimit > wage {
		return WageRequirementNotMetError{f.fc.WageLimit, wage}
	}
	return nil
}

type AgeRequirementNotMetError struct {
	RequiredMinimumAge int
	ApplicantAge       int
}

func (e AgeRequirementNotMetError) Error() string {
	return fmt.Sprintf(
		"applicant does not meet the minimum age requirement: required %d years, provided %d years",
		e.RequiredMinimumAge, e.ApplicantAge)
}

func (f *Factory) validateMinimumAge(birthDate time.Time, today time.Time) error {
	age := today.Year() - birthDate.Year()
	if today.YearDay() > birthDate.YearDay() {
		age--
	}

	if age < f.fc.MinimumAge {
		return AgeRequirementNotMetError{f.fc.MinimumAge, age}
	}
	return nil
}

func (f *Factory) RegisterNewCustomer(
	id string,
	nik string,
	fullname string,
	legalName string,
	placeOfBirth string,
	dateOfBirth time.Time,
	wage float64,
	photoURL string,
	kTPURL string,
	processDate time.Time,
) (*Customer, error) {
	if err := f.validateMinimumWage(wage); err != nil {
		return nil, err
	}
	if err := f.validateMinimumAge(dateOfBirth, processDate); err != nil {
		return nil, err
	}
	return &Customer{
		ID:        id,
		NIK:       nik,
		Fullname:  fullname,
		LegalName: legalName,
		PlaceAndDateOfBirth: PlaceAndDateOfBirth{
			Place: placeOfBirth,
			Date:  dateOfBirth,
		},
		Wage:     wage,
		PhotoURL: photoURL,
		KTPURL:   kTPURL,
		CreateAt: processDate,
		UpdateAt: processDate,
	}, nil
}
