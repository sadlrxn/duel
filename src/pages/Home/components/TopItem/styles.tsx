import { Text } from "components/Text";
import styled from "styled-components";

import { Box, Flex } from "components/Box";
import { Button } from "components/Button";

export const StyledText = styled(Text)`
  font-family: "Inter";
  font-style: normal;
  font-weight: 600;
  font-size: 25px;
  line-height: 30px;

  color: #ffffff;

  margin-right: 10px;
`;

export const StyledUserInfo = styled(Box)`
  position: relative;
  background: #202e3f;
  border-radius: 9px;

  &:hover {
    background: #26374b;
  }
`;

export const StyledContainer = styled(Flex)<{ isWinComp: boolean }>`
  background: ${({ isWinComp }) =>
    isWinComp
      ? "linear-gradient(270deg, #4b4430 0%, rgba(30, 43, 57, 0) 50.63%)"
      : "linear-gradient(90deg, rgba(32, 46, 63, 0) 50%, rgba(79, 255, 139, 0.2) 100%)"};

  justify-content: space-between;
  padding: 15px 25px;
  border-radius: 9px;

  &:before {
    position: absolute;
    background: ${({ isWinComp }) => (isWinComp ? "#ffbe5c" : "#4FFF8B")};
    content: "";
    display: block;
    right: 0;
    height: 50%;
    width: 2px;
    top: 25%;
  }
  flex-direction: column;
  align-items: center;
  gap: 12px;
  .width_700 & {
    flex-direction: row;
  }
`;

export const No = styled(Text)`
  font-family: "Inter";
  font-style: normal;
  font-weight: 400;
  font-size: 16px;
  line-height: 19px;
  text-align: center;

  color: #768bad;

  margin-right: 15px;
`;

export const Avatar = styled.div`
  width: 30px;
  height: 30px;

  margin-left: 10px;
  margin-right: 10px;
  img {
    width: 100%;
    border-radius: 8px;
  }
`;

export const StyleButton = styled(Button)`
  background: transparent;

  border: 1px solid #768bad;
  border-radius: 34px;

  font-family: "Inter";
  font-style: normal;
  font-weight: 400;
  font-size: 14px;
  line-height: 17px;
  /* identical to box height */

  text-align: center;

  /* Secondary Text */

  color: #768bad;

  padding: 4px 12px;

  &:hover {
    background: rgba(118, 139, 173, 0.15);
  }
`;
