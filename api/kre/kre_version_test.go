package kre_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/konstellation-io/kli/api/kre/version"
)

func TestVersionList(t *testing.T) {
	srv, client := gqlMockServer(t, "{\"runtimeId\":\"test-runtime\"}", `
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

	runtime := "test-runtime"
	list, err := c.List(runtime)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, list[0], version.Version{
		Name:   "test-v1",
		Status: "STOPPED",
	})
}

func TestVersionStart(t *testing.T) {
	runtime := "test-runtime"
	versionName := "123456"
	comment := "test start comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"runtimeId": "%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, runtime, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Start(runtime, versionName, comment)
	require.NoError(t, err)
}

func TestVersionStop(t *testing.T) {
	runtime := "test-runtime"
	versionName := "123456"
	comment := "test stop comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"runtimeId":"%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, runtime, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Stop(runtime, versionName, comment)
	require.NoError(t, err)
}

func TestVersionPublish(t *testing.T) {
	runtime := "test-runtime"
	versionName := "123456"
	comment := "test publish comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"runtimeId": "%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, runtime, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Publish(runtime, versionName, comment)
	require.NoError(t, err)
}

func TestVersionUnpublish(t *testing.T) {
	runtime := "test-runtime"
	versionName := "123456"
	comment := "test unpublish comment"
	expectedVariables := fmt.Sprintf(`
		{
			"input": {
					"runtimeId": "%s",
					"versionName": "%s",
					"comment": "%s"
			}
		}
 `, runtime, versionName, comment)

	srv, client := gqlMockServer(t, expectedVariables, "")
	defer srv.Close()

	c := version.New(client)

	err := c.Unpublish(runtime, versionName, comment)
	require.NoError(t, err)
}

func TestVersionGetConfig(t *testing.T) {
	runtime := "test-runtime"
	versionName := "test-v1"
	expectedVariables := fmt.Sprintf(`
		{
			"runtimeId":"%s",
			"versionName": "%s"
		}`, runtime, versionName)
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

	config, err := c.GetConfig(runtime, versionName)
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
					"runtimeId": "test-runtime",
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

	runtime := "test-runtime"
	completed, err := c.UpdateConfig(runtime, "test-v1", configVars)
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

	runtime := "test-runtime"
	versionName, err := c.Create(runtime, krtFile.Name())
	require.NoError(t, err)
	require.Equal(t, versionName, "test-v1")
}
