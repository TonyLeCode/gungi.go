package revised

import (
	"errors"
	"log"
	"strings"

	"github.com/whitemonarch/gungi-server/server/gungi/utils"
)

// Piece must be colorless
func coloredPiece(piece int, color int) int {
	if color == 0 {
		return piece
	} else {
		return piece + 13
	}
}

func (r *Revised) MakeMove(piece int, fromCoord int, moveType int, toCoord int) error {
	err := r.ValidateMove(piece, fromCoord, moveType, toCoord)
	if err != nil {
		return err
	}

	fromFile, fromRank := utils.CoordsToNotation(fromCoord)
	toFile, toRank := utils.CoordsToNotation(toCoord)
	var str string

	switch moveType {
	// Don't combine move and stack, possibly change later
	case MOVE:
		r.RemovePiece(fromCoord)
		r.PlacePiece(piece, toCoord)

		str = utils.EncodeSingleChar(piece) + fromFile + fromRank + "-" + toFile + toRank
		// Ex. L4D-4D
	case STACK:
		r.RemovePiece(fromCoord)
		r.PlacePiece(piece, toCoord)

		str = utils.EncodeSingleChar(piece) + fromFile + fromRank + "-" + toFile + toRank
		// Ex. L4D-4D
	case ATTACK:
		str = utils.EncodeSingleChar(piece) + fromFile + fromRank + "x" + utils.EncodeSingleChar(r.BoardSquares[toCoord].GetTop()) + toFile + toRank

		r.RemovePiece(fromCoord)
		r.RemovePiece(toCoord)
		r.PlacePiece(piece, toCoord)
		// Ex. L4Dxy5D
	case PLACE:
		r.Hand[piece]--
		r.PlacePiece(piece, toCoord)

		str = utils.EncodeSingleChar(piece) + toFile + toRank
		// Ex. L4D
	case READY:
		if r.TurnColor == 0 {
			str = "w-r"
		} else {
			str = "b-r"
		}
		r.Ready[r.TurnColor] = true
		if r.Ready[utils.OppositeColor(r.TurnColor)] {
			r.TurnColor = BLACK
		} else {
			r.TurnColor = utils.OppositeColor(r.TurnColor)
		}
	}
	r.History = append(r.History, str)

	if piece%13 == MARSHAL {
		r.MarshalCoords[r.TurnColor] = toCoord
	}

	if r.Ready[0] == r.Ready[1] {
		r.TurnColor = utils.OppositeColor(r.TurnColor)
	}
	r.TurnNumber++

	return nil
}

