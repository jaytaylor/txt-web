package main

import (
	"bytes"
	"os"

	"github.com/jaytaylor/txt-web/interfaces"

	"github.com/gigawattio/web/cli"
	webinterfaces "github.com/gigawattio/web/interfaces"
	cliv2 "gopkg.in/urfave/cli.v2"
)

const (
	AppName = "txt-web"
)

func webServiceProvider(ctx *cliv2.Context) (webinterfaces.WebService, error) {
	webService := interfaces.NewWebService(ctx.String("bind"))
	return webService, nil
}

func getVersion() string {
	buf := &bytes.Buffer{}
	DisplayVersion(buf, "\n")
	v := buf.String()
	return v
}

func main() {
	options := cli.Options{
		AppName:            AppName,
		Usage:              "HTML2Text web service",
		Version:            getVersion(),
		WebServiceProvider: webServiceProvider,
		ExitOnError:        true,
	}
	c, err := cli.New(options)
	if err != nil {
		cli.ErrorExit(os.Stderr, err, 1)
	}
	c.Main()
}
