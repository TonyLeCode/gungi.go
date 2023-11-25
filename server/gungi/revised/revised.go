package revised

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/TonyLeCode/gungi.go/server/gungi/utils"
)

func getColor(piece int) int {
	if piece <= 12 {
		return 0
	} else {
		return 1
	}
}

type Square struct {
	coord int
	stack []int
	prev  *Square
	next  *Square
}

func (s Square) IsOutOfBounds() bool {
	return s.stack[0] == -1
}

func (s Square) IsEmpty() bool {
	return len(s.stack) == 0
}

func (s Square) IsStackable() bool {
	return len(s.stack) > 0 && len(s.stack) < 3
}

func (s Square) GetTop() int {
	return s.stack[len(s.stack)-1]
}

func (s *Square) Push(piece int) {
	s.stack = append(s.stack, piece)
}

func (s *Square) Pop() int {
	popped := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return popped
}

type ListRef struct {
	Head [2]*Square
	Tail [2]*Square
}

func (s *ListRef) Remove(square *Square, color int) {
	if s.Head[color] == s.Tail[color] && s.Head[color] == square {
		s.Head[color] = nil
		s.Tail[color] = nil
		return
	}
	if s.Head[color] == square {
		s.Head[color] = square.next
		square.next.prev = nil
		return
	}
	if s.Tail[color] == square {
		s.Tail[color] = square.prev
		square.prev.next = nil
		return
	}
	square.prev.next = square.next
	square.next.prev = square.prev
}

func (s *ListRef) Push(square *Square, color int) {
	if s.Head[color] == nil && s.Tail[color] == nil {
		s.Head[color] = square
		s.Tail[color] = square
		return
	}
	square.prev = s.Tail[color]
	square.prev.next = square
	s.Tail[color] = square

}

// Compares altered square previous piece then changes ref if necessary
func (s *ListRef) AlteredSquare(square *Square, prevPiece int) {
	prevColor := getColor(prevPiece)
	if square.IsEmpty() {
		s.Remove(square, prevColor)
		return
	}

	currentColor := getColor(square.GetTop())

	if currentColor != prevColor {
		s.Remove(square, prevColor)
		s.Push(square, currentColor)
	}
}

type Hand [26]int

func (h *Hand) print() {
	toPrintW := strings.Builder{}
	toPrintB := strings.Builder{}
	for i := 0; i < 13; i++ {
		p := h[i]
		toPrintW.WriteString(utils.EncodeSingleChar(i))
		toPrintW.WriteString(":")
		toPrintW.WriteString(strconv.Itoa(p))
		toPrintW.WriteString("  ")
	}
	for i := 13; i < 26; i++ {
		p := h[i]
		toPrintB.WriteString(utils.EncodeSingleChar(i))
		toPrintB.WriteString(":")
		toPrintB.WriteString(strconv.Itoa(p))
		toPrintB.WriteString("  ")
	}

	fmt.Println("W hand: ")
	fmt.Println(toPrintW.String())
	fmt.Println("B hand: ")
	fmt.Println(toPrintB.String())
}

type Revised struct {
	BoardSquares  [utils.BOARD_SQUARE_NUM]Square
	ListRef       ListRef
	Hand          Hand
	History       []string
	MarshalCoords [2]int
	Ready         [2]bool
	TurnColor     int
	TurnNumber    int
}

func (r Revised) GetTurnColor() int {
	return r.TurnColor
}

func (r Revised) PieceCount() (int, int) {
	wCount := 0
	bCount := 0

	for listIndex := 0; listIndex < len(r.ListRef.Head); listIndex++ {
		currStack := r.ListRef.Head[listIndex]
		for currStack != nil {
			for i := 0; i < len(currStack.stack); i++ {
				if getColor(currStack.stack[i]) == 0 {
					wCount++
				} else {
					bCount++
				}
			}
			currStack = currStack.next
		}
	}
	return wCount, bCount
}

