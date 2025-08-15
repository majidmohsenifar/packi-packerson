package api

import (
	"job-test/internal/pack"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func valid(c *gin.Context) int {
	xStr := c.Query("x")
	if xStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "missing query param 'x'"})
		return 0
	}

	x, err := strconv.Atoi(xStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "'x' must be an integer"})
		return 0
	}
	if x < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "'x' must be >= 0"})
		return 0
	}

	// Define a maximum allowed value for 'x' to prevent excessive computation
	// This is a safeguard against performance issues or excessive memory usage.
	const maxAllowed = 1_000_000
	if x > maxAllowed {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "'x' is too large"})
		return 0
	}
	return x

}

type SetSizesReq struct {
	Sizes []int
}

// @Summary Set sizes
// @Description Set sizes for packs
// @Tags pack
// @Accept  json
// @Produce  json
// @Param sizes body int true "X is number"
// @Router /set-sizes [post]
func setSizes(c *gin.Context) {
	var reqBody SetSizesReq
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid sizes format"})
		return
	}
	newSizes := pack.SetSizes(reqBody.Sizes)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": map[string][]int{
			"sizes": newSizes,
		},
	})
}

// @Summary Get Correct
// @Description Give correct
// @Tags pack
// @Accept  json
// @Produce  json
// @Param x query int true "X is number"
// @Router /correct [get]
func correct(c *gin.Context) {
	x := valid(c)
	if x == 0 {
		return // Early return if validation fails
	}
	packs := pack.Correct(x)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"ordered": x,
			"packs":   packs,
		},
	})
}

// @Summary Get Incorrect
// @Description Give incorrect
// @Tags pack
// @Accept  json
// @Produce  json
// @Param x query int true "X is number"
// @Router /incorrect [get]
func incorrect(c *gin.Context) {
	x := valid(c)
	if x == 0 {
		return // Early return if validation fails
	}
	packs := pack.InCorrect(x)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"ordered": x,
			"packs":   packs,
		},
	})
}

func InitApi(r *gin.Engine) {
	// static files
	r.Static("/static", "./static") // For CSS, JS, images
	// Serve index.html at /
	r.GET("/", func(c *gin.Context) { // Serve index.html at /
		c.File("./static/index.html")
	})

	r.POST("/set-sizes", setSizes)
	r.GET("/correct", correct)
	r.GET("/incorrect", incorrect)
}
