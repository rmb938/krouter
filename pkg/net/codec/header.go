package codec

type RequestHeader struct {
	Key           int16
	Version       int16
	CorrelationId int32
	ClientID      *string
}

func (h *RequestHeader) Decode(pd *RawDecoder) (err error) {
	if h.Key, err = pd.Int16(); err != nil {
		return err
	}

	if h.Version, err = pd.Int16(); err != nil {
		return err
	}

	if h.CorrelationId, err = pd.Int32(); err != nil {
		return err
	}

	if h.ClientID, err = pd.NullableString(); err != nil {
		return err
	}

	// if h.Key == 18 || h.Version >= 2 {
	// 	// tagged fields
	// 	taggedFields, err := pd.UVarInt()
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	fmt.Printf("Tagged Fields: %v\n", taggedFields)
	//
	// }

	return
}

type ResponseHeader struct {
	Length        int32
	CorrelationId int32
	Version       int16
}

func (h *ResponseHeader) Encode(pe *RawEncoder) (err error) {
	length := h.Length + 4

	// if h.Version >= 1 {
	// 	length += 1
	// }

	pe.Int32(length)
	pe.Int32(h.CorrelationId)

	// if h.Version >= 1 {
	// 	// tagged fields
	// 	pe.UVarInt(0)
	// }

	return nil
}
