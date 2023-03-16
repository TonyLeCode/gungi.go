package gungi

import (
	"fmt"
	"log"
)

func IsSlidingPiece(piece int, tier int) {

}

func (b *Board) IsInCheck() {

}

type PossibleMove struct {
	Stack   *LLStack
	ToCoord int
}

// Generates legal moves
func (b *Board) GetLegalMoves() {

	// See is Marshal is in check
	// See if piece is pinned
	// Restrain hand placements
	// See if moving out of stack puts Marshal in check
	// Check if pawn is already in same file

	moveList := []PossibleMove{}

	// Loop through enemy sliding pieces to check possibility of pins
	currentStackNode := b.StackList[GetOppositeColor(b.TurnColor)].Head
	var slidingPieceList []LLStack
	for currentStackNode != nil {
		stack := currentStackNode.Value.(LLStack)

		// add sliding pieces for pin check
		switch stack.Stack.Top.Value.(int) % 13 {
		case MUSKETEER, SAMURAI, CANNON, SPY, TACTICIAN:
			slidingPieceList = append(slidingPieceList, stack)
		}

		pseudoMoves := b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length)
		log.Println(pseudoMoves)
		log.Println(moveList)

		currentStackNode = currentStackNode.Next
	}

	// Loop through pieces of current player
	currentStackNode = b.StackList[b.TurnColor].Head
	for currentStackNode != nil {
		stack := currentStackNode.Value.(LLStack)

		pseudoMoves := b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length)
		log.Println(pseudoMoves)
		log.Println(moveList)

		currentStackNode = currentStackNode.Next
	}
	fmt.Println(currentStackNode)
	fmt.Println(slidingPieceList)

}

func (b *Board) GetPseudoRangingPiece(coordinate int, offset int) []int {
	squares := []int{}

	i := coordinate + offset
	currSquare := b.BoardSquares[i]
	for currSquare == nil {
		squares = append(squares, i)
		i += offset
		currSquare = b.BoardSquares[i]
	}

	if currSquare.Value != -1 {
		squares = append(squares, i)
	}

	return squares
}

// Generates psuedo-legal moves for a piece at a coordinate
func (b *Board) GetPseudoLegalMoves(piece int, coordinate int, tier int) []int {
	offsets := []int{}
	squares := []int{}
	switch piece % 13 {
	case PAWN:
		switch tier {
		case 1:
			offsets = append(offsets, -12)
		case 2, 3:
			offsets = append(offsets, -11, -12, -13)
		}
	case LIEUTENANT_GENERAL:
		switch tier {
		case 1:
			offsets = append(offsets, -11, -12, -13, 11, 13)
		case 2:
			offsets = append(offsets, -11, -12, -13, 11, 12, 13)
		case 3:
			offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
		}
	case MAJOR_GENERAL:
		switch tier {
		case 1:
			offsets = append(offsets, -13, -11)
		case 2:
			offsets = append(offsets, -11, -12, -13, 11, 13)
		case 3:
			offsets = append(offsets, -11, -12, -13, -1, 1, 12)
		}
	case GENERAL:
		switch tier {
		case 1:
			offsets = append(offsets, -11, -12, -13, -1, 1, 12)
		case 2:
			offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
		case 3:
			offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13, -25, -24, -23)
		}
	case FORTRESS:
		// Cannot stack
		offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
	case KNIGHT:
		switch tier {
		case 1:
			offsets = append(offsets, -1, 1, -25, -23)
		case 2:
			offsets = append(offsets, -25, -23, -14, -10)
		case 3:
			offsets = append(offsets, -25, -23, -14, -10, 25, 23, 14, 10)
		}
	case ARCHER:
		switch tier {
		case 1:
			offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
		case 2:
			offsets = append(offsets, -26, -25, -24, -22, -10, -14, -2, 26, 25, 24, 22, 10, 14, 2)
		case 3:
			offsets = append(offsets, -39, -38, -37, -36, -35, -34, -33, -27, -15, -3, -21, -9, 39, 38, 37, 36, 35, 34, 33, 27, 15, 3, 21, 9)
		}
	case MUSKETEER:
		switch tier {
		case 1:
			offsets = append(offsets, -12)
		case 2:
			offsets = append(offsets, -12, -24)
		case 3:
			if GetColor(piece) == 0 {
				squares = append(squares, b.GetPseudoRangingPiece(coordinate, -12)...)
			} else {
				squares = append(squares, b.GetPseudoRangingPiece(coordinate, 12)...)
			}
		}
	case SAMURAI:
		switch tier {
		case 1:
			offsets = append(offsets, -11, -13, 11, 13)
		case 2:
			offsets = append(offsets, -26, -22, 26, 22)
		case 3:
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 13)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -13)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -11)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 11)...)
		}
	case CANNON:
		switch tier {
		case 1:
			offsets = append(offsets, -12, -1, 1, 12)
		case 2:
			offsets = append(offsets, -24, -12, -2, -1, 1, 2, 12, 24)
		case 3:
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 1)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -1)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -12)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 12)...)
		}
	case SPY:
		switch tier {
		case 1:
			offsets = append(offsets, -12)
		case 2:
			offsets = append(offsets, -11, -13, 11, 13)
		case 3:
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 13)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -13)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -11)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 11)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 1)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -1)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, -12)...)
			squares = append(squares, b.GetPseudoRangingPiece(coordinate, 12)...)
		}
	case TACTICIAN:
		switch tier {
		case 1:
			offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
		case 2, 3:
			piece := b.BoardSquares[coordinate].Value.(*LLStack).Stack.Top.Prev.Value.(int)
			sameColorPiece := piece % 13
			if GetColor(piece) == 1 {
				sameColorPiece += 13
			}
			squares = append(squares, b.GetPseudoLegalMoves(sameColorPiece, coordinate, tier)...)
		}
	case MARSHAL:
		offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
	}

	// Reverse direction for black
	for _, offset := range offsets {
		if GetColor(piece) == 1 {
			squares = append(squares, coordinate-offset)
		} else {
			squares = append(squares, coordinate+offset)
		}
	}

	return squares
}
