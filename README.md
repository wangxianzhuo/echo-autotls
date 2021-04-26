# echo-autotls
HTTPS server using certificates automatically installed from https://letsencrypt.org

## Archived

The problem [#1231](https://github.com/labstack/echo/issues/1231) has been solved [#1305](https://github.com/labstack/echo/pull/1305)


## Usage

```go
func main() {
	e := echo.New()
	e.Pre(middleware.HTTPSRedirect())

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<h1>Welcome to Echo!</h1>
			<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
		`)
	})
	e.Logger.Fatal(e.StartServer(autotls.DefaultManager("example.com").StartAutoTLS(":443")))

	// m := autotls.AutoTLSManager{}
	// m.Prompt = autocert.AcceptTOS
	// m.Cache = autocert.DirCache("/var/www/.cache")
	// e.Logger.Fatal(e.StartServer(m.StartAutoTLS(":443")))
}
```
