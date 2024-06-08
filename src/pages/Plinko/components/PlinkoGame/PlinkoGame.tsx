import React, { useMemo, useRef } from 'react';
import styled from 'styled-components';

import blueGlowImg from 'assets/imgs/plinko/blue.png';
import lightBlueGlowImg from 'assets/imgs/plinko/lightblue.png';
import purpleGlowImg from 'assets/imgs/plinko/purple.png';
import pinkGlowImg from 'assets/imgs/plinko/pink.png';
import orangeGlowImg from 'assets/imgs/plinko/orange.png';
import yellowGlowImg from 'assets/imgs/plinko/yellow.png';

import { config } from 'pages/Plinko/config';
import { getMultipliers } from 'pages/Plinko/utils';
import { usePlinko } from 'hooks';

import Background from './Background';
import Defs from './Defs';
import { useAppSelector } from 'state';
import { Ball } from '../Ball';

const glowImgs = [
  lightBlueGlowImg,
  lightBlueGlowImg,
  blueGlowImg,
  purpleGlowImg,
  pinkGlowImg,
  orangeGlowImg,
  yellowGlowImg
];

export default function PlinkoGame() {
  const { gameRows, gameMode } = usePlinko();
  const { balls } = useAppSelector(state => state.plinko);

  const svgRef = useRef<SVGSVGElement>(null);

  const [width, height] = useMemo(() => {
    const width = config.pinGap * (4 + gameRows);
    const height =
      config.startTop +
      config.pinGap * (gameRows - 1) +
      config.multiHeight * 1.75 +
      10;
    return [width, height];
  }, [gameRows]);

  const path = useMemo(() => {
    const height = config.startTop + config.pinGap * gameRows;
    const topLength = config.pinGap * 3 - config.startTop; //pinGap * 2 + pinGap / 2 * 2 - startTop
    const bottomLength = topLength + height;
    const borderRadius = 8;

    const path: string[] = [];
    path.push(`M${-topLength / 2 + borderRadius} 8`);
    path.push(`h${topLength - borderRadius * 2}`);
    path.push(`q${borderRadius} 0 ${borderRadius * 1.5} ${borderRadius}`);
    path.push(`l${height / 2 - borderRadius} ${height - borderRadius * 2}`);
    path.push(
      `q${borderRadius / 2} ${borderRadius} -${
        borderRadius / 2
      } ${borderRadius}`
    );
    path.push(`h -${bottomLength - borderRadius * 2}`);
    path.push(`q -${borderRadius} 0 -${borderRadius / 2} -${borderRadius}`);
    path.push(`l ${height / 2 - borderRadius} -${height - borderRadius * 2}`);
    path.push(
      `q${borderRadius / 2} -${borderRadius} ${
        borderRadius * 1.5
      } -${borderRadius}`
    );
    path.push('z');

    return path.join(' ');
  }, [gameRows]);

  const pins = useMemo(() => {
    let pins: any[] = [];

    for (let i = 0; i < gameRows; i++) {
      const linePins = config.startPins + i;

      for (let j = 0; j < linePins; j++) {
        const y = 8 + config.startTop + config.pinSize + i * config.pinGap;
        const x = -((linePins - 1) / 2 - j) * config.pinGap;

        pins.push(
          <>
            <circle
              className={'pins-' + i + '-' + j + '-bounce'}
              cx={x}
              cy={y}
              r={config.pinSize}
              fill="url(#pin-bounce-gradient)"
              stroke="none"
              key={'pin_' + i + '_' + j + '_bounce'}
            />
            <circle
              className={'pins-' + i + '-' + j}
              cx={x}
              cy={y}
              r={config.pinSize}
              fill="#9A7BB1B2"
              filter="url(#pin-shadow-1)"
              stroke="none"
              key={'pin_' + i + '_' + j}
            />
          </>
        );
      }
    }

    return pins;
  }, [gameRows]);

  const multipliers = useMemo(() => {
    const multipliers = getMultipliers(gameMode, gameRows).map(
      (item, index) => {
        const { multiplier, gradient } = item;
        const x =
          -((gameRows + 1) / 2 - index) * config.pinGap +
          (config.pinGap - config.multiWidth) / 2;
        const y = height - config.multiHeight - 10;
        return (
          <g
            key={multiplier + '_' + gradient + '_' + index}
            transform={`translate(${x} ${y})`}
          >
            <image
              href={glowImgs[gradient]}
              pointerEvents="none"
              width={config.multiWidth}
              y={-45}
              className={`multiplier_${index}_image`}
              opacity={0}
            />
            <g className={`multiplier_${index}`}>
              <rect
                width={config.multiWidth}
                height={config.multiHeight}
                rx={4}
                fill={`url(#multiplier-gradient-${gradient})`}
              />
              <rect
                width={config.multiWidth}
                height={config.multiHeight}
                rx={4}
                fill="white"
                fillOpacity={0.2}
              />
              <text
                fontSize="11px"
                fontWeight={700}
                fontFamily="Inter"
                x={config.multiWidth / 2}
                y={15}
                textAnchor="middle"
                fill="white"
              >
                {multiplier}
              </text>
              <text
                fontSize="11px"
                fontWeight={700}
                fontFamily="Inter"
                x={config.multiWidth / 2}
                y={26}
                textAnchor="middle"
                fill="white"
              >
                x
              </text>
            </g>
          </g>
        );
      }
    );

    return multipliers;
  }, [gameMode, gameRows, height]);

  return (
    <Svg
      viewBox={`${-width / 2} 0 ${width} ${height}`}
      width="100%"
      ref={svgRef}
    >
      <Background path={path} />
      {balls.map(ball => {
        return (
          <Ball
            svgRef={svgRef}
            ball={ball}
            key={'plinko_ball_' + ball.roundId}
          />
        );
      })}
      <g>{pins}</g>
      <g>{multipliers}</g>
      <Defs />
    </Svg>
  );
}

const Svg = styled.svg``;
