package projects

import "log"

func validateAddProjectPayload(payload AddProjectRequest) error {
	if payload.Collaborated == nil {
		return &CollaboratedRequiredError{}
	}
	log.Println("==collab", *payload.Collaborated)
	return nil
}
