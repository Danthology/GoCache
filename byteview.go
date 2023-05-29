package gocache

type Byteview struct {
	v []byte
}

func (this Byteview) Len() int {
	return len(this.v)
}

func (this Byteview) String() string {
	return string(this.v)
}

func (this Byteview) ByteSlice() []byte {
	return CloneValue(this)
}

func CloneValue(this Byteview) []byte {
	temp := make([]byte, this.Len())
	copy(temp, this.v)
	return temp
}
