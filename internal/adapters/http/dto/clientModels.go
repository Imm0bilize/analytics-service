package dto

type NumAgreedTasksResponse struct {
	Num int `json:"num_agreed_tasks"`
}

type NumRejectedTaskResponse struct {
	Num int `json:"num_rejected_tasks"`
}

type TotalTimeResponse struct {
	Time string `json:"total_time"`
}
