package gungi

import (
	"log"

	"github.com/TonyLeCode/gungi.go/server/ds"
)

func IsSlidingPiece(piece int, tier int) {

}

func (b *Board) IsInCheck() {

}

type PossibleMove struct {
	Stack   *LLStack
	ToCoord int
}

func ColorOffset(turnColor int, offset int) int {
	if turnColor == 0 {
		return offset
	} else {
		return offset * -1
	}
}

type PseudoMove struct {
	coordinate int
	moveList   []int
}

type XRay struct {
	coordinate int
	path       []XRaySquares
}

// Generates legal moves
func (b *Board) GenerateLegalMoves() {

	// See is Marshal is in check
	// See if piece is pinned
	// Restrain hand placements
	// See if moving out of stack puts Marshal in check
	// Check if pawn is already in same file

	// moveList := []PossibleMove{}
	// enemyMoveList := []PseudoMove{}
	var enemyXRaySquares XRay
	inCheck := false
	pinnedPiece := -1
	inDoubleCheck := false

	marshalSquare := b.BoardSquares[b.MarshalCoords[b.TurnColor]].Value.(*LLStack).Stack
	marshalHashmap := make(map[int]bool)
	for _, move := range b.GetPseudoLegalMoves(marshalSquare.Top.Value.(int), b.MarshalCoords[b.TurnColor], marshalSquare.Length) {
		marshalHashmap[move] = true
	}

	currentStackNode := b.StackList[GetOppositeColor(b.TurnColor)].Head
	for currentStackNode != nil {
		stack := currentStackNode.Value.(*LLStack)
		piece := stack.Stack.Top.Value.(int) % 13
		coord := stack.Coordinate
		// add sliding pieces for pin check
		if stack.Stack.Length == 3 {
			switch piece {
			case MUSKETEER, SAMURAI, CANNON, SPY:
				// TODO enemy sliding piece
				// if no piece in between, marshal in check
				// if in check, marshal must evade out of line, place/move a piece in between, or capture
				// if enemy piece in between, no restriction
				// if ally piece in between, it is pinned, piece can only move in line or capture attacking piece
				moves, tempXRayMoves, tempInCheck, tempInPath := b.CheckEnemyRanging(piece, coord)

				if tempInCheck {
					if inCheck {
						inDoubleCheck = true
					} else {
						inCheck = true
					}
				}
				if tempInPath {
					piecesInbetween := []int{}
					for _, move := range tempXRayMoves {
						if move.inBetween && move.occupied && GetColor(b.BoardSquares[move.coordinate].Value.(*LLStack).Stack.Top.Value.(int)) == b.TurnColor {
							piecesInbetween = append(piecesInbetween, move.coordinate)
						}
					}
					if len(piecesInbetween) <= 1 {
						enemyXRaySquares = XRay{
							coordinate: coord,
							path:       tempXRayMoves,
						}
					}
					if len(piecesInbetween) == 1 && GetColor(b.BoardSquares[piecesInbetween[0]].Value.(*LLStack).Stack.Top.Value.(int)) == b.TurnColor {
						pinnedPiece = piecesInbetween[0]
					}
					// enemyXRaySquares = append(enemyXRaySquares, XRay{
					// 	coordinate: coord,
					// 	path:       tempXRayMoves,
					// })
				}

				for _, move := range moves {
					if marshalHashmap[move] {
						delete(marshalHashmap, move)
					}
				}

				// enemyMoveList = append(enemyMoveList, PseudoMove{
				// 	coordinate: coord,
				// 	moveList:   moves,
				// })
			case TACTICIAN:
				lowerStack := stack.Stack.Top.Prev
				lowerPiece := lowerStack.Value.(int)
				if lowerPiece%13 == TACTICIAN {
					lowerPiece = lowerStack.Prev.Value.(int)
				}
				switch lowerPiece {
				case MUSKETEER, SAMURAI, CANNON, SPY:
					moves, tempXRayMoves, tempInCheck, tempInPath := b.CheckEnemyRanging(lowerPiece, coord)

					if tempInCheck {
						if inCheck {
							inDoubleCheck = true
						} else {
							inCheck = true
						}
					}
					if tempInPath {
						piecesInbetween := []int{}
						for _, move := range tempXRayMoves {
							if move.inBetween && move.occupied && GetColor(b.BoardSquares[move.coordinate].Value.(*LLStack).Stack.Top.Value.(int)) == b.TurnColor {
								piecesInbetween = append(piecesInbetween, move.coordinate)
							}
						}
						if len(piecesInbetween) <= 1 {
							enemyXRaySquares = XRay{
								coordinate: coord,
								path:       tempXRayMoves,
							}
						}
						if len(piecesInbetween) == 1 && GetColor(b.BoardSquares[piecesInbetween[0]].Value.(*LLStack).Stack.Top.Value.(int)) == b.TurnColor {
							pinnedPiece = piecesInbetween[0]
						}
					}

					for _, move := range moves {
						if marshalHashmap[move] {
							delete(marshalHashmap, move)
						}
					}

					// enemyMoveList = append(enemyMoveList, PseudoMove{
					// 	coordinate: coord,
					// 	moveList:   moves,
					// })
				default:
					moves := b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length)
					for _, move := range moves {
						if marshalHashmap[move] {
							delete(marshalHashmap, move)
						}
						if b.MarshalCoords[b.TurnColor] == move {
							if inCheck {
								inDoubleCheck = true
							} else {
								inCheck = true
							}
						}
					}

					// enemyMoveList = append(enemyMoveList, PseudoMove{
					// 	coordinate: coord,
					// 	moveList:   b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length),
					// })
				}
			default:
				moves := b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length)
				for _, move := range moves {
					if marshalHashmap[move] {
						delete(marshalHashmap, move)
					}
					if b.MarshalCoords[b.TurnColor] == move {
						if inCheck {
							inDoubleCheck = true
						} else {
							inCheck = true
						}
					}
				}
				// enemyMoveList = append(enemyMoveList, PseudoMove{
				// 	coordinate: coord,
				// 	moveList:   b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length),
				// })
			}
		} else {
			moves := b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length)
			for _, move := range moves {
				if marshalHashmap[move] {
					delete(marshalHashmap, move)
				}
				if b.MarshalCoords[b.TurnColor] == move {
					if inCheck {
						inDoubleCheck = true
					} else {
						inCheck = true
					}
				}
			}
			// enemyMoveList = append(enemyMoveList, PseudoMove{
			// 	coordinate: coord,
			// 	moveList:   b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length),
			// })
		}
		b.CheckEnemyMoves()
		// pseudoMoves := b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length)
		currentStackNode = currentStackNode.Next
	}

	// If marshal in check or in path of xray, remove moves from marshal
	xrayHashmap := make(map[int]bool)
	if len(enemyXRaySquares.path) > 0 {
		for _, move := range enemyXRaySquares.path {
			xrayHashmap[move.coordinate] = true
			if inCheck {
				if marshalHashmap[move.coordinate] {
					delete(marshalHashmap, move.coordinate)
				}
			}
		}
	}
	xrayInbetweenHashmap := make(map[int]bool)
	if inCheck && len(enemyXRaySquares.path) > 0 {
		for _, move := range enemyXRaySquares.path {
			if move.inBetween {
				xrayInbetweenHashmap[move.coordinate] = true
			}
		}
	}

	// TODO Handle capturing attacking piece
	// TODO Handle uncovering attack piece

	moveList := []PseudoMove{}

	// Loop through pieces of current player
	currentStackNode = b.StackList[b.TurnColor].Head
	for currentStackNode != nil {
		stack := currentStackNode.Value.(*LLStack)

		if stack.Stack.Top.Value.(int)%13 != MARSHAL {
			tempMoves := []int{}
			for _, move := range b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length) {
				if inCheck {
					if inDoubleCheck {

					} else if len(enemyXRaySquares.path) > 0 && xrayInbetweenHashmap[move] {
						tempMoves = append(tempMoves, move)
					}
				} else if !(!xrayHashmap[move] && stack.Coordinate == pinnedPiece) {
					// !(if pinned piece, and move is in path)
					tempMoves = append(tempMoves, move)
				}
			}
			moveList = append(moveList, PseudoMove{
				coordinate: stack.Coordinate,
				moveList:   tempMoves,
			})
		}
		// log.Println(pseudoMoves)
		// log.Println(moveList)
		currentStackNode = currentStackNode.Next
	}

	// Look for x-ray attacks for pin potential
	// Traverse enemy stacklist for ranging piece
	log.Println("marshal moves: ", marshalHashmap)
	log.Println("enemy xray: ", enemyXRaySquares)
	log.Println("enemy xray hashmap: ", xrayHashmap)
	log.Println("inCheck: ", inCheck)
	log.Println("inDoubleCheck: ", inDoubleCheck)
	log.Println("current player moves: ", moveList)

}

