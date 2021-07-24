package server

import (
	"../core/xhtml"
	nodeManager "../node_auto_discovery/node_alive_manager"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
)

func PageLogin(ctx *fasthttp.RequestCtx) {
	fmt.Println(GetRequestText(ctx))

	CheckThenRun(ctx, "text/html", func() {
		if bytes, err := ioutil.ReadFile("web/login.html"); err == nil {
			content := string(bytes)
			WriteContent(ctx, content)
		} else {
			ctx.Error("can't read the static files", 501)
		}
	})
}

func PageDataManager(ctx *fasthttp.RequestCtx) {
	fmt.Println(GetRequestText(ctx))

	structName := fmt.Sprintf("%s", ctx.UserValue("StructName"))
	CheckThenRun(ctx, "text/html", func() {
		if bytes, err := ioutil.ReadFile("web/data_manager.html"); err == nil {
			content := string(bytes)

			elements := make(map[string]string)
			elements["PageTitle"] = fmt.Sprintf("%s Manager - TopNoch", structName)
			elements["TablePageTitle"] = "Data Manager"
			elements["TableCardTitle"] = structName
			elements["ModalFields"] = xhtml.CreateModalFields(structName)
			elements["StructName"] = structName
			elements["ColumnDefs"] = xhtml.CreateColumnDefs(structName)

			elements["RequestHeaders"] = xhtml.CreateRequestHeaders(structName)

			content = PutElements(content, elements)
			WriteContent(ctx, content)
		} else {
			ctx.Error("can't read the static xfile", 401)
		}
	})
}

func PageDiscovery(ctx *fasthttp.RequestCtx) {
	fmt.Println(GetRequestText(ctx))

	pageTitle := "Discovery - TopNoch"
	tablePageTitle := "Device Discoverer"
	tableCardTitle := "Discovering devices on the xnetwork"
	CheckThenRun(ctx, "text/html", func() {
		if bytes, err := ioutil.ReadFile("web/discovery.html"); err == nil {
			content := string(bytes)

			elements := make(map[string]string)
			elements["PageTitle"] = pageTitle
			elements["TablePageTitle"] = tablePageTitle
			elements["TableCardTitle"] = tableCardTitle

			content = PutElements(content, elements)
			WriteContent(ctx, content)
		} else {
			ctx.Error("can't read the static xfile", 401)
		}
	})
}

func PageNodeRealtime(ctx *fasthttp.RequestCtx) {
	fmt.Println(GetRequestText(ctx))

	nodeKey := string(ctx.Request.Header.Peek("NodeKey"))

	fmt.Println("node realtime page")

	pageTitle := "Node Realtime - TopNoch"
	inPageTitle := fmt.Sprintf("Node Realtime Data (%s)", nodeKey)

	fmt.Println("nodeKey", nodeKey)
	nodeManager.NodeAdd(nodeKey)
	nodeManager.NodeGetOne(nodeKey).AliveCheckStart(false)

	CheckThenRun(ctx, "text/html", func() {
		if bytes, err := ioutil.ReadFile("web/node_realtime.html"); err == nil {
			content := string(bytes)

			elements := make(map[string]string)
			elements["PageTitle"] = pageTitle
			elements["InPageTitle"] = inPageTitle

			content = PutElements(content, elements)
			WriteContent(ctx, content)
		} else {
			ctx.Error("can't read the static xfile", 401)
		}
	})
}
