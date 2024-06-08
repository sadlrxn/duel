import { useRef } from 'react';
import styled, { keyframes, css } from 'styled-components';
import { CSSTransition } from 'react-transition-group';

import { Avatar, Chip, Flex, FlexProps } from 'components';
import { CrashCashEvent } from 'api/types/crash';
import { convertBalanceToChip } from 'utils/balance';

interface CashoutUserProps extends FlexProps {
  cashOut: CrashCashEvent;
  top?: number;
}

const animation = (x: number, y: number) => {
  return keyframes`
  0% { transform: translate(0px, 0px); opacity: 1; }
  80% { transform: translate(${x * 0.8}px, ${y * 0.8}px); opacity: 1; }
  100% { transform: translate(${x}px, ${y}px); opacity: 0; }
`;
};

export default function CashoutUser({ cashOut, ...props }: CashoutUserProps) {
  const nodeRef = useRef<any>();

  return (
    <CSSTransition
      timeout={{
        enter: 0,
        exit: 1400
      }}
      nodeRef={nodeRef}
      unmountOnExit
      className="crash_cashout_user"
      {...props}
    >
      <Container
        xDis={-30 - Math.random() * 30}
        yDis={50 + Math.random() * 30}
        ref={nodeRef}
        {...props}
      >
        <Avatar
          userId={cashOut.user.id}
          image={cashOut.user.avatar}
          name={cashOut.user.avatar}
          size="38px"
          border="1.19px solid #B2D1FF"
          borderRadius="100%"
          useProxy
        />
        {/* <Span fontSize="16px" color="textWhite" fontWeight={700}>
        {cashOut.user.name}
      </Span> */}
        <Chip
          price={convertBalanceToChip(cashOut.amount)}
          decimal={2}
          fontSize="14px"
          fontWeight={700}
          color="success"
          size={14}
        />
      </Container>
    </CSSTransition>
  );
}

const Container = styled(Flex)<{ xDis: number; yDis: number }>`
  flex-direction: column;
  align-items: center;
  gap: 10px;

  position: absolute;
  top: 60px;
  left: calc(100% - 180px);
  z-index: 2;
  opacity: 0;
  pointer-events: none;

  .width_700 & {
    left: calc(100% - 200px);
  }

  .width_900 & {
    left: calc(100% - 240px);
  }

  ${({ xDis, yDis }) => css`
    animation: ${animation(xDis, yDis)} 1.5s linear;
  `}
`;
