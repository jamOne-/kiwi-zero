package gomoku

import "github.com/jamOne-/kiwi-zero/game"

func ConvertGomokuFnToGeneralFeatuersFn(gomokuFn func(gomoku *GomokuGame) game.Features) game.GameToFeaturesFn {
	return func(g game.Game) game.Features {
		gomokuGame := g.(*GomokuGame)
		return gomokuFn(gomokuGame)
	}
}

func GomokuToBoard1(gomoku *GomokuGame) game.Features {
	features := make([][][]float32, TOTAL_SIZE)
	for row := 0; row < TOTAL_SIZE; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	for i, field := range gomoku.Board {
		features[i][0][0] = float32(field)
	}

	return features
}
