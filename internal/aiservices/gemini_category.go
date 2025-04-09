package aiservices

func SubcategoryRetrieval() (map[string]any, error) {
	structure, jsonErr := ReadJsonStructure("./internal/aiservices/jobs_guideds.json")
	if jsonErr != nil {
		return nil, jsonErr
	}

	return structure, nil
}
