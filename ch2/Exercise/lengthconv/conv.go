package lengthconv

func FToM(f Feet) Meters {
	return Meters(f / Feet(FtBase))
}

func MToF(m Meters) Feet {
	return Feet(m * Meters(FtBase))
}
