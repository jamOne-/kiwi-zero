package edaxPlayer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/jamOne-/kiwi-zero/reversi"

	"github.com/jamOne-/kiwi-zero/game"
)

type EdaxPlayer struct {
	alpha   float64
	beta    float64
	depth   int
	probcut int

	edaxIn  io.WriteCloser
	edaxOut *bufio.Reader
}

func NewEdaxPlayer(
	alpha float64,
	beta float64,
	depth int,
	probcut int,
) *EdaxPlayer {
	edaxCmd := exec.Command("/home/dominik/edax-linux/edax-4.4", "--cassio")
	edaxIn, _ := edaxCmd.StdinPipe()
	edaxOut, _ := edaxCmd.StdoutPipe()
	player := &EdaxPlayer{alpha, beta, depth, probcut, edaxIn, bufio.NewReader(edaxOut)}

	err := edaxCmd.Start()
	if err != nil {
		log.Fatalf("Edax start: %v", err)
	}

	player.writeCommand("ENGINE-PROTOCOL init")
	player.readUntilReady()

	return player
}

func (player *EdaxPlayer) SelectMove(g game.Game) game.Move {
	// Edax can't handle positions without move possibilities
	possibilities := g.GetPossibleMoves()
	if len(possibilities) == 1 {
		return possibilities[0]
	}

	reversiGame := g.(*reversi.ReversiGame)
	algorithm := "midgame-search"

	board := serializeReversi(reversiGame)
	command := fmt.Sprintf(
		"ENGINE-PROTOCOL %s %s %f %f %d %d",
		algorithm,
		board,
		player.alpha,
		player.beta,
		player.depth,
		player.probcut,
	)

	player.writeCommand(command)

	output := player.readUntilReady()
	re := regexp.MustCompile(`move (..)`)
	match := re.FindStringSubmatch(output)

	if match == nil {
		fmt.Println(fmt.Sprintf("Edax didn't return valid line! %s\n", output))
		return reversi.PASS_MOVE
	}

	edaxMove := strings.ToUpper(match[1])
	if edaxMove == "PA" {
		return reversi.PASS_MOVE
	}

	column := edaxMove[0] - 'A'
	row := edaxMove[1] - '1'
	move := column + row*8

	return game.Move(move)
}

func (player *EdaxPlayer) ClosePlayer() {
	player.writeCommand("ENGINE-PROTOCOL quit")
}

func (player *EdaxPlayer) writeCommand(command string) {
	_, err := player.edaxIn.Write([]byte(command + "\n"))

	if err != nil {
		log.Fatalf("Write: %v", err)
	}
}

func (player *EdaxPlayer) readUntilReady() string {
	output := ""
	line, err := "", error(nil)

	for line != "ready.\n" {
		line, err = player.edaxOut.ReadString('\n')

		if err != nil {
			log.Fatalf("Edax read: %v", err)
		}

		output += line
	}

	return output
}

func serializeReversi(reversi *reversi.ReversiGame) string {
	serialized := ""

	for _, field := range reversi.Board {
		serialized += fieldToEdaxField(field)
	}

	serialized += fieldToEdaxField(reversi.Turn)

	return serialized
}

func fieldToEdaxField(field game.Field) string {
	if field == game.BLACK {
		return "X"
	} else if field == game.WHITE {
		return "O"
	} else {
		return "-"
	}
}
