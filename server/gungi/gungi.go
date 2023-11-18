package gungi

import "github.com/TonyLeCode/gungi.go/server/gungi/revised"

type Ruleset interface {
	InitializeBoard()
	PrintBoard()
	SetBoardFromFen(fen string) error
	ValidateMove(piece int, fromCoord int, moveType int, toCoord int) error
	GetLegalMoves() (string, map[int][]int)
	MakeMove(piece int, fromCoord int, moveType int, toCoord int) error
	//Ready(playerColor string)
	//Resign(playerColor string)
	UndoMove()
	BoardToFen() string
	SetHistory(history []string)
	SerializeHistory() string
	ConvertInputCoord(coord int) int
	ConvertOutputCoord(coord int) int
}

type Game struct {
	Ruleset Ruleset
}

func CreateBoard(rules string) Game {
	newBoard := Game{}
	switch rules {
	case "revised":
		newBoard.Ruleset = &revised.Revised{}
	default:
		newBoard.Ruleset = &revised.Revised{}
	}
	newBoard.InitializeBoard()
	return newBoard
}

func NewBoard(r Ruleset) Game {
	newBoard := Game{Ruleset: r}
	newBoard.InitializeBoard()
	return newBoard
}

func (g *Game) InitializeBoard() {
	g.Ruleset.InitializeBoard()
}

func (g *Game) PrintBoard() {
	g.Ruleset.PrintBoard()
}

func (g *Game) SetBoardFromFen(fen string) error {
	return g.Ruleset.SetBoardFromFen(fen)
}

func (g *Game) ValidateMove(piece int, fromCoord int, moveType int, toCoord int) error {
	return g.Ruleset.ValidateMove(piece, fromCoord, moveType, toCoord)
}

func (g *Game) GetLegalMoves() (string, map[int][]int) {
	return g.Ruleset.GetLegalMoves()
}

func (g *Game) MakeMove(piece int, fromCoord int, moveType int, toCoord int) error {
	return g.Ruleset.MakeMove(piece, fromCoord, moveType, toCoord)
}

func (g *Game) UndoMove() {
	g.Ruleset.UndoMove()
}

func (g *Game) BoardToFen() string {
	return g.Ruleset.BoardToFen()
}

func (g *Game) SetHistory(history []string) {
	g.Ruleset.SetHistory(history)
}

func (g *Game) SerializeHistory() string {
	return g.Ruleset.SerializeHistory()
}

func (g *Game) ConvertInputCoord(coord int) int {
	return g.Ruleset.ConvertInputCoord(coord)
}

func (g *Game) ConvertOutputCoord(coord int) int {
	return g.Ruleset.ConvertOutputCoord(coord)
}
