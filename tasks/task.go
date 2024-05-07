package tasks

type Task struct {
	Interval     int
	DatasourceID int64
	pointId      int64
}

func NewTask(interval int, pointId int64, datasourceId int64) *Task {
	return &Task{
		Interval:     interval,
		pointId:      pointId,
		DatasourceID: datasourceId,
	}
}
