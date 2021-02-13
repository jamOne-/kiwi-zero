package games

import (
	"github.com/jamOne-/kiwi-zero/connectFour"
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/gomoku"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
)

type FeaturesFnInfo struct {
	Shape string
	Fn    game.GameToFeaturesFn
}

var FEATURES_FNS_DICT = map[string]map[string]*FeaturesFnInfo{
	"reversi":        REVERSI_FEATURES_FNS,
	"reversirandom":  REVERSI_FEATURES_FNS,
	"connect4":       CONNECT4_FEATURES_FNS,
	"connect4random": CONNECT4_FEATURES_FNS,
	"gomoku":         GOMOKU_FEATURES_FNS,
	"gomokurandom":   GOMOKU_FEATURES_FNS,
}

var REVERSI_FEATURES_FNS = map[string]*FeaturesFnInfo{
	"board3":         &FeaturesFnInfo{"(8,8,3)", OneHotBoard3},
	"boardmoves":     &FeaturesFnInfo{"(8,8,4)", reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoardMoves)},
	"paddedmoves":    &FeaturesFnInfo{"(10,10,5)", reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoardPaddedMoves)},
	"board1features": &FeaturesFnInfo{"(72,1,1)", reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToFeaturesExtended)},
	"board1":         &FeaturesFnInfo{"(64,1,1)", reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToFeatures)},
}

var CONNECT4_FEATURES_FNS = map[string]*FeaturesFnInfo{
	"board3": &FeaturesFnInfo{"(6,7,3)", OneHotBoard3},
	"board1": &FeaturesFnInfo{"(42,1,1)", connectFour.ConvertConnect4FnToGeneralFeatuersFn(connectFour.Connect4ToBoard1)},
	"b3turn": &FeaturesFnInfo{"(6,7,5)", connectFour.ConvertConnect4FnToGeneralFeatuersFn(connectFour.Connect4ToBoardTurn)},
}

var GOMOKU_FEATURES_FNS = map[string]*FeaturesFnInfo{
	"board3": &FeaturesFnInfo{"(8,8,3)", OneHotBoard3},
	"board1": &FeaturesFnInfo{"(64,1,1)", gomoku.ConvertGomokuFnToGeneralFeatuersFn(gomoku.GomokuToBoard1)},
}

func OneHotBoard3(game game.Game) game.Features {
	return game.OneHotBoard()
}
