import styled from 'styled-components';

import { Box, Flex, Button } from 'components';

export const ToggleAnimation = styled(Button)`
  background: transparent;
  position: absolute;
  right: 0;
  top: 0;
  height: 20px;
  z-index: 10;
  border-radius: 5px 15px 5px 5px;
  border: 2px solid #242f42;

  color: #fff9;

  opacity: 0.6;
  &:hover {
    opacity: 1;
  }
`;
ToggleAnimation.defaultProps = {
  variant: 'secondary'
};

export const GameWrapper = styled(Box)`
  position: relative;
  padding: 36px 42px;

  margin-top: -2px;
  border: solid #7389a9;
  border-width: 2px 2px 0 2px;
  border-radius: 0px 15px 0px 0px;
  color: #7389a9;
  overflow: hidden;
  min-height: 360px;
  transition: all 0.5s;
`;

export const Marker = styled.div`
  position: absolute;
  top: 0;
  left: 50%;
  height: 100%;
  width: 1px;
  background: linear-gradient(
    0deg,
    rgba(79, 255, 139, 0) 0%,
    ${({ theme }) => theme.colors.success} 100%
  );
`;

export const Sword = styled(Flex)`
  position: absolute;
  justify-content: center;
  align-items: center;
  width: 51px;
  height: 108px;
  left: 50%;
  top: 200px;
  transform: translateX(-50%);
  background: radial-gradient(
    50% 50% at 50% 50%,
    #4fff8b19 0%,
    rgba(79, 255, 139, 0) 100%
  );
  transition: margin-top 0.8s;
`;
