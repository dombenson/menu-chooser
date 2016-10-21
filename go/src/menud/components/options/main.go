package options

import (
	"database/sql"
	"encoding/json"
)

type Option interface {
	Name() string
	Description() string
	ID() int
	CourseID() int
	json.Marshaler
}

type option struct {
	id             int
	name           string
	courseid       int
	description    string
	descriptionSet bool
}

const GetOptionsSQL = "SELECT `optionid`,`courseid`,`name`,`description` FROM `options` WHERE `courseid` = ? ORDER BY name ASC"
const GetOptionSQL = "SELECT `optionid`,`courseid`,`name`,`description` FROM `options` WHERE `optionid` = ?"

func MakeOption(rows *sql.Rows) (newOption Option, err error) {
	retOption := &option{}
	newOption = retOption
	var tmpDesc *string
	err = rows.Scan(&retOption.id, &retOption.courseid, &retOption.name, &tmpDesc)
	if tmpDesc != nil {
		retOption.descriptionSet = true
		retOption.description = *tmpDesc
	} else {
		retOption.descriptionSet = false
		retOption.description = ""
	}
	return
}

func (this *option) Name() string {
	return this.name
}
func (this *option) Description() string {
	return this.description
}
func (this *option) ID() int {
	return this.id
}
func (this *option) CourseID() int {
	return this.courseid
}

func (this *option) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CourseID    int    `json:"courseId"`
	}{
		ID:          this.id,
		Name:        this.name,
		Description: this.description,
		CourseID:    this.courseid,
	})
}
