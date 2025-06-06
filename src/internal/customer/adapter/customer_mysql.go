package adapter

import (
	"context"
	"time"
	"xyz_multifinance/src/internal/customer/domain/customer"
	mysqlutil "xyz_multifinance/src/internal/shared/mysql"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlCustomer struct {
	CustomerID   string    `db:"customer_id"`
	NIK          string    `db:"nik"`
	FullName     string    `db:"full_name"`
	LegalName    string    `db:"legal_name"`
	PlaceOfBirth string    `db:"place_of_birth"`
	DateOfBirth  time.Time `db:"date_of_birth"`
	Wages        float64   `db:"wages"`
	KTPPhotoURL  string    `db:"ktp_photo_url"`
	PhotoURL     string    `db:"photo_url"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func unmarshallToDomain(record mysqlCustomer) *customer.Customer {
	return &customer.Customer{
		ID:        record.CustomerID,
		NIK:       record.NIK,
		Fullname:  record.FullName,
		LegalName: record.LegalName,
		PlaceAndDateOfBirth: customer.PlaceAndDateOfBirth{
			Place: record.PlaceOfBirth,
			Date:  record.DateOfBirth,
		},
		Wage:     record.Wages,
		PhotoURL: record.PhotoURL,
		KTPURL:   record.KTPPhotoURL,
		CreateAt: record.CreatedAt,
		UpdateAt: record.UpdatedAt,
	}
}

type MysqlCustomerRepository struct {
	db *sqlx.DB
}

func NewMysqlCustomerRepository(db *sqlx.DB) *MysqlCustomerRepository {
	return &MysqlCustomerRepository{
		db: db,
	}
}

func (m MysqlCustomerRepository) GetCustomer(
	ctx context.Context,
	customerID string,
) (*customer.Customer, error) {
	row := m.db.QueryRowContext(
		ctx, `
		SELECT 
			customer_id, nik, full_name, legal_name, place_of_birth,
			date_of_birth, wages, ktp_photo_url, photo_url,
			created_at, updated_at
		FROM Customers tl 
		WHERE tl.customer_id = ?`, customerID)

	var customerDb mysqlCustomer
	if err := row.Scan(
		&customerDb.CustomerID, &customerDb.NIK, &customerDb.FullName, &customerDb.LegalName, &customerDb.PlaceOfBirth,
		&customerDb.DateOfBirth, &customerDb.Wages, &customerDb.KTPPhotoURL, &customerDb.PhotoURL,
		&customerDb.CreatedAt, &customerDb.UpdatedAt,
	); err != nil {
		return nil, err
	}

	customer := unmarshallToDomain(customerDb)
	return customer, nil
}

func (m MysqlCustomerRepository) Create(
	ctx context.Context,
	customer *customer.Customer,
) error {
	for {
		err := m.create(ctx, customer)
		if val, ok := errors.Cause(err).(*mysql.MySQLError); ok && val.Number == mysqlutil.MySQLDeadlockErrorCode {
			continue
		}
		return err
	}
}

func (m MysqlCustomerRepository) create(
	ctx context.Context,
	customer *customer.Customer,
) error {
	createdDbCustomer := mysqlCustomer{
		CustomerID:   customer.ID,
		NIK:          customer.NIK,
		FullName:     customer.Fullname,
		LegalName:    customer.LegalName,
		PlaceOfBirth: customer.PlaceAndDateOfBirth.Place,
		DateOfBirth:  customer.PlaceAndDateOfBirth.Date,
		Wages:        customer.Wage,
		KTPPhotoURL:  customer.KTPURL,
		PhotoURL:     customer.PhotoURL,
		CreatedAt:    customer.CreateAt,
		UpdatedAt:    customer.UpdateAt,
	}
	_, err := m.db.NamedExecContext(
		ctx,
		`INSERT INTO 
			Customers (
				customer_id, nik, full_name, legal_name, place_of_birth,
				date_of_birth, wages, ktp_photo_url, photo_url,
				created_at, updated_at
			)
		VALUES
			(
				:customer_id, :nik, :full_name, :legal_name, :place_of_birth,
				:date_of_birth, :wages, :ktp_photo_url, :photo_url,
				:created_at, :updated_at
			)
		`, createdDbCustomer)
	if err != nil {
		return errors.Wrap(err, "unable to create customer")
	}
	return nil
}
