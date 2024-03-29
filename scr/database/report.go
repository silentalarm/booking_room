package database

import "database/sql"

type Report struct {
	ID       int
	NickName string
	IDIntra  string
	Comment  string
	Decided  bool
	ReportID string
}

func InsertNewReport(db *sql.DB, nickName, idIntra, comment string) {
	_, err := db.Exec(
		"INSERT INTO report (nickname, idintra, comment) values ($1, $2, $3)",
		nickName, idIntra, comment)
	if err != nil {
		panic(err)
	}
}

func GetReportList(db *sql.DB) ([]Report, error) {
	rows, err := db.Query("SELECT * FROM report")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []Report{}

	for rows.Next() {
		row := Report{}
		err := rows.Scan(
			&row.ID,
			&row.NickName,
			&row.IDIntra,
			&row.Comment,
			&row.Decided,
			&row.ReportID)

		if err != nil {
			continue
		}

		reports = append(reports, row)
	}
	return reports, nil
}

func ReportReady(db *sql.DB, reportID string) {
	_, err := db.Exec("UPDATE report SET decided=true WHERE reportid=$1", reportID)
	if err != nil {
		panic(err)
	}
}

func ReportDelete(db *sql.DB, reportID string) {
	_, err := db.Exec("DELETE FROM report WHERE reportid=$1", reportID)
	if err != nil {
		panic(err)
	}
}
