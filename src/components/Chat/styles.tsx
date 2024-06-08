import styled from 'styled-components';
import TextareaAutosize from 'react-textarea-autosize';

import { Box, Flex } from 'components/Box';
import { Button } from 'components/Button';
import { Span } from 'components/Text';
import { chatWidth } from 'theme';

export const BubbleBtn = styled(Button)<{
  stroke?: boolean;
  fill?: boolean;
}>`
  border-radius: 6px;
  background: #273d5b;
  width: 42px;
  height: 28px;

  &:hover {
    background: #1c2c41;
    svg {
      ${({ stroke }) => stroke && 'stroke: #4FFF8B;'}
      ${({ fill }) => fill && 'fill: #4FFF8B;'}
    }
  }
`;
BubbleBtn.defaultProps = {
  variant: 'secondary'
};

export const BubbleContainer = styled(Flex)`
  display: none;
  position: absolute;
  right: 6px;
  top: -3px;
  gap: 5px;
`;

export const ChatContainer = styled(Box)<{ show: boolean }>`
  position: absolute;
  box-sizing: border-box;
  height: calc(100vh - 65px);
  top: 0;
  right: 0;
  ${({ theme }) => theme.mediaQueries.xl} {
    position: relative;
  }
  transition: all 300ms;
  width: ${`${chatWidth}px`};
  background: #0d141c;
  display: flex;
  flex: none;
  flex-direction: column;
  justify-content: space-between;
  padding: 17px 0px 17px 17px;
  z-index: 20;
  transition: margin 0.5s ease;

  margin-right: ${({ show }) => (show ? '0px' : `-${chatWidth}px`)};
`;

export const ToggleBtn = styled(Button)<{ show: boolean }>`
  width: 46px;
  height: 46px;
  background: #242f42;
  border-radius: 5px;
  transition: 0.5s;

  transform: ${({ show }) =>
    show ? 'translateX(0px) rotate(0deg)' : 'translateX(-50px) rotate(180deg)'};
`;

export const ChatStatus = styled(Flex)`
  align-items: center;
  flex-grow: 1;
  background: #141f2e;
  padding: 11px 16px;
  border-radius: 5px;
  z-index: 10;
`;

export const HistoryTime = styled(Button)`
  color: #788ca9;
  font-weight: 600;
  font-size: 10px;
  line-height: 12px;
  position: absolute;
  bottom: 6px;
  right: 6px;
  padding: 2px 5px;
  border-radius: 6px;

  opacity: 0.8;
  pointer-events: none;
  display: none;
  background: rgba(39, 61, 91, 0.7);
  backdrop-filter: blur(2px);
`;

HistoryTime.defaultProps = {
  variant: 'secondary'
};

export const MsgBox = styled(Box)`
  position: relative;
  display: grid;
  align-items: start;
  gap: 10px;

  background: #131e2d;
  color: #b2d1ff;
  border-radius: 8px;
  padding: 13px;
  margin-right: 10px;

  &:hover {
    background: linear-gradient(90deg, #1a293c 0%, #121d2a 100%);
    ${HistoryTime} {
      display: flex;
    }
    ${BubbleContainer} {
      display: flex;
    }
  }
`;

export const ReplyMsgBox = styled(Box)`
  position: relative;
  display: grid;
  align-items: start;
  gap: 10px;

  color: #b2d1ff;
  border-radius: 8px 8px 0px 0px;

  font-size: 12px;
`;

export const GlassRect = styled(Box)`
  position: absolute;
  top: 0px;
  left: 0px;
  width: 319px;
  height: 181px;
  z-index: 5;
  pointer-events: none;
  background: linear-gradient(180deg, #0d141c 32.91%, rgba(3, 6, 9, 0) 100%);
`;

export const DuelItemContainer = styled(Flex)`
  cursor: pointer;
  user-select: none;
  gap: 8px;
  padding: 11px 7px;
  border-radius: 6px;
  font-weight: 500;
  font-size: 14px;
  line-height: 17px;
  &:hover,
  &:focus {
    background: #26374e;
  }
`;

export const MsgInput = styled(TextareaAutosize)`
  background: #1b283a;
  border-radius: 8px;
  width: 100%;
  padding: 11px;
  margin-top: 15px;

  resize: none;
  border: none;
  outline: none;
  color: #fff;
  font-family: 'Inter';
  font-style: normal;
  font-weight: 500;
  font-size: 14px;
  line-height: 17px;

  ::placeholder {
    color: #4f617b;
  }
`;

export const RuleBtn = styled(Button)`
  background: transparent;
  color: #768bad;
  font-weight: 600;
  border-radius: 5px;
  padding: 5px 10px;
  &:hover {
    background: #1a212a;
  }
`;

export const SendBtn = styled(Button)`
  font-weight: 700;
  border-radius: 5px;
  padding: 5px 10px;

  &:hover {
    box-shadow: 0px 0px 10px #4fff8b;
  }
`;

export const ChartMsgContainer = styled(Flex)`
  height: 100%;
  overflow: hidden;
  overflow-y: auto;
  flex-direction: column;
  gap: 14px;
  margin-right: 7px;
  & > :first-child {
    margin-top: auto !important;
    /* use !important to prevent breakage from child margin settings */
  }
`;

export const NotifierWrapper = styled.div`
  font-weight: 500;
  &:hover {
    font-weight: 700;
    cursor: pointer;
    user-select: none;
  }
`;

export const MsgSpan = styled(Span)`
  overflow-wrap: anywhere;
  line-height: 17px;
`;

export const ChatWarningContainer = styled(Flex)`
  align-items: center;
  gap: 8px;
  color: #ffffff;
  opacity: 0.5;
`;
