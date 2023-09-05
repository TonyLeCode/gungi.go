package ruleset

import (
	"github.com/TonyLeCode/gungi.go/server/gungi"
	"github.com/TonyLeCode/gungi.go/server/gungi/revised"
)

func NewBoard(r string) gungi.Game {
	newBoard := gungi.Game{}
	if r == "revised" {
		newBoard.Ruleset = &revised.Revised{}
	}
	newBoard.InitializeBoard()
	return newBoard
}
