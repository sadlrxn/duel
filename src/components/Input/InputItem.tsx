import React from 'react';
import styled from 'styled-components';

import { ReactComponent as TriangleArrowIcon } from 'assets/imgs/icons/triangle-arrow.svg';

import { Box, Flex, FlexProps } from 'components/Box';
import { Button } from 'components/Button';
import { Label, Span } from 'components/Text';
import Input from './Input';

export interface InputItemProps extends FlexProps {
  label?: string;
  description?: string;
  status?: string;
  allInButtonBorderRadius?: string | number;
  allInButton?: boolean;
  minButton?: boolean;
  maxButton?: boolean;
  upDownButton?: boolean;
  upButtonLabel?: string;
  downButtonLabel?: string;
  labelColor?: string;
  labelFontWeight?: number;
  labelLineHeight?: string | number;
  statusFontWeight?: number;
  descriptionFontWeight?: number;
  size?: number;
  handleAllIn?: any;
  handleMin?: any;
  handleMax?: any;
  handleUp?: any;
  handleDown?: any;
  onInputChange?: any;
  inputValue?: string;
  inputColor?: string;
  inputFontWeight?: number;
  inputLineHeight?: string | number;
  inputWidth?: string | number;
  inputSecondValue?: string;
  inputSecondColor?: string;
  inputSecondFontWeight?: number;
  inputSecondLineHeight?: string | number;
  inputSecondWidth?: string | number;
  inputSecondIcon?: string;
  inputSecondIconSize?: string | number;
  placeholder?: string;
  placeholderColor?: string;
  inputBoxBackground?: string;
  inputBoxBorderRadius?: string | number;
  inputBoxHeight?: string | number;
  disabled?: boolean;
  gap?: number | string;
  iconUrl?: string;
  iconSize?: string | number;
  step?: number | string;
  type?: React.HTMLInputTypeAttribute;
  tabIndex?: number;
  readOnly?: boolean;
}

export default function InputItem({
  label = '',
  status = '',
  description = '',
  allInButtonBorderRadius = '7px',
  allInButton = false,
  minButton = false,
  maxButton = false,
  upDownButton = false,
  upButtonLabel = '',
  downButtonLabel = '',
  labelColor = '#4F617B',
  labelFontWeight = 400,
  labelLineHeight = 1,
  statusFontWeight = 700,
  descriptionFontWeight = 700,
  size = 12,
  handleAllIn,
  handleMin,
  handleMax,
  handleUp,
  handleDown,
  onInputChange,
  inputValue = '',
  inputColor = '#FFF',
  inputFontWeight = 500,
  inputLineHeight,
  inputWidth,
  inputSecondValue,
  inputSecondColor = 'white',
  inputSecondFontWeight = 700,
  inputSecondLineHeight,
  inputSecondWidth = 'max-content',
  inputSecondIcon,
  inputSecondIconSize = '12px',
  placeholder = '',
  placeholderColor = '#4F617B',
  inputBoxBackground = '#03060999',
  inputBoxBorderRadius = '11px',
  inputBoxHeight,
  disabled = false,
  gap,
  iconUrl,
  iconSize = '16px',
  step = 'any',
  type,
  tabIndex,
  readOnly,
  ...props
}: InputItemProps) {
  return (
    <Flex
      flexDirection="column"
      gap={gap ? gap : '0.5em'}
      style={{
        pointerEvents: disabled ? 'none' : 'all',
        opacity: disabled ? 0.7 : 1
      }}
      {...props}
    >
      {label !== '' && (
        <Flex alignItems="center" justifyContent="space-between" height="17px">
          <Flex alignItems="center">
            <Label
              ml="2px"
              color={labelColor}
              fontWeight={labelFontWeight}
              lineHeight={labelLineHeight}
              fontSize="1em"
            >
              {label}
            </Label>
            <Label
              ml="4px"
              color={'#95A3B9'}
              fontWeight={statusFontWeight}
              lineHeight={labelLineHeight}
              fontSize="1em"
            >
              {status}
            </Label>
          </Flex>
          {allInButton && (
            <BetAllButton
              borderRadius={allInButtonBorderRadius}
              fontSize="1em"
              fontWeight={500}
              height="100%"
              onClick={handleAllIn}
              type="button"
            >
              All-in
            </BetAllButton>
          )}
        </Flex>
      )}
      <InputBox
        background={inputBoxBackground}
        borderRadius={inputBoxBorderRadius}
        padding={`3px 3px 3px 11px`}
        height={inputBoxHeight ? inputBoxHeight : `3em`}
      >
        {iconUrl !== undefined && (
          <img
            src={iconUrl}
            style={{
              width: iconSize,
              minWidth: iconSize,
              height: iconSize,
              maxHeight: iconSize
            }}
            alt=""
          />
        )}
        <Input
          placeholder={placeholder}
          placeholderColor={placeholderColor}
          fontSize="1em"
          width={inputWidth}
          onChange={onInputChange}
          value={inputValue}
          color={inputColor}
          fontWeight={inputFontWeight}
          lineHeight={inputLineHeight}
          flexGrow={1}
          step={step}
          type={type}
          tabIndex={tabIndex}
          readOnly={readOnly}
        />
        {inputSecondValue && (
          <>
            <VerticalDivider />
            <Flex
              alignItems="center"
              gap={3}
              minWidth={inputSecondWidth}
              mx="5px"
            >
              {inputSecondIcon && (
                <img
                  src={inputSecondIcon}
                  alt=""
                  style={{
                    width: inputSecondIconSize,
                    minWidth: inputSecondIconSize,
                    height: inputSecondIconSize,
                    minHeight: inputSecondIconSize
                  }}
                />
              )}
              <Span
                color={inputSecondColor}
                fontWeight={inputSecondFontWeight}
                lineHeight={inputSecondLineHeight}
                minWidth="max-content"
              >
                {inputSecondValue}
              </Span>
            </Flex>
          </>
        )}
        {minButton && (
          <MinMaxButton type="button" onClick={handleMin}>
            Min
          </MinMaxButton>
        )}
        {maxButton && (
          <MinMaxButton type="button" onClick={handleMax}>
            Max
          </MinMaxButton>
        )}
        {upDownButton && (
          <UpDownButtonWrapper minWidth="max-content" height="100%">
            <UpperBtn type="button" onClick={handleUp}>
              <TriangleArrowIcon
                width={`${size / 2}px`}
                height={`${size / 2}px`}
              />
              <Span>{upButtonLabel}</Span>
            </UpperBtn>
            <DownBtn type="button" onClick={handleDown}>
              <TriangleArrowIcon
                width={`${size / 2}px`}
                height={`${size / 2}px`}
              />
              <Span>{downButtonLabel}</Span>
            </DownBtn>
          </UpDownButtonWrapper>
        )}
        {!maxButton && handleMax && (
          <MinMaxButton
            className="input_item_mobile_max_btn"
            type="button"
            onClick={handleMax}
          >
            Max
          </MinMaxButton>
        )}
      </InputBox>
      {description !== '' && (
        <Flex width="100%" justifyContent="end">
          <Label
            color={labelColor}
            fontWeight={descriptionFontWeight}
            lineHeight={labelLineHeight}
            fontSize="1em"
          >
            {description}
          </Label>
        </Flex>
      )}
    </Flex>
  );
}

