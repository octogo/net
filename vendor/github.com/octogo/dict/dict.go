package dict

// Dict is a goroutine-safe wrapper around a map[string]interface{}
type Dict struct {
	data map[string]interface{}

	chDel   chan string
	chSet   chan Pair
	chGet   chan lookup
	chList  chan lookup
	chClose chan struct{}
}

// New returns an initialized *Dict.
func New(d map[string]interface{}) *Dict {
	if d == nil {
		d = make(map[string]interface{})
	}

	return (&Dict{
		data: d,

		chDel:   make(chan string),
		chSet:   make(chan Pair),
		chGet:   make(chan lookup),
		chList:  make(chan lookup),
		chClose: make(chan struct{}),
	}).run()
}

func (d *Dict) run() *Dict {
	go func() {
		defer func() {
			close(d.chDel)
			close(d.chSet)
			close(d.chGet)
			close(d.chList)
			close(d.chClose)
		}()

		for {
			select {
			case <-d.chClose:
				return

			case k := <-d.chDel:
				delete(d.data, k)

			case p := <-d.chSet:
				d.data[p.Key] = p.Value

			case l := <-d.chGet:
				v, found := d.data[l.Key]
				if !found {
					l.Error <- errNotFound
				} else {
					l.Response <- v
				}
				l.Close()

			case l := <-d.chList:
				out := make(map[string]interface{})
				for k, v := range d.data {
					out[k] = v
				}
				l.Response <- out
				l.Close()
			}
		}
	}()

	return d
}

// Del deletes the given key.
func (d Dict) Del(k string) {
	d.chDel <- k
}

// Set sets the given key to the given value.
func (d Dict) Set(k string, v interface{}) {
	d.chSet <- Pair{
		Key:   k,
		Value: v,
	}
}

// Get takes a key and returns a value and a boolean indicating success.
func (d Dict) Get(k string) (interface{}, bool) {
	l := lookup{
		Key:      k,
		Response: make(chan interface{}),
		Error:    make(chan error),
	}

	d.chGet <- l

	select {
	case <-l.Error:
		return nil, false

	case v := <-l.Response:
		return v, true
	}
}

// GetBool takes a key and returns a bool or an error.
func (d Dict) GetBool(k string) (bool, error) {
	v, found := d.Get(k)
	if !found {
		return false, errNotFound
	}
	return v.(bool), nil
}

// GetInt takes a key and returns an int or an error.
func (d Dict) GetInt(k string) (int, error) {
	v, found := d.Get(k)
	if !found {
		return 0, errNotFound
	}
	return v.(int), nil
}

// GetUint takes a key and returns an uint or an error.
func (d Dict) GetUint(k string) (uint, error) {
	v, found := d.Get(k)
	if !found {
		return 0, errNotFound
	}
	return v.(uint), nil
}

// GetString takes a key and returns a string or an error.
func (d Dict) GetString(k string) (string, error) {
	v, found := d.Get(k)
	if !found {
		return "", errNotFound
	}
	return v.(string), nil
}

// GetByte takes a key and returns a byte or an error.
func (d Dict) GetByte(k string) (byte, error) {
	v, found := d.Get(k)
	if !found {
		return 0x0, errNotFound
	}
	return v.(byte), nil
}

// GetBytes takes a key and returns a byte or an error.
func (d Dict) GetBytes(k string) ([]byte, error) {
	v, found := d.Get(k)
	if !found {
		return []byte{}, errNotFound
	}
	return v.([]byte), nil
}

// GetDefault takes a key and a default value and returns the key's value or
// the given default value.
func (d Dict) GetDefault(k string, def interface{}) interface{} {
	v, found := d.Get(k)
	if found {
		return v
	}
	return def
}

// GetDefaultBool takes a key and a bool and returns the value of the key if
// found or the given bool.
func (d Dict) GetDefaultBool(k string, def bool) bool {
	return d.GetDefault(k, def).(bool)
}

// GetDefaultInt takes a key and an int and returns the value of the key if
// found or the given int.
func (d Dict) GetDefaultInt(k string, def int) int {
	return d.GetDefault(k, def).(int)
}

// GetDefaultString takes a key and a string and returns the value of the key
// if found or the given string.
func (d Dict) GetDefaultString(k string, def string) string {
	return d.GetDefault(k, def).(string)
}

// GetDefaultByte takes a key and a byte and a byte and returns the value of
// the key if found or the given byte.
func (d Dict) GetDefaultByte(k string, def byte) byte {
	return d.GetDefault(k, def).(byte)
}

// GetDefaultBytes takes a key and a []byte and returns the value of the key
// if found or the given []byte.
func (d Dict) GetDefaultBytes(k string, def []byte) []byte {
	return d.GetDefault(k, def).([]byte)
}

// Map returns a map[string]interface{}
func (d Dict) Map() map[string]interface{} {
	out := make(map[string]interface{})
	l := lookup{
		Response: make(chan interface{}),
		Error:    make(chan error),
	}
	d.chList <- l
	select {
	case <-l.Error:
		return out
	case v := <-l.Response:
		m, ok := v.(map[string]interface{})
		if !ok {
			return out
		}
		return m
	}
}

// Ingest takes another Dict or map[string]interface{} and copies its values
// over into this Dict. If any other type is passed it returns an error.
func (d Dict) Ingest(other interface{}) error {
	switch other.(type) {
	case map[string]interface{}: // pass
	case *Dict:
		other = other.(*Dict).Map()
	default:
		return errUnknownType
	}
	for k, v := range other.(map[string]interface{}) {
		d.Set(k, v)
	}
	return nil
}

// Close frees the underlying resources.
func (d Dict) Close() {
	d.chClose <- struct{}{}
}
