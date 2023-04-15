package gungi

import (
	"log"
	"strings"

	"github.com/TonyLeCode/gungi.go/server/utils"
)

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
	Coordinate int
	MoveList   []int
}

type XRay struct {
	Coordinate int
	Path       []XRaySquares
}

func (b *Board) MakeMove(piece int, fromCoord int, moveType int, toCoord int) bool {
	isValid := b.ValidateMove(piece, fromCoord, moveType, toCoord)
	log.Println("isValid: ", isValid)
	if isValid {
		switch moveType {
		case MOVE:
			b.RemovePiece(fromCoord)
			b.PlacePiece(piece, toCoord)

			fromFile, fromRank := CoordsToNotation(fromCoord)
			toFile, toRank := CoordsToNotation(toCoord)
			str := EncodeSingleChar(piece) + fromFile + fromRank + "-" + toFile + toRank
			b.History = append(b.History, str)
		case STACK:
			b.RemovePiece(fromCoord)
			b.PlacePiece(piece, toCoord)

			fromFile, fromRank := CoordsToNotation(fromCoord)
			toFile, toRank := CoordsToNotation(toCoord)
			str := EncodeSingleChar(piece) + fromFile + fromRank + "-" + toFile + toRank
			b.History = append(b.History, str)
		case ATTACK:
			b.RemovePiece(fromCoord)
			b.RemovePiece(toCoord)
			b.PlacePiece(piece, toCoord)

			fromFile, fromRank := CoordsToNotation(fromCoord)
			toFile, toRank := CoordsToNotation(toCoord)
			str := EncodeSingleChar(piece) + fromFile + fromRank + "x" + EncodeSingleChar(b.BoardSquares[toCoord].Value.(*LLStack).Stack.Top.Value.(int)) + toFile + toRank
			b.History = append(b.History, str)
		case PLACE:
			b.Hand[piece]--
			b.PlacePiece(piece, toCoord)

			toFile, toRank := CoordsToNotation(toCoord)
			str := EncodeSingleChar(piece) + toFile + toRank
			b.History = append(b.History, str)
		case READY:
			b.Ready[b.TurnColor] = true
			b.TurnColor = GetOppositeColor(b.TurnColor)
		}

		if b.Ready[0] == b.Ready[1] {
			b.TurnColor = GetOppositeColor(b.TurnColor)
		}
		b.TurnNumber++
	}

	return isValid
}

func (b *Board) UndoMove() {
	if len(b.History) == 0 {
		return
	}
	lastMove := b.History[len(b.History)-1]
	if strings.Contains(lastMove, "-") {
		// Move and Stack
		fromPiece := DecodeSingleChar(string(lastMove[0]))
		fromCoord := CoordsToSquare(LetterToFile(string(lastMove[2]))-1, RevertRank(string(lastMove[1]))-1)
		toCoord := CoordsToSquare(LetterToFile(string(lastMove[5]))-1, RevertRank(string(lastMove[4]))-1)

		b.RemovePiece(toCoord)
		b.PlacePiece(fromPiece, fromCoord)
		b.History = utils.RemoveIndexStr(b.History, len(b.History)-1)
		b.TurnColor = GetOppositeColor(b.TurnColor)
		b.TurnNumber--
	} else if strings.Contains(lastMove, "x") {
		// Attack
		fromPiece := DecodeSingleChar(string(lastMove[0]))
		toPiece := DecodeSingleChar(string(lastMove[4]))
		fromCoord := CoordsToSquare(LetterToFile(string(lastMove[2]))-1, RevertRank(string(lastMove[1]))-1)
		toCoord := CoordsToSquare(LetterToFile(string(lastMove[6]))-1, RevertRank(string(lastMove[5]))-1)

		b.RemovePiece(toCoord)
		b.PlacePiece(toPiece, toCoord)
		b.PlacePiece(fromPiece, fromCoord)
		b.History = utils.RemoveIndexStr(b.History, len(b.History)-1)
		b.TurnColor = GetOppositeColor(b.TurnColor)
		b.TurnNumber--
	} else if lastMove == "Ready" {
		b.TurnColor = GetOppositeColor(b.TurnColor)
		b.Ready[b.TurnColor] = false
		b.TurnNumber--
	} else if len(lastMove) == 3 {
		// Place from hand
		piece := DecodeSingleChar(string(lastMove[0]))
		toCoord := CoordsToSquare(LetterToFile(string(lastMove[2]))-1, RevertRank(string(lastMove[1]))-1)

		b.RemovePiece(toCoord)
		b.Hand[piece]++
		b.History = utils.RemoveIndexStr(b.History, len(b.History)-1)
		b.TurnColor = GetOppositeColor(b.TurnColor)
		b.TurnNumber--
	}

	log.Println(lastMove)
}

