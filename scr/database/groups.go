package database

import "database/sql"

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
