package gungi

// Takes index that would normally be on 9x9 and transposes it to playable area of 12x15
func IndexToSquare(index int) int {
	return CoordsToSquare(index%9, index/9)
}

// Takes file and rank (index) that would normally be on 9x9 and converts to one dimensional index on playable area of 12x15
func CoordsToSquare(file int, rank int) int {
	return (file + 37) + (rank * 12)
}

func DecodeSingleChar(f string) int {
	// white is lowercase, black is uppercase
	switch f {
	case "p":
		return 0
	case "l":
		return 1
	case "s":
		return 2
	case "g":
		return 3
	case "f":
		return 4
	case "k":
		return 5
	case "y":
		return 6
	case "b":
		return 7
	case "w":
		return 8
	case "c":
		return 9
	case "n":
		return 10
	case "t":
		return 11
	case "m":
		return 12
	case "P":
		return 13
	case "L":
		return 14
	case "S":
		return 15
	case "G":
		return 16
	case "F":
		return 17
	case "K":
		return 18
	case "Y":
		return 19
	case "B":
		return 20
	case "W":
		return 21
	case "C":
		return 22
	case "N":
		return 23
	case "T":
		return 24
	case "M":
		return 25
	default:
		return -1
	}
}
func EncodeSingleChar(i int) string {
	// white is lowercase, black is uppercase
	switch i {
	case 0:
		return "p"
	case 1:
		return "l"
	case 2:
		return "s"
	case 3:
		return "g"
	case 4:
		return "f"
	case 5:
		return "k"
	case 6:
		return "y"
	case 7:
		return "b"
	case 8:
		return "w"
	case 9:
		return "c"
	case 10:
		return "n"
	case 11:
		return "t"
	case 12:
		return "m"
	case 13:
		return "P"
	case 14:
		return "L"
	case 15:
		return "S"
	case 16:
		return "G"
	case 17:
		return "F"
	case 18:
		return "K"
	case 19:
		return "Y"
	case 20:
		return "B"
	case 21:
		return "W"
	case 22:
		return "C"
	case 23:
		return "N"
	case 24:
		return "T"
	case 25:
		return "M"
	default:
		return ""
	}
}

// converts "b" or "w" to turn color
func LetterToTurn(color string) int {
	switch color {
	case "w":
		return WHITE
	case "b":
		return BLACK
	default:
		return WHITE
	}

}