func (b *Board) ValidateMove(piece int, fromCoord int, moveType int, toCoord int) bool {
	if (!b.Ready[0] || !b.Ready[1]) && (moveType != PLACE && moveType != READY) {
		return false
	}

	if GetColor(piece) != b.TurnColor {
		return false
	}

	inCheck, inDoubleCheck, _, enemyXRaySquares, attackPiece, filteredMoves := b.GenerateLegalMoves()

	toSquare := b.BoardSquares[toCoord]

	isValid := false

	switch moveType {
	case MOVE:
		if b.BoardSquares[fromCoord] == nil {
			return false
		} else if b.BoardSquares[fromCoord].Value == -1 {
			return false
		} else if b.BoardSquares[fromCoord].Value.(*LLStack).Stack.Top.Value != piece {
			return false
		}
		for _, move := range filteredMoves[fromCoord] {
			if move == toCoord && (toSquare == nil || toSquare.Value.(*LLStack).Stack.Length == 0) {
				isValid = true
			}
		}
	case STACK:
		if b.BoardSquares[fromCoord] == nil {
			return false
		} else if b.BoardSquares[fromCoord].Value == -1 {
			return false
		} else if b.BoardSquares[fromCoord].Value.(*LLStack).Stack.Top.Value != piece {
			return false
		}
		for _, move := range filteredMoves[fromCoord] {
			if move == toCoord && toSquare != nil && (toSquare.Value.(*LLStack).Stack.Length > 0 && toSquare.Value.(*LLStack).Stack.Length < 3) {
				isValid = true
			}
		}
	case ATTACK:
		if b.BoardSquares[fromCoord] == nil {
			return false
		} else if b.BoardSquares[fromCoord].Value == -1 {
			return false
		} else if b.BoardSquares[fromCoord].Value.(*LLStack).Stack.Top.Value != piece {
			return false
		}
		for _, move := range filteredMoves[fromCoord] {
			if move == toCoord && toSquare != nil && toSquare.Value.(*LLStack).Stack.Length > 0 && GetColor(toSquare.Value.(*LLStack).Stack.Top.Value.(int)) != b.TurnColor {
				isValid = true
			}
		}
	case PLACE:
		if b.Hand[piece] == 0 {
			return false
		}

		invalidSquares := []int{}
		if b.TurnColor == 0 {
			invalidSquares = []int{37, 38, 39, 40, 41, 42, 43, 44, 45,
				49, 50, 51, 52, 53, 54, 55, 56, 57,
				61, 62, 63, 64, 65, 66, 67, 68, 69}
		} else if b.TurnColor == 1 {
			invalidSquares = []int{109, 110, 111, 112, 113, 114, 115, 116, 117,
				121, 122, 123, 124, 125, 126, 127, 128, 129,
				133, 134, 135, 136, 137, 138, 139, 140, 141}
		}
		for _, coord := range invalidSquares {
			if toCoord == coord {
				return false
			}
		}

		if piece%13 == PAWN {
			fileCoord := toCoord % 12
			for i := 0; i < 9; i++ {
				if b.BoardSquares[(12*i)+36+fileCoord] != nil {
					currentNode := b.BoardSquares[(12*i)+36+fileCoord].Value.(*LLStack).Stack.Bottom
					for currentNode != nil {
						if currentNode.Value.(int)%13 == PAWN && GetColor(currentNode.Value.(int)) == b.TurnColor {
							return false
						}
						currentNode = currentNode.Next
					}
				}
			}
		}

		if !inDoubleCheck && toSquare == nil || toSquare.Value.(*LLStack).Stack.Length < 3 {
			if attackPiece != -1 && toCoord == attackPiece {
				isValid = true
			} else if len(enemyXRaySquares.Path) != 0 && inCheck {
				for _, move := range enemyXRaySquares.Path {
					if move.Coordinate == toCoord && move.InBetween {
						isValid = true
					}
				}
				if toCoord == enemyXRaySquares.Coordinate {
					isValid = true
				}
			} else {
				isValid = true
			}
		}

	case READY: //
		if b.Ready[b.TurnColor] {
			isValid = false
		}
	}
	return isValid
}

