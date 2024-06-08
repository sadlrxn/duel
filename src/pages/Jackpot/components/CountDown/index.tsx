import { useEffect, useRef, useState, useCallback } from 'react';
import styled from 'styled-components';
import gsap from 'gsap';

import { useAppSelector } from 'state';
import { TGameStatus } from 'api/types/jackpot';
import { Flex, Label } from 'components';

const Bar = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
`;

const Progress = styled.div`
  position: relative;
  border: 2px solid ${({ theme }) => theme.colors.border};
  border-top: 0;
  border-bottom: 0;
  height: 5px;
  overflow: hidden;
  background-color: #182738;
  margin-bottom: 0px;
`;

const Counter = styled.div`
  background-color: #05090d;
  width: 75px;
  height: 30px;
  justify-content: center;
  align-items: center;
  border: 2px solid ${({ theme }) => theme.colors.border};
  border-bottom-left-radius: 5px;
  border-bottom-right-radius: 5px;
  border-top: 0;
  font-weight: 600;
  color: #627694;

  display: flex;
`;

interface CountDownProps {
  status: TGameStatus;
  time: number;
  countingTime?: number;
  betCount?: number;
}

export default function CountDown({
  status,
  time: updatedTime,
  countingTime = 40,
  betCount
}: CountDownProps) {
  const meta = useAppSelector(state => state.meta.jackpot[state.jackpot.room]);
  const barRef = useRef<any>();

  const [max, setMax] = useState(40);
  const [count, setCount] = useState(0);
  const [roll, setRoll] = useState(false);

  const setData = useCallback(() => {
    let max = countingTime,
      count = countingTime;
    let time = (Date.now() - updatedTime) / 1000;
    let roll = false;
    if (time < 0) time = 0;
    switch (status) {
      case 'started':
        max = countingTime;
        count = Math.ceil(time);
        if (count > max) {
          count -= max;
          max = meta.rollingTime - meta.winnerTime;
          roll = true;
        }
        break;
      case 'rolling':
        max = meta.rollingTime - meta.winnerTime;
        count = Math.ceil(time);
        roll = true;
        break;
      case 'rollend':
        max = meta.winnerTime;
        count = Math.ceil(time);
        roll = true;
        break;
    }

    count = max - count;
    if (count < 0) count = 0;
    setMax(max);
    setCount(count);
    setRoll(roll);
  }, [status, updatedTime, meta, countingTime]);

  useEffect(() => {
    setData();
    const interval = setInterval(() => setData(), 1000);
    return () => clearInterval(interval);
  }, [setData]);

  useEffect(() => {
    gsap.fromTo(
      barRef.current,
      {
        background: roll ? '#FFE87F' : '#4FFF8B',
        width: `${((max - count) / max) * 100}%`
      },
      {
        width: `${((max - count + 1) / max) * 100}%`,
        ease: 'none',
        duration: 1
      }
    );
  }, [count, max, roll]);

  return (
    <>
      <Progress>
        <Bar ref={barRef} />
      </Progress>
      <Flex justifyContent="space-between" alignItems="center">
        <Counter>0:{`${count < 10 ? '0' : ''}${count}`}</Counter>
        {betCount !== undefined && (
          <Label color="text" fontWeight={600}>
            Remaining Bets: {meta.betCountLimit - betCount}
          </Label>
        )}
      </Flex>
    </>
  );
}