func (b *Board) CheckEnemyMoves() {

}

type XRaySquares struct {
	coordinate int
	inBetween  bool
	occupied   bool
}

// Returns moves, xraysquares, checked, inpath
func (b *Board) CheckEnemyRanging(piece int, coord int) ([]int, []XRaySquares, bool, bool) {

	var moves []int
	inCheck := false
	inPath := false
	var xraySquares []XRaySquares

	switch piece % 13 {
	case MUSKETEER:
		tempMoves, tempBlocked, tempInPath := b.XRayRangingPiece(coord, ColorOffset(b.TurnColor, -12))
		if !tempBlocked && tempInPath {
			inCheck = true
		}
		if tempInPath {
			inPath = true
		}
		moves, xraySquares = XRayHandler(moves, xraySquares, tempMoves, tempBlocked, tempInPath)
		// if !tempBlocked && tempInPath {
		// 	// is inCheck by ranging piece
		// 	inCheck = true
		// 	inPath = true
		// 	for _, move := range tempMoves {
		// 		if move.inBetween == true {
		// 			moves = append(moves, move.coordinate)
		// 		}
		// 		xraySquares = append(xraySquares, move)
		// 	}
		// } else if tempBlocked && tempInPath {
		// 	blocked = true
		// 	for _, move := range tempMoves {
		// 		if move.inBetween == true {
		// 			moves = append(moves, move.coordinate)
		// 		}
		// 		xraySquares = append(xraySquares, move)
		// 	}
		// } else {
		// 	for _, move := range tempMoves {
		// 		moves = append(moves, move.coordinate)
		// 		if move.occupied == true {
		// 			break
		// 		}
		// 	}
		// }
	case CANNON:
		offsets := [4]int{-12, -1, 1, 12}
		for _, offset := range offsets {
			tempMoves, tempBlocked, tempInPath := b.XRayRangingPiece(coord, offset)
			if !tempBlocked && tempInPath {
				inCheck = true
			}
			if tempInPath {
				inPath = true
			}
			moves, xraySquares = XRayHandler(moves, xraySquares, tempMoves, tempBlocked, tempInPath)
		}
	case SPY:
		offsets := [8]int{-12, -1, 1, 12, -11, -13, 11, 13}
		for _, offset := range offsets {
			tempMoves, tempBlocked, tempInPath := b.XRayRangingPiece(coord, offset)
			if !tempBlocked && tempInPath {
				inCheck = true
			}
			if tempInPath {
				inPath = true
			}
			moves, xraySquares = XRayHandler(moves, xraySquares, tempMoves, tempBlocked, tempInPath)
		}
	case SAMURAI:
		offsets := [4]int{-11, -13, 11, 13}
		for _, offset := range offsets {
			tempMoves, tempBlocked, tempInPath := b.XRayRangingPiece(coord, offset)
			if !tempBlocked && tempInPath {
				inCheck = true
			}
			if tempInPath {
				inPath = true
			}
			moves, xraySquares = XRayHandler(moves, xraySquares, tempMoves, tempBlocked, tempInPath)
		}
	}

	// it's actually impossible to be double checked by two ranging pieces in a move
	// if checked by a ranging piece, marshal must move out of piece's line
	// if x-rayed, check if there are any pieces in between
	// if == 1, then no restrictions
	// if == 1 and piece is current player's piece, it can only move in piece's line
	// if >= 1, not pinned
	return moves, xraySquares, inCheck, inPath
}