// Generates legal moves
func (b *Board) GenerateLegalMoves() (bool, bool, bool, XRay, int, map[int][]int) {
	var enemyXRaySquares XRay
	inCheck := false
	pinnedPiece := -1
	inDoubleCheck := false
	attackPiece := -1
	checkmate := false

	marshalSquare := b.BoardSquares[b.MarshalCoords[b.TurnColor]].Value.(*LLStack).Stack
	marshalHashmap := make(map[int]bool)
	for _, move := range b.GetPseudoLegalMoves(marshalSquare.Top.Value.(int), b.MarshalCoords[b.TurnColor], marshalSquare.Length) {
		marshalHashmap[move] = true
	}

	tempAttackPiece, tempXRay, tempPinnedPiece, tempInCheck, tempInDoubleCheck := b.CheckEnemyMoves(&marshalHashmap, inCheck, inDoubleCheck)
	if len(tempXRay.Path) != 0 {
		enemyXRaySquares = tempXRay
	}
	if tempPinnedPiece != -1 {
		pinnedPiece = tempPinnedPiece
	}
	if tempAttackPiece != -1 {
		attackPiece = tempAttackPiece
	}
	inCheck = tempInCheck
	inDoubleCheck = tempInDoubleCheck

	// If marshal in check or in path of xray, remove moves from marshal
	xrayHashmap := make(map[int]bool)
	if len(enemyXRaySquares.Path) > 0 {
		for _, move := range enemyXRaySquares.Path {
			xrayHashmap[move.Coordinate] = true
			if inCheck {
				if marshalHashmap[move.Coordinate] {
					delete(marshalHashmap, move.Coordinate)
				}
			}
		}
	}
	xrayInbetweenHashmap := make(map[int]bool)
	if inCheck && len(enemyXRaySquares.Path) > 0 {
		for _, move := range enemyXRaySquares.Path {
			if move.InBetween {
				xrayInbetweenHashmap[move.Coordinate] = true
			}
		}
	}

	moveList := []PseudoMove{}

	// Loop through pieces of current player
	currentStackNode := b.StackList[b.TurnColor].Head
	for currentStackNode != nil {
		stack := currentStackNode.Value.(*LLStack)

		pieceUnder := -1
		snared := false
		if stack.Stack.Length > 1 {
			pieceUnder = stack.Stack.Top.Prev.Value.(int)
		}

		if pieceUnder != -1 && GetColor(pieceUnder) != b.TurnColor {
			moves := b.GetPseudoLegalMoves(pieceUnder, stack.Coordinate, stack.Stack.Length-1)

			for _, move := range moves {
				if b.MarshalCoords[b.TurnColor] == move {
					snared = true
				}
			}
		}

		if stack.Stack.Top.Value.(int)%13 != MARSHAL && !snared {
			tempMoves := []int{}
			for _, move := range b.GetPseudoLegalMoves(stack.Stack.Top.Value.(int), stack.Coordinate, stack.Stack.Length) {
				if inCheck {
					if inDoubleCheck {

					} else if len(enemyXRaySquares.Path) > 0 && xrayInbetweenHashmap[move] {
						tempMoves = append(tempMoves, move)
					} else if move == attackPiece || move == enemyXRaySquares.Coordinate {
						tempMoves = append(tempMoves, move)
					}
				} else if stack.Coordinate == pinnedPiece && move == enemyXRaySquares.Coordinate {
					tempMoves = append(tempMoves, move)
				} else if !(!xrayHashmap[move] && stack.Coordinate == pinnedPiece) && b.MarshalCoords[b.TurnColor] != move {
					// !(if pinned piece, and move is in path)
					tempMoves = append(tempMoves, move)
				}
			}
			moveList = append(moveList, PseudoMove{
				Coordinate: stack.Coordinate,
				MoveList:   tempMoves,
			})
		}
		currentStackNode = currentStackNode.Next
	}

	filteredMoves := make(map[int][]int)

	for _, moves := range moveList {
		filteredMoves[moves.Coordinate] = append(filteredMoves[moves.Coordinate], moves.MoveList...)
	}

	// filteredMoves = append(filteredMoves, marshalHashmap)
	for i := range marshalHashmap {
		filteredMoves[b.MarshalCoords[b.TurnColor]] = append(filteredMoves[b.MarshalCoords[b.TurnColor]], i)
	}

	for i, moves := range filteredMoves {
		if len(moves) == 0 {
			delete(filteredMoves, i)
		}
	}

	if len(filteredMoves) == 0 {
		checkmate = true
	}

	log.Println("marshal moves: ", marshalHashmap)
	log.Println("enemy xray: ", enemyXRaySquares)
	log.Println("enemy xray hashmap: ", xrayHashmap)
	log.Println("inCheck: ", inCheck)
	log.Println("inDoubleCheck: ", inDoubleCheck)
	log.Println("checkMate: ", checkmate)
	log.Println("current player moves: ", filteredMoves)
	return inCheck, inDoubleCheck, checkmate, enemyXRaySquares, attackPiece, filteredMoves
}

type XRaySquares struct {
	Coordinate int
	InBetween  bool
	Occupied   bool
}

