package cherryHandler

import (
	"context"
	"github.com/cherry-game/cherry/code"
	"github.com/cherry-game/cherry/const"
	"github.com/cherry-game/cherry/extend/reflect"
	facade "github.com/cherry-game/cherry/facade"
	"github.com/cherry-game/cherry/logger"
	cherryMessage "github.com/cherry-game/cherry/net/message"
	cherrySession "github.com/cherry-game/cherry/net/session"
	"reflect"
)

type (
	Handler struct {
		facade.AppContext
		name             string                        // unique name
		eventSlice       map[string][]facade.EventFunc // event func
		localHandlers    map[string]*facade.HandlerFn  // local invoke Handler functions
		remoteHandlers   map[string]*facade.HandlerFn  // remote invoke Handler functions
		handlerComponent *Component                    // handler component
	}
)

func (h *Handler) Name() string {
	return h.name
}

func (h *Handler) SetName(name string) {
	h.name = name
}

func (h *Handler) OnPreInit() {
	h.eventSlice = make(map[string][]facade.EventFunc)
	h.localHandlers = make(map[string]*facade.HandlerFn)
	h.remoteHandlers = make(map[string]*facade.HandlerFn)

	var found = false
	h.handlerComponent, found = h.App().Find(cherryConst.HandlerComponent).(*Component)
	if found == false {
		panic("handler component not found.")
	}
}

func (h *Handler) OnInit() {
}

func (h *Handler) OnAfterInit() {
}

func (h *Handler) OnStop() {
}

func (h *Handler) Events() map[string][]facade.EventFunc {
	return h.eventSlice
}

func (h *Handler) Event(name string) ([]facade.EventFunc, bool) {
	events, found := h.eventSlice[name]
	return events, found
}

func (h *Handler) LocalHandlers() map[string]*facade.HandlerFn {
	return h.localHandlers
}

func (h *Handler) LocalHandler(funcName string) (*facade.HandlerFn, bool) {
	invoke, found := h.localHandlers[funcName]
	return invoke, found
}

func (h *Handler) RemoteHandlers() map[string]*facade.HandlerFn {
	return h.remoteHandlers
}

func (h *Handler) RemoteHandler(funcName string) (*facade.HandlerFn, bool) {
	invoke, found := h.remoteHandlers[funcName]
	return invoke, found
}

func (h *Handler) Component() *Component {
	return h.handlerComponent
}

func (h *Handler) AddLocals(localFns ...interface{}) {
	for _, fn := range localFns {
		funcName := cherryReflect.GetFuncName(fn)
		if funcName == "" {
			cherryLogger.Warnf("get function name fail. fn=%v", fn)
			continue
		}
		h.AddLocal(funcName, fn)
	}
}

func (h *Handler) AddLocal(name string, fn interface{}) {
	f, err := getInvokeFunc(name, fn)
	if err != nil {
		cherryLogger.Warn(err)
		return
	}

	h.localHandlers[name] = f
}

func (h *Handler) AddRemotes(remoteFns ...interface{}) {
	for _, fn := range remoteFns {
		funcName := cherryReflect.GetFuncName(fn)
		if funcName == "" {
			cherryLogger.Warnf("get function name fail. fn=%v", fn)
			continue
		}
		h.AddRemote(funcName, fn)
	}
}

func (h *Handler) AddRemote(name string, fn interface{}) {
	invokeFunc, err := getInvokeFunc(name, fn)
	if err != nil {
		cherryLogger.Warn(err)
		return
	}

	h.remoteHandlers[name] = invokeFunc
}

func getInvokeFunc(name string, fn interface{}) (*facade.HandlerFn, error) {
	invokeFunc, err := cherryReflect.GetInvokeFunc(name, fn)
	if err != nil {
		return invokeFunc, err
	}

	if len(invokeFunc.InArgs) == 2 {
		if invokeFunc.InArgs[1] == reflect.TypeOf(&cherryMessage.Message{}) {
			invokeFunc.IsRaw = true
		}
	}

	return invokeFunc, err
}

func (h *Handler) AddEvent(eventName string, fn facade.EventFunc) {
	if eventName == "" {
		cherryLogger.Warn("eventName is nil")
		return
	}

	if fn == nil {
		cherryLogger.Warn("event function is nil")
		return
	}

	events := h.eventSlice[eventName]
	events = append(events, fn)

	h.eventSlice[eventName] = events
}

func (h *Handler) PostEvent(e facade.IEvent) {
	if h.handlerComponent == nil {
		cherryLogger.Errorf("handlerComponent is nil. [event = %s]", e)
		return
	}

	h.handlerComponent.PostEvent(e)
}

func (h *Handler) AddBeforeFilter(beforeFilters ...FilterFn) {
	if h.handlerComponent != nil {
		h.handlerComponent.AddBeforeFilter(beforeFilters...)
	}
}

func (h *Handler) AddAfterFilter(afterFilters ...FilterFn) {
	if h.handlerComponent != nil {
		h.handlerComponent.AddAfterFilter(afterFilters...)
	}
}

func (h *Handler) ResponseCode(ctx context.Context, session *cherrySession.Session, code int32) {
	codeResult := cherryCode.GetCodeResult(code)
	if cherryCode.IsOK(code) {
		session.Response(ctx, codeResult)
	} else {
		session.Response(ctx, codeResult, true)
	}
}

func (h *Handler) Response(ctx context.Context, session *cherrySession.Session, data interface{}) {
	session.Response(ctx, data)
}
