package session

import (
	"../../core/xlog"
	"github.com/fasthttp/session"
	"github.com/fasthttp/session/providers/memory"
	"github.com/valyala/fasthttp"
	"time"
)

var (
	sessionMap                     = make(map[string]time.Time)
	sessionServer *session.Session = nil
)

func FetchSessionStore(ctx *fasthttp.RequestCtx) *session.Store {
	if sessionServer == nil {
		config := session.NewDefaultConfig()
		config.EncodeFunc = session.MSGPEncode
		config.DecodeFunc = session.MSGPDecode
		sessionServer = session.New(config)

		provider, err := memory.New(memory.Config{})
		if err != nil {
			xlog.LogFatal(err.Error())
		}
		err = sessionServer.SetProvider(provider)
		if err != nil {
			xlog.LogFatal(err.Error())
		}
	}

	store, e := sessionServer.Get(ctx)
	if e != nil {
		xlog.LogFatal(e.Error())
	} else {
		store.Set("is_auth", false)
		id := string(store.GetSessionID())
		sessionMap[id] = time.Now().Add(store.GetExpiration())
		return store
	}
	return nil
}

func IsExpired(id string) bool {
	for k, v := range sessionMap {
		if k == id {
			if time.Now().Before(v) {
				return false
			} else {
				return true
			}
		}
	}
	return true
}
