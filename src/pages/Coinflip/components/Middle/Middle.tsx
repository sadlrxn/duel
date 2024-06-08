import React, {
  useCallback,
  useEffect,
  useState,
  useRef,
  useMemo
} from 'react';
import gsap from 'gsap';
import ClipLoader from 'react-spinners/ClipLoader';

import { toast } from 'utils/toast';
import { useSound } from 'hooks';

import {
  Chip,
  Span,
  FlexProps,
  Duel as DuelIcon,
  Ana as AnaIcon
} from 'components';
import { CoinflipGameStatus as Status } from 'api/types/coinflip';
import state, { useAppSelector } from 'state';
import { setRequest } from 'state/coinflip/actions';
import { sendMessage } from 'state/socket';
import { updateBalance } from 'state/user/actions';

import {
  StyledButton,
  Container,
  CounterText,
  Coin,
  Duel,
  Ana,
  Side
} from './Middle.styles';
import { convertBalanceToChip } from 'utils/balance';

interface ButtonProps {
  price: number;
  request?: boolean;
  onClick?: any;
  type?: 'join' | 'cancel';
}

function Button({
  type = 'join',
  price,
  request = false,
  onClick
}: ButtonProps) {
  return (
    <StyledButton
      outlined
      borderColor={type === 'join' ? 'success' : 'warning'}
      color={type === 'join' ? 'success' : 'warning'}
      backgroundColor="transparent"
      onClick={request ? null : onClick}
    >
      {request ? (
        <ClipLoader color="success" loading={request} size={20} />
      ) : (
        <>
          <Span>{type === 'join' ? 'Join' : 'Cancel'}</Span>
          <Chip price={price} fontWeight={700} color="white" />
        </>
      )}
    </StyledButton>
  );
}

interface MiddleProps extends FlexProps {
  creator: boolean;
  roundId: number;
  amount: number;
  status: Status;
  winner: 'duel' | 'ana';
  login?: boolean;
  play?: boolean;
  time: number;
  request?: boolean;
  rowRef: React.MutableRefObject<null>;
  setEnd: React.Dispatch<React.SetStateAction<boolean>>;
}

