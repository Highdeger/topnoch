package server

import (
	"../core/xlist"
	"../server/session"
	"fmt"
	"github.com/valyala/fasthttp"
	"strings"
)

var (
	allowedHosts   = []string{"localhost", "127.0.0.1"}
	allowedSources = []string{"127.0.0.1"}
)

func GetRequestText(ctx *fasthttp.RequestCtx) string {
	requestMethod := fmt.Sprintf("%s", ctx.Method())
	requestUri := fmt.Sprintf("%s", ctx.RequestURI())
	requestPath := fmt.Sprintf("%s", ctx.Path())
	requestHost := fmt.Sprintf("%s", ctx.Host())
	requestQueryArgs := fmt.Sprintf("%s", ctx.QueryArgs())
	//requestUserAgent := fmt.Sprintf("%s", ctx.UserAgent())
	//requestConnTime := fmt.Sprintf("%s", ctx.ConnTime())
	//requestTime := fmt.Sprintf("%s", ctx.Time())
	requestConnNum := fmt.Sprintf("%d", ctx.ConnRequestNum())
	requestRemoteIp := fmt.Sprintf("%s", ctx.RemoteIP())

	//return fmt.Sprintf("(%s) uri:'%s' path:'%s' host:'%s' args:'%s ' src:'%s' useragent:'%s' connnum:'%s'", requestMethod, requestUri, requestPath, requestHost, requestQueryArgs, requestRemoteIp, requestUserAgent, requestConnNum)
	return fmt.Sprintf("(%s) uri:'%s' path:'%s' host:'%s' args:'%s ' src:'%s' connnum:'%s'", requestMethod, requestUri, requestPath, requestHost, requestQueryArgs, requestRemoteIp, requestConnNum)
}

// contentType examples: "text/html" or "text/json"
func CheckThenRun(ctx *fasthttp.RequestCtx, contentType string, run func()) {
	store := session.FetchSessionStore(ctx)
	if session.IsExpired(string(store.GetSessionID())) {
		ctx.Error("session expired", fasthttp.StatusUnauthorized)
		return
	} else {
		requestHost := fmt.Sprintf("%s", ctx.Host())
		hostname := strings.Split(requestHost, ":")[0]
		requestRemoteIp := fmt.Sprintf("%s", ctx.RemoteIP())

		ctx.SetContentType(fmt.Sprintf("%s; charset=utf8", contentType))

		if xlist.ItemIndex(hostname, &allowedHosts) != -1 {
			if xlist.ItemIndex(requestRemoteIp, &allowedHosts) != -1 {
				run()
			} else {
				ctx.Error(fmt.Sprintf("Access from your IP '%s', is not allowed. If you really need access, please contact your xnetwork's admin.", requestRemoteIp), 403)
			}
		} else {
			ctx.Error(fmt.Sprintf("Access by the name '%s' for host, is not allowed, if you really need access, please contact your xnetwork's admin.", hostname), 403)
		}

		// set additional headers
		//ctx.Response.Header.Set("X-My-Header", "my-header-value")

		// set cookies
		//var cookie fasthttp.Cookie
		//cookie.SetKey("cookie-name")
		//cookie.SetValue("cookie-value")
		//ctx.Response.Header.SetCookie(&cookie)
	}
}

func PutElements(input string, elements map[string]string) (output string) {
	output = input
	for k, v := range elements {
		output = strings.ReplaceAll(output, fmt.Sprintf("/*{{%s}}*/", k), fmt.Sprintf("%s", v))
	}
	return
}

func WriteContent(ctx *fasthttp.RequestCtx, content string) {
	if _, err := fmt.Fprint(ctx, content); err != nil {
		// add error in coreLog (connection lost)
	} else {
		// add response in coreLog (responded by how much bytes)
	}
}

func JsonError(msg string) string {
	return fmt.Sprintf("{\"err\":\"%s\"}", msg)
}
