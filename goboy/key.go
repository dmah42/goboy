package goboy

type key struct {
	keys [2]byte
	colidx byte
}

func makeKey() key {
	var k key
	k.keys[0] = 0xF
	k.keys[1] = 0xF
	return k
}

func (k key) ReadByte() byte {
	switch k.colidx {
		case 0x00: return 0x00
		case 0x10: return k.keys[0]
		case 0x20: return k.keys[1]
		default: return 0x00
	}
}

func (k *key) WriteByte(value byte) {
	k.colidx = value & 0x30
}

func (k *key) Keydown(keycode byte) {
	switch keycode {
		case 39: k.keys[1] &= 0xE
		case 37: k.keys[1] &= 0xD
		case 38: k.keys[1] &= 0xB
		case 40: k.keys[1] &= 0x7

		case 90: k.keys[0] &= 0xE
		case 88: k.keys[0] &= 0xD
		case 32: k.keys[0] &= 0xB
		case 13: k.keys[0] &= 0x7
	}
}

func (k *key) Keyup(keycode byte) {
	switch keycode {
	  case 39: k.keys[1] |= 0x1
	  case 37: k.keys[1] |= 0x2
	  case 38: k.keys[1] |= 0x4
	  case 40: k.keys[1] |= 0x8

	  case 90: k.keys[0] |= 0x1
	  case 88: k.keys[0] |= 0x2
	  case 32: k.keys[0] |= 0x5
	  case 13: k.keys[0] |= 0x8
	}
}
