import { useEffect, useRef, useState, useCallback } from 'react';
import styled from 'styled-components';
import gsap from 'gsap';

import { useAppSelector } from 'state';
import { TGameStatus } from 'api/types/jackpot';

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
  margin-bottom: 38px;

  .width_700 & {
    margin-bottom: 0px;
  }
`;

const Counter = styled.div`
  display: flex;
  background-color: #05090d;
  width: 100px;
  height: 30px;
  justify-content: center;
  align-items: center;
  border: 2px solid ${({ theme }) => theme.colors.border};
  border-bottom-left-radius: 5px;
  border-bottom-right-radius: 5px;
  border-top: 0;
  font-weight: 600;
  color: #627694;
`;

const convertNumberToString = (val: number, maxLength: number = 2) => {
  return val.toString().padStart(maxLength, '0');
};

interface CountDownProps {
  status: TGameStatus;
  time: number;
  isHistory?: boolean;
}

export default function CountDown({
  status,
  time: updatedTime,
  isHistory
}: CountDownProps) {
  const meta = useAppSelector(state => state.meta.grandJackpot);
  const barRef = useRef<any>();

  const [max, setMax] = useState(40);
  const [count, setCount] = useState(0);
  const [roll, setRoll] = useState(false);
  const [countText, setCountText] = useState('00:00:00');

  const setData = useCallback(() => {
    let max = meta.bettingTime,
      count = meta.bettingTime;
    let time = (Date.now() - updatedTime) / 1000;
    let roll = false;
    if (time < 0) time = 0;
    switch (status) {
      case 'started':
        max = meta.bettingTime;
        count = Math.ceil(time);
        if (count > max) {
          count -= max;
          max = meta.rollingTime - meta.winnerTime;
          roll = true;
        }
        break;
      case 'counting':
        max = meta.countingTime;
        count = Math.ceil(time);
        if (count > max) {
          count -= max;
          max = meta.rollingTime - meta.winnerTime;
          roll = true;
        }
        break;
      case 'rolling':
        max = isHistory ? 15 : meta.rollingTime - meta.winnerTime;
        count = Math.ceil(time);
        if (isHistory) {
          // max -= 45;
          // count -= 45;
        }
        roll = true;
        break;
      case 'rollend':
        max = meta.winnerTime;
        count = Math.ceil(time);
        if (isHistory) count = max;
        roll = true;
        break;
    }

    count = max - count;
    if (count < 0) count = 0;
    setMax(max);
    setCount(count);
    setRoll(roll);
    setCountText(
      `${convertNumberToString(
        Math.floor(count / 3600)
      )}:${convertNumberToString(
        Math.floor(count / 60) % 60
      )}:${convertNumberToString(Math.floor(count % 60))}`
    );
  }, [status, updatedTime, meta, isHistory]);

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
      <Counter>{countText}</Counter>
    </>
  );
}
