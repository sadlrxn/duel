import React, { useRef, useMemo, useEffect } from 'react';
import styled from 'styled-components';
import gsap from 'gsap';
import { CustomEase } from 'gsap/CustomEase';
import { motion } from 'framer-motion';

import { Box, BoxProps } from 'components/Box';
import { JackpotRoundData } from 'api/types/jackpot';
import { useAppSelector } from 'state';
import { useSound } from 'hooks';

import Card from './Card';
import { generateCandidateData } from '../util';

gsap.registerPlugin(CustomEase);

const CARD_COUNT = 20;
const RADIUS = Math.floor((CARD_COUNT * 210) / (2 * 3.141592));

export interface RollProps {
  roundData: JackpotRoundData;
  isHistory?: boolean;
  isGrand?: boolean;
}

export default function Roll({
  roundData,
  isHistory = false,
  isGrand = false,
  ...props
}: RollProps) {
  const rollRef = useRef(null);
  const cardRef = useRef(0);
  const allMeta = useAppSelector(state => state.meta);
  const room = useAppSelector(state => state.jackpot.room);
  const { id: userId } = useAppSelector(state => state.user);
  const { candidates, rolltime, winner, status, time, roundId } = roundData;
  const { id: winnerId } = winner;

  const roll = useMemo(() => status === 'rolling', [status]);
  const meta = useMemo(
    () => (isGrand ? allMeta.grandJackpot : allMeta.jackpot[room]),
    [isGrand, allMeta, room]
  );

  const { jackpotPlay } = useSound();

  const { candidateData, cardCount, rotation } = useMemo(
    () => generateCandidateData(candidates, winnerId, roundId),
    [candidates, winnerId, roundId]
  );

  useEffect(() => {
    if (!roll || cardCount === 0) return;
    if (!rollRef) return;

    const q = gsap.utils.selector(rollRef);

    const animateTime =
      (isHistory ? 15 : meta.rollingTime - meta.winnerTime) -
      (Date.now() - time) / 1000;

    let progress = (rolltime - animateTime) / rolltime;
    if (progress > 1) progress = 0.99;

    const tl = gsap
      .timeline()
      .set(rollRef.current, { rotationY: 0 })
      .set(q('.grand_jackpot_card'), {
        rotateY: i => i * (360 / CARD_COUNT),
        transformOrigin: `50% 50% ${RADIUS}px`,
        z: -RADIUS,
        backfaceVisibility: 'hidden',
        opacity: i => (i < 14 || i >= 100 - 6 ? 1 : 0)
      })
      .to(rollRef.current, {
        rotateY: `-${360 * 5 * (isHistory ? 1 : 4) + rotation}`,
        duration: rolltime > 1.5 ? rolltime - 1.5 : 0.01,
        ease: CustomEase.create('easeName', '0.36, 0.9, 0.3, 1'),
        delay: 1,
        onStart: () => {
          cardRef.current = 0;
        },
        onUpdate: () => {
          const rotation = -Number(
            gsap.getProperty(rollRef.current, 'rotateY')
          );
          const index =
            Math.floor(rotation / (180 / CARD_COUNT) + 1) % (cardCount * 2);
          if (index % 2 === 0 && index !== cardRef.current) jackpotPlay.tick();
          cardRef.current = index;
          let min = Math.floor(index / 2) - 6;
          let max = min + CARD_COUNT;
          if (min < 0) min += cardCount;
          if (max >= cardCount) max -= cardCount;
          gsap.set(q('.grand_jackpot_card'), {
            opacity: i => {
              return (max > min ? i >= min && i < max : i >= min || i < max)
                ? i === Math.floor(index / 2)
                  ? 1
                  : 0.7
                : 0;
            }
          });
        }
      })
      .add(() => {
        if (userId && userId === winnerId) jackpotPlay.win();
        else jackpotPlay.rollend();
      }, '+=0.5');

    tl.progress(progress);

    return () => {
      tl.kill();
    };
  }, [
    cardCount,
    jackpotPlay,
    meta,
    roll,
    rolltime,
    rotation,
    time,
    userId,
    winnerId,
    isHistory
  ]);

  return (
    <Container
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      {...props}
    >
      <CardWrapper ref={rollRef}>
        {candidateData.map((candidate, index) => {
          return (
            <Card
              key={`grand_jackpot_${candidate.name}_${index}`}
              className="grand_jackpot_card"
              name={candidate.name}
              avatar={candidate.avatar}
              percent={candidate.percent}
              count={cardCount}
            />
          );
        })}
      </CardWrapper>
    </Container>
  );
}

const CardWrapper = styled(Box)`
  width: 100%;
  height: 100%;
  position: absolute;

  transform-style: preserve-3d;
  user-select: none;
`;

const Container = styled(motion.div)<BoxProps>`
  width: 200px;
  height: 300px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  perspective: 800px;

  transform-style: preserve-3d;
  user-select: none;
`;
