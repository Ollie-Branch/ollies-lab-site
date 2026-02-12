package main


import (
	"text/template"
	"net/http"
	// "fmt"
	"bytes"
    "github.com/gin-gonic/gin"
)

// The template that does raw response when not using HTMX
type RawResponse struct {
	Content string
}

func main() {
    // Create a default gin router
    router := gin.Default()
	resp := RawResponse{"<p> \n \t Paragraph inside the content field \n </p>"}
	buf := new(bytes.Buffer)
	t1, err := template.New("t1").Parse("<h1> \n Test Content \n </h1> \n {{.Content}}")
	if err != nil {
		panic(err);
	}
	err = t1.Execute(buf, resp)
	if err != nil {
		panic(err);
	}

	// outputString := fmt.Sprintf("%s", buf)


    // Define a route for the root path
    router.GET("/", func(c *gin.Context) {
        c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
		})

    // Start the server on port 8080
    router.Run(":8080")
}

