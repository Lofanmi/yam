package curl

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/Lofanmi/yam/api"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/spf13/cast"
)

const (
	optionHttpGet       = 80    // CURLOPT_HTTPGET
	optionPost          = 47    // CURLOPT_POST
	optionHttpVersion   = 84    // CURLOPT_HTTP_VERSION
	optionCustomRequest = 10036 // CURLOPT_CUSTOMREQUEST
	optionURL           = 10002 // CURLOPT_URL
	optionPostFields    = 10015 // CURLOPT_POSTFIELDS
	optionHttpHeader    = 10023 // CURLOPT_HTTPHEADER
)

type curlSetOption struct {
	request api.PHPRequest
}

func NewSetOption(request api.PHPRequest) api.HookCallback {
	return &curlSetOption{request: request}
}

func (s *curlSetOption) Before(data api.ExecuteData, returnValue api.ZVal) {}

func (s *curlSetOption) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) != 3 {
		return
	}
	_handle, _options, value := args[0], args[1], args[2]
	if (!_handle.IsResource()) || (!_options.IsInt()) {
		return
	}
	id := _handle.AsResource().ID()
	if id <= 0 {
		return
	}
	curl := s.request.Curl(id)

	options := _options.AsInt()
	parseOption(curl, options, value)
}

func parseOption(curl api.Curl, options int, value api.ZVal) {
	switch options {
	case optionHttpGet:
		if value.IsBool() && value.AsBool() {
			curl.SetMethod(http.MethodGet)
		}
	case optionPost:
		if value.IsBool() && value.AsBool() {
			curl.SetMethod(http.MethodPost)
		}
	case optionHttpVersion:
		if value.IsInt() {
			switch value.AsInt() {
			case 1:
				curl.SetHttpVersion(semconv.HTTPFlavorHTTP10.Value.AsString())
			case 2:
				curl.SetHttpVersion(semconv.HTTPFlavorHTTP11.Value.AsString())
			}
		}
	case optionCustomRequest:
		curl.SetMethod(parseMethodByOptionCustomRequest(value))
	case optionURL:
		if value.IsString() {
			us := value.AsString()
			if u, err := url.Parse(us); err == nil {
				curl.SetUrl(us)
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
	case optionPostFields:
		var body string
		if value.IsString() {
			body = value.AsString()
		} else if value.IsArray() {
			var query url.Values
			for key, val := range value.AsArray().ToIIMap() {
				query.Add(cast.ToString(key), cast.ToString(val))
			}
			body = query.Encode()
		}
		curl.SetBody(body)
		curl.SetBodySize(len(body))
	case optionHttpHeader:
		if value.IsArray() {
			for _, val := range value.AsArray().ToIIMap() {
				curl.AppendHeader(cast.ToString(val))
			}
		}
	}
}

func parseMethodByOptionCustomRequest(value api.ZVal) (method string) {
	m := strings.ToUpper(value.AsString())
	switch m {
	case http.MethodGet:
		fallthrough
	case http.MethodHead:
		fallthrough
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodConnect:
		fallthrough
	case http.MethodOptions:
		fallthrough
	case http.MethodTrace:
		method = m
	}
	return
}
