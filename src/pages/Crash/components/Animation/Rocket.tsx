import { useCallback, useState, useRef, useEffect } from 'react';
import { Player, PlayerEvent } from '@lottiefiles/react-lottie-player';

import { Box } from 'components';
import Rocket3D from './Rocket3D';

export default function Rocket({
  visible = false,
  explosion = false,
  angle = 90,
  speed = 1,
  ...props
}: any) {
  const [show, setShow] = useState(true);

  const explosionRef = useRef<any>(null);

  const handleEvent = useCallback((event: PlayerEvent) => {
    if (event === 'loop') {
      setShow(false);
      explosionRef.current.stop();
    } else if (event === 'ready') {
      explosionRef.current.stop();
    }
  }, []);

  useEffect(() => {
    if (visible && explosion) {
      setShow(true);
      explosionRef.current.play();
    } else {
      setShow(false);
    }
  }, [explosion, visible]);

  return (
    <Box
      {...props}
      style={{
        pointerEvents: 'none',
        transformOrigin: '50% 0px',
        transform: ` translateX(-50%) rotate(${90 - angle}deg) translateY(-22%)`
      }}
    >
      <Player
        ref={explosionRef}
        src="/assets/Explosion.json"
        loop
        onEvent={handleEvent}
        style={{
          position: 'absolute',
          width: '100%',
          left: '50%',
          top: '50%',
          transform: 'translate(-50%, -60%) scale(2.5)',
          opacity: show ? 1 : 0
        }}
      />
      <Rocket3D
        speed={speed}
        style={{
          opacity: visible && !explosion ? 1 : 0,
          pointerEvents: 'none'
        }}
      />
      {/* <Player
        src="/assets/Rocket.json"
        loop
        autoplay
        speed={speed * 4}
        style={{
          position: 'absolute',
          left: '50%',
          top: '50%',
          transform: 'translate(-50%, -50%)',
          width: '100%',
          height: '100%',
          opacity: visible && !explosion ? 1 : 0
        }}
      /> */}
    </Box>
  );
}
