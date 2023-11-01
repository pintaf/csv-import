package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tableflow/go/pkg/db"
	"tableflow/go/pkg/evaluator"
	"tableflow/go/pkg/types"
)

// getWorkspaceAPIKey
//
//	@Summary		Get workspace API key
//	@Description	Get the current API key of a workspace
//	@Tags			Workspace
//	@Success		200	{object}	string
//	@Failure		400	{object}	types.Res
//	@Router			/admin/v1/workspace/{id}/api-key [get]
//	@Param			id	path	string	true	"Workspace ID"
func getWorkspaceAPIKey(c *gin.Context, getWorkspaceUser func(*gin.Context, string) (string, error)) {
	workspaceID := c.Param("id")
	if len(workspaceID) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Res{Err: "No workspace ID provided"})
		return
	}
	_, err := getWorkspaceUser(c, workspaceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Res{Err: err.Error()})
		return
	}
	apiKey, err := db.GetAPIKey(workspaceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Res{Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiKey)
}

// regenerateWorkspaceAPIKey
//
//	@Summary		Regenerate workspace API key
//	@Description	Regenerate the current API key of a workspace and return the new API key
//	@Tags			Workspace
//	@Success		200	{object}	string
//	@Failure		400	{object}	types.Res
//	@Router			/admin/v1/workspace/{id}/api-key [post]
//	@Param			id	path	string	true	"Workspace ID"
func regenerateWorkspaceAPIKey(c *gin.Context, getWorkspaceUser func(*gin.Context, string) (string, error)) {
	workspaceID := c.Param("id")
	if len(workspaceID) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Res{Err: "No workspace ID provided"})
		return
	}
	_, err := getWorkspaceUser(c, workspaceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Res{Err: err.Error()})
		return
	}
	apiKey, err := db.RegenerateAPIKey(workspaceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Res{Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiKey)
}

// getWorkspaceDataTypeValidations
//
//	@Summary		Get datatype validations
//	@Description	Get a map of available data types and allowed validations
//	@Tags			Workspace
//	@Success		200	{object}	map[string][]string
//	@Failure		400	{object}	types.Res
//	@Router			/admin/v1/workspace/{id}/datatype-validations [get]
//	@Param			id	path	string	true	"Workspace ID"
func getWorkspaceDataTypeValidations(c *gin.Context, getWorkspaceUser func(*gin.Context, string) (string, error), getAllowedValidateTypes func(string) map[string]bool) {
	workspaceID := c.Param("id")
	if len(workspaceID) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.Res{Err: "No workspace ID provided"})
		return
	}
	_, err := getWorkspaceUser(c, workspaceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, types.Res{Err: err.Error()})
		return
	}

	allowedValidateTypes := getAllowedValidateTypes(workspaceID)
	dataTypeValidations := make(map[string][]string)

	for dataType, validations := range evaluator.DataTypeValidations {
		var newValidations []string
		for _, validation := range validations {
			if allowedValidateTypes == nil || allowedValidateTypes[validation] {
				newValidations = append(newValidations, validation)
			}
		}
		dataTypeValidations[dataType] = newValidations
	}

	c.JSON(http.StatusOK, &dataTypeValidations)
}
