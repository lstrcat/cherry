package cherryProfile

import (
	cfacade "github.com/cherry-game/cherry/facade"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type (
	Config struct {
		jsoniter.Any
	}
)

func Wrap(val interface{}) *Config {
	return &Config{
		Any: jsoniter.Wrap(val),
	}
}

func (p *Config) GetConfig(path ...interface{}) cfacade.ProfileJSON {
	return &Config{
		Any: p.Any.Get(path...),
	}
}

func (p *Config) GetString(path interface{}, defaultVal ...string) string {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}
	return result.ToString()
}

func (p *Config) GetBool(path interface{}, defaultVal ...bool) bool {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}

		return false
	}

	return result.ToBool()
}

func (p *Config) GetInt(path interface{}, defaultVal ...int) int {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}

	return result.ToInt()
}

func (p *Config) GetInt32(path interface{}, defaultVal ...int32) int32 {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}

	return result.ToInt32()
}

func (p *Config) GetInt64(path interface{}, defaultVal ...int64) int64 {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}

	return result.ToInt64()
}

func (p *Config) GetDuration(path interface{}, defaultVal ...time.Duration) time.Duration {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}

	return time.Duration(result.ToInt64())
}

func (p *Config) Marshal(value interface{}) error {
	if p.LastError() != nil {
		return p.LastError()
	}
	return jsoniter.UnmarshalFromString(p.ToString(), value)
}