const Middle: React.FC<MiddleProps> = ({
  creator,
  roundId,
  amount,
  status,
  winner = 'duel',
  login = false,
  play = false,
  time = 0,
  request = false,
  rowRef,
  setEnd,
  ...props
}) => {
  const sides = Array(30).fill('');
  const { coinPlay, countPlay, coinStop, countStop } = useSound();

  const balance = useAppSelector(state => state.user.balance);

  const textRef = useRef(null);
  const coinRef = useRef(null);
  const svgRef = useRef(null);
  const circleRef = useRef(null);
  const [animating, setAnimating] = useState(false);

  const handleRequest = useCallback(() => {
    if (!login) {
      toast.info('Please sign in.');
      return;
    }
    if (!creator && amount > balance) {
      toast.error('Insufficient funds');
      return;
    }
    const content = creator
      ? JSON.stringify([{ eventType: 'cancel', roundId }])
      : JSON.stringify([
          {
            eventType: 'bet',
            roundId,
            amount
          }
        ]);
    if (!creator)
      state.dispatch(updateBalance({ type: -1, usdAmount: amount }));
    state.dispatch(sendMessage({ type: 'event', room: 'coinflip', content }));
    state.dispatch(setRequest({ roundId, status: true }));
  }, [roundId, creator, login, balance, amount]);

  const [color, deg] = useMemo(
    () => [
      winner === 'ana' ? '#422554' : '#25544D',
      winner === 'ana' ? -90 : 90
    ],
    [winner]
  );

  const animation = useCallback(
    (time: number) => {
      if (!coinRef || !textRef) return gsap.timeline();

      const firstAngle = 720; //720
      const secondAngle = 1440; //1440

      const rotateValue = Math.random() * 60 - 30;
      const result = winner === 'ana';

      let count = { value: 0 };
      let cur = 0;

      const tl = gsap
        .timeline()
        .set(coinRef.current, { display: 'none' })
        .set([svgRef.current, textRef.current], { display: 'block' })
        .add('countdown')
        .fromTo(
          count,
          {
            value: 3.5,
            onStart: () => {
              cur = 3;
            }
          },
          {
            value: 0.5,
            duration: 3,
            ease: 'none',
            roundProps: 'value',
            onUpdate: () => {
              if (Math.floor(count.value) === cur) return;
              cur = Math.floor(count.value);
              gsap.fromTo(
                textRef.current,
                {
                  textContent: `${cur}`,
                  scale: 1.4
                },
                {
                  scale: 1,
                  duration: 0.4,
                  onStart: () => {
                    play &&
                      time < (3 - cur) * 1 + 0.2 &&
                      //@ts-ignore
                      countPlay[`c${cur}`]();
                  }
                }
              );
            }
          }
        )
        .fromTo(
          circleRef.current,
          { strokeDashoffset: 0 },
          {
            strokeDashoffset: 2 * 3.14 * 19,
            ease: 'none',
            duration: 3
          },
          '<'
        )
        .set(
          [svgRef.current, textRef.current],
          { display: 'none' },
          'countdown+=3'
        )
        .add('flip')
        .fromTo(
          coinRef.current,
          { display: 'block', rotationY: '0deg', scale: 0.8 },
          {
            rotationY: `${firstAngle + (Math.random() < 0.5 ? 180 : 0) - 30}`,
            ease: 'ease-in',
            rotation: rotateValue / 2,
            scale: 1.2,
            duration: 1,
            onStart: () => {
              setAnimating(true);
              play && time < 3.2 && coinPlay.start && coinPlay.start();
            }
          }
        )
        .to(coinRef.current, {
          rotationY: '+=60',
          ease: 'ease-in',
          duration: 1
        })
        .add('result')
        .to(coinRef.current, {
          rotationY: `${secondAngle + (result ? 180 : 0)}`,
          scale: 1,
          rotation: rotateValue,
          ease: 'ease-out',
          duration: 0.5,
          onComplete: () => {
            play && time < 5.7 && coinPlay && coinPlay.end();
            setAnimating(false);
          }
        })
        .set(
          rowRef.current,
          {
            background: `linear-gradient(${deg}deg, ${color}00 30%, ${color}00 50%, ${color}00 70%)`
          },
          'result+=0.1'
        )
        .to(
          rowRef.current,
          {
            background: `linear-gradient(${deg}deg, ${color}00 30%, ${color} 50%, ${color}00 70%)`,
            duration: 0.3
          },
          'result+=0.2'
        )
        .to(
          rowRef.current,
          {
            background: `linear-gradient(${deg}deg, ${color}00 -50%, ${color}, ${color}00 50%)`,
            duration: 0.4,
            onStart: () => {
              setEnd(true);
            }
          },
          '>=+0.2'
        );
      return tl;
    },
    [rowRef, setEnd, winner, color, deg, coinPlay, countPlay, play]
  );

  useEffect(() => {
    let tl: gsap.core.Timeline | null = null;
    switch (status) {
      case 'created':
        if (!coinRef || !textRef || !svgRef) break;
        gsap.set([coinRef.current, textRef.current, svgRef.current], {
          display: 'none'
        });
        break;
      case 'joined':
        let progress = (Date.now() - time) / (6.1 * 1000);
        if (progress < 0) progress = 0;
        if (progress > 1) progress = 0.9999;
        tl = animation(progress * 6.1);
        tl.progress(progress);
        break;
      case 'ended':
        if (!coinRef || !rowRef || !svgRef) break;
        setEnd(true);
        gsap.set([textRef.current, svgRef.current], { display: 'none' });
        gsap.set(coinRef.current, {
          rotationY: `${winner === 'ana' ? 180 : 0}deg`
        });
        gsap.set(rowRef.current, {
          background: `linear-gradient(${deg}deg, ${color}00 -50%, ${color}, ${color}00 50%)`
        });
        break;
    }
    return () => {
      tl && tl.kill();
      coinStop.end();
      coinStop.start();
      countStop.c1();
      countStop.c2();
      countStop.c3();
    };
  }, [
    animation,
    color,
    deg,
    rowRef,
    setEnd,
    status,
    time,
    winner,
    coinStop,
    countStop
  ]);

  return (
    <Container {...props}>
      {status === 'created' && (
        <Button
          type={creator ? 'cancel' : 'join'}
          price={convertBalanceToChip(amount)}
          request={request}
          onClick={handleRequest}
        />
      )}
      <svg ref={svgRef} width={40} height={40}>
        <circle
          cx={20}
          cy={20}
          r={19}
          stroke="#556680"
          strokeWidth={2}
          fill="transparent"
        />
        <circle
          ref={circleRef}
          cx={20}
          cy={20}
          r={19}
          stroke="#4FFF8B"
          strokeWidth={2}
          fill="transparent"
          strokeDasharray={19 * 3.14 * 2}
          strokeDashoffset={0}
          strokeLinecap="round"
          transform="scale(-1, 1) rotate(-90) translate(-40, -40)"
        />
      </svg>
      <CounterText ref={textRef} />
      <Coin ref={coinRef}>
        <Duel>
          <DuelIcon size={28} />
        </Duel>
        <Ana>
          <AnaIcon size={28} />
        </Ana>
        {animating === true &&
          sides.map((side, index) => {
            return <Side key={`${side}_${index}`} index={index} />;
          })}
      </Coin>
    </Container>
  );
};

export default React.memo(Middle);
