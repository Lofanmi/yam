package curl

import (
	"net/url"

	"github.com/Lofanmi/yam/api"
)

type curlInit struct {
	request api.PHPRequest
}

func NewInit(request api.PHPRequest) api.HookCallback {
	return &curlInit{request: request}
}

func (s *curlInit) Before(data api.ExecuteData, returnValue api.ZVal) {}

func (s *curlInit) After(data api.ExecuteData, returnValue api.ZVal) {
	if !returnValue.IsResource() {
		return
	}
	id := returnValue.AsResource().ID()
	if id <= 0 {
		return
	}
	curl := s.request.Curl(id)
	args := data.Args()
	if len(args) == 1 && args[0].IsString() {
		us := args[0].AsString()
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
}
