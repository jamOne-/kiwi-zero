package connectFour

import "github.com/jamOne-/kiwi-zero/game"

func ConvertConnect4FnToGeneralFeatuersFn(connect4Fn func(c4game *ConnectFourGame) game.Features) game.GameToFeaturesFn {
	return func(g game.Game) game.Features {
		connect4game := g.(*ConnectFourGame)
		return connect4Fn(connect4game)
	}
}

func Connect4ToBoard1(c4game *ConnectFourGame) game.Features {
	features := make([][][]float32, TOTAL_SIZE)
	for row := 0; row < TOTAL_SIZE; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	for i, field := range c4game.Board {
		features[i][0][0] = float32(field)
	}

	return features
}
