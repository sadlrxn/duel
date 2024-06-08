import React from 'react';
import styled from 'styled-components';

import backgroundImg from 'assets/imgs/plinko/background.png';

import { Box } from 'components';

export default function Background() {
  return <Container />;
}

const Container = styled(Box)`
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;

  background-image: url(${backgroundImg});
  background-position: top right;
  background-size: cover;
  z-index: -1;
`;
