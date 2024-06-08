import React from 'react';
import styled from 'styled-components';
import { useAppSelector } from 'state';

import { Text, Snow1, Snow3, Snow6, SnowMan, Santa } from 'components';

export default function Dream() {
  const { isHoliday } = useAppSelector(state => state.user);

  return (
    <div style={{ position: 'relative', paddingTop: '10px' }}>
      {isHoliday && (
        <>
          <>
            <Snow1 position="absolute" size={60} top={-12} left={53} />
            <Snow6 position="absolute" size={150} top={-15} right={51} />
          </>
          <>
            <Snow1 position="absolute" size={60} bottom={-9} left={-7} />
            <SnowMan position="absolute" size={70} bottom={13} left={-12} />
            <Santa position="absolute" size={60} bottom={25} right={-7} />
            <Snow3 position="absolute" size={70} bottom="1px" right={-10} />
          </>
        </>
      )}
      <StyledDreamText className="dreamtower_dreamtext">DREAM</StyledDreamText>
      <StyledTowerText className="dreamtower_towertext">Tower</StyledTowerText>
    </div>
  );
}

const StyledDreamText = styled(Text)`
  font-family: 'Righteous';
  font-size: 64px;
  line-height: 1;
  text-align: center;
  color: transparent;
  position: relative;
  top: 5px;

  -webkit-text-stroke: 1px white;

  text-shadow: 0 0 3px #0680d100, 0 0 5px #0680d100, 0 0 10px #0680d100,
    0 0 20px #0680d1, 0 0 30px #0680d1, 0 0 40px #0680d1, 0 0 50px #0680d1;
`;

const StyledTowerText = styled(Text)`
  font-family: 'Mr Dafoe';
  font-size: 50px;
  line-height: 1;
  font-style: italic;
  text-align: center;
  color: transparent;
  rotate: -5deg;

  -webkit-text-stroke: 1px white;

  text-shadow: 0 0 3px #ff4da6, 0 0 5px #ff4da6, 0 0 10px #ff4da6,
    0 0 20px #ff4da6, 0 0 30px #ff4da6, 0 0 40px #ff4da6, 0 0 50px #ff4da6;
`;
