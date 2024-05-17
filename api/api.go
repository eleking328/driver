package api

import (
	"github.com/eleking328/driver-sdk/datasource"
)

type Driver interface {
	CreateChannel(datasourceId int64, properties []byte) (err error)

	ParsePointsProperties(datasourceId int64, ponits []byte) (map[int][]int64, error)

	Write(channelId int64, point []byte) (err error)

	Read(channelId int64, points []int64) (res *datasource.DataSourceNotify, err error)

	SubPointsData(channelId int64, subCh chan *datasource.DataSourceNotify) error
}
