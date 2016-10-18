package config

import (
	"fmt"
	"github.com/dombenson/go-ini"
)

const iniFilePath = "config.ini"

const dbSection = "db"

const userKey = "user"
const passKey = "pass"
const hostKey = "host"
const dbKey = "database"

var iniFile ini.Getter

func Init() {
	var err error
	iniFile, err = ini.LoadFile(iniFilePath)
	if err != nil {
		panic(err)
	}
}

func getWithDefault(section, key, defValue string) (value string) {
	value, ok := iniFile.Get(section, key)
	if !ok {
		value = defValue
	}
	return
}

func DBConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", getWithDefault(dbSection, userKey, "menud"), getWithDefault(dbSection, passKey, ""), getWithDefault(dbSection, hostKey, "127.0.0.1"), getWithDefault(dbSection, dbKey, "menud"))
}
