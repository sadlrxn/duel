import { useRef, useCallback, useEffect, useState } from 'react';
import ReactCanvasConfetti from 'react-canvas-confetti';

function randomInRange(min: number, max: number) {
  return Math.random() * (max - min) + min;
}

function getAnimationSettings(_angle: number, _originX: number) {
  return {
    particleCount: 1,
    startVelocity: 0,
    ticks: 50,
    gravity: 0.4,
    origin: {
      x: Math.random(),
      y: Math.random()
    },
    colors: ['#ffffff'],
    shapes: ['circle'],
    scalar: randomInRange(0.4, 1)
  };
}

export default function Snow() {
  const refAnimationInstance = useRef<any>();
  const [intervalId, setIntervalId] = useState<NodeJS.Timer | undefined>(
    undefined
  );

  const getInstance = useCallback((instance: any) => {
    refAnimationInstance.current = instance;
  }, []);

  const nextTickAnimation = useCallback(() => {
    if (refAnimationInstance && refAnimationInstance.current) {
      refAnimationInstance.current(getAnimationSettings(60, 0));
      refAnimationInstance.current(getAnimationSettings(120, 1));
    }
  }, [refAnimationInstance]);

  const startAnimation = useCallback(() => {
    if (!intervalId) {
      setIntervalId(setInterval(nextTickAnimation, 16));
    }
  }, [nextTickAnimation, intervalId]);

  // const pauseAnimation = useCallback(() => {
  //   clearInterval(intervalId);
  //   setIntervalId(undefined);
  // }, [intervalId]);

  // const stopAnimation = useCallback(() => {
  //   clearInterval(intervalId);
  //   setIntervalId(undefined);
  //   refAnimationInstance.current && refAnimationInstance.current.reset();
  // }, [intervalId]);

  useEffect(() => {
    startAnimation();
    return () => {
      clearInterval(intervalId);
    };
  }, [intervalId, startAnimation]);

  return (
    <>
      <ReactCanvasConfetti
        refConfetti={getInstance}
        style={{
          position: 'fixed',
          left: 0,
          top: 0,
          width: '100vw',
          height: '100vh'
        }}
      />
    </>
  );
}
