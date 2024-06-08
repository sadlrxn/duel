import { useRef, useEffect, useMemo } from 'react';
import styled from 'styled-components';
import gsap from 'gsap';
import { CustomEase } from 'gsap/CustomEase';
import { motion } from 'framer-motion';

import { useAppSelector } from 'state';
import { useSound as useSoundHook } from 'hooks';

import { Candidate, JackpotRoundData } from 'api/types/jackpot';
import { Box, BoxProps } from 'components/Box';
import SpinCard from './Card';

gsap.registerPlugin(CustomEase);

export interface RollProps {
  roundData: JackpotRoundData;
  isGrand?: boolean;
}

export default function Roll({
  roundData,
  isGrand = false,
  ...props
}: RollProps) {
  const rollRef = useRef(null);
  const cardRef = useRef(0);
  const allMeta = useAppSelector(state => state.meta);
  const room = useAppSelector(state => state.jackpot.room);
  const { id: userId } = useAppSelector(state => state.user);
  const { candidates, rolltime, winner, status, time } = roundData;
  const { id: winnerId } = winner;

  const roll = useMemo(() => status === 'rolling', [status]);
  const meta = useMemo(
    () => (isGrand ? allMeta.grandJackpot : allMeta.jackpot[room]),
    [isGrand, allMeta, room]
  );

  const { jackpotPlay } = useSoundHook();

  const candidateData = useMemo(() => {
    return candidates
      .reduce((prev: Candidate[], candidate) => {
        const array = Array(candidate.count).fill(candidate);
        return [...prev, ...array];
      }, [])
      .sort(() => {
        if (Math.random() > 0.5) return 1;
        if (Math.random() < 0.5) return -1;
        return 0;
      });
  }, [candidates]);

  const cardCount = useMemo(() => candidateData.length, [candidateData]);

  const rotation = useMemo(() => {
    if (winnerId === 0) return 0;
    const winner = candidates.find(candidate => candidate.id === winnerId);
    let index = Math.floor(Math.random() * winner!.count!);
    let i = cardCount - 1;
    for (; i >= 0; i--) {
      if (candidateData[i].id === winnerId) index--;
      if (index < 0) break;
    }
    index = (cardCount - i) % cardCount;
    const degree = 360 / cardCount;
    const rotate =
      degree * index + Math.random() * (degree * 0.8) - degree * 0.4;
    return rotate;
  }, [candidates, winnerId, cardCount, candidateData]);

  const radius = useMemo(() => {
    return (cardCount * 120) / (2 * 3.141592) - 72;
  }, [cardCount]);

  const size = useMemo(() => radius + 80, [radius]);

  useEffect(() => {
    if (!roll || cardCount === 0) return;
    if (!rollRef) return;

    const q = gsap.utils.selector(rollRef);

    const animateTime =
      meta.rollingTime - meta.winnerTime - (Date.now() - time) / 1000;

    let progress = (rolltime - animateTime) / rolltime;
    if (progress > 1) progress = 0.98;

    const tl = gsap
      .timeline()
      .set(rollRef.current, {
        rotate: 0,
        width: size * 2,
        height: size * 2
      })
      .set(q('svg'), {
        rotate: -90 + 360 / (cardCount * 2)
      })
      .set(q('.jackpot_card'), {
        transform: i =>
          `translate(-50%, -50%) rotateZ(${
            (360 * i) / cardCount
          }deg) translate(0, -${radius + 20}px)`
      })
      .to(rollRef.current, {
        rotate: `${360 * (isGrand ? 8 : 2) + rotation}`,
        duration: rolltime > 1 ? rolltime - 1 : 0.01,
        ease: CustomEase.create('easeName', '0.36, 0.9, 0.3, 1'),
        delay: 1,
        onStart: () => {
          cardRef.current = 0;
        },
        onUpdate: () => {
          const rotation = Number(gsap.getProperty(rollRef.current, 'rotate'));
          const index = Math.floor(rotation / (180 / cardCount));
          if (index % 2 === 1 && index !== cardRef.current) jackpotPlay.tick();
          cardRef.current = index;
        },
        onComplete: () => {
          if (userId && userId === winnerId) jackpotPlay.win();
          else jackpotPlay.rollend();
        }
      });

    tl.progress(progress);

    return () => {
      tl.kill();
    };
  }, [
    roll,
    cardCount,
    rotation,
    time,
    radius,
    size,
    rolltime,
    userId,
    winnerId,
    jackpotPlay,
    meta,
    isGrand
  ]);

  return (
    <Container
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      {...props}
    >
      <CardWrapper ref={rollRef}>
        <svg width={size * 2} height={size * 2}>
          <circle
            cx={size}
            cy={size}
            r={radius}
            stroke="#17263980"
            strokeWidth={124}
            fill="transparent"
          />
          <circle
            cx={size}
            cy={size}
            r={radius}
            stroke="#1A2A3E"
            strokeWidth={144}
            fill="transparent"
            strokeDasharray={`${(3.141592 * 2 * radius) / cardCount - 3} 3`}
          />
        </svg>
        {candidateData.map((candidate, index) => {
          return (
            <StyledCard
              key={`jackpot_${candidate.name}_${index}`}
              className="jackpot_card"
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

const StyledCard = styled(SpinCard)`
  position: absolute;
  left: 50%;
  top: 50%;
`;

const CardWrapper = styled(Box)`
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
`;

const Container = styled(motion.div)<BoxProps>`
  overflow: hidden;
  max-height: 330px;
`;
