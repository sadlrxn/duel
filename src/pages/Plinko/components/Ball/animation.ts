import gsap, { Power2 } from 'gsap';

import { config } from '../../config';

export const ballStarAnim = (target: gsap.TweenTarget) => {
  const duration = 0.2 + 0.05 * (Math.random() - 0.5);

  return gsap.timeline().fromTo(
    target,
    {
      x: 0,
      y: 8 - config.ballSize
    },
    {
      y: 8 + config.startTop - config.ballSize - config.pinSize,
      duration,
      ease: Power2.easeIn
    }
  );
};

export const ballBounceAnim = (
  parent: gsap.TweenTarget,
  target: gsap.TweenTarget,
  i: number,
  j: number,
  left: boolean
) => {
  const linePins = config.startPins + i;
  const x = -config.pinGap * ((linePins - 1) / 2 - j);
  const y = config.startTop + config.pinGap * i;

  const dx = (config.pinGap / 2) * (left ? -1 : 1);

  const ran = 0.3 + 0.2 * Math.random();
  const dy = (config.pinGap - config.pinSize * 2) * ran;

  const duration = 0.2 + 0.05 * (Math.random() - 0.5);

  const q = gsap.utils.selector(parent);

  return gsap
    .timeline()
    .set(target, { x, y })
    .to(target, {
      onStart: () => {
        gsap.set(q(`.pins-${i}-${j}`), {
          fill: '#FFFFFF',
          filter: 'url(#pin-shadow-2)'
        });
        gsap.fromTo(
          q(`.pins-${i}-${j}-bounce`),
          {
            transformOrigin: 'center center',
            scale: 1,
            opacity: 1
          },
          {
            scale: 3,
            opacity: 0,
            duration: 0.2,
            onComplete: () => {
              gsap.set(q(`.pins-${i}-${j}-bounce`), { scale: 1 });
              gsap.set(q(`.pins-${i}-${j}`), {
                fill: '#9A7BB1B2',
                filter: 'url(#pin-shadow-1)'
              });
            }
          }
        );
      },
      x: x + dx * 0.4,
      duration,
      ease: 'none'
    })
    .to(
      target,
      {
        y: y - dy,
        duration,
        ease: Power2.easeOut
      },
      '<'
    )
    .to(target, {
      x: x + dx,
      duration: duration * 1.5,
      ease: 'none'
    })
    .to(
      target,
      {
        y: y + config.pinGap,
        duration: duration * 1.5,
        ease: Power2.easeIn
      },
      '<'
    );
};
