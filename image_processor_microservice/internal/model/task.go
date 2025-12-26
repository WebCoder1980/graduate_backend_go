package model

import "time"

type TaskInfo struct {
	Id        int64     `json:"id"`
	CreatedDT time.Time `json:"created_dt"`
	Width     string    `json:"width"`
	Height    string    `json:"height"`
	Format    string    `json:"format"`
	Quality   float32   `json:"quality"`
}

type TaskResponse struct {
	TaskInfo
	CommonStatusId int64       `json:"common_status_id"`
	Images         []ImageInfo `json:"images"`
	CreatedDT      time.Time   `json:"created_dt"`
}
