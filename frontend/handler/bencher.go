package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/gobench/backend/model"
	"github.com/gophergala2016/gobench/common"
	"github.com/labstack/echo"
	"net/http"
)

// ApiNextTaskHandler implements method for getting next task from the queue
func (h *handler) ApiNextTaskHandler(c *echo.Context) error {

	enc := json.NewDecoder(c.Request().Body)

	var taskReq common.TaskRequest
	if err := enc.Decode(&taskReq); err != nil {
		c.Request().Body.Close()
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong JSON. Expected gobench.common.TaskRequest")
	}

	c.Request().Body.Close()

	// checks existing test environment by authKey received by client
	ok, err := h.back.Model.TestEnvironment.Exist(taskReq.AuthKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong authKey!")
	}

	// retrives single task
	taskRow, err := h.back.Model.Task.Next(taskReq.AuthKey)

	if err != nil {
		if err == model.ErrNotFound {
			return echo.NewHTTPError(http.StatusNoContent)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	task := common.TaskResponse{Id: taskRow.Id.String(), PackageUrl: taskRow.PackageUrl, Type: "benchmark"}

	return c.JSON(http.StatusOK, task)
}

// ApiSubmitTaskResult implements saving of benchmark result
func (h *handler) ApiSubmitTaskResult(c *echo.Context) error {

	var tr common.TaskResult

	enc := json.NewDecoder(c.Request().Body)
	err := enc.Decode(&tr)
	c.Request().Body.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong JSON. Expected gobench.common.TaskResult")
	}

	// checks existing test environment by authKey received by client
	ok, err := h.back.Model.TestEnvironment.Exist(tr.AuthKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong authKey!")
	}

	// retrives single task to update
	taskRow, err := h.back.Model.Task.Next(tr.AuthKey)
	if err != nil {
		if err == model.ErrNotFound {
			return echo.NewHTTPError(http.StatusNoContent)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	task := common.TaskResponse{Id: taskRow.Id.String(), PackageUrl: taskRow.PackageUrl}
	return c.JSON(http.StatusOK, &task)
}
