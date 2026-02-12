package main


import (
	"os"
	"text/template"
	"net/http"
	"fmt"
	// "io"
	"bytes"
    "github.com/gin-gonic/gin"
)

// The template that does raw response when not being requested by HTMX
type RawResponse struct {
	Content string
}

func ReturnContentfulPage(c *gin.Context, skeleton_path string, partial_path string) {
	buf := new(bytes.Buffer)
	skel_html, err := os.ReadFile(skeleton_path)
	if err != nil {
		fmt.Print(err)
	}
	partial_html, err := os.ReadFile(partial_path)
	if err != nil {
		fmt.Print(err)
	}
	tmpl, err := template.New("tmpl").Parse(string(skel_html))
	if err != nil {
		panic(err)
	}

	resp := RawResponse{string(partial_html)}
	err = tmpl.Execute(buf, resp)
	if err != nil {
		panic(err)
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

func ReturnPartial(c *gin.Context, partial_path string) {
	partial_html, err := os.ReadFile(partial_path)
	if err != nil {
		fmt.Print(err)
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(partial_html))
}

func main() {
    // Create a default gin router
    router := gin.Default()

	router.Static("/assets", "./assets")
	router.Static("/styles", "./styles")
	// I only do this in case I need to load extra content within a page at this
	// subdirectory. Hope this doesn't break anything tho.
	router.Static("/hypertext", "./hypertext")

	router.GET("/", func(c *gin.Context) {
		hx_header := c.Request.Header.Get("Hx-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/index.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/index.html")
		}
	})
	router.GET("/lab-reports", func(c *gin.Context) {
		hx_header := c.Request.Header.Get("Hx-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/lab-reports/reports-index.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/lab-reports/reports-index.html")
		}
	})
	router.GET("/projects", func(c *gin.Context) {
		hx_header := c.Request.Header.Get("Hx-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/projects/projects-index.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/projects/projects-index.html")
		}
	})
	router.GET("/tunes-town", func(c *gin.Context) {
		hx_header := c.Request.Header.Get("Hx-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/tunes-town.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/tunes-town.html")
		}
	})
	router.GET("/gallery", func(c *gin.Context) {
		hx_header := c.Request.Header.Get("Hx-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/gallery.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/tunes-town.html")
		}
	})
	router.GET("/museum", func(c *gin.Context) {
		hx_header := c.Request.Header.Get("Hx-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/museum.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/museum.html")
		}
	})
	router.GET("/adventure-map", func(c *gin.Context) {
		hx_header := c.Request.Header.Get("Hx-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/adventure-map.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/adventure-map.html")
		}
	})

    // Start the server on port 6767
    router.Run(":6767")
}

