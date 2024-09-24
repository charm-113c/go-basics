package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (int, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// Right now we're using database/sql, but you may want to abstract this,
	// If things go well, maintainability will be key, and GORM is a far better choice
	// for that, despite the performance trade-off.
	connStr := "user=postgres dbname=postgres password=bankingserver sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("DB is online")

	return &PostgresStore{
		db: db,
	}, nil
}

func (st *PostgresStore) CreateAccount(acc *Account) (int, error) {
	query := `INSERT INTO Account
		(firstName, lastName, accNumber, balance, createdAt)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id int
	err := st.db.QueryRow(query,
		acc.FirstName,
		acc.LastName,
		acc.AccNumber,
		acc.Balance,
		acc.CreatedAt).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (st *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := st.db.Query("SELECT * FROM Account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account

	for rows.Next() {
		acc := new(Account)
		acc, err = scanIntoAccount(rows)	
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	acc := new(Account)
	err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.AccNumber, &acc.Balance, &acc.CreatedAt)
	if err != nil { return nil, err }

	return acc, nil
}

func (st *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := st.db.Query("SELECT * FROM Account WHERE id = $1", id)
	if err != nil { return nil, err }
	defer rows.Close()

	for rows.Next() { 
		return scanIntoAccount(rows)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (st *PostgresStore) DeleteAccount(id int) error {
	// So we should consider soft deletions too, huh
	_, err := st.db.Exec("DELETE FROM Account WHERE id = $1", id)
	return err
}

func (st *PostgresStore) UpdateAccount(acc *Account) error {
	return nil
}

// This is the same pattern to initialize things (here the DB). In our case it's rather unnecssary, but if additional operations are necessary this will help keep things clean
func (st *PostgresStore) Init() error {
	return st.createAccountTable()
}
func (st *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS Account (
		id SERIAL PRIMARY KEY,
		firstName VARCHAR(50),
		lastName VARCHAR(50),
		accNumber SERIAL UNIQUE,
		balance INT,
		createdAt timestamp
	)`

	_, err := st.db.Exec(query)
	return err
}
