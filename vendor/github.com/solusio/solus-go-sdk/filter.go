package solus

import "strconv"

type filter struct {
	data map[string]string
}

func (f *filter) add(k, v string) {
	if f.data == nil {
		f.data = map[string]string{}
	}

	f.data[k] = v
}

func (f *filter) addInt(k string, v int) {
	f.add(k, strconv.Itoa(v))
}
