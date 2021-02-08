package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type BlockListDB struct {
	DB *sql.DB
}

type IPDetail struct {
	Id           string
	IPAddress    string
	UpdatedAt    string
	CreatedAt    string
	Response     string
	ResponseCode string
}

func NewBlocklistDB(dbPath string) (*BlockListDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening db file %v", err)
	}

	return &BlockListDB{
		DB: db,
	}, nil

}

func (db *BlockListDB) RetrieveByIp(ip string) ([]IPDetail, bool, error) {

	resultSet, err := db.DB.Query("SELECT * FROM BlockList where ip_address = ?", ip)
	if err != nil {
		return nil, false, fmt.Errorf("error retrieving result set %v", err)
	}

	defer resultSet.Close()

	var results []IPDetail
	for resultSet.Next() {

		var id string
		var ipAddress string
		var response string
		var responseCode string
		var createdAt string
		var updatedAt string

		if err := resultSet.Scan(&id, &response, &responseCode, &updatedAt, &createdAt, &ipAddress); err != nil {
			return nil, false, fmt.Errorf("error getting result for id due to %v", err)
		}

		if id == "" {
			continue
		}

		bl := IPDetail{
			Id:           id,
			IPAddress:    ipAddress,
			UpdatedAt:    updatedAt,
			CreatedAt:    createdAt,
			Response:     response,
			ResponseCode: responseCode,
		}

		results = append(results, bl)

	}

	if len(results) < 1 {
		return nil, false, nil
	}

	return results, true, nil

}

func (db *BlockListDB) UpdateTimestamp(ip, date string) error {

	prepare, err := db.DB.Prepare(
		"UPDATE BlockList SET created_at = ? where ip_address = ?",
	)

	if err != nil {
		return fmt.Errorf("error preparing statement %v", err)
	}

	if _, err := prepare.Exec(date, ip); err != nil {
		return fmt.Errorf("error executing update %v", err)
	}

	return nil
}

func (db *BlockListDB) CreateIpDetail(ipDetail IPDetail) error {

	prepare, err := db.DB.Prepare("INSERT INTO BlockList VALUES (?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("error inserting ipdetail intop database %v", err)
	}

	if _, err := prepare.Exec(ipDetail.Id, ipDetail.Response,
		ipDetail.ResponseCode, ipDetail.UpdatedAt, ipDetail.CreatedAt,
		ipDetail.IPAddress); err != nil {

		return fmt.Errorf("error creating ip detail record %v", err)
	}

	return nil

}