func RemoveIndexStr(s []string, index int) []string {
	ret := make([]string, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func (r *Revised) UndoMove() {
	//TODO regex and validation
	if len(r.History) == 0 {
		return
	}

	lastMove := r.History[len(r.History)-1]
	if lastMove == "w-r" || lastMove == "b-r" {
		//ready
		r.Ready[utils.OppositeColor(r.TurnColor)] = false
		r.History = RemoveIndexStr(r.History, len(r.History)-1)
		r.TurnColor = utils.OppositeColor(r.TurnColor)
		r.TurnNumber--
	} else if strings.Contains(lastMove, "x") {
		//attack
		fromPiece := utils.DecodeSingleChar(string(lastMove[0]))
		toPiece := utils.DecodeSingleChar(string(lastMove[4]))
		fromCoord := utils.CoordsToSquare(utils.LetterToFile(string(lastMove[2]))-1, utils.RevertRank(string(lastMove[1]))-1)
		toCoord := utils.CoordsToSquare(utils.LetterToFile(string(lastMove[6]))-1, utils.RevertRank(string(lastMove[5]))-1)

		r.RemovePiece(toCoord)
		r.PlacePiece(toPiece, toCoord)
		r.PlacePiece(fromPiece, fromCoord)

		if fromPiece%13 == MARSHAL {
			r.MarshalCoords[utils.OppositeColor(r.TurnColor)] = fromCoord
		}

		r.History = RemoveIndexStr(r.History, len(r.History)-1)
		r.TurnColor = utils.OppositeColor(r.TurnColor)
		r.TurnNumber--
	} else if strings.Contains(lastMove, "-") {
		//stack or move
		fromPiece := utils.DecodeSingleChar(string(lastMove[0]))
		fromCoord := utils.CoordsToSquare(utils.LetterToFile(string(lastMove[2]))-1, utils.RevertRank(string(lastMove[1]))-1)
		toCoord := utils.CoordsToSquare(utils.LetterToFile(string(lastMove[5]))-1, utils.RevertRank(string(lastMove[4]))-1)

		r.RemovePiece(toCoord)
		r.PlacePiece(fromPiece, fromCoord)

		if fromPiece%13 == MARSHAL {
			r.MarshalCoords[utils.OppositeColor(r.TurnColor)] = fromCoord
		}

		r.History = RemoveIndexStr(r.History, len(r.History)-1)
		r.TurnColor = utils.OppositeColor(r.TurnColor)
		r.TurnNumber--
	} else if len(lastMove) == 3 {
		//place
		piece := utils.DecodeSingleChar(string(lastMove[0]))
		toCoord := utils.CoordsToSquare(utils.LetterToFile(string(lastMove[2]))-1, utils.RevertRank(string(lastMove[1]))-1)

		r.RemovePiece(toCoord)
		r.Hand[piece]++

		if piece%13 == MARSHAL {
			r.MarshalCoords[utils.OppositeColor(r.TurnColor)] = 0
		}

		r.History = RemoveIndexStr(r.History, len(r.History)-1)
		if r.Ready[0] == r.Ready[1] {
			r.TurnColor = utils.OppositeColor(r.TurnColor)
		}
		r.TurnNumber--
	}

}

// Note: http input should be coordinates on 9x9, rulesets should convert themselves

// TODO Figure out how to handle checks in drafting phase
func (r *Revised) ValidateMove(piece int, fromCoord int, moveType int, toCoord int) error {
	if moveType < 0 || moveType > 4 {
		return errors.New("invalid move type")
	}
	if !(r.Ready[0] && r.Ready[1]) && !(moveType == PLACE || moveType == READY) {
		return errors.New("can only place or ready in drafting phase")
	}
	if getColor(piece) != r.TurnColor && piece != -1 {
		return errors.New("wrong color")
	}
	if r.MarshalCoords[r.TurnColor] == 0 && (piece != coloredPiece(MARSHAL, r.TurnColor) || moveType != PLACE) {
		return errors.New("must place marshal in drafting phase")
	}

	isValid := false

	checkStatus, moveList := r.GetLegalMoves()

	switch moveType {
	case MOVE:
		if fromCoord < 0 || fromCoord > utils.BOARD_SQUARE_NUM || toCoord < 0 || toCoord > utils.BOARD_SQUARE_NUM {
			return errors.New("invalid coordinates")
		}
		if piece < 0 || piece > 25 {
			return errors.New("invalid piece")
		}
		fromSquare := r.BoardSquares[fromCoord]
		toSquare := r.BoardSquares[toCoord]
		if fromSquare.IsEmpty() {
			return errors.New("invalid move piece")
		} else if fromSquare.IsOutOfBounds() {
			return errors.New("code error, out of bounds")
		} else if fromSquare.GetTop() != piece {
			return errors.New("invalid move piece")
		}

		for _, move := range moveList[fromCoord] {
			if move == toCoord && toSquare.IsEmpty() {
				isValid = true
			}
		}
	case STACK:
		if fromCoord < 0 || fromCoord > utils.BOARD_SQUARE_NUM || toCoord < 0 || toCoord > utils.BOARD_SQUARE_NUM {
			return errors.New("invalid coordinates")
		}
		if piece < 0 || piece > 25 {
			return errors.New("invalid piece")
		}
		fromSquare := r.BoardSquares[fromCoord]
		toSquare := r.BoardSquares[toCoord]

		if toSquare.IsEmpty() {
			return errors.New("cannot stack on empty square")
		}

		if fromSquare.GetTop()%13 == FORTRESS {
			return errors.New("fortress cannot stack")
		}
		if toSquare.GetTop()%13 == MARSHAL {
			return errors.New("cannot stack on marshal")
		}

		if fromSquare.IsEmpty() {
			return errors.New("invalid move piece")
		} else if fromSquare.IsOutOfBounds() {
			return errors.New("code error, out of bounds")
		} else if fromSquare.GetTop() != piece {
			return errors.New("invalid move piece")
		}

		for _, move := range moveList[fromCoord] {
			if move == toCoord && toSquare.IsStackable() {
				isValid = true
			}
		}
	case ATTACK:
		if fromCoord < 0 || fromCoord > utils.BOARD_SQUARE_NUM || toCoord < 0 || toCoord > utils.BOARD_SQUARE_NUM {
			return errors.New("invalid coordinates")
		}
		if piece < 0 || piece > 25 {
			return errors.New("invalid piece")
		}
		fromSquare := r.BoardSquares[fromCoord]
		toSquare := r.BoardSquares[toCoord]
		if fromSquare.IsEmpty() {
			return errors.New("invalid move piece")
		} else if fromSquare.IsOutOfBounds() {
			return errors.New("code error, out of bounds")
		} else if fromSquare.GetTop() != piece {
			return errors.New("invalid move piece")
		}

		if fromSquare.GetTop()%13 == FORTRESS && len(toSquare.stack) != 1 {
			return errors.New("fortress cannot stack")
		}

		for _, move := range moveList[fromCoord] {
			if move == toCoord && !toSquare.IsEmpty() && getColor(toSquare.GetTop()) != r.TurnColor {
				isValid = true
			}
		}
	case PLACE:
		if toCoord < 0 || toCoord > utils.BOARD_SQUARE_NUM {
			return errors.New("invalid coordinates")
		}
		if piece < 0 || piece > 25 {
			return errors.New("invalid piece")
		}
		if r.Hand[piece] == 0 {
			return errors.New("piece not in hand")
		}

		var count [2]int
		count[0], count[1] = r.PieceCount()
		if count[r.TurnColor] >= 26 {
			return errors.New("max pieces already on board")
		}

		if !r.Ready[r.TurnColor] {
			validSquares := []int{}
			if r.TurnColor == 0 {
				validSquares = []int{109, 110, 111, 112, 113, 114, 115, 116, 117,
					121, 122, 123, 124, 125, 126, 127, 128, 129,
					133, 134, 135, 136, 137, 138, 139, 140, 141}
			} else if r.TurnColor == 1 {
				validSquares = []int{37, 38, 39, 40, 41, 42, 43, 44, 45,
					49, 50, 51, 52, 53, 54, 55, 56, 57,
					61, 62, 63, 64, 65, 66, 67, 68, 69}
			}

			validCoord := false
			for _, coord := range validSquares {
				if toCoord == coord {
					validCoord = true
					break
				}
			}

			if !validCoord {
				return errors.New("must place in your territory during drafting phase")
			}
		}

		if checkStatus == "double-checked" {
			return errors.New("cannot place while double checked")
		}

		if toCoord == r.MarshalCoords[0] || toCoord == r.MarshalCoords[1] {
			return errors.New("cannot place on marshal")
		}

		invalidSquares := []int{}
		if r.TurnColor == 0 {
			invalidSquares = []int{37, 38, 39, 40, 41, 42, 43, 44, 45,
				49, 50, 51, 52, 53, 54, 55, 56, 57,
				61, 62, 63, 64, 65, 66, 67, 68, 69}
		} else if r.TurnColor == 1 {
			invalidSquares = []int{109, 110, 111, 112, 113, 114, 115, 116, 117,
				121, 122, 123, 124, 125, 126, 127, 128, 129,
				133, 134, 135, 136, 137, 138, 139, 140, 141}
		}
		for _, coord := range invalidSquares {
			if toCoord == coord {
				return errors.New("cannot drop in enemy territory")
			}
		}

		if piece%13 == PAWN {
			toFile := toCoord % 12
			for i := 0; i < 9; i++ {
				toSquare := r.BoardSquares[(12*i)+36+toFile]
				if !toSquare.IsEmpty() {
					for _, p := range toSquare.stack {
						if p == piece {
							return errors.New("cannot drop another pawn in this file")
						}
					}
				}
			}
		}

		toSquare := r.BoardSquares[toCoord]
		if len(toSquare.stack) >= 3 {
			return errors.New("cannot drop on full stack")
		}
		if piece%13 == FORTRESS && len(toSquare.stack) != 0 {
			return errors.New("cannot stack with fortress")
		}

		if checkStatus == "checked" {
			for _, move := range moveList[-1] {
				if move == toCoord {
					isValid = true
				}
			}
		} else {
			isValid = true
		}

	case READY:
		if r.Ready[r.TurnColor] {
			log.Println(r.Ready)
			return errors.New("already readied")
		}
		isValid = true
	}

	if isValid {
		return nil
	} else {
		return errors.New("invalid move")
	}
}

// Coord is the coordinate of the piece to move
// MoveList is the list of destination
type Moves struct {
	Coord    int
	MoveList []int
}

func (r *Revised) GetLegalMoves() (string, map[int][]int) {
	if r.MarshalCoords[r.TurnColor] == 0 {
		return "", map[int][]int{}
	}

	marshalCoord := r.MarshalCoords[r.TurnColor]
	// log.Println("marshal coord", marshalCoord)
	marshalMoves := make(map[int]bool)
	for _, coord := range r.GetPseudoLegalMoves(coloredPiece(MARSHAL, r.TurnColor), marshalCoord, 1) {
		marshalMoves[coord] = true
	}

	attackCoord, xraySquares, xrayCoord, pinnedCoord, checkStatus := r.CheckEnemyMoves(&marshalMoves)
	// log.Println("check enemy moves", attackCoord, xraySquares, xrayCoord, pinnedCoord, checkStatus)

	xrayMap := make(map[int]bool)
	xrayBetweenMap := make(map[int]bool)
	if len(xraySquares) > 0 {
		for _, move := range xraySquares {
			xrayMap[move.Coordinate] = true
			if move.InBetween {
				xrayBetweenMap[move.Coordinate] = true
			}
		}
	}

	moveList := []Moves{}

	currSquare := r.ListRef.Head[r.TurnColor]
	for currSquare != nil {
		stack := currSquare.stack
		snared := false

		if len(stack) > 1 {
			pieceUnder := stack[len(stack)-2]
			if getColor(pieceUnder) != r.TurnColor {
				for _, move := range r.GetPseudoLegalMoves(pieceUnder, currSquare.coord, len(stack)-1) {
					if r.MarshalCoords[r.TurnColor] == move {
						snared = true
					}
					if currSquare.GetTop()%13 == MARSHAL && marshalMoves[move] {
						if pieceUnder%13 == FORTRESS && r.BoardSquares[move].IsEmpty() {
							delete(marshalMoves, move)
						} else if pieceUnder%13 != FORTRESS {
							delete(marshalMoves, move)
						}
					}
				}
			}
		}

		if currSquare.GetTop()%13 != MARSHAL && !snared {
			moves := []int{}
			for _, move := range r.GetPseudoLegalMoves(currSquare.GetTop(), currSquare.coord, len(stack)) {
				if checkStatus == "double-checked" {
					continue
				}
				if checkStatus == "checked" {
					if len(xraySquares) > 0 && (xrayBetweenMap[move] || move == xrayCoord) {
						// log.Println(1)
						moves = append(moves, move)
					} else if move == attackCoord {
						// log.Println(2)
						moves = append(moves, move)
					}
				} else if currSquare.coord == pinnedCoord && xrayBetweenMap[move] {
					moves = append(moves, move)
				} else if currSquare.coord != pinnedCoord && move != r.MarshalCoords[r.TurnColor] {
					moves = append(moves, move)
				}
			}
			moveList = append(moveList, Moves{Coord: currSquare.coord, MoveList: moves})
		}
		currSquare = currSquare.next
	}

	filteredMoves := make(map[int][]int)
	//TODO fix filtered moves when in check

	if checkStatus == "checked" {
		// filteredMoves[-1] means piece placement from hand
		if attackCoord != -1 {
			if len(r.BoardSquares[attackCoord].stack) >= 0 && len(r.BoardSquares[attackCoord].stack) < 3 {
				filteredMoves[-1] = append(filteredMoves[-1], attackCoord)
			}
		} else {
			wList := []int{}
			rCheck := false
			for _, move := range xraySquares {
				if move.InBetween {
					wList = append(wList, move.Coordinate)
				}
				if move.Attacked && move.Coordinate == r.MarshalCoords[r.TurnColor] {
					rCheck = true
				}
			}
			if rCheck {
				if len(r.BoardSquares[xrayCoord].stack) >= 0 && len(r.BoardSquares[xrayCoord].stack) < 3 {
					wList = append(wList, xrayCoord)
				}
				filteredMoves[-1] = wList
			}
		}
	}

	for _, moves := range moveList {
		if checkStatus == "checked" {
			for _, move := range moves.MoveList {
				if move == attackCoord {
					filteredMoves[moves.Coord] = append(filteredMoves[moves.Coord], attackCoord)
				}
			}
		} else {
			filteredMoves[moves.Coord] = append(filteredMoves[moves.Coord], moves.MoveList...)
		}
	}
	for move := range marshalMoves {
		filteredMoves[r.MarshalCoords[r.TurnColor]] = append(filteredMoves[r.MarshalCoords[r.TurnColor]], move)
	}
	for i, moves := range filteredMoves {
		if len(moves) == 0 {
			delete(filteredMoves, i)
		}
	}

	isEmpty := true
	if r.TurnColor == 0 {
		for i := 0; i < 13; i++ {
			if r.Hand[i] > 0 {
				isEmpty = false
			}
		}
	} else {
		for i := 13; i < 26; i++ {
			if r.Hand[i] > 0 {
				isEmpty = false
			}
		}
	}

	if checkStatus == "checked" || checkStatus == "double-checked" {
		if isEmpty && len(filteredMoves) == 1 && len(filteredMoves[-1]) != 0 {
			checkStatus = "checkmated"
		} else if len(filteredMoves) == 0 {
			checkStatus = "checkmated"
		}
	} else if len(filteredMoves) == 0 && isEmpty {
		checkStatus = "stalemate"
	}

	return checkStatus, filteredMoves
}

type xraySquares struct {
	Coordinate int
	InBetween  bool
	Attacked   bool
	Occupied   bool
}

func (r *Revised) CheckEnemyRanging(piece int, coord int) ([]xraySquares, bool, bool) {
	// enemyColor := utils.OppositeColor(r.TurnColor)
	check := false
	inPath := false
	squares := []xraySquares{}

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
		tempInPath := false // Marshal is in path
		tempBlocked := false
		tempSquares := []xraySquares{}

		i := coord + offset
		currSquare := r.BoardSquares[i]
		for {
			if currSquare.IsEmpty() {
				tempSquares = append(tempSquares, xraySquares{
					Coordinate: i,
					InBetween:  !tempInPath,
					Attacked:   !tempBlocked,
					Occupied:   false,
				})
			} else if currSquare.IsOutOfBounds() {
				break
			} else {
				piece := currSquare.GetTop()
				attacked := !tempBlocked
				if piece%13 == MARSHAL && getColor(piece) == r.TurnColor {
					tempInPath = true
				} else if !tempInPath {
					tempBlocked = true
				}
				tempSquares = append(tempSquares, xraySquares{
					Coordinate: i,
					InBetween:  !tempInPath,
					Attacked:   attacked,
					Occupied:   true,
				})
			}
			i += offset
			currSquare = r.BoardSquares[i]
		}

		if !tempBlocked && tempInPath {
			check = true
		}
		if tempInPath {
			squares = tempSquares
			inPath = true
		}
	}
	return squares, check, inPath
}

func (r *Revised) CheckEnemyMoves(marshalMoves *map[int]bool) (int, []xraySquares, int, int, string) {
	// log.Println("marshal moves: ", marshalMoves)
	// log.Println("enemy moves: ")
	enemyColor := utils.OppositeColor(r.TurnColor)
	currSquare := r.ListRef.Head[enemyColor]
	var pinnedCoord int = -1
	var xrayCoord int = -1
	var xraySquares []xraySquares
	var checkStatus string
	var attackCoord int = -1

	for currSquare != nil {
		piece := currSquare.GetTop()
		tier := len(currSquare.stack)

		if piece%13 == TACTICIAN {
			for i := len(currSquare.stack) - 2; i >= 0; i-- {
				piece = currSquare.stack[i]
				if piece%13 != TACTICIAN {
					break
				}
			}
			piece = coloredPiece(piece%13, enemyColor)
		}

		if tier == 3 && piece%13 >= 7 && piece%13 <= 10 {
			moves, check, inPath := r.CheckEnemyRanging(piece, currSquare.coord)
			if check {
				if checkStatus == "checked" {
					checkStatus = "double-checked"
				} else {
					checkStatus = "checked"
				}
			}
			if inPath {
				piecesBetween := []int{}

				for _, move := range moves {
					if move.InBetween && move.Occupied {
						piecesBetween = append(piecesBetween, move.Coordinate)
					}
				}

				if len(piecesBetween) <= 1 {
					xraySquares = moves
					xrayCoord = currSquare.coord
				}

				if len(piecesBetween) == 1 && getColor(r.BoardSquares[piecesBetween[0]].GetTop()) == r.TurnColor {
					pinnedCoord = piecesBetween[0]
				}

				for _, move := range moves {
					if (*marshalMoves)[move.Coordinate] && move.Attacked {
						delete(*marshalMoves, move.Coordinate)
					}
				}
			}
		} else {
			moves := r.GetPseudoLegalMoves(piece, currSquare.coord, len(currSquare.stack))
			// log.Println(moves)
			for _, moveCoord := range moves {
				if (*marshalMoves)[moveCoord] && !(piece%13 == FORTRESS && !r.BoardSquares[moveCoord].IsEmpty()) {
					delete(*marshalMoves, moveCoord)
				}

				if r.MarshalCoords[r.TurnColor] == moveCoord {
					if piece%13 == FORTRESS && len(r.BoardSquares[r.MarshalCoords[r.TurnColor]].stack) != 1 {
						continue
					}

					attackCoord = currSquare.coord

					if checkStatus == "checked" {
						checkStatus = "double-checked"
					} else {
						checkStatus = "checked"
					}
				}
			}
		}
		currSquare = currSquare.next
	}
	return attackCoord, xraySquares, xrayCoord, pinnedCoord, checkStatus
}

// Will stop if out of bounds or a piece is in the way. However, you should double check if attacking or stacking is possible including stacking on marshal.
func (r *Revised) GetPseudoRanging(coord int, offset int) []int {
	squares := []int{}

	i := coord + offset
	currSquare := r.BoardSquares[i]
	for currSquare.IsEmpty() {
		squares = append(squares, i)
		i += offset
		currSquare = r.BoardSquares[i]
	}

	if !currSquare.IsOutOfBounds() {
		squares = append(squares, i)
	}

	return squares
}

// Generates psuedo-legal moves for a piece at a coordinate.
// Still must determine out of bound squares, fully stacked pieces, marshal, fortress, is checked, and pinned.
// Does not differentiate from attacking and stacking.
// Does not generate moves from hand.
func (r *Revised) GetPseudoLegalMoves(piece int, coord int, tier int) []int {
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
			offsets = append(offsets, -26, -25, -24, -23, -22, -10, -14, -2, 26, 25, 24, 23, 22, 10, 14, 2)
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
			if getColor(piece) == 0 {
				squares = append(squares, r.GetPseudoRanging(coord, -12)...)
			} else {
				squares = append(squares, r.GetPseudoRanging(coord, 12)...)
			}
		}
	case SAMURAI:
		switch tier {
		case 1:
			offsets = append(offsets, -11, -13, 11, 13)
		case 2:
			offsets = append(offsets, -26, -22, 26, 22)
		case 3:
			squares = append(squares, r.GetPseudoRanging(coord, 13)...)
			squares = append(squares, r.GetPseudoRanging(coord, -13)...)
			squares = append(squares, r.GetPseudoRanging(coord, -11)...)
			squares = append(squares, r.GetPseudoRanging(coord, 11)...)
		}
	case CANNON:
		switch tier {
		case 1:
			offsets = append(offsets, -12, -1, 1, 12)
		case 2:
			offsets = append(offsets, -24, -12, -2, -1, 1, 2, 12, 24)
		case 3:
			squares = append(squares, r.GetPseudoRanging(coord, 1)...)
			squares = append(squares, r.GetPseudoRanging(coord, -1)...)
			squares = append(squares, r.GetPseudoRanging(coord, -12)...)
			squares = append(squares, r.GetPseudoRanging(coord, 12)...)
		}
	case SPY:
		switch tier {
		case 1:
			offsets = append(offsets, -12)
		case 2:
			offsets = append(offsets, -11, -13, 11, 13)
		case 3:
			squares = append(squares, r.GetPseudoRanging(coord, 13)...)
			squares = append(squares, r.GetPseudoRanging(coord, -13)...)
			squares = append(squares, r.GetPseudoRanging(coord, -11)...)
			squares = append(squares, r.GetPseudoRanging(coord, 11)...)
			squares = append(squares, r.GetPseudoRanging(coord, 1)...)
			squares = append(squares, r.GetPseudoRanging(coord, -1)...)
			squares = append(squares, r.GetPseudoRanging(coord, -12)...)
			squares = append(squares, r.GetPseudoRanging(coord, 12)...)
		}
	case TACTICIAN:
		// Movement becomes the piece under it but on tactician's tier
		switch tier {
		case 1:
			offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
		case 2, 3:
			square := r.BoardSquares[coord]
			var piece int
			for i := len(square.stack) - 2; i >= 0; i-- {
				piece = square.stack[i]
				if piece%13 != TACTICIAN {
					break
				}
			}
			sameColorPiece := coloredPiece(piece%13, r.TurnColor)
			squares = append(squares, r.GetPseudoLegalMoves(sameColorPiece, coord, tier)...)
		}
	case MARSHAL:
		offsets = append(offsets, -11, -12, -13, -1, 1, 11, 12, 13)
	}

	for _, offset := range offsets {
		var offsetCoord int

		if getColor(piece) == BLACK {
			offsetCoord = coord - offset
		} else {
			offsetCoord = coord + offset
		}

		square := r.BoardSquares[offsetCoord]
		if square.IsEmpty() || !square.IsOutOfBounds() {
			squares = append(squares, offsetCoord)
		}
	}

	return squares
}
