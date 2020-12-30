package database

import (
	"database/sql"
)

type Group struct {
	ID      int
	Club    string
	Name    string
	Owner   string
	GroupID int
}

type ByID []Group

func (a ByID) Len() int           { return len(a) }
func (a ByID) Less(i, j int) bool { return a[i].GroupID < a[j].GroupID }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func CreateGroup(db *sql.DB, groupName, clubName, ownerName string) error {
	_, err := db.Exec(
		"INSERT INTO clubgroups (groupname, clubname, owner) values ($1, $2, $3)",
		clubName, groupName, ownerName)
	if err != nil {
		return err
	}
	return nil
}

func SetGroupOwner(db *sql.DB, ownerName, clubName, groupName string) error {
	_, err := db.Exec(
		"UPDATE clubgroups SET owner=$1 WHERE clubname=$2 AND groupname=$3",
		ownerName, clubName, groupName)
	if err != nil {
		return err
	}
	return nil
}

func IsUserGroupOwner(db *sql.DB, nickName, clubName, groupName string) bool {
	err := db.QueryRow("SELECT owner FROM clubgroups WHERE clubname=$1 and groupname=$2",
		clubName, groupName).Scan(&nickName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	return true
}

func GetUserGroupOwner(db *sql.DB, clubName, groupName string) string {
	var nick string
	err := db.QueryRow("SELECT owner FROM clubgroups WHERE clubname=$1 and groupname=$2",
		clubName, groupName).Scan(&nick)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return nick
	}
	return nick
}

func GetClubGroups(db *sql.DB, clubName string) ([]Group, error) {
	rows, err := db.Query("SELECT * FROM clubgroups WHERE clubname=$1",
		clubName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []Group{}

	for rows.Next() {
		row := Group{}

		err := acceptGoupRows(rows, &row)
		if err != nil {
			return nil, err
		}

		groups = append(groups, row)
	}
	return groups, nil
}

func acceptGoupRows(rows *sql.Rows, row *Group) error {
	err := rows.Scan(
		&row.ID,
		&row.Name,
		&row.Club,
		&row.Owner,
		&row.GroupID)

	return err
}
