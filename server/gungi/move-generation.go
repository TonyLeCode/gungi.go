package gungi

import "fmt"

func IsSlidingPiece(piece int, tier int) {

}

func (b *Board) IsInCheck() {
	currentStackNode := b.StackList[GetOppositeColor(b.TurnColor)].Head
	var slidingPieceList []int
	for currentStackNode != nil {

	}
	fmt.Println(currentStackNode)

}

func GenerateAllMoves() {
	// See is Marshal is in check
	// See if piece is pinned
	// Restrain hand placements
	// See if moving out of stack puts Marshal in check
	// Check if pawn is already in same file

}

func (piece *Piece) GetMove() {
	switch int(*piece) {
	case WHITE_PAWN, BLACK_PAWN:
		fmt.Println("hello")
	case WHITE_LIEUTENANT_GENERAL, BLACK_LIEUTENANT_GENERAL:
		fmt.Println("hello")
	case WHITE_MAJOR_GENERAL, BLACK_MAJOR_GENERAL:
		fmt.Println("hello")
	case WHITE_GENERAL, BLACK_GENERAL:
		fmt.Println("hello")
	case WHITE_FORTRESS, BLACK_FORTRESS:
		fmt.Println("hello")
	case WHITE_KNIGHT, BLACK_KNIGHT:
		fmt.Println("hello")
	case WHITE_ARCHER, BLACK_ARCHER:
		fmt.Println("hello")
	case WHITE_MUSKETEER, BLACK_MUSKETEER:
		fmt.Println("hello")
	case WHITE_SAMURAI, BLACK_SAMURAI:
		fmt.Println("hello")
	case WHITE_CANNON, BLACK_CANNON:
		fmt.Println("hello")
	case WHITE_SPY, BLACK_SPY:
		fmt.Println("hello")
	case WHITE_TACTICIAN, BLACK_TACTICIAN:
		fmt.Println("hello")
	case WHITE_MARSHAL, BLACK_MARSHAL:
		fmt.Println("hello")
	default:
		fmt.Println("false")
	}
}
