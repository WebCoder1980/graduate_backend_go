package model

type ImageInfo struct {
	Id       int64
	Filename string
	Format   string
	TaskId   int64
	Position int
	StatusId int64
}

type ImageStatus struct {
	TaskId   int64
	Position int
	//StatusId int64 TODO
}
