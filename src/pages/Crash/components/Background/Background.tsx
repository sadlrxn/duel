import { useEffect, useRef } from 'react';
import styled from 'styled-components';

import { Box } from 'components';
import { useAppSelector } from 'state';
import { CRASH_BACK_TIME } from 'pages/Crash/config';

import { backgroundAnim } from './animation';

export default function Background() {
  const status = useAppSelector(state => state.crash.status);
  const globalTime = useAppSelector(state => state.crash.time);
  const usingAnimation = useAppSelector(state => state.user.crashAnimation);

  const containerRef = useRef<HTMLDivElement>(null);

  const progressRef = useRef(0);

  useEffect(() => {
    let tl = backgroundAnim(containerRef);
    let progress = 0;
    let interval: NodeJS.Timer;
    const duration = backgroundAnim(containerRef).duration();

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
        if (progress > 1) progress = 1;
        if (progress < 0) progress = 0;
        tl.progress(progress).play();

        interval = setInterval(() => {
          tl.pause();
          progress = (Date.now() - globalTime) / (duration * 1000);
          if (progress > 1) progress = 1;
          if (progress < 0) progress = 0;
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
    <>
      <Container ref={containerRef}></Container>
    </>
  );
}

const Container = styled(Box)`
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  z-index: -3;
  background: linear-gradient(0deg, #2875ad 0%, #00315d 100%);
`;