// Returns moves, xraysquares, checked, inpath
func (b *Board) CheckEnemyRanging(piece int, coord int) ([]int, []XRaySquares, bool, bool) {

	var moves []int
	inCheck := false
	inPath := false
	var xraySquares []XRaySquares

	offsets := []int{}

	switch piece % 13 {
	case MUSKETEER:
		offsets = []int{-12}
	case CANNON:
		offsets = []int{-12, -1, 1, 12}
	case SPY:
		offsets = []int{-12, -1, 1, 12, -11, -13, 11, 13}
	case SAMURAI:
		offsets = []int{-11, -13, 11, 13}
	}

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

	// it's actually impossible to be double checked by two ranging pieces in a move
	// if checked by a ranging piece, marshal must move out of piece's line
	// if x-rayed, check if there are any pieces in between
	// if == 1, then no restrictions
	// if == 1 and piece is current player's piece, it can only move in piece's line
	// if >= 1, not pinned
	return moves, xraySquares, inCheck, inPath
}

// 	}

// 	for _, offset := range offsets {
// 		tempMoves, tempBlocked, tempInPath := b.XRayRangingPiece(coord, offset)
// 		if !tempBlocked && tempInPath {
// 			inCheck = true
// 		}
// 		if tempInPath {
// 			inPath = true
// 		}
// 		moves, xraySquares = XRayHandler(moves, xraySquares, tempMoves, tempBlocked, tempInPath)
// 	}

// Returns moves, xraysquares, checked, inpath
func XRayHandler(moves []int, xraySquares []XRaySquares, tempMoves []XRaySquares, tempBlocked bool, tempInPath bool) ([]int, []XRaySquares) {

	if !tempBlocked && tempInPath {
		// is checked by ranging piece
		for _, move := range tempMoves {
			if move.InBetween {
				moves = append(moves, move.Coordinate)
			}
			xraySquares = append(xraySquares, move)
		}
	} else if tempBlocked && tempInPath {
		// marshal in path but blocked (not in check)
		for _, move := range tempMoves {
			if move.InBetween {
				moves = append(moves, move.Coordinate)
			}
			xraySquares = append(xraySquares, move)
		}
	} else {
		// otherwise, just add moves normally
		for _, move := range tempMoves {
			moves = append(moves, move.Coordinate)
			if move.Occupied {
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
				Coordinate: i,
				InBetween:  !inPath,
				Occupied:   false,
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
				Coordinate: i,
				InBetween:  !inPath,
				Occupied:   true,
			})
		}
		i += offset
		currSquare = b.BoardSquares[i]
	}

	return xraySquares, blocked, inPath
}

func (b *Board) CheckEnemyMoves(marshalHashmap *map[int]bool, inCheck bool, inDoubleCheck bool) (int, XRay, int, bool, bool) {
	var enemyXRaySquares XRay
	pinnedPiece := -1
	attackPiece := -1

	currentStackNode := b.StackList[GetOppositeColor(b.TurnColor)].Head
	for currentStackNode != nil {
		// log.Println("node: ", currentStackNode.Value)
		stack := currentStackNode.Value.(*LLStack)
		piece := stack.Stack.Top.Value.(int) % 13
		coord := stack.Coordinate

		if stack.Stack.Length == 3 && piece == TACTICIAN {
			piece = stack.Stack.Top.Prev.Value.(int) % 13
			if piece == TACTICIAN {
				piece = stack.Stack.Bottom.Value.(int) % 13
			}
		}

		if stack.Stack.Length == 3 && piece >= 7 && piece <= 10 {
			// ranging
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
					if move.InBetween && move.Occupied && GetColor(b.BoardSquares[move.Coordinate].Value.(*LLStack).Stack.Top.Value.(int)) == b.TurnColor {
						piecesInbetween = append(piecesInbetween, move.Coordinate)
					}
				}
				if len(piecesInbetween) <= 1 {
					enemyXRaySquares = XRay{
						Coordinate: coord,
						Path:       tempXRayMoves,
					}
				}
				if len(piecesInbetween) == 1 && GetColor(b.BoardSquares[piecesInbetween[0]].Value.(*LLStack).Stack.Top.Value.(int)) == b.TurnColor {
					pinnedPiece = piecesInbetween[0]
				}
			}

			for _, move := range moves {
				if (*marshalHashmap)[move] {
					delete(*marshalHashmap, move)
				}
			}
		} else {
			// every other piece
			moves := b.GetPseudoLegalMoves(piece, stack.Coordinate, stack.Stack.Length)
			for _, move := range moves {
				if (*marshalHashmap)[move] {
					delete(*marshalHashmap, move)
				}
				if b.MarshalCoords[b.TurnColor] == move {
					attackPiece = stack.Coordinate
					if inCheck {
						inDoubleCheck = true
					} else {
						inCheck = true
					}
				}
			}
		}
		currentStackNode = currentStackNode.Next
	}
	return attackPiece, enemyXRaySquares, pinnedPiece, inCheck, inDoubleCheck
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
				sameColorPiece = piece.Value.(int) % 13
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
