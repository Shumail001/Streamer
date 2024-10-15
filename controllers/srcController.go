package controllers

import (
	"Streamer/database"
	"Streamer/media/source-pipelines"
	"Streamer/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

var validate = validator.New()
var sources = make(map[string]*gst.Pipeline)

func StartPipeline() gin.HandlerFunc {
	return func(c *gin.Context) {
		var source models.RtspSrcModel

		// Bind the Json
		if err := c.ShouldBindJSON(&source); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Json Format"})
			return
		}

		// validate the struct
		if err := validate.Struct(source); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Initialize Object Box

		ob := database.InitObjectBox()
		if ob == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to init ObjectBox"})
			return
		}
		defer ob.Close()

		// Open the box for rtsp source model
		box := models.BoxForRtspSrcModel(ob)
		if box == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the box for source-pipelines model"})
			return
		}

		// Start the pipeline
		srcPipeline, _, _, er := source_pipelines.CreateSourcePipeline(source.Src)
		if er != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
			return
		}
		// Add in the database

		id, err := box.Put(&models.RtspSrcModel{
			Id:  source.Id,
			Src: source.Src,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result, err := box.Get(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source-pipelines"})
			return
		}

		// Start the Pipeline
		sources[strconv.FormatUint(id, 10)] = srcPipeline
		if err := srcPipeline.SetState(gst.StatePlaying); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start pipeline"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Source Start Successfully", "Response": result})

	}
}

func AddSource() gin.HandlerFunc {
	return func(c *gin.Context) {
		var source models.SrcModel

		// Bind JSON
		if err := c.ShouldBindJSON(&source); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the Struct Model
		switch source.Protocol {
		case models.UDP, models.RTP, models.SRT:
			if err := validate.StructExcept(&source, "Url", "Path"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		case models.RTMP, models.RTSP, models.DASH, models.HLS:
			if err := validate.StructExcept(&source, "Address", "Port", "Path"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Protocol not supported"})
			return
		}

		// Initialize Object box
		ob := database.InitObjectBox()
		if ob == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "object box init error"})
			return
		}

		defer ob.Close()

		// Open the box for Source Model
		box := models.BoxForSrcModel(ob)
		if box == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the box for source-pipelines model"})
			return
		}

		// Start the pipeline
		srcPipeline, _, _, er := source_pipelines.CreateSourcePipeline(string(source.Protocol))
		if er != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
			return
		}

		// Add source-pipelines in Database
		source.CreatedAt = time.Now().String()
		source.UpdatedAt = time.Now().String()
		var id uint64
		var err error
		switch source.Protocol {
		case models.UDP, models.RTP, models.SRT:
			id, err = box.Put(&models.SrcModel{
				Id:        source.Id,
				Protocol:  source.Protocol,
				Address:   source.Address,
				Port:      source.Port,
				CreatedAt: source.CreatedAt,
				UpdatedAt: source.UpdatedAt,
			})
		case models.RTMP, models.DASH, models.RTSP, models.HLS:
			id, err = box.Put(&models.SrcModel{
				Id:        source.Id,
				Protocol:  source.Protocol,
				Url:       source.Url,
				Port:      source.Port,
				CreatedAt: source.CreatedAt,
				UpdatedAt: source.UpdatedAt,
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid protocol: %s", source.Protocol)})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source-pipelines"})
			return
		}

		result, err := box.Get(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source-pipelines"})
			return
		}
		// Start the pipeline
		sources[strconv.FormatUint(id, 10)] = srcPipeline
		// Start the pipeline
		if err := srcPipeline.SetState(gst.StatePlaying); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start pipeline"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Source Created Successfully", "Response": result})
	}
}

func GetSourceById() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("source_id")

		// Initialize object box

		ob := database.InitObjectBox()
		if ob == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "object box init error"})
			return
		}
		defer ob.Close()

		// Open the box for Source Model
		box := models.BoxForSrcModel(ob)
		if box == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the box for source-pipelines model"})
			return
		}

		// Convert Source ID from string to uint64
		sourceID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Source ID"})
			return
		}

		result, err := box.Get(sourceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get source-pipelines by id"})
			return
		}

		if result == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not found source-pipelines with this id"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": result})

	}
}

func GetSourceList() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initialize Database
		ob := database.InitObjectBox()
		if ob == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "object box init error"})
			return
		}

		defer ob.Close()
		box := models.BoxForSrcModel(ob)
		if box == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the box for source-pipelines model"})
			return
		}

		result, err := box.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get source-pipelines list"})
			return
		}
		if result == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not found source-pipelines list"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": result})
	}
}

func RemoveSourceById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("source_id")

		// Initialize database
		ob := database.InitObjectBox()
		if ob == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "object box init error"})
			return
		}
		defer ob.Close()
		box := models.BoxForSrcModel(ob)
		if box == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the box for source-pipelines model"})
			return
		}

		sourceId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Source ID"})
			return
		}

		result, err := box.Get(sourceId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get source-pipelines by id"})
			return
		}
		if result == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not found source-pipelines with this id"})
			return
		}

		err = box.Remove(result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove source-pipelines by id"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Source Removed Successfully"})
	}
}

func UpdateSourceById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("source_id")

		// Initialize Database
		ob := database.InitObjectBox()
		if ob == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "object box init error"})
			return
		}
		defer ob.Close()

		// Open the Box for Source Model
		box := models.BoxForSrcModel(ob)
		if box == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the box for source-pipelines model"})
			return
		}

		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Source ID"})
			return
		}

		var updateSource models.SrcModel

		// Bind the Json
		if err := c.ShouldBindJSON(&updateSource); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the updated Source
		switch updateSource.Protocol {
		case models.UDP, models.RTP, models.SRT:
			if err := validate.StructExcept(&updateSource, "Url", "Path"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		case models.RTMP, models.RTSP, models.DASH, models.HLS:
			if err := validate.StructExcept(&updateSource, "Address", "Port", "Path"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Protocol not supported"})
			return
		}

		// Get the Source id
		source, err := box.Get(parsedId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get source-pipelines by this id"})
			return
		}

		if source == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "source-pipelines not found"})
			return
		}
		//Update the source-pipelines fields
		switch updateSource.Protocol {
		case models.UDP, models.RTP, models.SRT:
			source.Protocol = updateSource.Protocol
			source.Address = updateSource.Address
			source.Port = updateSource.Port
		case models.RTMP, models.DASH, models.RTSP, models.HLS:
			source.Protocol = updateSource.Protocol
			source.Url = updateSource.Url
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid protocol: %s", source.Protocol)})
			return

		}
		// Set the updated time
		source.UpdatedAt = time.Now().String()

		// update the source-pipelines data in database
		updatedID, err := box.Put(source)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source-pipelines"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Source with ID %v updated successfully", updatedID), "Source": updatedID})

	}
}
