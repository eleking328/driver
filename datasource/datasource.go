package datasource

type DataPointBase struct {
	Dsec        string      `json:"describe"`
	CloudId     int64       `json:"id"`
	Name        string      `json:"name"`
	Properties  interface{} `json:"properties"`
	ReadAccess  bool        `json:"read"`
	WriteAccess bool        `json:"write"`
}

type DSEntity struct {
	DataPoint  []DataPointBase `json:"datapoint"`
	Dsec       string          `json:"describe"`
	DriverId   int64           `json:"driverId"`
	CloudId    int64           `json:"id"`
	Name       string          `json:"name"`
	Properties interface{}     `json:"properties"`
	CreateTime string          `db:"create_time" json:"createTime"`
}
