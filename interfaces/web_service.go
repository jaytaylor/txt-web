package interfaces

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gigawattio/errorlib"
	"github.com/gigawattio/web"
	"github.com/gigawattio/web/route"
	"github.com/jaytaylor/html2text"
	"github.com/nbio/hitch"
	"github.com/parnurzeal/gorequest"
)

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
	index := func(w http.ResponseWriter, req *http.Request) {
		web.RespondWithHtml(w, 200, `<html><head><title>hello world</title></head><body>hello world</body></html>`)
	}
	routes := []route.RouteMiddlewareBundle{
		route.RouteMiddlewareBundle{
			Middlewares: []func(http.Handler) http.Handler{service.LoggerMiddleware},
			RouteData: []route.RouteDatum{
				{"get", "/", index},
				{"get", "/v1/:url", service.txt},
				// {"post", "/v1/tesseract/*url", service.tesseractUrl},
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
	// generics.GenericObjectEndpoint(w, req, func() (object interface{}, err error) {
	obj, err := func() (interface{}, error) {
		url := normalizeURL(hitch.Params(req).ByName("url"))
		r := gorequest.New()
		resp, data, errs := r.Get(url).End()
		if err := errorlib.Merge(errs); err != nil {
			return nil, err
		}
		if resp.StatusCode/100 != 2 {
			return nil, fmt.Errorf("expected status code in the 2xx range but actual=%v", resp.StatusCode)
		}
		txt, err := html2text.FromString(data, html2text.Options{PrettyTables: true})
		if err != nil {
			return nil, err
		}
		return []byte(txt), nil
		// return []byte("<html><body>hi " + url), nil
	}()
	if err != nil {
		web.RespondWithText(w, http.StatusInternalServerError, err.Error())
		return
	}
	web.RespondWithText(w, http.StatusOK, obj)
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
