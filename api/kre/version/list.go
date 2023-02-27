package version

// Version represents a Version entity in KRE.
type Version struct {
	Name   string
	Status string
}

// List contains a list of  Version.
type List []Version

// List calls to KRE API and returns a list of Version entities.
func (c *versionClient) List(runtime string) (List, error) {
	query := `
		query GetVersions($runtimeId: ID!) {
			versions(runtimeId: $runtimeId) {
				name
				status
			}
		}
	`

	vars := map[string]interface{}{
		"runtimeId": runtime,
	}

	var respData struct {
		Versions List
	}

	err := c.client.MakeRequest(query, vars, &respData)

	return respData.Versions, err
}
