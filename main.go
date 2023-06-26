package main
/* Derived from :
  https://go.dev/doc/tutorial/web-service-gin
	and from:
	https://betterprogramming.pub/how-to-generate-html-with-golang-templates-5fad0d91252
*/

import (
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album. Notice the change from the base tutorial: `form replaces `json
// binding:"required" for ID means ID cannot be a zero value (="" for a string)
// To test that, try entering blank ID when adding new album
// with output:  Error Key: 'album.ID' Error:Field validation for 'ID' failed on the 'required' tag
type album struct {
	ID     string `form:"id" binding:"required"`
	Title  string `form:"title"`
	Artist string `form:"artist"`
	Price  float64 `form:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "You reached the album management site...",
		})
	})
	router.GET("/addalbum", func(c *gin.Context) {
		c.HTML(http.StatusOK, "addalbum.html", gin.H{
			"content": "Add new album :",
			"action": "/albums",
		})
	})
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum)
	router.POST("/albums/:id", putAlbum) // use POST method for update: no gin API for PUT method
	router.GET("/update/:id", updateAlbum)

	router.Run(":8082")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
	  "content": template.HTML(prettyPrint(albums)),
	})
}

// postAlbums adds an album from form received in the request body.
func postAlbum(c *gin.Context) {
	var newAlbum album
  log.Println("Prix="+c.PostForm("price"))
    // Call Bind to bind the received form to
    // newAlbum.
    if err := c.Bind(&newAlbum); err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
			  "content": template.HTML("Error "+err.Error()),
			})
			return
    }

    // Add the new album to the slice.
    albums = append(albums, newAlbum)

		c.HTML(http.StatusOK, "index.html", gin.H{
		  "content": template.HTML("<p>New album :</p>"+prettyPrint(newAlbum)),
		})
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.HTML(http.StatusOK, "index.html", gin.H{
			  "content": template.HTML("<p>Album with id "+id+" :</p>"+prettyPrint(a)),
			})
			return
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": "album "+string(id)+" not found",
	})
}

func putAlbum(c *gin.Context) { // Using method POST, no good support for PUT in Gin
	id := c.Param("id")
	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for i, a := range albums {
		if a.ID == id {
			// Use GinGonic `Bind` function to bind the data in the form to albums[i]
			// All validation and conversion (eg converting price from string to float64)
			// will be carried out and errors (eg conversion errors) will be reported
			if err := c.Bind(&albums[i]); err != nil {
				c.HTML(http.StatusOK, "index.html", gin.H{
				  "content": template.HTML("Error "+err.Error()),
				})
				return
	    }
			albums[i].ID = id // should not change id !!
			c.HTML(http.StatusOK, "index.html", gin.H{
			  "content": template.HTML("<p>Updated album with id "+id+" :</p>"+prettyPrint(albums[i])),
			})
			return
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": "album "+string(id)+" not found",
	})
}

func updateAlbum(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.HTML(http.StatusOK, "addalbum.html", gin.H{
			  "content": template.HTML("<p>Update album with id "+id+" :</p>"),
				"action": "/albums/"+id,
				"id" : a.ID,
				"title" : a.Title,
				"artist" : a.Artist,
				"price" : fmt.Sprintf("%.2f", a.Price),
			})
			return
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": "album "+string(id)+" not found",
	})
}

func prettyPrint(v any) string {
	b,_ := json.MarshalIndent(v, "  ", "\t")
	res := strings.ReplaceAll(string(b), "\n", "<br />")
	return strings.ReplaceAll(string(res), "\t", "&nbsp;&nbsp;&nbsp;&nbsp;")
}
