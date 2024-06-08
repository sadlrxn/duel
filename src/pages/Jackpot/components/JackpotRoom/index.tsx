import { useState, useCallback, useMemo } from 'react';
import { AnimatePresence } from 'framer-motion';

import sword from 'assets/imgs/icons/sword.svg';

import { Content } from '../Content';
import CountDown from '../CountDown';

import { GameWrapper, Marker, Sword } from './styles';
import { JackpotRoundData } from 'api/types/jackpot';
import { useAppSelector } from 'state';
import { Roll2D, Roll3D, Box } from 'components';
import { convertBalanceToChip } from 'utils/balance';

export interface JackpotRoomProps {
  roundData: JackpotRoundData;
  isHistory?: boolean;
  betCount?: number;
  time?: number;
}

const JackpotRoom = ({
  roundData: game,
  isHistory = false,
  betCount = 5,
  time = Date.now()
}: JackpotRoomProps) => {
  const animation = useAppSelector(state => state.user.jackpotAnimation);
  const [show, setShow] = useState(false);

  const handleShow = useCallback(() => setShow(true), []);
  const handleClose = useCallback(() => setShow(false), []);

  const [roll, win, winnerName] = useMemo(() => {
    return [
      game.status === 'rolling' && game.candidates.length > 0,
      game.status === 'rollend' || isHistory,
      game.winner.name
    ];
  }, [game, isHistory]);

  return (
    <>
      <GameWrapper
        position="relative"
        background={
          win
            ? 'linear-gradient(0deg, #ffe24b30 0%, #ffe24b00 100%)'
            : 'linear-gradient(180deg, #05090d 0%, #0b141e 100%)'
        }
      >
        {/* <ToggleAnimation onClick={() => {
          state.dispatch(toggleJackpotAnimation());
        }}>
          {animation}
        </ToggleAnimation> */}
        {show && (
          <Content
            variant="secondary"
            win={win}
            nfts={game.nfts}
            nftsToShow={4}
            winnerName={winnerName}
            usdBetAmount={convertBalanceToChip(game.usdBetAmount)}
            nftBetAmount={convertBalanceToChip(game.nftBetAmount)}
            onClose={handleClose}
            usdProfit={game.usdProfit}
            usdFee={game.usdFee}
            nftProfit={game.nftProfit}
            nftFee={game.nftFee}
            roundId={game.roundId}
            time={time}
          />
        )}
        {roll ? (
          <>
            <Box overflow="hidden" maxWidth="100%" height="300px" mt="50px">
              <AnimatePresence exitBeforeEnter>
                {animation === '2D' ? (
                  <Roll2D roundData={game} key="roll2d" />
                ) : (
                  <Roll3D roundData={game} isHistory key="roll3d" />
                )}
              </AnimatePresence>
            </Box>
            <Marker />
            <Sword mt={animation === '2D' ? '50px' : '80px'}>
              <img src={sword} alt="" width={21} height={42} />
            </Sword>
          </>
        ) : (
          <Content
            win={win}
            nfts={game.nfts}
            nftsToShow={4}
            winnerName={winnerName}
            usdBetAmount={convertBalanceToChip(game.usdBetAmount)}
            nftBetAmount={convertBalanceToChip(game.nftBetAmount)}
            onClick={handleShow}
            usdProfit={game.usdProfit}
            usdFee={game.usdFee}
            nftProfit={game.nftProfit}
            nftFee={game.nftFee}
            roundId={game.roundId}
            time={time}
          />
        )}
      </GameWrapper>
      <CountDown
        status={game.status}
        time={game.time}
        countingTime={game.countingTime ?? 40}
        betCount={isHistory ? undefined : betCount}
      />
    </>
  );
};

export default JackpotRoom;
