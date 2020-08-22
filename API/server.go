package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	/*
		r.GET("/kill/:pid", func(c *gin.Context) {
			pidString := c.Param("pid")

			pid, err := strconv.Atoi(pidString)

			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())

			} else {
				fmt.Printf("Se procedera a matar al proceso: %d \n", pid)

				err = killProcess(pid)

				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
				} else {
					c.AbortWithStatus(http.StatusNoContent)
				}
			}
		})

		r.GET("/cpu", func(c *gin.Context) {
			idle, total := getCPUUsage()
			c.JSON(http.StatusOK, gin.H{"idle": idle, "total": total})
		})

		r.GET("/ram", func(c *gin.Context) {
			used, total := getRAMUsage()
			c.JSON(http.StatusOK, gin.H{
				"used":    used,
				"total":   total,
				"usedMb":  float64(used) / 1024,
				"totalMb": float64(total) / 1024,
				"usage":   float64(used*100) / float64(total),
			})
		})
	*/
	r.GET("/processes", func(c *gin.Context) {
		processes := getLinuxProcesses()
		c.JSON(http.StatusOK, processes)
	})

	_ = r.Run(":8080") // listen and serve on 0.0.0.0:8080

}
