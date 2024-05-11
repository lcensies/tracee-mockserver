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

var (
	counter           models.EventCounter = map[string]int{}
	fileEventsCounter                     = map[FileEventsSummary]int{}
)

type FileEventsSummary struct {
	TargetFile   string `json:"pathname"`
	ProcessImage string `json:"image"`
	ProcessName  string `json:"process"`
	EventName    string `json:"event"`
	DeviceType   string `json:"dev_type"`
}

var fileEvents map[string]string = map[string]string{
	"vfs_write":          "vfs_write",
	"vfs_read":           "vfs_read",
	"vfs_readv":          "vfs_readv",
	"security_file_open": "security_file_open",
}

func getEventTargetFile(e *trace.Event) string {
	for _, arg := range e.Args {
		if arg.ArgMeta.Name == "pathname" {
			return arg.Value.(string)
		}
	}
	return ""
}

func updateEventsSummary(e *trace.Event) {
	counter[e.EventName] += 1

	if _, ok := fileEvents[e.EventName]; !ok {
		return
	}

	eventKey := FileEventsSummary{
		TargetFile:   getEventTargetFile(e),
		EventName:    e.EventName,
		ProcessImage: e.Executable.Path,
		ProcessName:  e.ProcessName,
	}

	fileEventsCounter[eventKey] += 1
}

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

	updateEventsSummary(&e)
}

func HandleEventsCount(c *gin.Context) {
	c.JSON(http.StatusOK, counter)
}

func HandleEventsCountReset(c *gin.Context) {
	log.Info().Msg("Clearing statistics")
	counter = map[string]int{}
	fileEventsCounter = map[FileEventsSummary]int{}
}

func HandleFileEventsCount(c *gin.Context) {
	type fileEventsSummary struct {
		FileEventsSummary
		Count int `json:"count"`
	}
	var summaries []fileEventsSummary
	for ioKey, count := range fileEventsCounter {
		summaries = append(summaries,
			fileEventsSummary{
				ioKey,
				count,
			},
		)
	}
	c.JSON(http.StatusOK, summaries)
}
