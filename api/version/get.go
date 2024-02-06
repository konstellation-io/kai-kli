package version

import (
	"github.com/konstellation-io/kli/internal/entity"
	"github.com/konstellation-io/kli/internal/services/configuration"
)

// Get returns a version of a product given its tag.
func (c *Client) Get(server *configuration.Server, productID string, versionTag *string) (*entity.Version, error) {
	query := `
		query Version($productID: ID!, $tag: String) {
			version(productID: $productID, tag: $tag) {
				tag
        description
        creationDate
        creationAuthor
        publicationDate
        publicationAuthor
        status
        error
				publishedTriggers {
					trigger
					url
				}
        config {
					key
					value
        }
        workflows {
					name
					type
					config {
							key
							value
					}
					processes {
						name
						type
						image
						replicas
						gpu
						subscriptions
						status
						config {
							key
							value
						}
						objectStore {
							name
							scope
						}
						secrets {
							key
							value
						}
						networking {
							targetPort
							destinationPort
							protocol
						}
						resourceLimits {
							cpu {
								request
								limit
							}
							memory {
								request
								limit
							}
						}
					}
        }
			}
		}
		`
	vars := map[string]interface{}{
		"productID": productID,
	}

	if versionTag != nil {
		vars["tag"] = *versionTag
	}

	var respData struct {
		Version *entity.Version `json:"version"`
	}

	err := c.gqlClient.MakeRequest(server, query, vars, &respData)

	return respData.Version, err
}
