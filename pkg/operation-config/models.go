package operationConfig

type OperationConfig struct {
	Id              int  `json:"id,omitempty"`
	AllowNewProject bool `json:"allow_new_project,omitempty"`
}
