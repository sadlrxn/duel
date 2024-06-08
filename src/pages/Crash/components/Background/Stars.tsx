import { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import gsap from 'gsap';

// import starsBg from 'assets/imgs/crash/stars.png';
import starsBg from 'assets/imgs/crash/stars2.svg';

import { CRASH_BACK_TIME } from 'pages/Crash/config';
import { calculateAngleWithSpeed, calculateSpeed } from 'pages/Crash/utils';
import { Box } from 'components';
import { useAppSelector } from 'state';

import { starAnimation as opacityAnimation } from './animation';

export default function Stars() {
  const {
    status,
    time: globalTime,
    serverTimeElapsed,
    startedTime
  } = useAppSelector(state => state.crash);
  const { eventInterval: duration, preparingDuration } = useAppSelector(
    state => state.meta.crash
  );
  const usingAnimation = useAppSelector(state => state.user.crashAnimation);

  const [speed, setSpeed] = useState(1);
  const [angle, setAngle] = useState(45);
  const [animDuration, setAnimDuration] = useState(0.2);

  const [deltaX, setDeltaX] = useState(0);
  const [deltaY, setDeltaY] = useState(0);

  const starsRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let speed = 1;
    let angle = 45;
    let animDuration = 1;
    if (usingAnimation) {
      switch (status) {
        case 'bet':
          speed = 0;
          animDuration = 0;
          break;
        case 'ready':
          speed = 1;
          animDuration = preparingDuration;
          break;
        case 'play':
          const timeElapsed = (Date.now() - startedTime) / 1000;
          speed = calculateSpeed(timeElapsed);
          angle = calculateAngleWithSpeed(speed);
          animDuration = 1;
          break;
        case 'explosion':
          speed = 0;
          animDuration = 0;
          break;
        case 'back':
          speed = -4;
          animDuration = CRASH_BACK_TIME;
          break;
      }
    } else speed = 0;
    setSpeed(speed);
    setAngle(angle);
    setAnimDuration(animDuration);
  }, [
    status,
    serverTimeElapsed,
    startedTime,
    duration,
    preparingDuration,
    usingAnimation
  ]);

  useEffect(() => {
    const deltaX = Math.cos((angle / 180) * Math.PI) * 300;
    const deltaY = Math.sin((angle / 180) * Math.PI) * 300;
    setDeltaX(deltaX);
    setDeltaY(deltaY);
    // console.info('set delta x & y');
  }, [angle]);

  useEffect(() => {
    if (!usingAnimation) return;
    const tl = gsap.timeline().to(starsRef.current, {
      backgroundPositionX: `-=${deltaX * speed * animDuration}px`,
      backgroundPositionY: `+=${deltaY * speed * animDuration}px`,
      duration: animDuration,
      ease: 'none',
      modifiers: {
        backgroundPositionX: gsap.utils.unitize(
          gsap.utils.wrap(-1920 * 2, 1920 * 2),
          'px'
        ),
        backgroundPositionY: gsap.utils.unitize(
          gsap.utils.wrap(-1620 * 2, 1620 * 2),
          'px'
        )
      }
    });

    return () => {
      tl.kill();
    };
  }, [animDuration, deltaX, deltaY, speed, serverTimeElapsed, usingAnimation]);

  const progressRef = useRef(0);

  useEffect(() => {
    let interval: NodeJS.Timer;
    let tl = opacityAnimation(containerRef);
    let progress = 0;
    const duration = tl.duration();

    if (!usingAnimation) {
      tl.progress(0.01).pause();
      return;
    }

    switch (status) {
      case 'bet':
        tl.progress(0.01).pause();
        break;
      case 'ready':
        tl.progress(0.01).pause();
        break;
      case 'play':
        progress = (Date.now() - globalTime) / (duration * 1000);
        if (progress > 1) progress = 0.998;
        if (progress < 0) progress = 0.01;
        tl.progress(progress).play();

        interval = setInterval(() => {
          tl.pause();
          progress = (Date.now() - globalTime) / (duration * 1000);
          if (progress > 1) progress = 0.998;
          if (progress < 0) progress = 0.01;
          tl.progress(progress).play();
        }, 1000);
        break;
      case 'explosion':
        tl.progress(progressRef.current).pause();
        break;
      case 'back':
        let backDuration = CRASH_BACK_TIME - (Date.now() - globalTime) / 1000;
        if (backDuration < 0) backDuration = 0;
        tl.progress(progressRef.current)
          .duration(backDuration)
          .reverse()
          .then(() => {
            progressRef.current = 0;
            tl.progress(0.01).pause();
          });
        break;
    }

    return () => {
      progressRef.current = tl.progress();
      tl.kill();
      clearInterval(interval);
    };
  }, [globalTime, status, usingAnimation]);

  return (
    <Container ref={containerRef}>
      <StarBg ref={starsRef} />
    </Container>
  );
}

const StarBg = styled(Box)`
  width: 100%;
  height: 100%;
  background-image: url(${starsBg});
  background-repeat: repeat;
`;

const Container = styled(Box)`
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  z-index: -2;
`;
