import React, { ChangeEvent, useRef, useEffect } from 'react';
import { ReactComponent as TriangleArrowIcon } from 'assets/imgs/icons/triangle-arrow.svg';
import coinSvg from 'assets/imgs/coins/coin.svg';
import blueCoinSvg from 'assets/imgs/coins/coin-blue.svg';
import {
  InputContainer,
  Image,
  Multipliers,
  Amount,
  StyledButton,
  StyledInput
} from './BetInput.styles';
import { useAppSelector } from 'state';

interface IProps {
  value: string;
  className?: string;
  tabIndex?: number;

  onChange?: (event: ChangeEvent<HTMLInputElement>) => void;
  duplicateBet?: () => void;
  divideBetInHalf?: () => void;
}

export default function BetInput({
  divideBetInHalf,
  duplicateBet,
  onChange,
  tabIndex,
  value
}: IProps) {
  const { betBalanceType } = useAppSelector(state => state.user);

  const inputRef = useRef(null);

  useEffect(() => {
    if (!inputRef) return;
    //@ts-ignore
    inputRef.current.focus();
  }, []);

  return (
    <>
      <Amount>
        <InputContainer>
          {betBalanceType === 'coupon' ? (
            <Image src={blueCoinSvg} alt="Coin" />
          ) : (
            <Image src={coinSvg} alt="Coin" />
          )}

          <StyledInput
            ref={inputRef}
            tabIndex={tabIndex}
            type="number"
            placeholder="0.00"
            pattern="^[0-9]*\.?[0-9]*$"
            step="any"
            onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) => {
              if (
                e.key === 'e' ||
                e.key === 'E' ||
                e.key === '+' ||
                e.key === '-'
              )
                e.preventDefault();
            }}
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
