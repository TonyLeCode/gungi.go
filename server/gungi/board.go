package gungi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/TonyLeCode/gungi.go/ds"
)

type Piece int

type BoardSquares [BOARD_SQUARE_NUM]*ds.Node

// 0 = white, 1 = black
type StackList [2]ds.LinkedList

type Board struct {
	BoardSquares BoardSquares
	StackList    StackList
	Hand         [26]int

	TurnColor  int
	TurnNumber int

	MoveList []int
	History  []History
}

type History struct {
	Move     int
	MoveType int
}

var outsideSquare = ds.Node{
	Value: -1,
}

// Also resets board
func (b *Board) InitializeBoard() {
	b.StackList[0] = ds.LinkedList{}
	b.StackList[1] = ds.LinkedList{}

	// Board Squares
	for i := 0; i < len(b.BoardSquares); i++ {
		b.BoardSquares[i] = &outsideSquare
	}
	for i := 0; i < 81; i++ {
		b.BoardSquares[IndexToSquare(i)] = nil
	}

	b.Hand = [26]int{
		9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1,
		9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1,
	}

	b.TurnColor = 0
	b.TurnNumber = 0

	b.History = []History{}
}

func (b *Board) SetBoardFromFen(fen string) error {
	b.InitializeBoard()
	// index 0 holds piece position, index 1 holds hand piece count, index 2 is color turn
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
			newStack := ds.Stack{}
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
					newStack.Push(DecodeSingleChar(string(piece)))
				}
			}
			if newStack.Length > 0 {
				newNode := b.StackList[GetColor(newStack.Top.Value.(int))].Push(&newStack)
				b.BoardSquares[CoordsToSquare(fileIndex, i)] = newNode
			}
			fileIndex += 1
		}
	}
	// assign pieces to hand
	handSplit := strings.Split(fields[1], "/")
	piecesSplitW := strings.Split(handSplit[0], "")
	piecesSplitB := strings.Split(handSplit[1], "")
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

// Returns color from piece
func GetColor(piece int) int {
	if piece <= 12 {
		return 0
	} else {
		return 1
	}
}

// Switches index of a stack if the color of the top piece changes
func (p *StackList) ShiftStack(node *ds.Node, prevPiece int) {
	stack := node.Value.(*ds.Stack)
	piece := stack.Top.Value.(int)

	pieceColor := GetColor(piece)
	prevPieceColor := GetColor(prevPiece)

	if pieceColor != prevPieceColor {
		p[prevPieceColor].Remove(node)
		p[pieceColor].Push(node)
	}
}

// Places piece on top of a stack, or creates one if empty
func (b *Board) PlacePiece(piece int, coordinate int) error {
	square := b.BoardSquares[coordinate]

	if square == nil {
		newStack := ds.Stack{}
		newStack.Push(piece)
		node := b.StackList[GetColor(piece)].Push(&newStack)
		b.BoardSquares[coordinate] = node
	} else if square.Value == -1 {
		return errors.New("invalid square")
	} else {
		stack := square.Value.(*ds.Stack)
		prevTop := stack.Top.Value.(int)
		stack.Push(piece)
		b.StackList.ShiftStack(square, prevTop)
	}

	return nil
}

// Removes stack from StackList
func (p *StackList) RemoveStack(node *ds.Node) {
	stack := node.Value.(*ds.Stack)

	if stack.Length > 1 {
		stack.Pop()
	} else {
		piece := stack.Top.Value.(int)
		pieceColor := GetColor(piece)
		p[pieceColor].Remove(node)
	}
}

// Removes top piece from a stack, or removes it if it's the only piece
func (b *Board) RemovePiece(coordinate int) error {
	square := b.BoardSquares[coordinate]

	if square == nil {
		return errors.New("square already empty")
	} else if square.Value == -1 {
		return errors.New("square out of bounds")
	}

	stack := square.Value.(*ds.Stack)
	if stack.Length == 1 {
		b.BoardSquares[coordinate] = nil
	}
	b.StackList.RemoveStack(square)
	return nil
}

// TODO
func (b *Board) BoardToFen() {

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
				square := b.BoardSquares[index]
				if square == nil {
					fmt.Print("   -")
				} else if square.Value == -1 {
					fmt.Print("  -1")
				} else {
					stack := square.Value.(*ds.Stack)
					if stack.Length != 0 {
						val := EncodeSingleChar(stack.Top.Value.(int))
						fmt.Print("   ", val)
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
	fmt.Println()
	fmt.Print("Hand: ")
	for i, piece := range b.Hand {
		fmt.Print(EncodeSingleChar(i), ": ", piece, ", ")
	}
	fmt.Println()
}