func (r Revised) ConvertInputCoord(coord int) int {
	return utils.IndexToSquare(coord)
}
func (r Revised) ConvertOutputCoord(coord int) int {
	return utils.SquareToIndex(coord)
}

// Also resets board
func (r *Revised) InitializeBoard() {
	outsideSquare := Square{stack: []int{-1}}

	// Board Squares
	for i := 0; i < len(r.BoardSquares); i++ {
		r.BoardSquares[i] = outsideSquare
	}
	for i := 0; i < 81; i++ {
		r.BoardSquares[utils.IndexToSquare(i)] = Square{coord: utils.IndexToSquare(i)}
	}

	r.ListRef = ListRef{}
	r.Hand = [26]int{
		9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1,
		9, 4, 4, 6, 2, 2, 2, 1, 2, 2, 2, 1, 1,
	}
	r.History = []string{}
	r.MarshalCoords = [2]int{}
	r.Ready = [2]bool{}
	r.TurnColor = 0
	r.TurnNumber = 0
}

func (r *Revised) SetBoardFromFen(fen string) error {
	r.InitializeBoard()

	fields := strings.Split(fen, " ")
	if len(fields) != 4 {
		return errors.New("incorrect fen length")
	}

	rowSplit := strings.Split(fields[0], "/")
	if len(rowSplit) != 9 {
		return errors.New("incorrect fen length")
	}

	for i, row := range rowSplit {
		columnSplit := strings.Split(row, ",")

		fileIndex := 0
		for _, column := range columnSplit {
			newSquare := Square{}
			for _, piece := range column {
				if piece >= '1' && piece <= '9' {
					skipNum, err := strconv.Atoi(string(piece))
					if err != nil {
						return errors.New("invalid fen squares")
					}
					fileIndex += skipNum - 1
				} else {
					switch utils.DecodeSingleChar(string(piece)) {
					case BLACK_MARSHAL:
						r.MarshalCoords[1] = utils.CoordsToSquare(fileIndex, i)
					case WHITE_MARSHAL:
						r.MarshalCoords[0] = utils.CoordsToSquare(fileIndex, i)
					}
					newSquare.Push(utils.DecodeSingleChar(string(piece)))
				}
			}

			if !newSquare.IsEmpty() {
				coord := utils.CoordsToSquare(fileIndex, i)
				newSquare.coord = coord
				r.BoardSquares[coord] = newSquare
				color := getColor(newSquare.GetTop())
				r.ListRef.Push(&r.BoardSquares[coord], color)
			}
			fileIndex += 1
		}
	}

	//assign pieces to hand
	handSplit := strings.Split(fields[1], "/")
	piecesSplitW := strings.Split(handSplit[0], "")
	piecesSplitB := strings.Split(handSplit[1], "")
	piecesConcat := append(piecesSplitW, piecesSplitB...)
	for i, pieceAmount := range piecesConcat {
		pieceNum, err := strconv.Atoi(pieceAmount)
		if err != nil {
			return errors.New("invalid fen stockpile")
		}
		r.Hand[i] = pieceNum
	}

	r.TurnColor = utils.LetterToTurn(fields[2])

	ready := strings.Split(fields[3], "")
	if ready[0] == "0" {
		r.Ready[0] = false
	} else {
		r.Ready[0] = true
	}
	if ready[1] == "0" {
		r.Ready[1] = false
	} else {
		r.Ready[1] = true
	}

	return nil
}

// TODO place and remove marshal should change marshal coords
func (r *Revised) PlacePiece(piece int, coord int) {
	square := &r.BoardSquares[coord]

	if square.IsEmpty() {
		square.Push(piece)
		r.ListRef.Push(square, getColor(piece))
	} else if square.IsStackable() {
		prevPiece := square.GetTop()
		square.Push(piece)
		r.ListRef.AlteredSquare(square, prevPiece)
	}
}

