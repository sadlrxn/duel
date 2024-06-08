import React, { useEffect, useState } from 'react';
import { useClickAnyWhere, useWindowSize } from 'usehooks-ts';
import styled from 'styled-components';

import { ModalProps } from 'components/Modal';
import { duelEmojis } from 'config';
import { useMatchBreakpoints } from 'hooks';
import { chatWidth } from 'theme';

import DuelEmojiItem from './DuelEmojiItem';

interface EmojiModalProps extends ModalProps {
  inputRef: React.RefObject<HTMLInputElement | HTMLTextAreaElement>;
}

export default function EmojiModal({ onDismiss, inputRef }: EmojiModalProps) {
  const [rendered, setRendered] = useState(false);
  const { width } = useWindowSize();
  const { isMobile } = useMatchBreakpoints();

  useEffect(() => {
    setRendered(true);
  }, []);

  useClickAnyWhere(() => {
    rendered && onDismiss && onDismiss();
  });

  return (
    <StyledModal style={{ left: isMobile ? `10px` : `${width - chatWidth}px` }}>
      {duelEmojis.map(emoji => (
        <DuelEmojiItem
          key={emoji.id}
          emoji={emoji}
          onDismiss={onDismiss}
          inputRef={inputRef}
        />
      ))}
    </StyledModal>
  );
}

const StyledModal = styled.div`
  position: fixed;
  bottom: 15px;
  z-index: 60;
  color: #b2d1ff;
  max-height: calc(100vh - 90px);
  background: linear-gradient(180deg, #0c1725 0%, #18283e 100%);
  border: 1px solid #26374e;
  box-shadow: 0px 4px 22px rgba(0, 0, 0, 0.4);
  border-radius: 6px;
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
  overflow: hidden auto;
  padding: 8px;
  ${({ theme }) => theme.mediaQueries.sm} {
    transform: translateX(-100%);
  }
`;
