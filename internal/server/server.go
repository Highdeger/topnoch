package server

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"reflect"
	"strings"
)

func RunMainServer(address *string, compress *bool) {
	addr := *address
	if strings.HasPrefix(*address, ":") {
		addr = "localhost" + *address
	}

	for {
		r := router.New()
		r.ServeFiles("/manage/assets/{filepath:*}", "web/assets")
		r.ServeFiles("/assets/{filepath:*}", "web/assets")
		r.ServeFiles("/{filepath:*}", "web/assets/images")

		r.Handle(fasthttp.MethodGet, "/", PageLogin)
		r.Handle(fasthttp.MethodGet, "/manage/{StructName:*}", PageDataManager)
		r.Handle(fasthttp.MethodGet, "/discovery", PageDiscovery)
		r.Handle(fasthttp.MethodGet, "/node/realtime", PageNodeRealtime)

		r.Handle(fasthttp.MethodGet, "/all", ApiAllData)
		r.Handle(fasthttp.MethodGet, "/one", ApiOneData)
		r.Handle(fasthttp.MethodPost, "/add", ApiAddData)
		r.Handle(fasthttp.MethodDelete, "/delete", ApiDeleteData)
		r.Handle(fasthttp.MethodPut, "/edit", ApiEditData)

		r.Handle(fasthttp.MethodGet, "/authenticate", ApiAuthenticate)

		r.Handle(fasthttp.MethodGet, "/get/fields/all", ApiGetFieldsAll)
		r.Handle(fasthttp.MethodGet, "/get/fields/edit", ApiGetFieldsEdit)
		r.Handle(fasthttp.MethodPost, "/discovery/start", ApiDiscoveryStart)
		r.Handle(fasthttp.MethodGet, "/terminate", ApiTerminate)

		r.Handle(fasthttp.MethodGet, "/param/all", ApiParamGetAll)
		r.Handle(fasthttp.MethodGet, "/param/last", ApiParamGetLast)

		handler := r.Handler
		if *compress {
			handler = fasthttp.CompressHandler(handler)
		}

		log.Printf("Server is running on %s\n", addr)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Fatalf("Error in ListenAndServe: (Type: %s) -> %s", reflect.TypeOf(err), err)
		}
	}
}
