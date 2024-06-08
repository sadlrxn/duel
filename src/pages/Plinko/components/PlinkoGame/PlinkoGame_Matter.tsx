import { useEffect, useMemo, useState, useCallback } from 'react';
import {
  Engine,
  Events,
  World,
  Runner,
  Render,
  Body,
  Bodies,
  Composite,
  IEventCollision
} from 'matter-js';
import styled from 'styled-components';
import gsap from 'gsap';

import { Box, Flex } from 'components';
import state, { useAppSelector } from 'state';
import * as plinkoActions from 'state/plinko/actions';
import { usePlinko } from 'hooks';

import { config, multipliers as totalMultipliers } from '../../config';

const random = (min: number, max: number) => {
  return min + Math.random() * (max - min);
};

export default function PlinkoGame() {
  const { balls } = useAppSelector(state => state.plinko);
  const { gameRows: lines, gameMode: level } = usePlinko();

  const [_, setLastMultipliers] = useState<number[]>([]);
  const [ballIdsInGame, setBallIdsInGame] = useState<number[]>([]);

  const engine = useMemo(() => Engine.create(), []);

  const [width, height] = useMemo(() => {
    const width = config.pinGap * (3 + lines);
    const height =
      config.startTop + config.pinGap * (lines - 1) + config.multiHeight + 14;
    return [width, height];
  }, [lines]);

  const multipliers = useMemo(() => {
    //@ts-ignore
    return totalMultipliers[`${level.toLowerCase()}_${lines}`];
  }, [lines, level]);

  useEffect(() => {
    const element = document.getElementById('plinko_game_body');

    engine.gravity.y = config.engineGravity;

    const render = Render.create({
      element: element!,
      engine,
      bounds: {
        max: {
          x: width,
          y: height
        },
        min: {
          x: 0,
          y: 0
        }
      },
      options: {
        background: 'transparent',
        hasBounds: true,
        width: width,
        height: height,
        wireframes: false
      }
    });

    const runner = Runner.create();
    Runner.run(runner, engine);
    Render.run(render);
    return () => {
      World.clear(engine.world, true);
      Engine.clear(engine);
      render.canvas.remove();
      render.textures = {};
    };
  }, [engine, height, width]);

  const pins = useMemo(() => {
    const pins: Body[] = [];

    for (let l = 0; l < lines; l++) {
      const linePins = config.startPins + l;
      const lineWidth = linePins * config.pinGap;

      for (let i = 0; i < linePins; i++) {
        const pinX =
          width / 2 - lineWidth / 2 + i * config.pinGap + config.pinGap / 2;

        const pinY = config.startTop + config.pinSize / 2 + config.pinGap * l;

        const pin = Bodies.circle(pinX, pinY, config.pinSize, {
          label: `pin-${l}-${i}`,
          render: {
            fillStyle: '#F5DCFF'
          },
          isStatic: true,
          restitution: 0.6,
          friction: 0.8
        });

        pins.push(pin);
      }
    }

    return pins;
  }, [lines, width]);

  const leftWall = useMemo(
    () =>
      Bodies.rectangle(-3, height, width * 3, 2, {
        angle: 90,
        render: {
          visible: true
        },
        isStatic: true
      }),
    [height, width]
  );

  const rightWall = useMemo(
    () =>
      Bodies.rectangle(width + 3, height, width * 3, 2, {
        angle: -90,
        render: {
          visible: true
        },
        isStatic: true
      }),
    [height, width]
  );

  const floor = useMemo(
    () =>
      Bodies.rectangle(width / 2, height + 100, width * 2, 20, {
        label: 'multiplier-1',
        render: {
          visible: false
        },
        isStatic: true
      }),
    [height, width]
  );

  const multipliersBodies = useMemo(() => {
    const multipliersBodies: Body[] = [];

    let lastMultiplierX = width / 2 - (config.pinGap * lines) / 2;

    multipliers.forEach((multiplier: any, index: number) => {
      const multiplierBody = Bodies.rectangle(
        lastMultiplierX,
        height - config.multiHeight / 2,
        config.multiWidth,
        config.multiHeight,
        {
          label: `multiplier-${multiplier}-${index}`,
          isStatic: true,
          chamfer: {
            radius: [5, 5, 5, 5]
          },
          render: {
            fillStyle: '#182738',
            lineWidth: 0
          }
        }
      );

      lastMultiplierX += config.pinGap;
      multipliersBodies.push(multiplierBody);
    });

    return multipliersBodies;
  }, [height, lines, multipliers, width]);

  const onCollideWithMultiplier = useCallback(
    (ball: Body, multiplier: Body) => {
      ball.collisionFilter.group = 2;
      setBallIdsInGame(prev => {
        const index = prev.findIndex(id => id === ball.id);
        if (index === -1) return prev;
        prev.splice(index, 1);
        return prev;
      });
      World.remove(engine.world, ball);
      state.dispatch(plinkoActions.removeBall(ball.id));

      const ballValue = ball.label.split('-')[1];
      const multiplierValue = multiplier.label.split(
        '-'
      )[1] as unknown as number;

      setLastMultipliers(prev => [multiplierValue, prev[0], prev[1], prev[2]]);

      if (+ballValue < 0) return;

      const target = multiplier;
      const posY = target.position.y;

      const posValue = {
        value: 0
      };

      gsap
        .timeline({
          onUpdate: () => {
            const val = posValue.value / 10 ** 3;

            Body.setPosition(target, {
              x: target.position.x,
              y: posY + (val * config.multiHeight) / 6
            });
            // target.position.y = posY + (val * multiplierConfig.height) / 10;
          }
        })
        .set(posValue, { value: 0 })
        .to(posValue, {
          value: 10 ** 3,
          roundProps: 'value',
          duration: 0.1
        })
        .to(posValue, {
          value: 0,
          roundProps: 'value',
          duration: 0.1
        });

      // const newBalance = +ballValue * multiplierValue;
    },
    [engine.world]
  );

  const onCollideWithPin = useCallback(
    (ball: Body, pin: Body) => {
      const newPin = Bodies.circle(
        pin.position.x,
        pin.position.y,
        pin.circleRadius!,
        {
          isStatic: true,
          collisionFilter: {
            group: -1,
            category: 2,
            mask: 0
          },
          render: {
            fillStyle: '#fff'
          }
        }
      );

      World.add(engine.world, newPin);

      const scale = {
        value: 0
      };

      gsap.fromTo(
        scale,
        { value: 0 },
        {
          value: 10 ** 6,
          duration: 0.5,
          roundProps: 'value',
          onUpdate: () => {
            const val = scale.value / 10 ** 6;
            Body.scale(newPin, 1.05, 1.05); // 1.05
            newPin.render.opacity = 1 - val;
          },
          onComplete: () => {
            World.remove(engine.world, newPin);
          }
        }
      );
    },
    [engine.world]
  );

  const onBodyCollision = useCallback(
    (event: IEventCollision<Engine>) => {
      const pairs = event.pairs;
      for (const pair of pairs) {
        const { bodyA, bodyB } = pair;
        if (
          bodyA.label.includes('ball') &&
          bodyB.label.includes('multiplier')
        ) {
          onCollideWithMultiplier(bodyA, bodyB);
        } else if (
          bodyB.label.includes('ball') &&
          bodyA.label.includes('multiplier')
        ) {
          onCollideWithMultiplier(bodyB, bodyA);
        }

        if (bodyA.label.includes('ball') && bodyB.label.includes('pin')) {
          onCollideWithPin(bodyA, bodyB);
        } else if (
          bodyB.label.includes('ball') &&
          bodyA.label.includes('pin')
        ) {
          onCollideWithPin(bodyB, bodyA);
        }
      }
    },
    [onCollideWithMultiplier, onCollideWithPin]
  );

  useEffect(() => {
    const composite = Composite.add(engine.world, [
      ...pins,
      ...multipliersBodies,
      leftWall,
      rightWall,
      floor
    ]);

    return () => {
      Composite.clear(composite, false);
    };
  }, [engine.world, floor, leftWall, pins, rightWall, multipliersBodies]);

  useEffect(() => {
    Events.on(engine, 'collisionStart', onBodyCollision);

    return () => {
      Events.off(engine, 'collisionStart', onBodyCollision);
    };
  }, [engine, onBodyCollision]);

  const addBall = useCallback(
    (id: number, ballValue: number) => {
      const minBallX = width / 2 - config.pinSize * 3 + config.pinGap;
      const maxBallX =
        width / 2 - config.pinSize * 3 - config.pinGap + config.pinGap / 2;

      const ballX = random(minBallX, maxBallX);
      const ballColor = '#00ff00';
      const ball = Bodies.circle(ballX, 0, config.ballSize, {
        restitution: 0.6,
        friction: 0.8,
        label: `ball-${ballValue}`,
        id,
        // frictionAir: 0.05,
        density: 1,
        collisionFilter: {
          group: -1
        },
        render: {
          fillStyle: ballColor
        },
        isStatic: false
      });

      setBallIdsInGame(prev => {
        return [...prev, id];
      });
      Composite.add(engine.world, ball);
    },
    [engine.world, width]
  );

  useEffect(() => {
    balls.forEach(ball => {
      if (ballIdsInGame.findIndex(id => id === ball.roundId) === -1)
        addBall(ball.roundId, ball.betAmount);
    });
  }, [addBall, ballIdsInGame, balls]);

  return (
    <Flex width="100%" justifyContent="center" height="100%">
      <Container width={width} height={height}>
        <CanvasWrapper id="plinko_game_body" />
      </Container>
    </Flex>
  );
}

const Container = styled(Flex)`
  flex-direction: column;
  align-items: center;
`;

const CanvasWrapper = styled(Box)`
  display: flex;
  justify-content: center;
  width: 100%;

  canvas {
    width: 100%;
    max-width: 760px;
  }
`;
