import { useRef, useEffect, useCallback, useMemo, useState } from 'react';
import styled from 'styled-components';
import { toast } from 'react-toastify';
import { ClipLoader } from 'react-spinners';
import gsap from 'gsap';

import { Box, Button, Flex } from 'components';
import { useAppDispatch, useAppSelector } from 'state';
import { updateBalance } from 'state/user/actions';
import {
  setAccumulated,
  setAutoBetAmount,
  setAutoBetCount,
  setAutoPath,
  setAutoStatus,
  updateStatus
} from 'state/dreamtower/actions';
import { api } from 'services';
import { useSound, useCustomSWR } from 'hooks';

import Row from './Row';
import Payout from './Payout';
import { towerResetAnimation, towerWinAnimation } from '../animation';
import { convertBalanceToChip } from 'utils/balance';

interface Props {
  dreamRef: React.MutableRefObject<any>;
}

export default function AutoTower({ dreamRef }: Props) {
  const user = useAppSelector(state => state.user);
  const game = useAppSelector(state => state.dreamtower.game);
  const auto = useAppSelector(state => state.dreamtower.auto);
  const meta = useAppSelector(state => state.meta.dreamtower);
  const dispatch = useAppDispatch();
  const [initialBetAmount, setInitialBetAmount] = useState(0);
  const [currentBetAmount, setCurrentBetAmount] = useState(0);

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

  const buttonRef = useRef<any>(null);
  const [request, setRequest] = useState(false);

  const handleBet = useCallback(
    async (initialBet: number) => {
      if (
        auto.betCount === 0 ||
        (auto.stopProfit !== undefined &&
          auto.accumulated >= auto.stopProfit) ||
        (auto.stopLoss !== undefined && auto.accumulated <= -auto.stopLoss)
      ) {
        return;
      }

      if (auto.betAmount > user.balance) {
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
          usdAmount: auto.betAmount,
          wagered: auto.betAmount,
          balanceType: user.betBalanceType
        })
      );

      if (auto.betCount && auto.betCount > 0) {
        dispatch(setAutoBetCount(auto.betCount - 1));
      }

      const betData = {
        betAmount: Math.floor(auto.betAmount),
        bets: game.bets,
        difficulty: game.difficulty.level,
        paidBalanceType: user.betBalanceType
      };

      try {
        setRequest(true);
        let res = await api.post('/dreamtower/bet', betData);
        setCurrentBetAmount(auto.betAmount);
        setRequest(false);

        dispatch(updateStatus(res.data));
        if (res.data.status !== 'loss') {
          const profit =
            maxWinning > res.data.multiplier * auto.betAmount
              ? Math.floor(res.data.multiplier * auto.betAmount)
              : Math.floor(maxWinning);

          dispatch(
            updateBalance({
              type: 1,
              usdAmount: profit,
              balanceType: betData.paidBalanceType
            })
          );

          dispatch(
            setAccumulated(
              auto.accumulated +
                Math.floor((res.data.multiplier - 1) * auto.betAmount)
            )
          );
          if (auto.changeBetOnWin === 0) {
            dispatch(setAutoBetAmount(initialBet));
          } else if (auto.changeBetOnWin) {
            dispatch(
              setAutoBetAmount(
                Math.floor((auto.betAmount * (100 + auto.changeBetOnWin)) / 100)
              )
            );
          }
        } else {
          dispatch(setAccumulated(auto.accumulated - auto.betAmount));
          if (auto.changeBetOnLoss === 0) {
            dispatch(setAutoBetAmount(initialBet));
          } else if (auto.changeBetOnLoss) {
            dispatch(
              setAutoBetAmount(
                Math.floor(
                  (auto.betAmount * (100 + auto.changeBetOnLoss)) / 100
                )
              )
            );
          }
        }
      } catch (error: any) {
        dispatch(
          updateBalance({
            type: 1,
            usdAmount: betData.betAmount,
            wagered: betData.betAmount,
            balanceType: betData.paidBalanceType
          })
        );

        if (error.response.status === 503) {
          toast.error('This function is blocked by admin.');
        } else toast.error(error.response.data.message);
        setRequest(false);
        dispatch(setAutoStatus(''));
        return;
      }
    },
    [
      auto,
      user.balance,
      user.betBalanceType,
      game.betAmount,
      game.bets,
      game.difficulty.level,
      meta.minAmount,
      meta.maxAmount,
      dispatch
    ]
  );

  const handleStart = useCallback(async () => {
    if (request) return;
    if (user.name === '') {
      toast.info('Please sign in.');
      return;
    }
    setInitialBetAmount(game.betAmount);
    dispatch(setAutoStatus('running'));
    handleBet(game.betAmount);
  }, [dispatch, handleBet, request, user.name, game.betAmount]);

  const handleStop = useCallback(async () => {
    if (request) return;
    dispatch(setAutoStatus(''));
  }, [dispatch, request]);

  const handleClickSquare = async (height: any, index: any) => {
    dispatch(setAutoPath({ height, index }));
  };

  useEffect(() => {
    let tl: gsap.core.Tween | undefined = undefined;
    if (buttonRef)
      gsap.set(buttonRef.current, {
        x: '-100%'
      });
    if (auto.status === 'running') {
      if (buttonRef && game.roundId) {
        tl = gsap.to(buttonRef.current, {
          x: 0,
          duration: game.status === 'win' ? 6 : 3,
          ease: 'none',
          onComplete: () => {
            handleBet(initialBetAmount);
          }
        });
      }
    }
    return () => {
      tl && tl.kill();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [game.roundId, game.status, auto.status]);

  const endedHeight = useMemo(() => {
    if (game.status === '') return -1;
    for (var i = 0; i < game.tower.length; i++) {
      if (game.tower[i][game.bets[i]] === 0) return i;
    }
    return game.tower.length - 1;
  }, [game.tower, game.bets, game.status]);

  useEffect(() => {
    if (!dreamRef) return;
    if (game.status === 'win') {
      const q = gsap.utils.selector(dreamRef);

      const tl = towerWinAnimation({
        textTargets: [q('.dreamtower_dreamtext'), q('.dreamtower_towertext')],
        rowTargets: [...q('.dreamtower_row')],
        duration: 0.15
      });
      towerPlay.win && towerPlay.win();
      return () => {
        tl.kill();
        towerStop.win();
      };
    } else {
      const q = gsap.utils.selector(dreamRef);

      const rows = [...q('.dreamtower_row')];
      var tr1: typeof rows = [],
        tr2: typeof rows = [],
        tr3: typeof rows = [],
        tr4: typeof rows = [];
      rows.forEach((row, index) => {
        if (auto.status === '' && index === game.bets.length) {
          tr4.push(row);
        } else if (
          auto.status === 'running' &&
          (game.status === 'win' || game.status === 'cashout')
        ) {
          tr2.push(row);
        } else if (game.status === 'loss' && index === endedHeight) {
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
        rowTargets4: tr4
      });
      return () => {
        tl.kill();
      };
    }
  }, [
    game.status,
    game.bets,
    endedHeight,
    towerPlay,
    towerStop,
    auto.status,
    dreamRef
  ]);

  return (
    <>
      {(game.status === 'win' || game.status === 'cashout') && (
        <Payout
          multiplier={game.multiplier}
          profit={
            maxWinning > game.multiplier * currentBetAmount
              ? Math.floor(game.multiplier * currentBetAmount)
              : Math.floor(maxWinning)
          }
          chipType={user.betBalanceType}
        />
      )}
      <Flex py="5px" flexDirection={'column-reverse'}>
        {game.tower.map((v, i) => (
          <Row
            roundId={game.roundId}
            key={i}
            value={v}
            towerMode="auto"
            isNext={auto.status === '' && i === game.bets.length}
            isHighlight={
              (game.status === 'loss' && i === endedHeight) ||
              (auto.status === '' && i === game.bets.length) ||
              (auto.status === 'running' &&
                (game.status === 'win' || game.status === 'cashout'))
            }
            isClickable={auto.status === '' && i <= game.bets.length}
            isUnderBroken={i <= endedHeight}
            selectedIndex={game.bets.length > i ? game.bets[i] : undefined}
            nextMultiplier={game.nextMultiplier!}
            blocksInRow={game.difficulty.blocksInRow}
            handleClickSquare={(index: any) => {
              handleClickSquare(i, index);
            }}
          />
        ))}
      </Flex>
      <Box px={'20px'}>
        <StyledCashButton
          onClick={
            request
              ? undefined
              : auto.status === 'running'
              ? handleStop
              : game.bets.length > 0
              ? handleStart
              : null
          }
          //@ts-ignore
          style={{ position: 'relative' }}
          overflow="hidden"
          disabled={game.bets.length === 0}
        >
          <div
            style={{
              position: 'absolute',
              width: '100%',
              height: '100%',
              top: 0,
              background: '#00000040'
            }}
            ref={buttonRef}
          />
          {request ? (
            <ClipLoader color="#fff" size={20} />
          ) : (
            <>
              {auto.status === 'running'
                ? 'Stop Auto Bet'
                : game.bets.length > 0
                ? 'Start Auto Bet'
                : 'Choose Your Path'}
            </>
          )}
        </StyledCashButton>
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
