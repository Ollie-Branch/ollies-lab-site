package main


import (
	"os"
	"text/template"
	"net/http"
	"fmt"
	"bytes"
    "github.com/gin-gonic/gin"
)

// The template that does raw response when not being requested by HTMX
type ContentfulResponse struct {
	Title 	string
	Content	string
}

func ReturnContentfulPage(c *gin.Context, skeleton_path string, partial_path string, doc_title string) {
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

	resp := ContentfulResponse{doc_title, string(partial_html)}
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
	// I don't even know what this controls, but if it is what allows the
	// server to be reverse proxied I set it to localhost and the server's ip
	router.SetTrustedProxies([]string{"127.0.0.1", "93.188.161.205"})

	router.Static("/assets", "./assets")
	router.Static("/styles", "./styles")
	router.Static("/scripts", "./scripts")
	// I only do this in case I need to load extra content within a page at this
	// subdirectory. Hope this doesn't break anything tho.
	router.Static("/hypertext", "./hypertext")

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./assets/favicon.ico")
	})

	/*
	 * Some things to know for HTMX implementers using Gin:
	 * 		Gin's templating system assumes the template is user-generated or
	 * 		comes from an otherwise unsafe source, so we have to do a workaround
	 * 		in the ReturnContentfulPage function to just template the raw text
	 * 		and return it as HTML so HTMX can actually use it. If you're just
	 * 		using HTMX to return fragments of plain text, this isn't a problem,
	 * 		but if you want rich HTML in your HTMX requests you have to do this.
	 *
	 * 		The default behavior is preferred in 99% of cases, because you don't
	 * 		want users uploading <script> tags to your site and injecting
	 * 		malware for other users or your backend (depending on how your
	 * 		backend works and what the script can access), but for HTMX and the
	 * 		backend returning HTML fragments for the page to load, this isn't
	 * 		a feasible way to load content with HTMX.
	 *
	 * 		If I ever add a feature for user-generated text to be uploaded to
	 * 		the site, I will use the default templating system which escapes
	 * 		HTML to prevent cross site scripting and other types of attacks.
	 *
	 * 		Caching is also your biggest enemy, when your browser caches an
	 * 		endpoint, HTMX will also use that cached page in its request to
	 * 		avoid having to make requests to the server. This can result in
	 * 		either HTMX using the whole page as a fragment, or the user not
	 * 		receiving the whole page when starting from a non-root endpoint (if
	 * 		they have already visited that page before using a link/button in
	 * 		HTMX).
	 *
	 * 		The HTMX docs have some extra info on how to deal with caching here:
	 * 		https://htmx.org/docs/#caching
	 */
	router.GET("/", func(c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/index.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/index.html", "Ollie's Lab")
		}
	})
	router.GET("/lab-reports", func(c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/lab-reports/reports-index.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/lab-reports/reports-index.html", "\"Lab Reports\"")
		}
	})
	router.GET("/lab-reports/simple-site-htmx-go", func (c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/lab-reports/simple-site-htmx-go.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/lab-reports/simple-site-htmx-go.html", "Quick and Dirty HTMX+Go")
		}
	})
	router.GET("/projects", func(c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/projects/projects-index.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/projects/projects-index.html", "Projects")
		}
	})
	router.GET("/tunes-town", func(c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/tunes-town.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/tunes-town.html", "Tunes Town")
		}
	})
	router.GET("/gallery", func(c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/gallery.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/gallery.html", "Gallery")
		}
	})
	router.GET("/museum", func(c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/museum.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/museum.html", "Museum")
		}
	})
	router.GET("/adventure-map", func(c *gin.Context) {
		// fix for HTMX cache issues
		// https://htmx.org/docs/#caching
		c.Header("Vary", "HX-Request")
		hx_header := c.Request.Header.Get("HX-Request")
		if(hx_header == "true") {
			ReturnPartial(c, "./hypertext/adventure-map.html")
		} else {
			ReturnContentfulPage(c, "./hypertext/base.html", "./hypertext/adventure-map.html", "Adventure Map")
		}
	})

    // Start the server on port 6767
    router.Run(":6767")
}