// Returns moves, xraysquares, checked, inpath
func XRayHandler(moves []int, xraySquares []XRaySquares, tempMoves []XRaySquares, tempBlocked bool, tempInPath bool) ([]int, []XRaySquares) {

	if !tempBlocked && tempInPath {
		// is checked by ranging piece
		for _, move := range tempMoves {
			if move.inBetween {
				moves = append(moves, move.coordinate)
			}
			xraySquares = append(xraySquares, move)
		}
	} else if tempBlocked && tempInPath {
		// marshal in path but blocked (not in check)
		for _, move := range tempMoves {
			if move.inBetween {
				moves = append(moves, move.coordinate)
			}
			xraySquares = append(xraySquares, move)
		}
	} else {
		// otherwise, just add moves normally
		for _, move := range tempMoves {
			moves = append(moves, move.coordinate)
			if move.occupied {
				break
			}
		}
	}
	return moves, xraySquares
}

// Generates ranging move xray. Return xray path {coord, inbetween, occupied}, ranging piece blocked, if marshal in path
func (b *Board) XRayRangingPiece(coordinate int, offset int) ([]XRaySquares, bool, bool) {
	xraySquares := []XRaySquares{}
	inPath := false
	blocked := false

	i := coordinate + offset
	currSquare := b.BoardSquares[i]
	for {
		if currSquare == nil {
			// get moves in line
			xraySquares = append(xraySquares, XRaySquares{
				coordinate: i,
				inBetween:  !inPath,
				occupied:   false,
			})
		} else if currSquare.Value == -1 {
			break
		} else {
			// gets pieces in between, does not distinguish enemy and ally pieces
			square := currSquare.Value.(*LLStack)
			piece := square.Stack.Top.Value.(int)
			// sees if marshal is in line
			if piece%13 == MARSHAL && GetColor(piece) == b.TurnColor {
				inPath = true
			} else {
				blocked = true
			}
			xraySquares = append(xraySquares, XRaySquares{
				coordinate: i,
				inBetween:  !inPath,
				occupied:   true,
			})
		}
		i += offset
		currSquare = b.BoardSquares[i]
	}

	return xraySquares, blocked, inPath
}

