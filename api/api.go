package api

import (
	"git.cddpi.com/iot/iot-edge-driver/datasource"
)

type Driver interface {
	CreateChannel(datasourceId int64, properties []byte) (err error)

	ParsePointsProperties(datasourceId int64, ponits []byte) (map[int][]int64, error)

	Write(channelId int64) (err error)

	Read(channelId int64, points []int64) (res *datasource.DataSourceNotify, err error)

	SubPointsData(channelId int64, subCh chan *datasource.DataSourceNotify) error
}
