package headers

// Size gets the size of the header and ignores errors
func Size()uint64{
	dummyHeader := Header{}
	var size int
	size, _ = dummyHeader.MarshalLen()
	return uint64(size)
}