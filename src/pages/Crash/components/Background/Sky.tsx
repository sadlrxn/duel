import { useRef, useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';
import gsap from 'gsap';

import satelliteImg from 'assets/imgs/crash/satellite.png';

import planet1Img from 'assets/imgs/crash/planet1.png';
import planet2Img from 'assets/imgs/crash/planet2.png';
import planet3Img from 'assets/imgs/crash/planet3.png';
import planet4Img from 'assets/imgs/crash/planet4.png';
import planet5Img from 'assets/imgs/crash/planet5.png';
import planet6Img from 'assets/imgs/crash/planet6.png';
import planet7Img from 'assets/imgs/crash/planet7.png';
import planet8Img from 'assets/imgs/crash/planet8.png';
import planet9Img from 'assets/imgs/crash/planet9.png';
import planet10Img from 'assets/imgs/crash/planet10.png';

import { CRASH_BACK_TIME } from 'pages/Crash/config';
import { Box } from 'components';
import { useAppSelector } from 'state';

import { saturnAnimation, satelliteAnimation } from './animation';

const imgs = {
  planet1: planet1Img,
  planet2: planet2Img,
  planet3: planet3Img,
  planet4: planet4Img,
  planet5: planet5Img,
  planet6: planet6Img,
  planet7: planet7Img,
  planet8: planet8Img,
  planet9: planet9Img,
  planet10: planet10Img
};

export default function Sky() {
  const status = useAppSelector(state => state.crash.status);
  const globalTime = useAppSelector(state => state.crash.time);
  const usingAnimation = useAppSelector(state => state.user.crashAnimation);

  const [image, setImage] = useState(imgs['planet1']);

  const saturnRef = useRef<HTMLImageElement>(null);
  const satelliteRef = useRef<HTMLImageElement>(null);

  const saturnAnimEnd = useRef(0);
  const satelliteAnimEnd = useRef(0);

  const nextSaturnAnim = useRef(40);
  const nextSatelliteAnim = useRef(27);

  const saturnProgRef = useRef(0);
  const satelliteProgRef = useRef(0);

  const clearSaturnAnim = useCallback(() => {
    saturnAnim.current = null;
    //@ts-ignore
    const image = imgs[`planet${1 + Math.floor(Math.random() * 10)}`];
    setImage(image);
    if (!saturnRef || !saturnRef.current) return;
    gsap.set(saturnRef.current, { left: '-1000px' });
  }, []);

  const clearSatelliteAnim = useCallback(() => {
    satelliteAnim.current = null;
    if (!satelliteRef || !satelliteRef.current) return;
    gsap.set(satelliteRef.current, { left: '-1000px' });
  }, []);

  const onSaturnAnimComplete = useCallback(() => {
    saturnAnimEnd.current = Date.now();
    nextSaturnAnim.current = 27 + Math.random() * 5;
    // nextSaturnAnim.current = -100;
    clearSaturnAnim();
  }, [clearSaturnAnim]);

  const onSatelliteAnimComplete = useCallback(() => {
    satelliteAnimEnd.current = Date.now();
    // nextSatelliteAnim.current = 8 + Math.random() * 6;
    nextSatelliteAnim.current = -100;
    clearSatelliteAnim();
  }, [clearSatelliteAnim]);

  const saturnAnim = useRef<gsap.core.Timeline | null>();
  const satelliteAnim = useRef<gsap.core.Timeline | null>();

  const animPlay = useCallback(() => {
    if (saturnAnim.current) saturnAnim.current.pause();
    if (satelliteAnim.current) satelliteAnim.current.pause();

    const saturnAnimStart =
      saturnAnimEnd.current + nextSaturnAnim.current * 1000;
    const satelliteStart =
      satelliteAnimEnd.current + nextSatelliteAnim.current * 1000;

    const saturnProg = (Date.now() - saturnAnimStart) / 30000;
    const satelliteProg = (Date.now() - satelliteStart) / 5000;

    if (
      saturnProg >= 0 &&
      saturnProg < 1 &&
      (!saturnAnim || !saturnAnim.current)
    ) {
      saturnAnim.current = saturnAnimation(saturnRef, onSaturnAnimComplete);
    }

    if (
      satelliteProg >= 0 &&
      satelliteProg < 1 &&
      (!satelliteAnim || !satelliteAnim.current)
    ) {
      satelliteAnim.current = satelliteAnimation(
        satelliteRef,
        onSatelliteAnimComplete
      );
    }

    if (saturnProg < 1 && saturnProg >= 0)
      saturnAnim.current?.progress(saturnProg).play();
    if (satelliteProg < 1 && satelliteProg >= 0)
      satelliteAnim.current?.progress(satelliteProg).play();
  }, [onSatelliteAnimComplete, onSaturnAnimComplete]);

  useEffect(() => {
    let interval: NodeJS.Timer | null = null;

    if (status !== 'play') {
      saturnAnimEnd.current = 0;
      satelliteAnimEnd.current = 0;

      nextSaturnAnim.current = 40;
      nextSatelliteAnim.current = 27;
    }

    switch (status) {
      case 'bet':
        clearSaturnAnim();
        clearSatelliteAnim();
        break;
      case 'ready':
        break;
      case 'play':
        if (saturnAnimEnd.current === 0) {
          saturnAnimEnd.current = globalTime;
          satelliteAnimEnd.current = globalTime;
        }

        animPlay();
        interval = setInterval(() => {
          animPlay();
        }, 300);

        break;
      case 'explosion':
        break;
      case 'back':
        let backDuration = CRASH_BACK_TIME - (Date.now() - globalTime) / 1000;
        if (backDuration < 0) backDuration = 0;

        if (saturnAnim.current)
          saturnAnim.current
            .progress(saturnProgRef.current)
            .duration(backDuration)
            .reverse()
            .then(() => {
              saturnProgRef.current = 0;
              saturnAnim.current = null;
            });

        if (satelliteAnim.current)
          satelliteAnim.current
            .progress(satelliteProgRef.current)
            .duration(backDuration)
            .reverse()
            .then(() => {
              satelliteProgRef.current = 0;
              satelliteAnim.current = null;
            });

        break;
    }

    return () => {
      if (saturnAnim.current)
        saturnProgRef.current = saturnAnim.current.progress();
      if (satelliteAnim.current)
        satelliteProgRef.current = satelliteAnim.current.progress();

      saturnAnim.current?.pause();
      satelliteAnim.current?.pause();

      if (interval) clearInterval(interval);
    };
  }, [
    onSatelliteAnimComplete,
    onSaturnAnimComplete,
    globalTime,
    status,
    animPlay,
    clearSaturnAnim,
    clearSatelliteAnim
  ]);

  return (
    <Container>
      <Saturn
        src={image}
        alt=""
        ref={saturnRef}
        style={{ opacity: usingAnimation ? 1 : 0 }}
      />
      <Satellite
        ref={satelliteRef}
        style={{ opacity: usingAnimation ? 1 : 0 }}
      />
    </Container>
  );
}

const Saturn = styled.img`
  position: absolute;
  width: 260px;
  height: 260px;
  left: -400px;
`;

const Satellite = styled.img`
  position: absolute;
  width: 100px;
  height: 40px;
  left: -400px;
`;
Satellite.defaultProps = {
  src: satelliteImg,
  alt: ''
};

const Container = styled(Box)`
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
`;
