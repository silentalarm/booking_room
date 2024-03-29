package database

import (
	"database/sql"
	"errors"
)

type Group struct {
	ID      int
	Club    string
	Name    string
	Owner   string
	GroupID int
	Color   string
	Size    int
}

type ByID []Group

func (a ByID) Len() int           { return len(a) }
func (a ByID) Less(i, j int) bool { return a[i].GroupID < a[j].GroupID }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func InsertGroup(db *sql.DB, groupName, clubName, ownerName, colorGroup string, groupID int) error {
	_, err := db.Exec(
		"INSERT INTO clubgroups (groupname, clubname, owner, groupid, color) values ($1, $2, $3, $4, $5)",
		groupName, clubName, ownerName, groupID, colorGroup)
	if err != nil {
		return err
	}
	return nil
}

func CreateMainGroup(db *sql.DB, clubName, ownerName string) error {
	err := InsertGroup(db, "main", clubName, ownerName, "#ffffff", 0)
	if err != nil {
		return err
	}
	return nil
}

func CreateNewGroup(db *sql.DB, groupName, clubName, ownerName, colorGroup string) error {
	isExist := GroupIsExist(db, groupName, clubName)
	if isExist == true {
		return errors.New("groupName: group is exist")
	}

	lasId, err := GetLastGroupID(db, clubName)
	if err != nil {
		return err
	}
	lasId += 1

	err = InsertGroup(db, groupName, clubName, ownerName, colorGroup, lasId)
	if err != nil {
		return err
	}
	return nil
}

func AddGroup(db *sql.DB, groupName, clubName, ownerName, colorGroup string) error {
	newGroupID, err := GetLastGroupID(db, groupName)
	if err != nil {
		return err
	}

	err = InsertGroup(db, groupName, clubName, ownerName, colorGroup, newGroupID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteGroup(db *sql.DB, groupName, clubName string) error {
	_, err := db.Exec("DELETE FROM clubgroups WHERE groupname=$1 and clubname=$2",
		groupName, clubName)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllClubGroups(db *sql.DB, clubName string) error {
	_, err := db.Exec("DELETE FROM clubgroups WHERE clubname=$1",
		clubName)
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

func SetUserGroup(db *sql.DB, groupName, nickName, clubName string) error {
	isExist := GroupIsExist(db, groupName, clubName)
	if isExist == false {
		return errors.New("groupName: group does not exist")
	}

	_, err := db.Exec("UPDATE clubmembers SET groupname=$1 WHERE nickname=$2 and clubname=$3",
		groupName, nickName, clubName)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func SetUsersGroup(db *sql.DB, oldGroupName, newGroupName, clubName string) error {
	isExist := GroupIsExist(db, newGroupName, clubName)
	if isExist == false {
		return errors.New("groupName: group does not exist")
	}

	_, err := db.Exec("UPDATE clubmembers SET groupname=$1 WHERE groupname=$2 and clubname=$3",
		newGroupName, oldGroupName, clubName)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func ChangeUserGroup(db *sql.DB, groupName, nickName, clubName string) error {
	isExist := GroupIsExist(db, groupName, clubName)
	if isExist == false {
		return errors.New("groupName: group does not exist")
	}

	isOwner := IsUserGroupOwner(db, nickName, clubName, groupName)
	if isOwner == true {
		return errors.New("nickName: user is group owner")
	}

	_, err := db.Exec("UPDATE clubmembers SET groupname=$1 WHERE nickname=$2 and clubname=$3",
		groupName, nickName, clubName)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func UserLeaveGroup(db *sql.DB, nickName, clubName string) error {
	err := SetUserGroup(db, "main", nickName, clubName)
	if err != nil {
		panic(err)
	}
	return nil
}

func GroupIsExist(db *sql.DB, groupName, clubName string) bool {
	err := db.QueryRow("SELECT groupname FROM clubgroups WHERE groupname=$1 and clubname=$2",
		groupName, clubName).Scan(&groupName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	return true
}

func GetLastGroupID(db *sql.DB, clubName string) (int, error) {
	var groupID int
	err := db.QueryRow("SELECT MAX(groupid) FROM clubgroups WHERE clubname=$1",
		clubName).Scan(&groupID)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return -1, err
	}
	return groupID, nil
}

func IsUserGroupOwner(db *sql.DB, nickName, clubName, groupName string) bool {
	var owner string
	err := db.QueryRow("SELECT owner FROM clubgroups WHERE clubname=$1 and groupname=$2",
		clubName, groupName).Scan(&owner)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	if owner != nickName {
		return false
	}

	return true
}

func IsUserGroupsOwner(db *sql.DB, nickName, clubName string) bool {
	err := db.QueryRow("SELECT owner FROM clubgroups WHERE owner=$1 and clubname=$2",
		nickName, clubName).Scan(&nickName)
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

func GetUserGroup(db *sql.DB, nickName, clubName string) (string, error) {
	var name string
	err := db.QueryRow("SELECT groupname FROM clubmembers WHERE nickname=$1 and clubname=$2",
		nickName, clubName).Scan(&name)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return name, err
	}
	return name, nil
}

func GetGroupColor(db *sql.DB, clubName, groupName string) string {
	color := "#ffffff"
	err := db.QueryRow("SELECT color FROM clubgroups WHERE clubname=$1 and groupname=$2",
		clubName, groupName).Scan(&color)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return color
	}
	return color
}

func GetGroupSize(db *sql.DB, clubName, groupName string) (int, error) {
	var size int
	err := db.QueryRow("SELECT COUNT(*) FROM clubmembers WHERE clubname=$1 and groupname=$2",
		clubName, groupName).Scan(&size)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return size, err
	}

	return size, nil
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

		row.Size, err = GetGroupSize(db, clubName, row.Name)
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
		&row.GroupID,
		&row.Color)

	return err
}
