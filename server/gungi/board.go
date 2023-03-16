package gungi

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/TonyLeCode/gungi.go/server/ds"
)

type Piece int

type BoardSquares [BOARD_SQUARE_NUM]*ds.Node

// 0 = white, 1 = black
type StackList [2]ds.LinkedList

type LLStack struct {
	Stack      ds.Stack
	Coordinate int
}

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
	//TODO maybe make easier to read with regex
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
			// TODO make new struct for stacks that include coordinate index
			// TODO modify all *LLStack types to compensate
			// newStack := ds.Stack{}
			newStack := LLStack{
				Stack: ds.Stack{},
			}
			for _, piece := range column {
				// no separation between stacks
				if piece >= '1' && piece <= '9' {
					skipNum, err := strconv.Atoi(string(piece))
					if err != nil {
						return errors.New("invalid fen squares")
					}
					fileIndex += skipNum - 1
				} else {
					newStack.Stack.Push(DecodeSingleChar(string(piece)))
				}
			}
			if newStack.Stack.Length > 0 {
				newStack.Coordinate = CoordsToSquare(fileIndex, i)
				newNode := b.StackList[GetColor(newStack.Stack.Top.Value.(int))].Push(&newStack)
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

func (b *Board) BoardToFen() string {
	var fenString strings.Builder
	skipIndex := 0
	for i := 0; i < PLAYABLE_SQUARE_NUM; i++ {
		square := b.BoardSquares[IndexToSquare(i)]

		if i%9 == 0 && i != 0 {
			fenString.WriteString(strconv.Itoa(skipIndex) + "/")
			skipIndex = 0
		}

		if square == nil {
			skipIndex++
		} else if square.Value == -1 {
			log.Println("out of bounds")
		} else {
			stack := square.Value.(*LLStack).Stack
			var stackString strings.Builder

			pointer := stack.Bottom
			for pointer != nil {
				stackString.WriteString(EncodeSingleChar(pointer.Value.(int)))
				pointer = pointer.Next
			}

			if skipIndex != 0 {
				fenString.WriteString(strconv.Itoa(skipIndex) + ",")
				skipIndex = 0
			}
			fenString.WriteString(stackString.String())

			if i%8 != 0 || i == 0 {
				fenString.WriteString(",")
			}
		}
		if i == 80 && skipIndex != 0 {
			fenString.WriteString(strconv.Itoa(skipIndex))
		}
	}

	fenString.WriteString(" ")
	for i, amount := range b.Hand {
		if i == 13 {
			fenString.WriteString("/")
		}
		fenString.WriteString(strconv.Itoa(amount))
	}

	fenString.WriteString(" ")
	if b.TurnColor == 0 {
		fenString.WriteString("w")
	} else {
		fenString.WriteString("b")
	}

	return fenString.String()
}

// Returns color from piece
func GetColor(piece int) int {
	if piece <= 12 {
		return 0
	} else {
		return 1
	}
}

func GetOppositeColor(color int) int {
	if color == 0 {
		return 1
	}
	return 0
}

// Switches index of a stack if the color of the top piece changes
func (p *StackList) ShiftStack(node *ds.Node, prevPiece int) {
	stack := node.Value.(*LLStack).Stack
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
		newStack := LLStack{
			Stack:      ds.Stack{},
			Coordinate: coordinate,
		}
		newStack.Stack.Push(piece)
		node := b.StackList[GetColor(piece)].Push(&newStack)
		b.BoardSquares[coordinate] = node
	} else if square.Value == -1 {
		return errors.New("invalid square")
	} else {
		stack := square.Value.(*LLStack).Stack
		prevTop := stack.Top.Value.(int)
		stack.Push(piece)
		b.StackList.ShiftStack(square, prevTop)
	}

	return nil
}

// Removes stack from StackList
func (p *StackList) RemoveStack(node *ds.Node) {
	stack := node.Value.(*LLStack).Stack

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

	stack := square.Value.(*LLStack).Stack
	if stack.Length == 1 {
		b.BoardSquares[coordinate] = nil
	}
	b.StackList.RemoveStack(square)
	return nil
}

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
					stack := square.Value.(*LLStack).Stack
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

//9   37,  38,  39,  40,  41,  42,  43,  44,  45,

//8   49,  50,  51,  52,  53,  54,  55,  56,  57,

//7   61,  62,  63,  64,  65,  66,  67,  68,  69,

//6   73,  74,  75,  76,  77,  78,  79,  80,  81,

//5   85,  86,  87,  88,  89,  90,  91,  92,  93,

//4   97,  98,  99, 100, 101, 102, 103, 104, 105,

//3  109, 110, 111, 112, 113, 114, 115, 116, 117,

//2  121, 122, 123, 124, 125, 126, 127, 128, 129,

//1  133, 134, 135, 136, 137, 138, 139, 140, 141,

//    A    B    C    D    E    F    G    H    I
