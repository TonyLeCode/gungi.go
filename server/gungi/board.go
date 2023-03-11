package gungi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/TonyLeCode/gungi.go/ds"
)

type Piece int

type BoardSquares = [BOARD_SQUARE_NUM]ds.Stack

// 0 = white, 1 = black
type PieceList = [2]*ds.LinkedList

type Fen string

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

type Board struct {
	BoardSquares BoardSquares
	PieceList    PieceList
	Hand         [26]int

	MarshalSquare [2]int
	TurnColor     int
	TurnNumber    int

	History []History
}

type History struct {
	Move     int
	MoveType int
}

// Takes index that would normally be on 9x9 and transposes it to playable area of 12x15
func IndexToSquare(index int) int {
	return CoordsToSquare(index%9, index/9)
}

// Takes file and rank that would normally be on 9x9 and converts to one dimensional index on playable area of 12x15
func CoordsToSquare(file int, rank int) int {
	return (file + 37) + (rank * 12)
}

var outsideSquare = ds.Stack{
	Length: -1,
}

// Also resets board
func (b *Board) InitializeBoard() {

	// Board Squares
	for i := 0; i < len(b.BoardSquares); i++ {
		b.BoardSquares[i] = outsideSquare
	}
	for i := 0; i < 81; i++ {
		b.BoardSquares[IndexToSquare(i)] = ds.Stack{}
	}

	b.PieceList = PieceList{}
	b.Hand = [26]int{
		9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1,
		9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1,
	}

	b.MarshalSquare = [2]int{}
	b.TurnColor = 0
	b.TurnNumber = 0

	b.History = []History{}
}

func (b *Board) SetBoardFromFen(fen string) error {
	b.InitializeBoard()
	// index 0 holds piece position, index 1 holds hand piece count, index 3 is color turn
	fields := strings.Split(fen, " ")

	// forward slash splits into rows
	split := strings.Split(fields[0], "/")
	if len(split) != 9 {
		return errors.New("incorrect fen length")
	}

	// assign pieces to squares
	for i, row := range split {
		// comma splits into columns
		split2 := strings.Split(row, ",")
		// fileIndex must be separate because we will skip indices
		fileIndex := 0
		for _, column := range split2 {
			for _, piece := range column {
				// no separation between stacks
				if piece >= '1' && piece <= '9' {
					skipNum, err := strconv.Atoi(string(piece))
					if err != nil {
						return errors.New("invalid fen squares")
					}
					fileIndex += skipNum - 1
				} else {
					// TODO assign pieces to list
					b.BoardSquares[CoordsToSquare(fileIndex, i)].Push(DecodeSingleChar(string(piece)))
				}
			}
			fileIndex += 1
		}
	}

	// assign pieces to hand
	handSplit := strings.Split(fields[1], "/")
	piecesSplitW := strings.Split(handSplit[0], ",")
	piecesSplitB := strings.Split(handSplit[1], ",")
	piecesConcat := append(piecesSplitW, piecesSplitB...)
	for i, pieceAmount := range piecesConcat {
		pieceNum, err := strconv.Atoi(pieceAmount)
		if err != nil {
			return errors.New("invalid fen stockpile")
		}
		b.Hand[i] = pieceNum
	}

	b.TurnColor = LetterToTurn(fields[2])

	return nil
}

// func Mailbox() [81]int {
// 	// return [81]int{
// 	// 	37, 38, 39, 40, 41, 42, 43, 44, 45,
// 	// 	49, 50, 51, 52, 53, 54, 55, 56, 57,
// 	// 	61, 62, 63, 64, 65, 66, 67, 68, 69,
// 	// 	73, 74, 75, 76, 77, 78, 79, 80, 81,
// 	// 	85, 86, 87, 88, 89, 90, 91, 92, 93,
// 	// 	97, 98, 99, 100, 101, 102, 103, 104, 105,
// 	// 	109, 110, 111, 112, 113, 114, 115, 116, 117,
// 	// 	121, 122, 123, 124, 125, 126, 127, 128, 129,
// 	// 	133, 134, 135, 136, 137, 138, 139, 140, 141,
// 	// }
// 	return [81]int{}
// }

func (b *Board) PrintBoard() {
	for i := 0; i < 15; i++ {
		for j := 0; j < 12; j++ {
			index := i*12 + j
			if index < len(b.BoardSquares) {
				if b.BoardSquares[index].Length == -1 {
					fmt.Print("  -1")
				} else {
					square := b.BoardSquares[index]
					if square.Length != 0 {
						// val := square.Top.Value.(int)
						// padding := strings.Repeat(" ", 4-len(strconv.Itoa(val)))
						// fmt.Print(padding, val)
						val := EncodeSingleChar(square.Top.Value.(int))
						fmt.Print("   ", val)
					} else {
						fmt.Print("   -")
					}
				}
			} else {
				// double checks print
				fmt.Print(" z ")
			}
		}
		fmt.Println()
	}
}

func (b *Board) PrintHand() {

}
