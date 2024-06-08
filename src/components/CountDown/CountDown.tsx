import { useEffect, useState, useCallback } from 'react';
import { Span } from 'components/Text';

interface CountDownProps {
  endedAt: number;
}

const convertNumberToString = (val: number, maxLength: number = 2) => {
  return val.toString().padStart(maxLength, '0');
};

export default function CountDown({ endedAt }: CountDownProps) {
  const [count, setCount] = useState(0);
  const [countText, setCountText] = useState('00:00:00');

  const tick = useCallback(() => {
    let time = (endedAt - Date.now()) / 1000;

    if (time < 0) {
      setCount(0);
    } else {
      setCount(Math.ceil(time));
    }

    setCountText(
      `${convertNumberToString(
        Math.floor(count / 3600)
      )}:${convertNumberToString(
        Math.floor(count / 60) % 60
      )}:${convertNumberToString(Math.floor(count % 60))}`
    );
  }, [count, endedAt]);

  useEffect(() => {
    tick();
    const interval = setInterval(() => tick(), 1000);
    return () => clearInterval(interval);
  }, [tick]);

  return (
    <Span color="#CAFFFF" fontWeight={600} fontSize={12}>
      {countText}
    </Span>
  );
}