// X-rays pieces from marshal to see potentially pinned pieces. Returns if ranging offset in range of a piece, potentially pinned squares, and pseudo valid squares.
func (b *Board) GetCheckRangingPiece(offset int) (bool, []int, []int) {
	inPath := false
	potentiallyPinnedSquares := []int{} // includes enemy pieces to check if ranging piece is blocked
	squares := []int{}

	i := b.MarshalCoords[b.TurnColor] + offset
	currSquare := b.BoardSquares[i]
	endLoop := false
	for endLoop {
		if currSquare == nil {
			squares = append(squares, i)
		} else if currSquare.Value == -1 {
			endLoop = true
		} else {

			piece := currSquare.Value.(*ds.Stack).Top.Value.(int)

			switch piece % 13 {
			case SPY, CANNON, SAMURAI, TACTICIAN, MUSKETEER:
				if GetColor(piece) != b.TurnColor {
					inPath = true
					endLoop = true
				} else {
					squares = append(squares, i)
					potentiallyPinnedSquares = append(potentiallyPinnedSquares, piece)
				}
			default:
				squares = append(squares, i)
				potentiallyPinnedSquares = append(potentiallyPinnedSquares, piece)
			}

		}
	}

	// TODO might have to check line behind marshal too?
	// will check behind anyways, can probably just combine

	return inPath, potentiallyPinnedSquares, squares
}

// Will stop if out of bounds or a piece is in the way. However, you should double check if attacking or stacking is possible including stacking on marshal.
// Must create a separate function to x-ray for checking pins
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

// Generates psuedo-legal moves for a piece at a coordinate.
// Still must determine out of bound squares, fully stacked pieces, marshal, fortress, is checked, and pinned.
// Does not differentiate from attacking and stacking.
// Does not generate moves from hand.
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
		// Movement becomes the piece under it but on tactician's tier
		switch tier {
		case 1:
			offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
		case 2, 3:
			piece := b.BoardSquares[coordinate].Value.(*LLStack).Stack.Top.Prev
			sameColorPiece := piece.Value.(int) % 13
			if sameColorPiece == TACTICIAN && tier == 3 {
				piece = piece.Prev
			}
			if GetColor(piece.Value.(int)) == 1 {
				sameColorPiece += 13
			}
			squares = append(squares, b.GetPseudoLegalMoves(sameColorPiece, coordinate, tier)...)
		}
	case MARSHAL:
		offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
	}

	// Reverses direction for black
	for _, offset := range offsets {
		if GetColor(piece) == 1 {
			squares = append(squares, coordinate-offset)
		} else {
			squares = append(squares, coordinate+offset)
		}
	}

	tempSquares := []int{}
	for _, move := range squares {
		if b.BoardSquares[move] == nil {
			tempSquares = append(tempSquares, move)
		} else if b.BoardSquares[move].Value != -1 {
			tempSquares = append(tempSquares, move)
		}
	}

	return tempSquares
}
