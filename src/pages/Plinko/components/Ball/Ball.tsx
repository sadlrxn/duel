import React, { useEffect, useRef } from 'react';
import gsap from 'gsap';

import { PlinkoBall } from 'api/types/plinko';
import { useAppDispatch } from 'state';
import { removeBall } from 'state/plinko/actions';
import { usePlinko } from 'hooks';

import { ballBounceAnim, ballStarAnim } from './animation';

interface BallProps {
  ball: PlinkoBall;
  svgRef: React.RefObject<SVGSVGElement>;
}

export default function Ball({ svgRef, ball }: BallProps) {
  const dispatch = useAppDispatch();
  const { animSpeed } = usePlinko();

  const ballRef = useRef<SVGGElement>(null);
  const progress = useRef(0);

  useEffect(() => {
    if (!ballRef || !ballRef.current) return;
    if (!svgRef || !svgRef.current) return;
    const q = gsap.utils.selector(svgRef.current);
    const tl = gsap.timeline();
    let j = 1;

    tl.add(ballStarAnim(ballRef.current));

    for (let i = 0; i < ball.path.length; i++) {
      const left = ball.path[i] === 'L';
      tl.add(ballBounceAnim(svgRef.current, ballRef.current, i, j, left));
      j += left ? 0 : 1;
    }

    tl.addLabel('ball_anim_end')
      .fromTo(
        q(`.multiplier_${j - 1}_image`),
        { opacity: 0 },
        { opacity: 1, duration: 0.1 }
      )
      .to(
        q(`.multiplier_${j - 1}_image`),
        { opacity: 0, duration: 0.1 },
        '>+=0.2'
      )
      .fromTo(
        q(`.multiplier_${j - 1}`),
        { y: 0 },
        {
          y: 3,
          duration: 0.05
        },
        'ball_anim_end'
      )
      .to(q(`.multiplier_${j - 1}`), { y: 0, duration: 0.05 });

    tl.call(() => {
      //After Animation End
      dispatch(removeBall(ball.roundId));
    });

    const duration = tl.duration();

    tl.duration(duration * (1 / animSpeed)).progress(progress.current);

    return () => {
      progress.current = tl.progress();
      tl.kill();
    };
  }, [animSpeed, ball.path, ball.roundId, dispatch, svgRef]);

  return (
    <>
      <g
        className="plinko-ball"
        ref={ballRef}
        fill="url(#ball-gradient-1)"
        filter="url(#ball-shadow)"
      >
        <circle
          cx="0"
          cy="0"
          r="8"
          style={{
            boxShadow: '0px 0px 4px 1px rgba(255, 77, 1, 0.5)'
          }}
        />
        <circle cx="0" cy="0" r="7.5" stroke="white" strokeOpacity="0.15" />
      </g>
    </>
  );
}
