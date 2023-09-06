package utils

func OppositeColor(color int) int {
	if color == 0 {
		return 1
	}
	return 0
}

// Takes index that would normally be on 9x9 and transposes it to playable area of 12x15
func IndexToSquare(index int) int {
	return (index%9 + 37) + (index / 9 * 12)
}

// Takes file and rank (index) that would normally be on 9x9 and converts to one dimensional index on playable area of 12x15
func CoordsToSquare(file int, rank int) int {
	return (file + 37) + (rank * 12)
}

// Takes index of 12x15 and turns it into 9x9 file and rank
func SquareToCoords(square int) (int, int) {
	file := (square / 12) - 2
	rank := square % 12
	return file, rank
}

// Gets index of 12x15 and turns it into board notation (Rank, File)
func CoordsToNotation(square int) (string, string) {
	file, rank := SquareToCoords(square)
	return InvertRank(file), FileToLetter(rank)
}

// Index of 12x15 into 9x9 index
func SquareToIndex(index int) int {
	file, rank := SquareToCoords(index)
	return (file-1)*9 + (rank - 1)
}

func FileToLetter(rank int) string {
	switch rank {
	case 1:
		return "A"
	case 2:
		return "B"
	case 3:
		return "C"
	case 4:
		return "D"
	case 5:
		return "E"
	case 6:
		return "F"
	case 7:
		return "G"
	case 8:
		return "H"
	case 9:
		return "I"
	}
	return ""
}

func LetterToFile(rank string) int {
	switch rank {
	case "A":
		return 1
	case "B":
		return 2
	case "C":
		return 3
	case "D":
		return 4
	case "E":
		return 5
	case "F":
		return 6
	case "G":
		return 7
	case "H":
		return 8
	case "I":
		return 9
	}
	return 0
}

func RevertRank(file string) int {
	switch file {
	case "9":
		return 1
	case "8":
		return 2
	case "7":
		return 3
	case "6":
		return 4
	case "5":
		return 5
	case "4":
		return 6
	case "3":
		return 7
	case "2":
		return 8
	case "1":
		return 9
	}
	return 0
}

func InvertRank(file int) string {
	switch file {
	case 1:
		return "9"
	case 2:
		return "8"
	case 3:
		return "7"
	case 4:
		return "6"
	case 5:
		return "5"
	case 6:
		return "4"
	case 7:
		return "3"
	case 8:
		return "2"
	case 9:
		return "1"
	}
	return ""
}

func DecodeSingleChar(f string) int {
	// white is lowercase, black is uppercase
	switch f {
	case "P":
		return 0
	case "L":
		return 1
	case "S":
		return 2
	case "G":
		return 3
	case "F":
		return 4
	case "K":
		return 5
	case "Y":
		return 6
	case "B":
		return 7
	case "W":
		return 8
	case "C":
		return 9
	case "N":
		return 10
	case "T":
		return 11
	case "M":
		return 12
	case "p":
		return 13
	case "l":
		return 14
	case "s":
		return 15
	case "g":
		return 16
	case "f":
		return 17
	case "k":
		return 18
	case "y":
		return 19
	case "b":
		return 20
	case "w":
		return 21
	case "c":
		return 22
	case "n":
		return 23
	case "t":
		return 24
	case "m":
		return 25
	default:
		return -1
	}
}
func EncodeSingleChar(i int) string {
	// white is lowercase, black is uppercase
	switch i {
	case 0:
		return "P"
	case 1:
		return "L"
	case 2:
		return "S"
	case 3:
		return "G"
	case 4:
		return "F"
	case 5:
		return "K"
	case 6:
		return "Y"
	case 7:
		return "B"
	case 8:
		return "W"
	case 9:
		return "C"
	case 10:
		return "N"
	case 11:
		return "T"
	case 12:
		return "M"
	case 13:
		return "p"
	case 14:
		return "l"
	case 15:
		return "s"
	case 16:
		return "g"
	case 17:
		return "f"
	case 18:
		return "k"
	case 19:
		return "y"
	case 20:
		return "b"
	case 21:
		return "w"
	case 22:
		return "c"
	case 23:
		return "n"
	case 24:
		return "t"
	case 25:
		return "m"
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
