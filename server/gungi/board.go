package gungi

import (
	"fmt"

	"github.com/TonyLeCode/gungi.go/ds"
)

type Piece int

type BoardSquares = [BOARD_SQUARE_NUM]ds.Stack

type PieceList = [2]*ds.LinkedList

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

// Takes file and rank and converts to one dimensional index on playable area of 12x15
func CoordsToSquare(file int, rank int) int {
	return (file + 37) + (rank * 12)
}

//	var outsideVal = ds.StackNode{
//		Value: nil,
//	}
var outsideSquare = ds.Stack{
	Length: -1,
}

func (board *Board) InitializeBoard() {

	// Board Squares
	for i := 0; i < len(board.BoardSquares); i++ {
		board.BoardSquares[i] = outsideSquare
	}
	for i := 0; i < 81; i++ {
		newSquare := ds.Stack{}
		board.BoardSquares[IndexToSquare(i)] = newSquare
	}
	board.printBoard()

	// Hand
	board.Hand = [26]int{9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1, 9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1}
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

func (board *Board) printBoard() {
	for i := 0; i < 15; i++ {
		for j := 0; j < 12; j++ {
			index := i*12 + j
			if index < len(board.BoardSquares) {
				if board.BoardSquares[index].Length == -1 {
					fmt.Print(" -1")
				} else {
					// fmt.Print(arr[index].Stack.Top)
					fmt.Print("  0")
				}
			} else {
				// double checks print
				fmt.Print(" z ")
			}
		}
		fmt.Println()
	}
}
