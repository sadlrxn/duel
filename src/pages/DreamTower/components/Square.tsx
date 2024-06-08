import { FC, useMemo } from 'react';
import styled, { css } from 'styled-components';
import { BrokeIcon, Button, StarIcon, BellIcon } from 'components';
import { DreamtowerMode } from 'api/types/dreamtower';
import { useAppSelector } from 'state';

const Square: FC<{
  roundId?: number;
  towerMode?: DreamtowerMode;
  isNext?: boolean;
  isClickable: boolean;
  isStar: boolean;
  isSelected: boolean;
  isInPath: boolean;
  nextMultiplier: number;
  blocks: number;
  handleClick: () => void;
}> = ({
  roundId = 0,
  isNext = false,
  towerMode = 'manual',
  isClickable,
  isStar,
  isSelected,
  isInPath,
  nextMultiplier,
  blocks,
  handleClick
}) => {
  const isAuto = useMemo(() => towerMode === 'auto', [towerMode]);
  const isHoliday = useAppSelector(state => state.user.isHoliday);

  return (
    <StyledSquare
      roundId={roundId}
      isNext={isNext}
      isStar={isStar}
      isSelected={isSelected}
      isInPath={isInPath}
      disabled={!isClickable}
      blocks={blocks}
      towerMode={towerMode}
      onClick={isClickable ? handleClick : null}
    >
      {isNext && <>{nextMultiplier.toFixed(2)}x</>}
      {isStar &&
        isSelected &&
        (isHoliday ? (
          <BellIcon size={68} isAuto={isAuto} roundId={roundId} />
        ) : (
          <StarIcon isAuto={isAuto} roundId={roundId} />
        ))}
      {isStar &&
        !isSelected &&
        (isHoliday ? (
          <BellIcon
            size={68}
            isAuto={isAuto}
            isSelected={isSelected}
            roundId={roundId}
            opacity={0.15}
          />
        ) : (
          <StarIcon
            isAuto={isAuto}
            isSelected={isSelected}
            roundId={roundId}
            opacity={0.15}
          />
        ))}
      {!isStar && isSelected && <BrokeIcon />}
    </StyledSquare>
  );
};

const StyledSquare = styled(Button)<{
  isStar: boolean;
  isSelected: boolean;
  isNext: boolean;
  isInPath: boolean;
  blocks: number;
  towerMode: string;
  roundId: number;
}>`
  height: 60px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.25);
  /* background: #1c2c45; */
  gap: 4px;
  font-size: 12px;
  font-weight: 700;
  color: #ffffff;

  ${({ blocks: cols }) => {
    return css`
      width: ${240 / cols}px;
    `;
  }}

  ${({ disabled }) => {
    if (disabled)
      return css`
        opacity: 1;
        cursor: not-allowed;
        pointer-events: none;
      `;
  }}

  ${({ isNext, towerMode: mode }) => {
    if (isNext) {
      return mode === 'manual'
        ? `box-shadow: 0px 0px 20px rgba(219, 0, 255, 0.25); background: rgba(184, 0, 154, 0.25); border: 1px solid #CB72E9;`
        : `box-shadow: 0px 0px 20px rgba(0, 71, 255, 0.25); background: rgba(3, 13, 108, 0.25); border: 1px solid #4154FF;`;
    }
  }}
    
  ${({ isStar, isSelected, isInPath }) => {
    if (isSelected) {
      if (isStar)
        return `border: 1px solid #FBD92B; box-shadow: inset 0px 0px 20px rgba(255, 199, 0, 0.25);`;
      else return `border: 1px solid #D70606;`;
    } else {
      if (isInPath)
        return `box-shadow: 0px 0px 20px rgba(0, 71, 255, 0.25); background: rgba(3, 13, 108, 0.25); border: 1px solid #4154FF;`;
    }
  }}
`;

export default Square;
