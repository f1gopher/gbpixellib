package display

func (s *Screen) LCDEnable() bool {
	return s.memory.ReadBit(lcdcRegister, 7)
}

func (s *Screen) WindowTileMapStart() uint16 {
	if s.memory.ReadBit(lcdcRegister, 6) {
		return 0x9C00
	}

	return 0x9800
}

func (s *Screen) WindowEnable() bool {
	if !s.memory.ReadBit(lcdcRegister, 5) {
		return s.memory.ReadBit(lcdcRegister, 0)
	}

	return true
}

func (s *Screen) BgWindowTileDataArea() uint16 {
	if s.memory.ReadBit(lcdcRegister, 4) {
		return 0x8000
	}

	return 0x8800
}

func (s *Screen) BackgroundTileMapStart() uint16 {
	if s.memory.ReadBit(lcdcRegister, 3) {
		return 0x9C00
	}

	return 0x9800
}

func (s *Screen) ObjSize() byte {
	if s.memory.ReadBit(lcdcRegister, 2) {
		return 16
	}

	return 8
}

func (s *Screen) ObjEnable() bool {
	return s.memory.ReadBit(lcdcRegister, 1)
}

func (s *Screen) BgWindowEnablePriority() bool {
	return s.memory.ReadBit(lcdcRegister, 0)
}

func (s *Screen) LY() byte {
	return s.memory.ReadByte(lcdScanline)
}

func (s *Screen) LYC() byte {
	return s.memory.ReadByte(0xFF45)
}

func (s *Screen) LCDStatusStatInterruptLycLy() bool {
	return s.memory.ReadBit(lcdStatus, 6)
}

func (s *Screen) LCDStatusStatInterruptMode2Oam() bool {
	return s.memory.ReadBit(lcdStatus, 5)
}

func (s *Screen) LCDStatusStatInterruptMode1Vblank() bool {
	return s.memory.ReadBit(lcdStatus, 4)
}

func (s *Screen) LCDStatusStatInterruptMode0Hblank() bool {
	return s.memory.ReadBit(lcdStatus, 3)
}

func (s *Screen) LCDStatusLycLy() bool {
	return s.memory.ReadBit(lcdStatus, 2)
}

func (s *Screen) LCDStatusMode() lcdStatusMode {
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

func (s *Screen) LCDSInterrptEnabled() bool {
	return s.memory.ReadBit(lcdStatus, 6)
}

func (s *Screen) SCY() byte {
	return s.memory.ReadByte(0xFF42)
}

func (s *Screen) SCX() byte {
	return s.memory.ReadByte(0xFF43)
}

func (s *Screen) WY() byte {
	return s.memory.ReadByte(0xFF4A)
}

func (s *Screen) WX() byte {
	return s.memory.ReadByte(0xFF4B)
}

func (s *Screen) BGPIndex3Color() ScreenColor {
	return s.paletteColor(0xFF47, 6)
}

func (s *Screen) BGPIndex2Color() ScreenColor {
	return s.paletteColor(0xFF47, 4)
}

func (s *Screen) BGPIndex1Color() ScreenColor {
	return s.paletteColor(0xFF47, 2)
}

func (s *Screen) BGPIndex0Color() ScreenColor {
	return s.paletteColor(0xFF47, 0)
}

func (s *Screen) ObjPalette0Index3Color() ScreenColor {
	return s.paletteColor(0xFF48, 6)
}

func (s *Screen) ObjPalette0Index2Color() ScreenColor {
	return s.paletteColor(0xFF48, 4)
}

func (s *Screen) ObjPalette0Index1Color() ScreenColor {
	return s.paletteColor(0xFF48, 2)
}

func (s *Screen) ObjPalette0Index0Color() ScreenColor {
	return s.paletteColor(0xFF48, 0)
}

func (s *Screen) ObjPalette1Index3Color() ScreenColor {
	return s.paletteColor(0xFF49, 6)
}

func (s *Screen) ObjPalette1Index2Color() ScreenColor {
	return s.paletteColor(0xFF49, 4)
}

func (s *Screen) ObjPalette1Index1Color() ScreenColor {
	return s.paletteColor(0xFF49, 2)
}

func (s *Screen) ObjPalette1Index0Color() ScreenColor {
	return s.paletteColor(0xFF49, 0)
}

func (s *Screen) paletteColor(address uint16, offset byte) ScreenColor {
	value := s.memory.ReadByte(address)
	value = (value >> offset) & 0b00000011

	if value == 0x00 {
		return White
	} else if value == 0x01 {
		return LightGray
	} else if value == 0x02 {
		return DarkGray
	} else if value == 0x03 {
		return Black
	}

	panic("Invalid palette color")
}
