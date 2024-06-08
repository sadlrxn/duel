import React, {
  createContext,
  useMemo,
  useEffect,
  useState,
  useRef,
  useCallback
} from 'react';
import { shallowEqual } from 'react-redux';
import gsap from 'gsap';

import { useAppSelector, useAppDispatch } from 'state';
import {
  calculateSpeed,
  calculateAngle,
  calculateAngleWithSpeed
} from 'pages/Crash/utils';
import { CrashAutoBet, CrashBet } from 'api/types/crash';
import { sendMessage } from 'state/socket';
import { toggleCrashAnimation, updateBalance } from 'state/user/actions';
import { setBalanceType } from 'state/user/actions';
import { convertChipToBalance } from 'utils/balance';

export const CrashContext = createContext<any>({});

export const CrashProvider: React.FC<React.PropsWithChildren> = ({
  children
}) => {
  const {
    bets,
    status,
    serverTimeElapsed,
    startedTime,
    from,
    to,
    roundId,
    cashIns,
    cashOuts
  } = useAppSelector(state => state.crash);
  const {
    eventInterval: duration,
    maxCashOut,
    minBetAmount,
    maxBetAmount
  } = useAppSelector(state => state.meta.crash, shallowEqual);
  const user = useAppSelector(state => state.user, shallowEqual);

  const dispatch = useAppDispatch();

  const initialAutoBet = useMemo(() => {
    return {
      betAmount: 100000,
      paidBalanceType: user.betBalanceType,
      cashOutAt: 2.5,
      rounds: -1,
      pnl: 0,
      betId: -1,
      roundId: -1,
      profit: 0,
      bettedRounds: 0,
      isComplete: false,
      isBetted: false,
      onLoss: 0,
      onWin: 0,
      stopProfit: 0,
      stopLoss: 0
    } as CrashAutoBet;
  }, [user.betBalanceType]);

  const [liveMultiplier, setMultiplier] = useState(1);
  const [rocketSpeed, setSpeed] = useState(1);
  const [rocketAngle, setAngle] = useState(45);
  const [rulerMax, setMultipleMax] = useState(1);
  const [hRepeatBet, setRepeatBet] = useState(false);
  const [hRepeatBetEnabled, setRepeatBetEnabled] = useState(true);
  const [repeatBetRequest, setRepeatBetRequest] = useState(false);
  const [nextHookBet, setNextBet] = useState(false);
  const [nextBets, setNextBets] = useState<CrashBet[]>([]);
  const [hCurrentAutoBetIndex, setCurrentAutoBetIndex] = useState(0);
  const [hAutoBets, setAutoBets] = useState<CrashAutoBet[]>([initialAutoBet]);
  const [hAutoBetEnable, setAutoBetEnable] = useState(false);
  const [autoBetRequest, setAutoBetRequest] = useState(false);
  const [hCurrentAutoBet, setCurrentAutoBet] = useState(initialAutoBet);
  const [hShowStatus, setShowStatus] = useState(false);

  const roundIdRef = useRef<number>(-1);

  const [userRoundBets, userBetCount] = useMemo(() => {
    const userRoundBets = bets.filter(bet => bet.user.id === user.id);
    return [userRoundBets, userRoundBets.length];
  }, [bets, user.id]);

  const [userBets, userBetted] = useMemo(() => {
    const userBets = userRoundBets.filter(bet => !bet.payoutMultiplier);
    const userBetted = userRoundBets.length !== 0;
    return [userBets, userBetted];
  }, [userRoundBets]);

  const [totalBet, couponTotalBet] = useMemo(() => {
    let totalBet: number = 0;
    let couponTotalBet: number = 0;

    userBets.forEach(bet => {
      const balance = Math.min(maxCashOut, bet.betAmount * liveMultiplier);
      if (bet.paidBalanceType === 'chip') totalBet += balance;
      else if (bet.paidBalanceType === 'coupon') couponTotalBet += balance;
    });

    return [totalBet, couponTotalBet];
  }, [maxCashOut, liveMultiplier, userBets]);

  const handleToggleAnimation = useCallback(() => {
    dispatch(toggleCrashAnimation());
  }, [dispatch]);

  useEffect(() => {
    if (hAutoBetEnable) setRepeatBetEnabled(false);
  }, [hAutoBetEnable]);

  useEffect(() => {
    if (user.balances.coupon.balance >= convertChipToBalance(0.01)) {
      dispatch(setBalanceType('coupon'));
    } else {
      dispatch(setBalanceType('chip'));
    }
  }, [dispatch, user.balances.coupon.balance]);

  useEffect(() => {
    if (autoBetRequest) {
      let autoBets = [...hAutoBets];
      for (let i = 0; i < cashIns.length; i++) {
        const index = autoBets.findIndex(bet => {
          return (
            bet.betAmount === cashIns[i].amount &&
            bet.paidBalanceType === cashIns[i].balanceType &&
            bet.cashOutAt === cashIns[i].cashOutAt
          );
        });
        if (index === -1 || autoBets[index].isBetted) continue;
        autoBets[index].betId = cashIns[i].betId;
        autoBets[index].roundId = roundId;
        autoBets[index].bettedRounds = autoBets[index].bettedRounds + 1;
        autoBets[index].isBetted = true;
      }
      setAutoBets(() => [...autoBets]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [autoBetRequest, cashIns, roundId, status]);

  useEffect(() => {
    if (autoBetRequest) {
      let autoBets = [...hAutoBets];
      for (let i = 0; i < cashOuts.length; i++) {
        const index = autoBets.findIndex(bet => {
          return (
            bet.betId === cashOuts[i].betId &&
            bet.roundId === roundId &&
            bet.isBetted
          );
        });
        if (index === -1 || autoBets[index].profit !== 0) continue;

        autoBets[index].pnl += cashOuts[i].amount - autoBets[index].betAmount;
        autoBets[index].profit = cashOuts[i].amount;
      }
      setAutoBets(() => [...autoBets]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [autoBetRequest, cashOuts, roundId]);

  useEffect(() => {
    if (status === 'back' && autoBetRequest) {
      let autoBets = [...hAutoBets];
      for (let i = 0; i < autoBets.length; i++) {
        if (autoBets[i].isBetted) {
          if (autoBets[i].profit === 0) {
            autoBets[i].pnl -= autoBets[i].betAmount;
            autoBets[i].betAmount *= 1 + autoBets[i].onLoss / 100;
          } else {
            autoBets[i].betAmount *= 1 + autoBets[i].onWin / 100;
          }

          if (autoBets[i].betAmount > maxBetAmount)
            autoBets[i].betAmount = maxBetAmount;
          if (autoBets[i].betAmount < minBetAmount)
            autoBets[i].betAmount = minBetAmount;
        }

        autoBets[i].betAmount = Math.floor(autoBets[i].betAmount);
        autoBets[i].isBetted = false;
        autoBets[i].profit = 0;

        if (
          autoBets[i].bettedRounds === autoBets[i].rounds ||
          (autoBets[i].stopLoss && autoBets[i].pnl <= -autoBets[i].stopLoss) ||
          (autoBets[i].stopProfit && autoBets[i].pnl >= autoBets[i].stopProfit)
        ) {
          autoBets[i].isComplete = true;
        }
      }
      setAutoBets(() => [...autoBets]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [autoBetRequest, status]);

  useEffect(() => {
    if (userRoundBets.length !== 0 && !nextHookBet) {
      setNextBets(userRoundBets);
    }
  }, [nextHookBet, userRoundBets]);

  useEffect(() => {
    if (nextBets.length === 0 || nextBets[0].roundId < roundId - 1) {
      setRepeatBetEnabled(false);
      setRepeatBet(false);
    } else {
      setRepeatBetEnabled(true);
    }
  }, [nextBets, roundId]);

  useEffect(() => {
    if (roundId !== roundIdRef.current) {
      roundIdRef.current = roundId;
      setRepeatBetRequest(false);
      setAutoBetRequest(false);
    }

    if (hAutoBetEnable) {
      if (!autoBetRequest && !userBetted && !repeatBetRequest) {
        let userBalance = user.balance;
        const betBalanceType = user.betBalanceType;
        hAutoBets.forEach(bet => {
          if (bet.betAmount > userBalance) return;
          if (
            bet.paidBalanceType === 'chip' &&
            betBalanceType === 'coupon' &&
            userBalance > 0
          )
            return;
          if (bet.isComplete) return;
          userBalance -= bet.betAmount;

          dispatch(
            sendMessage({
              type: 'event',
              room: 'crash',
              content: JSON.stringify({
                type: 'cash-in',
                content: JSON.stringify({
                  amount: bet.betAmount,
                  balanceType: bet.paidBalanceType,
                  roundId,
                  cashOutAt: bet.cashOutAt
                })
              })
            })
          );

          dispatch(
            updateBalance({
              type: -1,
              usdAmount: bet.betAmount,
              wagered: bet.betAmount,
              balanceType: bet.paidBalanceType
            })
          );
        });
        setAutoBetRequest(true);
      }
    } else if (
      hRepeatBetEnabled &&
      !repeatBetRequest &&
      (nextHookBet || hRepeatBet) &&
      nextBets.length > 0 &&
      nextBets[0].roundId === roundId - 1
    ) {
      let userBalance = user.balance;
      const betBalanceType = user.betBalanceType;
      nextBets.forEach(bet => {
        if (bet.betAmount > userBalance) return;
        if (
          bet.paidBalanceType === 'chip' &&
          betBalanceType === 'coupon' &&
          userBalance > 0
        )
          return;

        userBalance -= bet.betAmount;

        dispatch(
          sendMessage({
            type: 'event',
            room: 'crash',
            content: JSON.stringify({
              type: 'cash-in',
              content: JSON.stringify({
                amount: bet.betAmount,
                balanceType: bet.paidBalanceType,
                roundId,
                cashOutAt: bet.cashOutAt
              })
            })
          })
        );

        dispatch(
          updateBalance({
            type: -1,
            usdAmount: bet.betAmount,
            wagered: bet.betAmount,
            balanceType: bet.paidBalanceType
          })
        );
      });
      setNextBet(false);
      setRepeatBetRequest(true);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    dispatch,
    nextBets,
    nextHookBet,
    repeatBetRequest,
    hRepeatBetEnabled,
    hAutoBetEnable,
    hRepeatBet,
    roundId,
    autoBetRequest,
    userBetted,
    hAutoBets
  ]);

  useEffect(() => {
    const timeElapsed = Date.now() - startedTime - serverTimeElapsed;
    let speed = 1;
    let angle = 45;
    if (status === 'play') {
      speed = calculateSpeed((timeElapsed + serverTimeElapsed) / 1000);
      angle = calculateAngleWithSpeed(speed);
    }
    setSpeed(speed);
    setAngle(angle);

    let progress = timeElapsed / duration;
    if (progress > 1) progress = 1;
    if (progress < 0) progress = 0;

    let roundValue = {
      value: 0
    };

    const tl = gsap
      .timeline()
      .fromTo(
        roundValue,
        {
          value: 0
        },
        {
          value: 10 ** 6,
          roundProps: 'value',
          duration: duration / 1000,
          ease: 'none',
          onUpdate: () => {
            const percent = roundValue.value / 10 ** 6;
            const multiplier = from + (to - from) * percent;
            setMultiplier(multiplier);

            if (status === 'play' || status === 'explosion') {
              const time = (Date.now() - startedTime) / 1000;
              const angle = calculateAngle(time);
              let rocketPercent = time / 15;
              if (rocketPercent > 1) rocketPercent = 1;
              const max = ((multiplier - 1) / rocketPercent) * 1.0 + 1;
              setMultipleMax(max);
              setSpeed(speed);
              setAngle(angle);
            } else {
              setMultipleMax(1);
            }
          }
        }
      )
      .progress(progress);

    return () => {
      tl.kill();
    };
  }, [duration, from, serverTimeElapsed, startedTime, status, to]);

  useEffect(() => {
    if (hCurrentAutoBetIndex < 0) setCurrentAutoBetIndex(0);
    if (hCurrentAutoBetIndex >= hAutoBets.length)
      setCurrentAutoBetIndex(hAutoBets.length - 1);
  }, [hCurrentAutoBetIndex, hAutoBets.length]);

  const handleAddSlot = useCallback(() => {
    if (hAutoBetEnable) return;
    if (hAutoBets.length < 5) {
      setAutoBets(prev => [...prev, initialAutoBet]);
    }
  }, [hAutoBetEnable, hAutoBets.length, initialAutoBet]);

  const handleRemove = useCallback(
    (index: number) => {
      if (hAutoBetEnable) return;
      if (hAutoBets.length > 1) {
        setAutoBets(prev => {
          let newBets = [...prev];
          newBets.splice(index, 1);
          return newBets;
        });
      }
    },
    [hAutoBetEnable, hAutoBets.length]
  );

  const handleReset = useCallback(() => {
    if (hAutoBetEnable) return;
    if (hAutoBets.length === 1) {
      setAutoBets([initialAutoBet]);
    }
  }, [hAutoBetEnable, hAutoBets.length, initialAutoBet]);

  const multiplier = useMemo(() => liveMultiplier, [liveMultiplier]);
  const speed = useMemo(() => rocketSpeed, [rocketSpeed]);
  const angle = useMemo(() => rocketAngle, [rocketAngle]);
  const multipleMax = useMemo(() => rulerMax, [rulerMax]);
  const repeatBet = useMemo(() => hRepeatBet, [hRepeatBet]);
  const repeatBetEnabled = useMemo(
    () => hRepeatBetEnabled,
    [hRepeatBetEnabled]
  );
  const autoBets = useMemo(() => hAutoBets, [hAutoBets]);
  const autoBetEnable = useMemo(() => hAutoBetEnable, [hAutoBetEnable]);
  const currentAutoBet = useMemo(() => hCurrentAutoBet, [hCurrentAutoBet]);
  const showStatus = useMemo(() => hShowStatus, [hShowStatus]);
  const usingAnimation = useMemo(
    () => user.crashAnimation,
    [user.crashAnimation]
  );
  const currentAutoBetIndex = useMemo(
    () => hCurrentAutoBetIndex,
    [hCurrentAutoBetIndex]
  );

  return (
    <CrashContext.Provider
      value={{
        multiplier,
        speed,
        angle,
        multipleMax,
        userBets,
        totalBet,
        couponTotalBet,
        userBetted,
        repeatBet,
        repeatBetEnabled,
        autoBets,
        autoBetEnable,
        currentAutoBet,
        currentAutoBetIndex,
        showStatus,
        usingAnimation,
        userBetCount: useMemo(() => userBetCount, [userBetCount]),
        setRepeatBet,
        setNextBet,
        setNextBets,
        setAutoBets,
        setCurrentAutoBetIndex,
        setAutoBetEnable,
        setCurrentAutoBet,
        setShowStatus,
        handleToggleAnimation,
        handleAddSlot,
        handleRemove,
        handleReset
      }}
    >
      {children}
    </CrashContext.Provider>
  );
};
