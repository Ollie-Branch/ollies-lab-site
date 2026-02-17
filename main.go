package main


import (
	"os"
	"text/template"
	"net/http"
	"fmt"
	"bytes"
	"log"
	"github.com/pelletier/go-toml"
    "github.com/gin-gonic/gin"
)

// The template that does raw response when not being requested by HTMX
type ContentfulResponse struct {
	Title 	string
	Content	string
}

type Config struct {
	TrustedProxies		[]string 	`toml:"trusted_proxies"`
	Page []struct {
		URL				string 		`toml:"url"`
		SkeletonPath	string 		`toml:"skeleton-path"`
		ContentPath		string 		`toml:"content-path"`
		Title			string 		`toml:"title"`
	}
}

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

func ServeHTMXContent(c *gin.Context, skel_path string, content_path string,
						title string) {
	// https://htmx.org/docs/#caching
	c.Header("Vary", "HX-Request")
	hx_header := c.Request.Header.Get("HX-Request")
	if(hx_header == "true") {
		ReturnPartial(c, content_path)
	} else {
		ReturnContentfulPage(c, skel_path, content_path, title)
	}
}

func main() {
	var config Config
	configFile, err := os.ReadFile("./config.toml")
	if err != nil {
		log.Fatal(err)
	}
	err = toml.Unmarshal([]byte(configFile), &config)

    router := gin.Default()
	// I don't even know what this controls, but if it is what allows the
	// server to be reverse proxied I set it to localhost and the server's ip
	// This is now set in the config file "config.toml"
	router.SetTrustedProxies(config.TrustedProxies)

	router.Static("/assets", "./assets")
	router.Static("/styles", "./styles")
	router.Static("/scripts", "./scripts")

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./assets/favicon.ico")
	})

	for i := range config.Page {
		router.GET(config.Page[i].URL, func(c *gin.Context) {
			ServeHTMXContent(c, config.Page[i].SkeletonPath,
				config.Page[i].ContentPath, config.Page[i].Title)
		})
	}

    router.Run(":6767")
}

