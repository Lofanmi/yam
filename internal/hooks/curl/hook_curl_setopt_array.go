package curl

import (
	"github.com/Lofanmi/yam/api"
)

type curlSetOptionArray struct {
	request api.PHPRequest
}

func NewSetOptionArray(request api.PHPRequest) api.HookCallback {
	return &curlSetOptionArray{request: request}
}

func (s *curlSetOptionArray) Before(data api.ExecuteData, returnValue api.ZVal) {
}

func (s *curlSetOptionArray) After(data api.ExecuteData, returnValue api.ZVal) {
	args := data.Args()
	if len(args) != 2 {
		return
	}
	_handle, value := args[0], args[1]
	if (!_handle.IsResource()) || (!value.IsArray()) {
		return
	}
	id := _handle.AsResource().ID()
	if id <= 0 {
		return
	}
	curl := s.request.Curl(id)

	m := value.AsArray().ToIZMap()
	for key, val := range m {
		if options, ok := key.(int); ok {
			parseOption(curl, options, val)
		}
	}
}
