package weekly_raffle

import (
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/sirupsen/logrus"
)

/**
* @External
* Initializes weekly raffle module.
*  - Initializes current round.
*  - If initializing current round is succeed, initializes socket.
 */
func Initialize(eventEmitter chan types.WSEvent) {
	go func() {
		if err := initRound(); err == nil {
			initSocket(eventEmitter)
		} else {
			log.LogMessage(
				"weekly_raffle_initializes",
				"failed to initialize weekly raffle round",
				"error",
				logrus.Fields{
					"error": err.Error(),
				},
			)
		}
	}()
}