const UpperBtn = styled(Button)`
  color: #768bad;
  font-weight: 600;
`;

const DownBtn = styled(Button)`
  color: #768bad;
  font-weight: 600;
  svg {
    transform: rotate(180deg);
  }
`;

const UpDownButtonWrapper = styled(Box)`
  display: flex;
  flex-direction: row-reverse;
  gap: 5px;

  ${UpperBtn} {
    background: #1a293d;
    border-radius: 5px;
    min-width: 44px;
    width: 44px;
    height: 100%;

    svg {
      display: none;
    }
    span {
      display: block;
    }
  }

  ${DownBtn} {
    background: #1a293d;
    border-radius: 5px;
    min-width: 44px;
    width: 44px;
    height: 100%;

    svg {
      display: none;
    }
    span {
      display: block;
    }
  }

  .width_700 & {
    display: block;

    ${UpperBtn} {
      border-radius: 7px 7px 0px 0px;
      min-width: auto;
      width: auto;
      height: 50%;
      background: #24354d;
      svg {
        display: block;
      }
      span {
        display: none;
      }
    }
    ${DownBtn} {
      border-radius: 0px 0px 7px 7px;
      min-width: auto;
      width: auto;
      height: 50%;
      background: #24354d;
      svg {
        display: block;
      }
      span {
        display: none;
      }
    }
  }
`;

const VerticalDivider = styled(Box)`
  width: 1px;
  height: 100%;
  background: #24354c;
`;

const MinMaxButton = styled(Button)`
  border-radius: 5px;
  background: #1a293d;
  color: #768bad;
  font-weight: 600;
  min-width: 44px;
  width: 44px;
  height: 100%;

  &.input_item_mobile_max_btn {
    display: flex;

    .width_700 & {
      display: none;
    }
  }
`;
MinMaxButton.defaultProps = {
  variant: 'secondary'
};

const BetAllButton = styled(Button)`
  background-color: #24354c;
  color: #526d90;

  transition: all 0.3s ease-in;
  line-height: 0;
  /* padding: 4px 10px; */

  &:hover {
    background-color: ${({ theme }) => theme.coinflip.greenDark};
    color: ${({ theme }) => theme.colors.success};
  }

  display: none;

  .width_700 & {
    display: flex;
  }
`;

const InputBox = styled(Flex)`
  align-items: center;
  flex-grow: 1;
  gap: 5px;

  &:focus-within {
    box-shadow: 0 0 0 1px #7389a980;
  }
`;
