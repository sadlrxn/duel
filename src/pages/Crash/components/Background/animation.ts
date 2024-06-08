import React from 'react';
import gsap from 'gsap';

export const backgroundAnim = (ref: React.RefObject<HTMLDivElement>) => {
  if (!ref || !ref.current) return gsap.timeline({ paused: true });
  return gsap
    .timeline({ paused: true })
    .set(ref.current, {
      background: 'linear-gradient(0deg, #2875ad 0%, #00315D 100% )'
    })
    .to(
      ref.current,
      {
        background: 'linear-gradient(0deg, #0459A5 0%, #001120 100% )',
        duration: 10,
        ease: 'none'
      },
      '3'
    )
    .to(ref.current, {
      background: 'linear-gradient(0deg, #230060 0%, #000306 100% )',
      duration: 10,
      ease: 'none'
    })
    .to(ref.current, {
      background: 'linear-gradient(0deg, #00021C 0%, #000000 100% )',
      duration: 7,
      ease: 'none'
    })
    .to(ref.current, {
      background: 'linear-gradient(0deg, #000000 0%, #000000 100% )',
      duration: 5,
      ease: 'none'
    });
};

export const starAnimation = (ref: React.RefObject<HTMLDivElement>) => {
  if (!ref || !ref.current) return gsap.timeline({ paused: true });
  return gsap
    .timeline({ paused: true })
    .set(ref.current, {
      opacity: 0.2
    })
    .to(
      ref.current,
      {
        opacity: 0.5,
        duration: 10,
        ease: 'none'
      },
      '3'
    )
    .to(ref.current, {
      opacity: 0.8,
      duration: 10,
      ease: 'none'
    })
    .to(ref.current, {
      opacity: 1,
      duration: 7,
      ease: 'none'
    });
};

export const saturnAnimation = (
  ref: React.RefObject<HTMLImageElement>,
  onComplete?: any
) => {
  if (!ref || !ref.current) return gsap.timeline({ paused: true });
  return gsap
    .timeline({ paused: true })
    .set(ref.current, {
      opacity: 1
    })
    .fromTo(
      ref.current,
      {
        left: `${-ref.current.clientWidth - 100}px`,
        top: `${-ref.current.clientHeight / 2}px`
      },
      {
        left: `calc(100% + ${ref.current.clientWidth + 100}px)`,
        top: `${25 + Math.random() * 10}%`,
        ease: 'none',
        duration: 30,
        onComplete
      }
    );
};

export const satelliteAnimation = (
  ref: React.RefObject<HTMLImageElement>,
  onComplete?: any
) => {
  if (!ref || !ref.current) return gsap.timeline({ paused: true });

  const y = 60 + Math.random() * 30;
  let start: number | string = `${-ref.current.clientWidth - 200}px`,
    end: number | string = `calc(100% + ${ref.current.clientWidth + 200}px)`;

  // if (Math.random() > 0.5) {
  let temp = start;
  start = end;
  end = temp;
  // }

  return gsap
    .timeline({ paused: true })
    .set(ref.current, {
      opacity: 1
    })
    .fromTo(
      ref.current,
      {
        left: start,
        top: `${y}%`,
        rotate: `${Math.random() * 180}`
      },
      {
        left: end,
        ease: 'none',
        duration: 5,
        rotate: `+=${60 + Math.random() * 60}`,
        onComplete
      }
    );
};
