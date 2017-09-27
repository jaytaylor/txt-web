package interfaces

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gigawattio/errorlib"
	"github.com/gigawattio/web"
	"github.com/gigawattio/web/route"
	"github.com/jaytaylor/html2text"
	"github.com/nbio/hitch"
	"github.com/parnurzeal/gorequest"
)

var AsyncRequestHandlerTimeout time.Duration = 5 * time.Second

type WebService struct {
	*web.WebServer
}

func NewWebService(addr string) *WebService {
	service := &WebService{}
	options := web.WebServerOptions{
		Addr:    addr,
		Handler: service.activateRoutes().Handler(),
	}
	service.WebServer = web.NewWebServer(options)
	return service
}

func (service *WebService) activateRoutes() *hitch.Hitch {
	// index := func(w http.ResponseWriter, req *http.Request) {
	// 	web.RespondWithHtml(w, 200, `<html><head><title>TXT-Web</title></head><body>Welcome to TXT-Web!</body></html>`)
	// }
	routes := []route.RouteMiddlewareBundle{
		route.RouteMiddlewareBundle{
			Middlewares: []func(http.Handler) http.Handler{service.LoggerMiddleware},
			RouteData: []route.RouteDatum{
				// {"get", "/", index},
				{"get", "/*url", service.txt},
			},
		},
	}
	h := route.Activate(routes)
	return h
}

func (service *WebService) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("method=%s url=%s remoteAddr=%s referer=%s\n", req.Method, req.URL.String(), req.RemoteAddr, req.Referer())
		next.ServeHTTP(w, req)
	})
}

func (service *WebService) txt(w http.ResponseWriter, req *http.Request) {
	url := hitch.Params(req).ByName("url")
	fmt.Println(url)

	if len(url) == 1 {
		web.RespondWithHtml(w, 200, `<html><head><title>TXT-Web</title></head><body>Welcome to TXT-Web!</body></html>`)
		return
	}

	if len(url) > 0 {
		url = url[1:]
	}
	url = normalizeURL(url)

	var (
		ch = asyncTXT(url)
		r  result
	)

	select {
	case r = <-ch:
	case <-time.After(AsyncRequestHandlerTimeout):
		r = result{"", fmt.Errorf("timed out after %v", AsyncRequestHandlerTimeout)}
	}

	if r.err != nil {
		web.RespondWithText(w, http.StatusInternalServerError, r.err.Error())
		return
	}
	web.RespondWithText(w, http.StatusOK, r.content)
}

type result struct {
	content string
	err     error
}

func asyncTXT(url string) chan result {
	ch := make(chan result)
	go func() {
		r := gorequest.New()
		resp, data, errs := r.Get(url).End()
		if err := errorlib.Merge(errs); err != nil {
			ch <- result{"", err}
			return
		}
		if resp.StatusCode/100 != 2 {
			err := fmt.Errorf("expected status code 2xx but actual=%v", resp.StatusCode)
			ch <- result{"", err}
			return
		}
		txt, err := html2text.FromString(data, html2text.Options{PrettyTables: true})
		if err != nil {
			ch <- result{"", err}
			return
		}
		ch <- result{txt, nil}
		return
	}()
	return ch
}

func normalizeURL(url string) string {
	if strings.HasPrefix(url, "//") {
		url = "http:" + url
	}
	if !strings.HasPrefix(strings.ToLower(url), "http://") && !strings.HasPrefix(strings.ToLower(url), "https://") {
		url = "http://" + url
	}
	return url
}
