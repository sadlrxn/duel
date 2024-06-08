import { useEffect, useMemo, FC } from 'react';
import { Box, Button, Flex } from 'components';
import styled from 'styled-components';
import Row from './Row';
import { useAppSelector } from 'state';
import gsap from 'gsap';
import { towerResetAnimation } from '../animation';
import Payout from './Payout';

const HistoryTower: FC<{
  goBackHandler: () => void;
  dreamRef: React.MutableRefObject<any>;
}> = ({ goBackHandler, dreamRef }) => {
  const review = useAppSelector(state => state.dreamtower.history.review);

  const endedHeight = useMemo(() => {
    if (review.status === '') return -1;
    for (var i = 0; i < review.tower.length; i++) {
      if (review.tower[i][review.bets[i]] === 0) return i;
    }
    return review.tower.length - 1;
  }, [review.tower, review.bets, review.status]);

  useEffect(() => {
    if (!dreamRef) return;
    // if (review.status === "win") {
    //   const q = gsap.utils.selector(dreamRef);
    //   const tl = towerWinAnimation({
    //     textTargets: [q(".dreamtower_dreamtext"), q(".dreamtower_towertext")],
    //     rowTargets: [...q(".dreamtower_row")],
    //     duration: 0.15,
    //   });
    //   return () => {
    //     tl.kill();
    //   };
    // } else {
    const q = gsap.utils.selector(dreamRef);

    const rows = [...q('.dreamtower_row')];
    var tr1: typeof rows = [],
      tr2: typeof rows = [],
      tr3: typeof rows = [];
    rows.forEach((row, index) => {
      if (review.status === 'win' || review.status === 'cashout') {
        tr2.push(row);
      } else if (review.status === 'loss' && index === endedHeight) {
        tr3.push(row);
      } else {
        tr1.push(row);
      }
    });
    const tl = towerResetAnimation({
      textTargets: [q('.dreamtower_dreamtext'), q('.dreamtower_towertext')],
      rowTargets1: tr1,
      rowTargets2: tr2,
      rowTargets3: tr3,
      rowTargets4: []
    });
    return () => {
      tl.kill();
    };
    // }
  }, [review.status, review.bets, endedHeight, dreamRef]);

  return (
    <>
      {(review.status === 'win' || review.status === 'cashout') && (
        <Payout
          multiplier={review.multiplier}
          profit={review.profit}
          chipType={review.paidBalanceType}
        />
      )}
      <Flex py="5px" flexDirection={'column-reverse'}>
        {review.tower.map((v, i) => (
          <Row
            key={i}
            value={v}
            roundId={review.roundId}
            isNext={false}
            isHighlight={
              (review.status === 'loss' && i === review.bets.length - 1) ||
              (review.status === 'playing' && i === review.bets.length) ||
              review.status === 'win' ||
              review.status === 'cashout'
            }
            isClickable={false}
            isUnderBroken={i <= endedHeight}
            selectedIndex={review.bets.length > i ? review.bets[i] : undefined}
            nextMultiplier={0}
            blocksInRow={review.difficulty.blocksInRow}
            handleClickSquare={() => {}}
          />
        ))}
      </Flex>
      <Box px={'20px'}>
        <StyledCashButton onClick={goBackHandler}>
          Back to Current Game
        </StyledCashButton>
      </Box>
    </>
  );
};

const StyledCashButton = styled(Button)`
  width: 100%;
  margin-top: 5px;
  padding: 18px 0px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 700;
  color: black;
`;

export default HistoryTower;
