package dto

type NumAgreedTasksDTO struct {
	Num int `json:"num_agreed_tasks"`
}

type NumRejectedTaskDTO struct {
	Num int `json:"num_rejected_tasks"`
}

type TotalTimeDTO struct {
	Time string `json:"total_time"`
}
