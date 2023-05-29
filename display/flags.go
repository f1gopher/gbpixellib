package display

func (s *Screen) lcdEnable() bool {
	return s.memory.ReadBit(lcdcRegister, 7)
}

func (s *Screen) windowTileMapStart() uint16 {
	if s.memory.ReadBit(lcdcRegister, 6) {
		return 0x9C00
	}

	return 0x9800
}

func (s *Screen) windowEnable() bool {
	return s.memory.ReadBit(lcdcRegister, 5)
}

func (s *Screen) bgWindowTileDataArea() uint16 {
	if s.memory.ReadBit(lcdcRegister, 4) {
		return 0x8000
	}

	return 0x8800
}

func (s *Screen) bgTileMapArea(bit byte) uint16 {
	if s.memory.ReadBit(lcdcRegister, bit) {
		return 0x9C00
	}

	return 0x9800
}

func (s *Screen) objSize() byte {
	if s.memory.ReadBit(lcdcRegister, 2) {
		return 16
	}

	return 8
}

func (s *Screen) objEnable() bool {
	return s.memory.ReadBit(lcdcRegister, 1)
}

func (s *Screen) bgWindowEnablePriority() bool {
	return s.memory.ReadBit(lcdcRegister, 0)
}

func (s *Screen) ly() byte {
	return s.memory.ReadByte(lcdScanline)
}

func (s *Screen) lyc() byte {
	return s.memory.ReadByte(0xFF45)
}

func (s *Screen) lcdStatus_STAT_Interrupt_LYC_LY() bool {
	return s.memory.ReadBit(lcdStatus, 6)
}

func (s *Screen) lcdStatus_STAT_Interrupt_Mode2_OAM() bool {
	return s.memory.ReadBit(lcdStatus, 5)
}

func (s *Screen) lcdStatus_STAT_Interrupt_Mode1_VBlank() bool {
	return s.memory.ReadBit(lcdStatus, 4)
}

func (s *Screen) lcdStatus_STAT_Interrupt_Mode0_HBlank() bool {
	return s.memory.ReadBit(lcdStatus, 3)
}

func (s *Screen) lcdStatus_LYC_LY() bool {
	return s.memory.ReadBit(lcdStatus, 2)
}

func (s *Screen) lcdStatus_Mode() lcdStatusMode {
	value := s.memory.ReadByte(lcdStatus)

	if value&0x0000 == 0x0000 {
		return hblank
	} else if value&0x0001 == 0x0001 {
		return vblank
	} else if value&0x0002 == 0x0002 {
		return searchOAM
	} else {
		return transferringToController
	}
}

func (s *Screen) scy() byte {
	return s.memory.ReadByte(0xFF42)
}

func (s *Screen) scx() byte {
	return s.memory.ReadByte(0xFF43)
}

func (s *Screen) wy() byte {
	return s.memory.ReadByte(0xFF4A)
}

func (s *Screen) wx() byte {
	return s.memory.ReadByte(0xFF4B)
}

func (s *Screen) bgpIndex3Color() screenColor {
	return s.bgpColor(6)
}

func (s *Screen) bgpIndex2Color() screenColor {
	return s.bgpColor(4)
}

func (s *Screen) bgpIndex1Color() screenColor {
	return s.bgpColor(2)
}

func (s *Screen) bgpIndex0Color() screenColor {
	return s.bgpColor(0)
}

func (s *Screen) bgpColor(offset byte) screenColor {
	value := s.memory.ReadByte(0xFF47)

	value = value >> offset

	if value&0x00 == 0x00 {
		return white
	} else if value&0x01 == 0x01 {
		return lightGray
	} else if value&0x02 == 0x02 {
		return darkGray
	} else {
		return black
	}
}

func (s *Screen) obp0() byte {
	return s.memory.ReadByte(0xFF48)
}

func (s *Screen) obp1() byte {
	return s.memory.ReadByte(0xFF49)
}
