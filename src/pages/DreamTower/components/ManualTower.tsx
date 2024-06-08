import { useEffect, useCallback, useState, useMemo } from 'react';
import styled from 'styled-components';
import ClipLoader from 'react-spinners/ClipLoader';
import { toast } from 'react-toastify';
import gsap from 'gsap';

import { ReactComponent as CoinBlueIcon } from 'assets/imgs/coins/coin-blue.svg';

import { Box, Button, CoinIcon, Flex } from 'components';
import state, { useAppDispatch, useAppSelector } from 'state';
import { updateBalance } from 'state/user/actions';
import { raise, reset, updateStatus } from 'state/dreamtower/actions';
import { useSound, useCustomSWR } from 'hooks';

import { api } from 'services';
import Row from './Row';
import { towerResetAnimation, towerWinAnimation } from '../animation';
import Payout from './Payout';
import { convertBalanceToChip } from 'utils/balance';

interface Props {
  dreamRef: React.MutableRefObject<any>;
}

export default function ManualTower({ dreamRef }: Props) {
  const dispatch = useAppDispatch();
  const user = useAppSelector(state => state.user);
  const game = useAppSelector(state => state.dreamtower.game);
  const meta = useAppSelector(state => state.meta.dreamtower);
  const [request, setRequest] = useState(false);

  const { data: maxWinData } = useCustomSWR({
    key: 'dreamtower_max_win',
    route: '/dreamtower/max-win',
    method: 'get'
  });

  const maxWinning = useMemo(() => {
    if (!maxWinData) return 0;
    return maxWinData;
  }, [maxWinData]);

  const { towerPlay, towerStop } = useSound();

  const handleBet = useCallback(async () => {
    if (user.name === '') {
      toast.info('Please sign in.');
      return;
    }

    if (game.betAmount > user.balance) {
      toast.error('Insufficient funds');
      return;
    }

    if (
      (game.betAmount < meta.minAmount && game.betAmount > 0) ||
      game.betAmount < 0
    ) {
      toast.warning(
        `Minimum bet amount is ${Math.floor(
          convertBalanceToChip(meta.minAmount)
        )} chip.`
      );
      return;
    } else if (game.betAmount > meta.maxAmount) {
      toast.warning(
        `Maximum bet amount is ${Math.floor(
          convertBalanceToChip(meta.maxAmount)
        )} chips.`
      );
      return;
    }

    dispatch(
      updateBalance({
        type: -1,
        balanceType: user.betBalanceType,
        wagered: game.betAmount,
        usdAmount: game.betAmount
      })
    );

    const betData = {
      betAmount: game.betAmount,
      bets: [],
      difficulty: game.difficulty.level,
      paidBalanceType: user.betBalanceType
    };

    setRequest(true);
    try {
      let res = await api.post('/dreamtower/bet', betData);
      state.dispatch(reset());
      state.dispatch(updateStatus(res.data));
    } catch (error: any) {
      dispatch(
        updateBalance({
          type: 1,
          balanceType: betData.paidBalanceType,
          usdAmount: betData.betAmount,
          wagered: betData.betAmount
        })
      );
      if (error.response.status === 503) {
        toast.error('This function is blocked by admin.');
      } else toast.error(error.response.data.message);
    }
    setRequest(false);
  }, [
    user.name,
    user.balance,
    user.betBalanceType,
    game.betAmount,
    game.difficulty.level,
    meta.minAmount,
    meta.maxAmount,
    dispatch
  ]);

  const handleCashout = useCallback(async () => {
    if (user.name === '') {
      toast.info('Please sign in.');
      return;
    }
    if (game.status !== 'playing') {
      // console.log(game.status);
      toast.error('Can only cashout from currently playing round.');
      return;
    }
    setRequest(true);
    try {
      let res = await api.post('/dreamtower/cashout', {
        roundId: game.roundId
      });

      const profit =
        maxWinning > res.data.multiplier * game.betAmount
          ? Math.floor(res.data.multiplier * game.betAmount)
          : Math.floor(maxWinning);

      dispatch(
        updateBalance({
          type: 1,
          usdAmount: profit,
          balanceType: res.data.paidBalanceType
        })
      );
      towerPlay.lastStar();
      state.dispatch(updateStatus(res.data));
    } catch (error: any) {
      if (error.response.status === 503) {
        toast.error('This function is blocked by admin.');
      } else toast.error(error.response.data.message);
    }
    setRequest(false);
  }, [
    user.name,
    game.status,
    game.roundId,
    game.betAmount,
    towerPlay,
    dispatch
  ]);

  const handleClickSquare = useCallback(
    async (index: any) => {
      setRequest(true);
      try {
        const res = await api.post('/dreamtower/raise', {
          roundId: game.roundId,
          bet: index,
          height: game.bets.length
        });
        dispatch(updateStatus(res.data));
        dispatch(raise({ status: res.data.status, bet: index }));
        if (res.data.status === 'win') {
          const profit =
            maxWinning > res.data.multiplier * game.betAmount
              ? Math.floor(res.data.multiplier * game.betAmount)
              : Math.floor(maxWinning);

          state.dispatch(
            updateBalance({
              type: 1,
              usdAmount: profit,
              balanceType: res.data.paidBalanceType
            })
          );
        }
      } catch (error: any) {
        if (error.response.status === 503) {
          toast.error('This function is blocked by admin.');
        } else toast.error(error.response.data.message);
        return;
      }
      setRequest(false);
    },
    [dispatch, game.betAmount, game.roundId, game.bets]
  );

  useEffect(() => {
    if (!dreamRef) return;
    if (game.status === 'win') {
      const q = gsap.utils.selector(dreamRef);
      const tl = towerWinAnimation({
        textTargets: [q('.dreamtower_dreamtext'), q('.dreamtower_towertext')],
        rowTargets: [...q('.dreamtower_row')],
        duration: 0.15
      });
      let timeout: NodeJS.Timeout;
      if (towerPlay.lastStar) {
        towerPlay.lastStar();
        timeout = setTimeout(() => {
          towerPlay.win && towerPlay.win();
        }, 3700);
      }
      return () => {
        tl.kill();
        towerStop.lastStar();
        towerStop.win();
        clearTimeout(timeout);
      };
    } else {
      if (game.bets.length > 0) {
        if (game.status === 'loss') {
          towerPlay.break && towerPlay.break();
        } else if (
          game.status === 'playing' &&
          game.bets.length < game.height
        ) {
          towerPlay.selectStar && towerPlay.selectStar();
        }
      }

      const q = gsap.utils.selector(dreamRef);

      const rows = [...q('.dreamtower_row')];
      var tr1: typeof rows = [],
        tr2: typeof rows = [],
        tr3: typeof rows = [];
      rows.forEach((row, index) => {
        if (
          (game.status === 'playing' && index === game.bets.length) ||
          game.status === 'win' ||
          game.status === 'cashout'
        ) {
          tr2.push(row);
        } else if (game.status === 'loss' && index === game.bets.length - 1) {
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
    }
  }, [game.status, game.bets, towerPlay, towerStop, game.height, dreamRef]);

  return (
    <>
      {(game.status === 'win' || game.status === 'cashout') && (
        <Payout
          multiplier={game.multiplier}
          profit={Math.floor(game.profit)}
          chipType={user.betBalanceType}
        />
      )}
      <Flex py="5px" flexDirection={'column-reverse'}>
        {game.tower.map((v, i) => (
          <Row
            key={i}
            value={v}
            isNext={game.status === 'playing' && i === game.bets.length}
            isHighlight={
              (game.status === 'loss' && i === game.bets.length - 1) ||
              (game.status === 'playing' && i === game.bets.length) ||
              game.status === 'win' ||
              game.status === 'cashout'
            }
            isClickable={game.status === 'playing' && i === game.bets.length}
            isUnderBroken={i < game.bets.length}
            selectedIndex={game.bets.length > i ? game.bets[i] : undefined}
            nextMultiplier={game.nextMultiplier!}
            blocksInRow={game.difficulty.blocksInRow}
            handleClickSquare={request ? () => {} : handleClickSquare}
          />
        ))}
      </Flex>
      <Box px={'20px'}>
        {request ? (
          <StyledCashButton>
            <ClipLoader color="#ffffff" loading={request} size={20} />
          </StyledCashButton>
        ) : game.status === 'playing' ? (
          <>
            <StyledCashButton
              onClick={request ? undefined : handleCashout}
              disabled={game.bets.length === 0}
            >
              Cashout
              {game.bets.length > 0 && (
                <>
                  {user.betBalanceType === 'coupon' ? (
                    <CoinBlueIcon />
                  ) : (
                    <CoinIcon />
                  )}
                  {maxWinning > game.multiplier * game.betAmount
                    ? convertBalanceToChip(
                        Math.floor(game.multiplier * game.betAmount)
                      ).toFixed(2)
                    : convertBalanceToChip(Math.floor(maxWinning)).toFixed(2)}
                </>
              )}
            </StyledCashButton>
          </>
        ) : (
          <StyledCashButton onClick={request ? undefined : handleBet}>
            Place Bet
          </StyledCashButton>
        )}
      </Box>
    </>
  );
}

const StyledCashButton = styled(Button)`
  width: 100%;
  margin-top: 5px;
  padding: 18px 0px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 700;
  color: black;
`;
