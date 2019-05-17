package web

import (
	"github.com/gin-gonic/gin"
	"github.com/orvice/monitor-server/enum"
	"github.com/orvice/monitor-server/internal/config"
	"github.com/orvice/monitor-server/internal/hub"
	"github.com/orvice/monitor-server/internal/mod"
	"net/http"
	"strconv"
)

type Ret struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func DataCollect(c *gin.Context) {
	if c.GetHeader(enum.PostKey) != config.PostKey {
		c.JSON(http.StatusUnauthorized, Ret{
			Code:    http.StatusUnauthorized,
			Message: "Wrong Key",
		})
		return
	}

	nid := c.Param("node_id")
	nodeID, err := strconv.Atoi(nid)
	if err != nil {
		c.JSON(http.StatusBadRequest, Ret{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, Ret{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	data := mod.Packet{
		NodeID:  int32(nodeID),
		Message: rawData,
	}
	hub.Manager.SendPacket(data)
	c.JSON(http.StatusOK, Ret{
	})
	return
}
