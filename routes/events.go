package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aquasecurity/tracee/types/trace"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/lcensies/tracee-mockserv/models"
)

// "github.com/aquasecurity/tracee/types/protocol"

var counter models.EventCounter = map[string]int{"vfs_file_write": 0}

func HandleEventsSink(c *gin.Context) {
	ByteBody, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(ByteBody))

	// jsonData, err := io.ReadAll(c.Request.Body)
	// if err != nil {
	// 	log.Error().Msgf("Failed to read request: %v", err)
	// 	return
	// }
	// fmt.Printf(string(jsonData))

	// var e models.MockEvent
	var e trace.Event
	if err := c.BindJSON(&e); err != nil {
		log.Error().Msgf("Failed to decode json: %v", err)
		return
	}
	log.Info().Msgf("Recevied %v", e.EventName)
	eventJson, _ := json.Marshal(e)
	log.Info().Msgf(string(eventJson))

	counter[e.EventName] += 1
}

func HandleEventsCount(c *gin.Context) {
	c.JSON(http.StatusOK, counter)
}

func HandleEventsCountReset(c *gin.Context) {
	log.Info().Msg("Clearing statistics")
	counter = map[string]int{"vfs_file_write": 0}
}
