import { FC } from 'react';
import { Flex } from 'components';
import Square from './Square';
import styled from 'styled-components';
import { DreamtowerMode } from 'api/types/dreamtower';

const Row: FC<{
  roundId?: number;
  value: number[];
  towerMode?: DreamtowerMode;
  isNext: boolean;
  isHighlight: boolean;
  isClickable: boolean;
  isUnderBroken: boolean;
  selectedIndex?: number;
  nextMultiplier: number;
  blocksInRow: number;
  handleClickSquare: (index: any) => void;
}> = ({
  roundId = 0,
  value,
  isNext,
  towerMode = 'manual',
  isHighlight,
  isClickable,
  selectedIndex,
  isUnderBroken,
  nextMultiplier,
  blocksInRow,
  handleClickSquare
}) => {
  return (
    <StyledRow className="dreamtower_row" highlight={isHighlight}>
      {value.map((v, i) => (
        <Square
          key={i}
          roundId={roundId}
          isNext={isNext}
          isClickable={isClickable}
          isStar={value[i] === 1}
          isInPath={selectedIndex === i}
          isSelected={isUnderBroken && selectedIndex === i}
          towerMode={towerMode}
          nextMultiplier={nextMultiplier}
          blocks={blocksInRow}
          handleClick={() => handleClickSquare(i)}
        />
      ))}
    </StyledRow>
  );
};

const StyledRow = styled(Flex)<{
  highlight: boolean;
}>`
  justify-content: space-between;
  padding: 5px 20px;
  border-width: 0px 2px 0px 2px;
  /* background: #e8caee; */
  /* box-shadow: 0px 0px 7px 3px rgba(255, 18, 246, 0.25); */

  ${({ highlight }) =>
    highlight &&
    `
        border-color: #e8caee;
        border-style: solid;
      `}
`;

export default Row;
