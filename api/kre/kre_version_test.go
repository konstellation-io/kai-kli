package kre_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/kre/version"
)

func TestVersionList(t *testing.T) {
	srv, client := gqlMockServer(t, "{\"productID\":\"test-product\"}", `
		{
				"data": {
						"versions": [
								{
										"name": "test-v1",
										"status": "STOPPED"
								}
						]
				}
		}
	`)
	defer srv.Close()

	c := version.New(client)

	product := "test-product"
	list, err := c.List(product)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, list[0], version.Version{
		Name:   "test-v1",
		Status: "STOPPED",
	})
}

func TestVersionStart(t *testing.T) {
	product := "test-product"
	versionName := "123456"
	comment := "test start comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"productID": "%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, product, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Start(product, versionName, comment)
	require.NoError(t, err)
}

func TestVersionStop(t *testing.T) {
	product := "test-product"
	versionName := "123456"
	comment := "test stop comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"productID":"%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, product, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Stop(product, versionName, comment)
	require.NoError(t, err)
}

func TestVersionPublish(t *testing.T) {
	product := "test-product"
	versionName := "123456"
	comment := "test publish comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"productID": "%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, product, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Publish(product, versionName, comment)
	require.NoError(t, err)
}

func TestVersionUnpublish(t *testing.T) {
	product := "test-product"
	versionName := "123456"
	comment := "test unpublish comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"productID": "%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, product, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Unpublish(product, versionName, comment)
	require.NoError(t, err)
}

func TestVersionGetConfig(t *testing.T) {
	product := "test-product"
	versionName := "test-v1"
	expectedVariables := fmt.Sprintf(`
		{
			"productID":"%s",
			"versionName": "%s"
		}`, product, versionName)
	srv, client := gqlMockServer(t, expectedVariables, `
		{
			"data": {
				"version": {
					"config": {
						"completed": false,
						"vars": [
							{
								"key": "KEY1",
								"value": "value",
								"type": "VARIABLE"
							},
							{
								"key": "KEY2",
								"value": "",
								"type": "VARIABLE"
							}
						]
					}
				}
			}
		}
	`)

	defer srv.Close()

	c := version.New(client)

	config, err := c.GetConfig(product, versionName)
	require.NoError(t, err)
	require.False(t, config.Completed)
	require.Len(t, config.Vars, 2)
	require.EqualValues(t, config, &version.Config{
		Completed: false,
		Vars: []*version.ConfigVariable{
			{
				Key:   "KEY1",
				Value: "value",
				Type:  version.ConfigVariableTypeVariable,
			},
			{
				Key:   "KEY2",
				Value: "",
				Type:  version.ConfigVariableTypeVariable,
			},
		},
	})
}

func TestVersionUpdateConfig(t *testing.T) {
	configVars := []version.ConfigVariableInput{
		{
			"key":   "KEY2",
			"value": "newValue",
		},
	}
	requestVars := `
		{
			"input": {
					"configurationVariables": [
						{ "key": "KEY2", "value": "newValue" }
					],
					"productID": "test-product",
					"versionName": "test-v1"
			}
		}
	`
	srv, client := gqlMockServer(t, requestVars, `
		{
			"data": {
				"updateVersionConfiguration": {
					"config": {
						"completed": true
					}
				}
			}
		}
	`)

	defer srv.Close()

	c := version.New(client)

	product := "test-product"
	completed, err := c.UpdateConfig(product, "test-v1", configVars)
	require.NoError(t, err)
	require.True(t, completed)
}

func TestVersionCreate(t *testing.T) {
	srv, client := gqlMockServer(t, "", `
		{
				"data": {
						"createVersion": {
								"name": "test-v1"
						}
				}
		}
	`)
	defer srv.Close()

	c := version.New(client)

	krtFile, err := os.CreateTemp("", ".test.krt")
	require.NoError(t, err)

	defer os.RemoveAll(krtFile.Name())

	product := "test-product"
	versionName, err := c.Create(product, krtFile.Name())
	require.NoError(t, err)
	require.Equal(t, versionName, "test-v1")
}
