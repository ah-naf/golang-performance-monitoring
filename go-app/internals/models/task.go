package models

type Task struct {
	ID          int    `json:"id" binding:""`
	Title       string `json:"title" binding:"required"`
	Description string `json:"desc" binding:"required"`
	Status      string `json:"status" binding:"required"`
}
