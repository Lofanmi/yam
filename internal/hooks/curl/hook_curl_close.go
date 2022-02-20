package curl

import (
	"net/url"
	"unsafe"

	"github.com/Lofanmi/yam/api"
	"github.com/Lofanmi/yam/internal/zend"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type curlClose struct {
	request api.PHPRequest
}

func NewClose(request api.PHPRequest) api.HookCallback {
	return &curlClose{request: request}
}

func getZValString(m map[string]api.ZVal, key string, ifEmpty string) string {
	if v, ok := m[key]; ok {
		if v != nil && v.IsString() {
			return v.AsString()
		}
	}
	return ifEmpty
}

func getZValInt(m map[string]api.ZVal, key string, ifEmpty int) int {
	if v, ok := m[key]; ok {
		if v != nil && v.IsInt() {
			return v.AsInt()
		}
	}
	return ifEmpty
}

func getZValFloat(m map[string]api.ZVal, key string, ifEmpty float64) float64 {
	if v, ok := m[key]; ok {
		if v != nil && v.IsFloat() {
			return v.AsFloat()
		}
	}
	return ifEmpty
}

func (s *curlClose) Before(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) != 1 || !args[0].IsResource() {
		return
	}
	id := args[0].AsResource().ID()
	if id <= 0 {
		return
	}
	curl := s.request.Curl(id)
	api.CurlGetInfo(data.Pointer(), func(ret unsafe.Pointer) {
		var m map[string]api.ZVal
		zv := zend.NewZVal(ret)
		if !zv.IsArray() {
			return
		}
		m = zv.AsArray().ToSZMap()
		if len(m) <= 0 {
			return
		}
		curl.SetUrl(getZValString(m, "url", ""))
		curl.SetContentType(getZValString(m, "content_type", ""))
		curl.SetHttpCode(getZValInt(m, "http_code", 0))
		curl.SetTotalTime(getZValFloat(m, "total_time", 0.0))
		curl.SetNameLookupTime(getZValFloat(m, "namelookup_time", 0.0))
		curl.SetConnectTime(getZValFloat(m, "connect_time", 0.0))
		curl.SetPreTransferTime(getZValFloat(m, "pretransfer_time", 0.0))
		curl.SetStartTransferTime(getZValFloat(m, "starttransfer_time", 0.0))
		switch getZValInt(m, "http_version", 0) {
		case 1:
			curl.SetHttpVersion(semconv.HTTPFlavorHTTP10.Value.AsString())
		case 2:
			curl.SetHttpVersion(semconv.HTTPFlavorHTTP11.Value.AsString())
		}
		curl.SetRemoteIP(getZValString(m, "primary_ip", ""))
		curl.SetRemotePort(getZValInt(m, "primary_port", 0))

		redirectURL := getZValString(m, "redirect_url", "")
		if redirectURL != "" {
			if u, err := url.Parse(redirectURL); err == nil {
				curl.SetUrl(redirectURL)
				curl.SetScheme(u.Scheme)
				curl.SetHost(u.Hostname())
				if u.Port() == "" {
					if u.Scheme == "https" {
						curl.SetPort("443")
					} else {
						curl.SetPort(u.Port())
					}
				} else {
					curl.SetPort(u.Port())
				}
				curl.SetPath(u.Path)
				curl.SetQuery(u.RawQuery)
			}
		}
	})
	curl.EndTrace(s.request, id)
}

func (s *curlClose) After(data api.ExecuteData, returnValue api.ZVal) {}
