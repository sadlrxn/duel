import styled from "styled-components";
import { lighten } from "polished";

import { Button } from "components/Button";
import { Input } from "components/Input";
import { getColor } from "utils/getThemeValue";

export const Label = styled.label`
  font-size: 10px;
  margin-bottom: 6px;

  ${({ theme }) => theme.mediaQueries.xxxl} {
    font-size: 13px;
    margin-bottom: 10px;
  }
`;

export const Labels = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;

  ${Label}:last-child {
    transition: all 0.15s ease-in;

    &:focus,
    &:hover {
      color: ${getColor("green")};
    }

    &:active {
      text-shadow: 0 0 10px ${getColor("green")};
    }
  }
`;

export const Image = styled.img`
  flex: none;
  width: 9px;
  height: 9px;

  .width_1100 & {
    width: 14px;
    height: 14px;
  }
`;

export const Multipliers = styled.div`
  // position: absolute;
  // top: 0;
  // right: 0px;
  height: 100%;

  display: flex;
  flex: none;
  padding: 3px;
  flex-direction: column;
`;

export const InputContainer = styled.div`
  position: relative;
  display: flex;
  align-items: center;
  gap: 0.3rem;
  width: 100%;
  height: 50px;
  background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0) 162.5%);
  border-radius: 9px;
  padding-left: 0.5rem;

  .width_1100 & {
    width: 140px;
    height: 50px;
  }
`;

export const Amount = styled.div`
  display: flex;
  flex-direction: column;
`;

export const StyledButton = styled(Button)`
  border: 0;
  width: 100%;
  min-width: 20px;
  height: 100%;
  background: #24354d;

  display: flex;
  justify-content: center;
  align-items: center;

  &:hover,
  &:focus {
    background: ${lighten(0.1, "#24354d")};
  }

  &:active {
    background: ${lighten(0.05, "#24354d")};
  }

  &:first-child {
    border-radius: 5px 5px 0 0;
  }

  &:last-child {
    border-radius: 0 0 5px 5px;

    svg {
      transform: rotate(180deg) scale(0.8);

      ${({ theme }) => theme.mediaQueries.xl} {
        transform: rotate(180deg) scale(1);
      }
    }
  }

  svg {
    transform: scale(0.8);

    ${({ theme }) => theme.mediaQueries.xl} {
      transform: scale(1);
    }
  }
`;

export const StyledInput = styled(Input)`
  &[type="number"]::-webkit-inner-spin-button,
  &[type="number"]::-webkit-outer-spin-button {
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    margin: 0;
  }

  &[type="number"] {
    width: 100%;
    font-family: Inter;
    font-style: normal;
    font-weight: normal;
    line-height: 24px;
    color: #96a8c2;
    text-align: left;
    font-size: 14px;
    padding: 11px 0;

    .width_1100 & {
      font-size: 20px;
      // padding: 18px 0px 18px 15px;
      padding: 18px 0;
    }
  }
`;

export const Text = styled.p`
  margin-bottom: 1.5rem;

  font: 600 22px "Inter";
  text-align: center;
  color: #fff;
`;

export const Popup = styled.section`
  max-width: 416px;
  margin: 0 auto;

  background: linear-gradient(180deg, #0f2035 0%, #1a293d 100%);
  border-radius: 12px;

  padding: 2rem 1rem;

  ${({ theme }) => theme.mediaQueries.sm} {
    padding: 2rem 3rem;
  }
`;

export const ButtonGroup = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
`;
