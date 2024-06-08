import React, { useEffect, useRef, useState, useCallback } from 'react';
import gsap from 'gsap';

import { Text } from 'components/Text';
import { Flex } from 'components/Box';
import { ReactComponent as ChevronIcon } from 'assets/imgs/icons/chevron.svg';
import { ReactComponent as UserIcon } from 'assets/imgs/icons/user.svg';
import { ReactComponent as DotIcon } from 'assets/imgs/icons/dot.svg';

import { useAppDispatch, useAppSelector } from 'state';
import { toggleChat } from 'state/chat/actions';

import ChatMsg from './ChatMsg';
import ChatInput from './ChatInput';
import {
  ChartMsgContainer,
  ChatContainer,
  ChatStatus,
  GlassRect,
  ToggleBtn
} from './styles';
import { useMatchBreakpoints } from 'hooks';

const Chat = () => {
  const dispatch = useAppDispatch();
  const { isMobile } = useMatchBreakpoints();
  const { id: currentUserId } = useAppSelector(state => state.user);
  const { show, msgs, users } = useAppSelector(state => state.chat);
  const containerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const [isHover, setIsHover] = useState(false);

  const toggleChatFunc = () => {
    dispatch(toggleChat());
  };

  const handleReplyMsgClick = useCallback(
    (msgId: number) => {
      const index = msgs.slice(-100).findIndex(msg => msg.id === msgId);
      if (index === -1) return;
      const msgElement =
        document.getElementsByClassName('chat_original_msg')[index];
      if (!msgElement) return;
      msgElement.scrollIntoView({ block: 'center', behavior: 'smooth' });
      gsap
        .timeline()
        .to(msgElement, { background: '#273D5B', duration: 0.2 }, '0.4')
        .to(msgElement, { background: '#131E2D', duration: 0.4 }, '>=0.3');
    },
    [msgs]
  );

  useEffect(() => {
    if (!('Notification' in window)) {
      console.info('This browser does not support desktop notification');
    } else {
      Notification.requestPermission();
    }
    if (containerRef.current)
      containerRef.current.scrollTo({
        top: containerRef.current.scrollHeight
      });
  }, []);

  const lastMsgIdRef = useRef<number>(0);

  useEffect(() => {
    if (containerRef.current && msgs.length > 0) {
      if (lastMsgIdRef.current !== msgs[msgs.length - 1].id)
        if (
          msgs[msgs.length - 1].author.id === currentUserId ||
          (!isMobile && !isHover) ||
          !show
        )
          containerRef.current.scrollTo({
            top: containerRef.current.scrollHeight
          });

      lastMsgIdRef.current = msgs[msgs.length - 1].id;
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [msgs, currentUserId, show]);

  return (
    <ChatContainer show={show}>
      <Flex gap={8} zIndex={10} pr="17px">
        <ToggleBtn show={show} onClick={toggleChatFunc}>
          <ChevronIcon />
        </ToggleBtn>

        <ChatStatus>
          <Text
            color={'#fff'}
            fontSize="20px"
            lineHeight={'24px'}
            fontWeight={600}
            mr="16px"
          >
            Chat
          </Text>

          <Flex mr="6px">
            <UserIcon />
            <DotIcon color="#4FFF8B" />
          </Flex>

          <Text color={'#4FFF8B'}>{users.length}</Text>
        </ChatStatus>
      </Flex>
      <GlassRect />
      <ChartMsgContainer
        ref={containerRef}
        onMouseEnter={() => setIsHover(true)}
        onMouseLeave={() => setIsHover(false)}
      >
        {msgs.slice(-100).map(msg => (
          <ChatMsg
            inputRef={inputRef}
            handleReplyMsgClick={handleReplyMsgClick}
            key={`chat_message_${msg.author.id}_${msg.time}`}
            msg={msg}
          />
        ))}
      </ChartMsgContainer>

      <ChatInput
        handleReplyMsgClick={handleReplyMsgClick}
        inputRef={inputRef}
      />
    </ChatContainer>
  );
};

export default Chat;
