import styled from 'styled-components';

import { Box, Flex } from 'components';
import towerBg from 'assets/imgs/dreamtower/tower.png';

import Logo from '../../Logo';
import Row from '../../Row';

interface TowerProps {
  tower: number[][];
  blocksInRow: number;
}

export default function Tower({ tower, blocksInRow = 4 }: TowerProps) {
  return (
    <StyledTowerBox>
      <Logo />
      <Flex py="10px" flexDirection={'column-reverse'}>
        {tower.map((v, i) => (
          <Row
            key={i}
            value={v}
            isNext={false}
            isHighlight={false}
            isClickable={false}
            isUnderBroken={false}
            selectedIndex={undefined}
            nextMultiplier={-1}
            blocksInRow={blocksInRow}
            handleClickSquare={() => {}}
          />
        ))}
      </Flex>
    </StyledTowerBox>
  );
}

const StyledTowerBox = styled(Box)`
  width: 320px;
  height: 780px;
  background-image: url(${towerBg});
  background-size: cover;
  position: relative;
`;
