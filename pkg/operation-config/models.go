package operationConfig

type OperationConfig struct {
	Id              int   `json:"id,omitempty"`
	AllowNewProject *bool `json:"allowNewProject,omitempty"`
}
