import { useMemo } from 'react';
import { AnimatePresence } from 'framer-motion';

import sword from 'assets/imgs/icons/sword.svg';

import {
  GameWrapper,
  Marker,
  Sword
} from 'pages/Jackpot/components/JackpotRoom/styles';
import { JackpotRoundData } from 'api/types/jackpot';
import { Box, Roll2D, Roll3D } from 'components';
import { useAppSelector } from 'state';

import { Content } from '../Content';
import CountDown from '../CountDown/CountDown';
import { convertBalanceToChip } from 'utils/balance';

export interface JackpotRoomProps {
  roundData: JackpotRoundData;
  handleShowNFT?: any;
  isHistory?: boolean;
}

const JackpotRoom = ({
  roundData: game,
  isHistory,
  handleShowNFT
}: JackpotRoomProps) => {
  const animation = useAppSelector(state => state.user.grandJackpotAnimation);
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
        background={
          win
            ? 'linear-gradient(0deg, #ffe24b30 0%, #ffe24b00 100%)'
            : 'linear-gradient(180deg, #05090d 0%, #0b141e 100%)'
        }
      >
        {/* <ToggleAnimation onClick={() => {
          state.dispatch(toggleGrandJackpotAnimation());
        }}>
          {animation}
        </ToggleAnimation> */}
        {roll ? (
          <>
            <Box overflow="hidden" maxWidth="100%" height="300px" mt="50px">
              <AnimatePresence exitBeforeEnter>
                {animation === '2D' ? (
                  <Roll2D roundData={game} isGrand key="roll2d_grand" />
                ) : (
                  <Roll3D
                    roundData={game}
                    isHistory={isHistory}
                    isGrand
                    key="roll3d_grand"
                  />
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
            nfts={game.nfts ?? []}
            nftsToShow={4}
            winnerName={winnerName}
            usdBetAmount={convertBalanceToChip(game.usdBetAmount)}
            nftBetAmount={convertBalanceToChip(game.nftBetAmount)}
            handleShowNFT={handleShowNFT}
          />
        )}
      </GameWrapper>
      <CountDown status={game.status} time={game.time} isHistory={isHistory} />
    </>
  );
};

export default JackpotRoom;
