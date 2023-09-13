package todos

type TodoModel struct {
	Title       string `json:"title" binding:"required"`
	State       bool   `json:"State" default:"false"`
	Description string `json:"description"`
}
