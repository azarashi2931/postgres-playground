package handler

import (
	"database/sql"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/koyashiro/postgres-playground/backend/runtime"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

type PlaygroundCreationResponse struct {
	Id string `json:"id"`
}

type Playground struct {
	Id string `json:"id"`
}

type ExecuteQueryRequest struct {
	Id    string `json:"id"`
	Port  int    `json:"port"`
	Query string `json:"query"`
}

type ExecuteQueryResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func PostPlayground(c echo.Context) error {
	id := uuid.NewString()
	dbm := runtime.NewDBManage(client.DefaultDockerHost)
	db, err := dbm.Create(id)
	if err != nil {
		c.Logger().Error(err)
		res := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, res)
	}
	c.Logger().Info(db)
	return c.JSON(http.StatusOK, db)
}

func GetPlayground(c echo.Context) error {
	id := c.Param("id")
	c.Logger().Info("show playground: " + id)
	playground := Playground{
		Id: id,
	}
	return c.JSON(http.StatusOK, playground)
}

func DeletePlayground(c echo.Context) error {
	id := c.Param("id")
	dbm := runtime.NewDBManage(client.DefaultDockerHost)
	if err := dbm.Destroy(id); err != nil {
		c.Logger().Error(err)
		res := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, res)
	}
	return c.JSON(http.StatusNoContent, nil)
}

func ExecuteQuery(c echo.Context) error {
	// id := c.Param("id")
	req := new(ExecuteQueryRequest)
	if err := c.Bind(req); err != nil {
		c.Logger().Error("invalid parameter")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid parameter"})
	}

	dbd, err := sql.Open("postgres", "user=playground password=password dbname=playground sslmode=disable port="+strconv.Itoa(req.Port))
	if err != nil {
		c.Logger().Error(err)
		res := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, res)
	}

	rows, err := dbd.Query(req.Query)
	if err != nil {
		c.Logger().Error(err)
		res := ErrorResponse{Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, res)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			c.Logger().Error(err)
		}
	}()

	var result string
	for rows.Next() {
		if err = rows.Scan(&result); err != nil {
			err = rows.Close()
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}

	res := ExecuteQueryResponse{
		Result: result,
	}
	return c.JSON(http.StatusOK, res)
}