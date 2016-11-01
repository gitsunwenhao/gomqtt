//header模块单元测试
package protocol

import (
	"testing"
)

func Test_Type(t *testing.T) {
	h := &header{
		remLen:   10000,
		typeFlag: 108,
		packetID: 65535,
	}
	h.SetType(CONNACK)

	if h.Type() != CONNACK {
		t.Errorf("test header type failed,expected: %08b, got: %08b", CONNACK, h.Type())
	}
}

func Test_Flags(t *testing.T) {
	h := &header{
		remLen:   10000,
		typeFlag: 108,
		packetID: 65535,
	}

	if h.Flags() != 12 {
		t.Errorf("test header flag failed,expected: %08b, got: %08b", 12, h.Flags())
	}
}

func Test_RemainLength(t *testing.T) {
	h := &header{
		remLen:   10000,
		typeFlag: 108,
		packetID: 65535,
	}

	var target = []struct {
		value int32
		suc   bool
	}{
		{-1, false},
		{0, true},
		{65535, true},
		{268435456, false},
	}

	for _, v := range target {
		h.SetRemainingLength(v.value)
		if (h.RemainingLength() == v.value) != v.suc {
			t.Errorf("test header remain length failed, expected : %d, got: %d", v.value, h.RemainingLength())
		}
	}
}

func Test_Len(t *testing.T) {
	h := &header{
		remLen:   10000,
		typeFlag: 108,
		packetID: 65535,
	}

	var target = []struct {
		value  int32
		expect int
	}{
		{127, 2},
		{128, 3},
		{16383, 3},
		{16384, 4},
		{2097151, 4},
		{2097152, 5},
		{268435455, 5},
	}

	for _, v := range target {
		h.SetRemainingLength(v.value)
		if h.Len() != v.expect {
			t.Errorf("test header  length failed, expected : %d, got: %d", v.expect, h.Len())
		}
	}
}

func Test_PackeID(t *testing.T) {
	h := &header{
		remLen:   10000,
		typeFlag: 108,
		packetID: 65535,
	}

	h.SetPacketID(65531)
	if h.PacketID() != 65531 {
		t.Errorf("test header packet ID failed, expected %d, got %d", 65531, h.PacketID())
	}

	var v = 65537 //0001 0000 0000 0000 0001
	h.SetPacketID(uint16(v))
	if h.PacketID() != 1 {
		t.Errorf("test header packet ID failed, expected %d, got %d", 1, h.PacketID())
	}
}

func Test_Encode(t *testing.T) {
	h := &header{
		remLen:   10000,
		typeFlag: 108,
		packetID: 65535,
	}

	//[16 144 78]
	h.SetType(CONNECT)
	expect := []byte{16, 144, 78}

	buf := make([]byte, 3)
	h.encode(buf)

	for i := 0; i < 3; i++ {
		if buf[i] != expect[i] {
			t.Errorf("test header encode, expected %v, got %v", expect, buf)
		}
	}
}
