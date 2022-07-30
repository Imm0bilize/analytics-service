package kafkaSchemes

type TaskAnalyticsCreateTypePayload struct {
	TaskID string `json:"task_id"`
}

type TaskAnalyticsCreateType struct {
	BaseTopic
	Payload TaskAnalyticsCreateTypePayload `json:"payload"`
}

type TaskAnalyticsAcceptRejectTypePayload struct {
	TaskID    string `json:"task_id"`
	Email     string `json:"email"`
	Time      string `json:"time"`
	TaskState string `json:"task_state"`
}

type TaskAnalyticsAcceptRejectType struct {
	BaseTopic
	Payload TaskAnalyticsAcceptRejectTypePayload `json:"payload"`
}
