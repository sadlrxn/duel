import React, { useEffect, useState, useCallback } from 'react';
import styled, { keyframes } from 'styled-components';
import gsap from 'gsap';
import { CustomEase } from 'gsap/CustomEase';

import { Box, Flex, Span } from 'components';
import { useAppSelector } from 'state';
import { useCrash, useSound } from 'hooks';

import { Rocket } from '../Animation';
import MultiplierRuler from './MultiplierRuler';
// import Cashout from './Cashout';

gsap.registerPlugin(CustomEase);

interface CrashMainProps {
  graphRef: React.RefObject<HTMLDivElement>;
}

export default function CrashMain({ graphRef }: CrashMainProps) {
  const {
    status,
    time: globalTime,
    startedTime
  } = useAppSelector(state => state.crash);
  const usingAnimation = useAppSelector(state => state.user.crashAnimation);
  const { bettingDuration } = useAppSelector(state => state.meta.crash);
  const { speed, angle, multiplier, multipleMax } = useCrash();
  const { crashPlay } = useSound();

  const [leftTime, setLeftTime] = useState(() => '0:00');

  const calcualteLeftTime = useCallback(() => {
    let leftTime = bettingDuration * 1000 - Math.floor(Date.now() - globalTime);
    if (leftTime < 0) leftTime = 0;
    if (leftTime > bettingDuration * 1000) leftTime = bettingDuration * 1000;
    else if (Math.floor(leftTime / 10) === 500) crashPlay.c4();
    else if (Math.floor(leftTime / 10) === 400) crashPlay.c3();
    else if (Math.floor(leftTime / 10) === 300) crashPlay.c2();
    else if (Math.floor(leftTime / 10) === 200) crashPlay.c1();
    else if (Math.floor(leftTime / 10) === 100) crashPlay.cend();
    setLeftTime(
      () =>
        Math.floor(leftTime / 1000) +
        ':' +
        ('00' + Math.floor((leftTime % 1000) / 10)).slice(-2)
    );
  }, [bettingDuration, crashPlay, globalTime]);

  const rocketAnimation = useCallback(() => {
    if (!graphRef || !graphRef.current) return gsap.timeline();

    const q = gsap.utils.selector(graphRef);
    let progress = (Date.now() - startedTime) / 15000;
    if (progress > 1) progress = 1;
    // if (progress < 0) progress = 0;

    return gsap
      .timeline()
      .fromTo(
        q('.crash_rocket'),
        {
          top: '100%',
          left: 'calc(100% - 120px)'
        },
        {
          top: '10%',
          duration: 15,
          ease: 'none'
        }
      )
      .fromTo(
        [q('.crash_multiplier'), q('.crash_cashout_users')],
        {
          top: '100%'
        },
        {
          top: '10%',
          duration: 15,
          ease: 'none'
        },
        '<'
      )
      .progress(progress);
  }, [graphRef, startedTime]);

  useEffect(() => {
    let interval: NodeJS.Timer | null = null;
    let tl: gsap.core.Timeline | null = null;
    let q: gsap.utils.SelectorFunc | null = null;
    let progress = 0;

    switch (status) {
      case 'bet':
        calcualteLeftTime();

        interval = setInterval(() => {
          calcualteLeftTime();
        }, 10);
        break;
      case 'ready':
        if (!graphRef || !graphRef.current) break;
        q = gsap.utils.selector(graphRef);

        progress = (Date.now() - globalTime) / 2000;
        if (progress > 1) progress = 1;
        if (progress < 0) progress = 0;

        gsap.set(q('.crash_multiplier'), { top: '100%' });

        tl = gsap
          .timeline()
          .fromTo(
            q('.crash_rocket'),
            {
              left: `-${graphRef.current.offsetLeft + 100}px`
            },
            {
              left: 'calc(100% - 120px)',
              duration: 2,
              ease: 'none'
            }
          )
          .fromTo(
            q('.crash_rocket'),
            {
              top: '120%'
            },
            {
              top: '100%',
              duration: 2,
              ease: CustomEase.create('prepareRocket', '1, 0.68, 0.94, 0.88')
            },
            '<'
          )
          .progress(progress);

        break;
      case 'play':
        if (!graphRef || !graphRef.current) break;
        tl = rocketAnimation();

        interval = setInterval(() => {
          if (tl) tl.kill();
          tl = rocketAnimation();
        }, 1000);
        break;
      case 'explosion':
        break;
      case 'back':
        break;
    }

    return () => {
      if (interval) clearInterval(interval);
      if (tl) tl.kill();
    };
  }, [
    bettingDuration,
    calcualteLeftTime,
    globalTime,
    graphRef,
    rocketAnimation,
    startedTime,
    status
  ]);

  return (
    <>
      {status !== 'bet' && (
        <Multiplier
          className="crash_multiplier"
          style={{
            opacity: status === 'back' ? 0 : 1,
            transition: 'opacity 1s'
          }}
        >
          <Flex flexDirection="column">
            <CustomSpan
              background={
                status === 'explosion' || status === 'back'
                  ? 'linear-gradient(180deg, #f45050 0%, #b33b3b 100%)'
                  : 'linear-gradient(180deg, #a2a8b3 0%, #a2a8b3 100%)'
              }
              style={{
                fontSize: '13px',
                fontWeight: 600,
                position: 'absolute',
                top: 0,
                transform: 'translate(0, -100%)',
                width: '170px',
                minWidth: '170px'
              }}
            >
              CURRENT MULTIPLIER
            </CustomSpan>
            <CustomSpan
              background={
                status === 'explosion' || status === 'back'
                  ? 'linear-gradient(180deg, #f45050 0%, #b33b3b 100%)'
                  : 'linear-gradient(180deg, #fff 0%, #fff 100%)'
              }
            >
              {multiplier.toFixed(2) + 'x'}
            </CustomSpan>
          </Flex>
          <Box
            height="1px"
            width="200px"
            background="linear-gradient(90deg, rgba(77, 169, 255, 0) 0%, rgba(77, 169, 255, 0.5) 11.98%, rgba(77, 169, 255, 0) 100%)"
          />
        </Multiplier>
      )}
      {status === 'bet' ? (
        <CounterWrapper
          width="100%"
          justifyContent="center"
          alignItems="center"
          height="500px"
        >
          <Span color="#A2A8B2" fontWeight={600} fontSize="1em">
            GAME STARTING IN
          </Span>
          <Span color="white" fontWeight={700} fontSize="3.5em" width="2.5em">
            {leftTime}
          </Span>
        </CounterWrapper>
      ) : (
        <MultiplierRuler max={multipleMax} />
      )}
      <CustomRocket
        className="crash_rocket"
        angle={angle}
        speed={speed}
        explosion={status === 'explosion'}
        visible={usingAnimation && status !== 'bet' && status !== 'back'}
      />
      {/* {status !== 'explosion' && status !== 'back' && (
        <CustomCashout multiplier={multiplier} />
      )} */}
    </>
  );
}

