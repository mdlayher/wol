// +build gofuzz

package wol

func Fuzz(data []byte) int {
	p := new(MagicPacket)
	if err := p.UnmarshalBinary(data); err != nil {
		return 0
	}

	if _, err := p.MarshalBinary(); err != nil {
		panic(err)
	}

	return 1
}
