// Code generated by "spawner -pool=true"; DO NOT EDIT.

package msg

type factory func() interface{}

func Spawner(name string, newPool ...bool) (interface{}, bool) {
	if len(newPool) > 0 && newPool[0] {
		p, ok := spawnerPools[name]
		if !ok {
			return nil, false
		}
		return p.Get(), true
	}
	f, ok := spawner[name]
	if !ok {
		return nil, ok
	}
	return f(), true
}

func Put(name string, x interface{}) {
	pool, ok := spawnerPools[name]
	if !ok {
		return
	}
	pool.Put(x)
}

var spawner = map[string]factory{
	"msg.LileiSay":     func() interface{} { return &LileiSay{} },
	"msg.HanMeimeiSay": func() interface{} { return &HanMeimeiSay{} },
}