const appearAnim = keyframes`
  from {
    opacity:0
  }
  to {
    opacity: 1
  }
`;

const CustomSpan = styled.span<{ background?: string }>`
  background: ${({ background }) => (background ? background : 'transparent')};
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;

  line-height: 1.3;
`;

const CounterWrapper = styled(Flex)`
  width: 100%;
  flex-direction: column;

  justify-content: center;
  align-items: center;

  height: 100px;

  font-size: 12px;

  .width_700 & {
    height: 500px;
  }

  .width_900 & {
    font-size: 14px;
  }

  .width_1100 & {
    font-size: 16px;
  }
`;

const CustomRocket = styled(Rocket)`
  position: absolute;
  right: 120px;

  width: 100px;
  height: 360px;

  /* right: 120px; */

  .width_700 & {
  }

  .width_1100 & {
    width: 150px;
    height: 500px;
  }
`;

const Multiplier = styled(Flex)`
  animation: ${appearAnim} 1s;

  position: absolute;
  right: 0;

  align-items: center;
  transform: translate(0, -50%);

  font-family: Termina;
  font-weight: 700;
  font-size: 36px;
  line-height: 43.2px;

  .width_700 & {
    font-size: 42px;
    line-height: 50px;
  }

  .width_900 & {
    font-size: 60px;
    line-height: 72px;
    /* width: 400px; */
  }
`;

// const CustomCashout = styled(Cashout)`
//   position: absolute;
//   top: 0px;
//   left: 0px;

//   .width_700 & {
//     top: 100px;
//     left: -20px;
//   }
// `;
