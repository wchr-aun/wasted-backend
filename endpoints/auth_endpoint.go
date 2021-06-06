package endpoints

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func Authentication(c *gin.Context) {
	firestoreCon := c.MustGet("firestoreCon").(*firestore.Client)
	uuid := c.MustGet("UUID").(string)

	defer firestoreCon.Close()

	doc, err := firestoreCon.Collection("users").Doc(uuid).Get(context.Background())
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			addUser(c, uuid)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fullname": doc.Data()["fullname"].(string),
	})
}

func addUser(c *gin.Context, uuid string) {
	c.JSON(http.StatusOK, gin.H{
		"hey": "ok",
	})
}
