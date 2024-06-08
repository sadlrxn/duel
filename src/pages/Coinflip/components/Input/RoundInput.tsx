import React, { ChangeEvent } from 'react';
import { ReactComponent as TriangleArrowIcon } from 'assets/imgs/icons/triangle-arrow.svg';
import {
  InputContainer,
  Multipliers,
  Amount,
  StyledButton,
  StyledInput
} from './BetInput.styles';

interface IProps {
  value: string;
  className?: string;
  isCount?: boolean;
  tabIndex?: number;

  onChange?: (event: ChangeEvent<HTMLInputElement>) => void;
  duplicateBet?: () => void;
  divideBetInHalf?: () => void;
}

export default function BetInput({
  divideBetInHalf,
  duplicateBet,
  onChange,
  value,
  isCount = false,
  tabIndex = 0
}: IProps) {
  return (
    <>
      <Amount>
        <InputContainer>
          <StyledInput
            tabIndex={tabIndex}
            type="number"
            placeholder={isCount ? '0' : '0.00'}
            pattern="^\d*(\.\d{0,2})?$"
            {...{ onChange, value }}
          />
          <Multipliers>
            <StyledButton type="button" onClick={duplicateBet}>
              <TriangleArrowIcon />
            </StyledButton>
            <StyledButton type="button" onClick={divideBetInHalf}>
              <TriangleArrowIcon />
            </StyledButton>
          </Multipliers>
        </InputContainer>
      </Amount>
    </>
  );
}
