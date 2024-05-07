package service

import (
	"github.com/eleking328/driver-sdk/datasource"
)

type DSManager struct {
	cache map[int64]datasource.DSEntity
}

func CreateDS(info datasource.DataSourceInfo) (err error) {
	return nil
}

func DeleteDS(dsID int64) (err error) {
	return nil
}

func IsWorking(dsID int64) bool {
	return false
}
