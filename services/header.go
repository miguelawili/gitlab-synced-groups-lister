package services

type Header struct {
	Key   string
	Value string
}

func NewHeader(key string, value string) Header {
	return Header{
		Key:   key,
		Value: value,
	}
}

func (h *Header) GetKey() string {
	return h.Key
}

func (h *Header) GetValue() string {
	return h.Value
}
