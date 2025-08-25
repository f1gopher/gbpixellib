package sound

func (s *Sound) extract(address uint16, mask uint8, shift uint8) uint8 {
	source := s.memory.ReadByte(address)
	selected := source & mask
	return selected >> shift
}

func (s *Sound) Sqaure1_SweepPeriod() uint8 {
	return s.extract(0xFF10, 0b01110000, 4)
}

func (s *Sound) Square1_Negate() bool {
	return s.memory.ReadBit(0xFF10, 3)
}

func (s *Sound) Square1_Shift() uint8 {
	return s.extract(0xFF10, 0b00000111, 0)
}

func (s *Sound) Square1_Duty() uint8 {
	return s.extract(0xFF11, 0b11000000, 6)
}

func (s *Sound) Square1_LengthLoad() uint8 {
	return 64 - s.extract(0xFF11, 0b00111111, 0)
}

func (s *Sound) Square1_StartingVolumne() uint8 {
	return s.extract(0xFF12, 0b11110000, 4)
}

func (s *Sound) Square1_EnvelopeAddMode() bool {
	return s.memory.ReadBit(0xFF12, 3)
}

func (s *Sound) Square1_Period() uint8 {
	return s.extract(0xFF12, 0b00000111, 0)
}

func (s *Sound) Square1_FrequencyLSB() uint8 {
	return s.memory.ReadByte(0xFF13)
}

func (s *Sound) Square1_Trigger() bool {
	return s.memory.ReadBit(0xFF14, 7)
}

func (s *Sound) Square1_LengthEnable() bool {
	return s.memory.ReadBit(0xFF14, 6)
}

func (s *Sound) Square1_FrequencyMSB() uint8 {
	return s.extract(0xFF14, 0b00000111, 0)
}

func (s *Sound) Square2_Duty() uint8 {
	return s.extract(0xFF16, 0b11000000, 6)
}

func (s *Sound) Square2_LengthLoad() uint8 {
	return 64 - s.extract(0xFF16, 0b00111111, 0)
}

func (s *Sound) Square2_StartingVolumne() uint8 {
	return s.extract(0xFF17, 0b11110000, 4)
}

func (s *Sound) Square2_EnvelopeAddMode() bool {
	return s.memory.ReadBit(0xFF17, 3)
}

func (s *Sound) Square2_Period() uint8 {
	return s.extract(0xFF17, 0b00000111, 0)
}

func (s *Sound) Square2_FrequencyLSB() uint8 {
	return s.memory.ReadByte(0xFF18)
}

func (s *Sound) Square2_Trigger() bool {
	return s.memory.ReadBit(0xFF19, 7)
}

func (s *Sound) Square2_LengthEnable() bool {
	return s.memory.ReadBit(0xFF19, 6)
}

func (s *Sound) Square2_FrequencyMSB() uint8 {
	return s.extract(0xFF19, 0b00000111, 0)
}

func (s *Sound) Wave_DACPower() bool {
	return s.memory.ReadBit(0xFF1A, 7)
}

func (s *Sound) Wave_LengthLoad() uint8 {
	return s.memory.ReadByte(0xFF1B)
}

func (s *Sound) Wave_VolumeCode() uint8 {
	return s.extract(0xFF1C, 0b01100000, 5)
}

func (s *Sound) Wave_FrequencyLSB() uint {
	return uint(s.memory.ReadByte(0xFF1D))
}

func (s *Sound) Wave_Trigger() bool {
	return s.memory.ReadBit(0xFF1E, 7)
}

func (s *Sound) Wave_LengthEnable() bool {
	return s.memory.ReadBit(0xFF1E, 6)
}

func (s *Sound) Wave_FrequencyMSB() uint8 {
	return s.extract(0xFF1E, 0b00000111, 0)
}

func (s *Sound) Noise_LengthLoad() uint8 {
	return 64 - s.extract(0xFF20, 0b00111111, 0)
}

func (s *Sound) Noise_StartingVolume() uint8 {
	return s.extract(0xFF21, 0b11110000, 4)
}

func (s *Sound) Noise_EnvelopeAddMode() bool {
	return s.memory.ReadBit(0xFF21, 3)
}

func (s *Sound) Noise_Period() uint8 {
	return s.extract(0xFF21, 0b00000111, 0)
}

func (s *Sound) Noise_ClockShift() uint8 {
	return s.extract(0xFF22, 0b11110000, 4)
}

func (s *Sound) Noise_WidthModeLFSR() bool {
	return s.memory.ReadBit(0xFF23, 3)
}

func (s *Sound) Noise_DivisorCode() uint8 {
	return s.extract(0xFF23, 0b00000111, 0)
}

func (s *Sound) Noise_Trigger() bool {
	return s.memory.ReadBit(0xFF24, 7)
}

func (s *Sound) Noise_LengthEnable() bool {
	return s.memory.ReadBit(0xFF24, 6)
}

func (s *Sound) ControlStatus_VinLEnable() bool {
	return s.memory.ReadBit(0xFF24, 7)
}

func (s *Sound) ControlStatus_LeftVol() uint8 {
	return s.extract(0xFF24, 0b01110000, 4)
}

func (s *Sound) ControlStatus_VinREnable() bool {
	return s.memory.ReadBit(0xFF24, 3)
}

func (s *Sound) ControlStatus_RightVol() uint8 {
	return s.extract(0xFF24, 0b00000111, 0)
}

// TODO - is left and right corretc way around?
func (s *Sound) ControlStatus_Channel4LeftEnabled() bool {
	return s.memory.ReadBit(0xFF25, 7)
}

func (s *Sound) ControlStatus_Channel3LeftEnabled() bool {
	return s.memory.ReadBit(0xFF25, 6)
}

func (s *Sound) ControlStatus_Channel2LeftEnabled() bool {
	return s.memory.ReadBit(0xFF25, 5)
}

func (s *Sound) ControlStatus_Channel1LeftEnabled() bool {
	return s.memory.ReadBit(0xFF25, 4)
}

func (s *Sound) ControlStatus_Channel4RightEnabled() bool {
	return s.memory.ReadBit(0xFF25, 3)
}

func (s *Sound) ControlStatus_Channel3RightEnabled() bool {
	return s.memory.ReadBit(0xFF25, 2)
}

func (s *Sound) ControlStatus_Channel2RightEnabled() bool {
	return s.memory.ReadBit(0xFF25, 1)
}

func (s *Sound) ControlStatus_Channel1RightEnabled() bool {
	return s.memory.ReadBit(0xFF25, 0)
}

func (s *Sound) ControlStatus_AllOn() bool {
	return s.memory.ReadBit(0xFF26, 7)
}

func (s *Sound) ControlStatus_Channel4Enabled() bool {
	return s.memory.ReadBit(0xFF26, 3)
}

func (s *Sound) ControlStatus_Channel3Enabled() bool {
	return s.memory.ReadBit(0xFF26, 2)
}

func (s *Sound) ControlStatus_Channel2Enabled() bool {
	return s.memory.ReadBit(0xFF26, 1)
}

func (s *Sound) ControlStatus_Channel1Enabled() bool {
	return s.memory.ReadBit(0xFF26, 0)
}