func (r *Revised) RemovePiece(coord int) {
	square := &r.BoardSquares[coord]

	if square.IsEmpty() {
		return
	}
	if square.IsOutOfBounds() {
		return
	}

	p := square.Pop()

	r.ListRef.AlteredSquare(square, p)
}

func (l *ListRef) Print(color int) {
	toPrint := strings.Builder{}
	current := l.Head[color]
	for current != nil {
		piece := utils.EncodeSingleChar(current.GetTop())
		toPrint.WriteString(piece + " ")
		current = current.next
	}
	fmt.Println(toPrint.String())
}

func (r *Revised) BoardToFen() string {
	var fenString strings.Builder
	skipIndex := 0

	for i := 0; i < utils.PLAYABLE_SQUARE_NUM; i++ {
		square := r.BoardSquares[utils.IndexToSquare(i)]

		if square.IsEmpty() {
			skipIndex++
		} else if square.IsOutOfBounds() {
			log.Println("out of bounds")
		} else {
			stackStr := strings.Builder{}

			for i := 0; i < len(square.stack); i++ {
				stackStr.WriteString(utils.EncodeSingleChar(square.stack[i]))
			}

			if skipIndex != i%9 && skipIndex != 0 {
				fenString.WriteString(",")
			}

			if skipIndex != 0 {
				fenString.WriteString(strconv.Itoa(skipIndex))
				skipIndex = 0
			}

			if i%9 != 0 || skipIndex != 0 {
				fenString.WriteString(",")
			}

			fenString.WriteString(stackStr.String())

		}
		if i%9 == 8 {
			if skipIndex != 0 {
				if skipIndex != 9 {
					fenString.WriteString(",")
				}
				fenString.WriteString(strconv.Itoa(skipIndex))
				skipIndex = 0
			}
			if i != 80 {
				fenString.WriteString("/")
			}
		}
	}

	fenString.WriteString(" ")
	for i, amount := range r.Hand {
		if i == 13 {
			fenString.WriteString("/")
		}
		fenString.WriteString(strconv.Itoa(amount))
	}

	fenString.WriteString(" ")
	if r.TurnColor == 0 {
		fenString.WriteString("w")
	} else {
		fenString.WriteString("b")
	}

	fenString.WriteString(" ")
	if r.Ready[0] {
		fenString.WriteString("1")
	} else {
		fenString.WriteString("0")
	}
	if r.Ready[1] {
		fenString.WriteString("1")
	} else {
		fenString.WriteString("0")
	}
	return fenString.String()
}

func (r *Revised) SetHistory(history []string) {
	r.History = history
}

func (r *Revised) SerializeHistory() string {
	return strings.Join(r.History, " ")
}

func (r *Revised) PrintBoard() {
	for i := 0; i < 15; i++ {
		for j := 0; j < 12; j++ {
			index := i*12 + j
			if index < len(r.BoardSquares) {
				square := r.BoardSquares[index]
				if square.IsEmpty() {
					fmt.Print("   -")
				} else if square.IsOutOfBounds() {
					fmt.Print("  -1")
				} else {
					piece := utils.EncodeSingleChar(square.GetTop())
					fmt.Print("   ", piece)
				}
			} else {
				// double checks print
				fmt.Print(" z ")
			}
		}
		fmt.Println()

	}
	// TODO Temporary for debugging
	// r.Hand.print()

	// fmt.Println("W list: ")
	// r.ListRef.print(0)
	// fmt.Println("B list: ")
	// r.ListRef.print(1)
	// fmt.Println("ready: ")
	// fmt.Println(r.Ready)
	// fmt.Println("marshal: ")
	// fmt.Println(r.MarshalCoords)
	// fmt.Println("turn: ")
	// fmt.Println(r.TurnColor)
	// fmt.Println("move history: ")
	// fmt.Println(r.History)
	// fmt.Println()
}
