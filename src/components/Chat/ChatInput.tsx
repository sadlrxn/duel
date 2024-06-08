import React, {
  useMemo,
  useState,
  useEffect,
  useCallback,
  useRef
} from 'react';
import { MentionsInput, Mention } from 'react-mentions';
import { CountdownCircleTimer } from 'react-countdown-circle-timer';

import frogImg from 'assets/imgs/chat/frog.webp';
import { Box, Flex } from 'components/Box';
import { Span } from 'components/Text';
import { Button } from 'components/Button';
import { useModal } from 'components/Modal';
import { ProfileModal, TipModal } from 'components/Modals';
import { duelEmojis } from 'config';
import { useAppDispatch, useAppSelector } from 'state';
import { sendMessage, setChatMsg, setChatError, setReply } from 'state/actions';
import { formatUserName } from 'utils/format';

import EmojiModal from './EmojiModal';
import DuelEmojiItem from './EmojiModal/DuelEmojiItem';
import CommandModal from './CommandModal';
import RulesModal from './RulesModal';
import { neverMatchingRegex, tipRegex, detailsRegex } from './config';
import {
  LockedChat,
  BannedUser,
  MutedUser,
  NotEnoughPeople,
  NotEnoughChip
} from './utils';
import { RuleBtn, SendBtn, DuelItemContainer } from './styles';
import ReplyMsg from './ReplyMsg';

import classNames from './example.module.css';
import { convertChipToBalance } from 'utils/balance';
interface Emoji {
  emoji: string;
  name: string;
  shortname: string;
  unicode: string;
  html: string;
  category: string;
  order: string;
}

export default function ChatInput({
  inputRef,
  handleReplyMsgClick
}: {
  inputRef: React.RefObject<HTMLInputElement>;
  handleReplyMsgClick?: any;
}) {
  const [emojis, setEmojis] = useState<Emoji[]>([]);
  const [showCommandModal, setShowCommandModal] = useState(false);
  const [isOnComposition, setIsOnComposition] = useState(false);
  const [focusIndex, setFocusIndex] = useState(0);
  const [keyPress, setKeyPress] = useState(0); // 1: Arrow Up 2: Arrow Down 0: Any Key

  const [animateWarningKey, setAnimateWarningKey] = useState(0);
  const {
    users,
    msg,
    error: chatError,
    maxLength,
    wagerLimit,
    msgs,
    replyTo,
    chatCooldown
  } = useAppSelector(state => state.chat);
  const {
    id: userId,
    role: userRole,
    statistics: { total_wagered },
    isAuthenticated
  } = useAppSelector(state => state.user);
  const dispatch = useAppDispatch();

  const [sendedTime, setSendedTime] = useState(
    Date.now() - chatCooldown * 1000
  );
  const [sendable, setSendable] = useState(true);

  const mentionItemRef = useRef<HTMLDivElement | null>(null);

  const focusToInput = useCallback(
    (event: KeyboardEvent) => {
      if (!inputRef || !inputRef.current) return;
      if (event.key === '/') {
        if (document.activeElement !== inputRef.current) {
          event.preventDefault();
          inputRef.current.focus();
        }
      }
    },
    [inputRef]
  );

  useEffect(() => {
    window.addEventListener('keypress', focusToInput);

    return () => {
      window.removeEventListener('keypress', focusToInput);
    };
  }, [focusToInput]);

  useEffect(() => {
    if (keyPress === 0) return;
    if (!mentionItemRef || !mentionItemRef.current) return;

    const parent = mentionItemRef.current.offsetParent;
    const parentTop = parent!.scrollTop;
    const parentBottom = parentTop + parent!.clientHeight;
    const itemTop = mentionItemRef.current.offsetTop;
    const itemBottom = itemTop + mentionItemRef.current.clientHeight;

    setKeyPress(0);
    if (parentTop <= itemTop && itemBottom <= parentBottom) return;
    if (keyPress === 1) mentionItemRef.current.scrollIntoView(true);
    else if (keyPress === 2) mentionItemRef.current.scrollIntoView(false);
  }, [mentionItemRef, focusIndex, keyPress]);

  const replyMsg = useMemo(() => {
    if (!replyTo) return undefined;
    return msgs.find(msg => msg.id === replyTo);
  }, [replyTo, msgs]);

  const parsedUserName = useMemo(() => {
    if (msg.match(tipRegex) || msg.match(detailsRegex)) {
      let parsedUserName = msg.split(' ')[1];
      if (parsedUserName[0] === '@') {
        const flag = parsedUserName[1] === '{' ? 1 : 0;
        parsedUserName = parsedUserName.slice(1 + flag, flag ? -1 : undefined);
      }
      return parsedUserName;
    }
  }, [msg]);

  const parsedTipAmount = useMemo(() => {
    if (msg.match(tipRegex)) {
      return msg.split(' ')[2] ? parseFloat(msg.split(' ')[2]) : 0;
    }
    return 0;
  }, [msg]);
  const wagerAmountToUnlockChat = useMemo(() => {
    return wagerLimit - total_wagered;
  }, [total_wagered, wagerLimit]);

  const [onPresentEmojiModal] = useModal(
    <EmojiModal inputRef={inputRef} />,
    false,
    false,
    false
  );

  const [onPresentTipModal] = useModal(
    <TipModal name={parsedUserName} amount={parsedTipAmount} />
  );
  const [onPresentProfileModal] = useModal(
    <ProfileModal name={parsedUserName} />
  );
  const [onPresentRulesModal] = useModal(<RulesModal />, false, false, false);

  const onChange = (e: any) => {
    dispatch(setChatMsg(e.target.value));
  };

  const handleKeyPress = useCallback(
    (e: any) => {
      switch (e.keyCode) {
        case 27: // ESC
          dispatch(setReply());
          break;
        case 38: //Arrow Up
          setKeyPress(1);
          break;
        case 40: //Arrow Down
          setKeyPress(2);
          break;
      }
    },
    [dispatch]
  );

  const remainLetters = useMemo(() => {
    return maxLength - msg.length;
  }, [msg, maxLength]);

  const handleSubmit = (e: any) => {
    const remainTime = chatCooldown - (Date.now() - sendedTime) / 1000;
    e.preventDefault();
    if ((userRole === 'user' || userRole === 'ambassador') && remainTime >= 0) {
      return;
    }
    if (msg.trim() === '') return;
    if (remainLetters < 0) return;
    // Check status
    if (wagerAmountToUnlockChat > 0) {
      setChatError('locked');
      setAnimateWarningKey(Math.random());
    }

    dispatch(setChatMsg(''));
    if (msg.startsWith('/tip')) {
      onPresentTipModal();
      return;
    }
    if (msg.startsWith('/details')) {
      onPresentProfileModal();
      return;
    }
    let message = msg;
    if (
      msg.startsWith('/mute') ||
      msg.startsWith('/unmute') ||
      msg.startsWith('/ban') ||
      msg.startsWith('/unban')
    ) {
      let parsedUserName = msg.split(' ')[1];
      if (parsedUserName && parsedUserName[0] === '@') {
        const flag = parsedUserName[1] === '{' ? 1 : 0;
        parsedUserName = parsedUserName.slice(1 + flag, flag ? -1 : undefined);
        message = msg.split(' ')[0] + ' ' + parsedUserName;
      }
    }
    if (msg.startsWith('/rain')) {
      const splits = msg.split(' ');
      message =
        splits[0] +
        ' ' +
        splits[1] +
        ' ' +
        convertChipToBalance(+splits[2]).toFixed(0);
    }
    dispatch(
      sendMessage({
        type: 'message',
        room: 'chat',
        content: JSON.stringify({
          message: message.trim(),
          replyTo: replyMsg?.deleted ? undefined : replyTo
        })
      })
    );
    dispatch(setReply());
    setSendedTime(Date.now());
    if (
      (userRole === 'user' || userRole === 'ambassador') &&
      chatCooldown !== 0
    ) {
      setSendable(false);
      setTimeout(() => {
        setSendable(true);
      }, chatCooldown * 1000);
    }
  };

  useEffect(() => {
    if (msg.startsWith('/')) {
      setShowCommandModal(true);
    } else {
      setShowCommandModal(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [msg]);

  useEffect(() => {
    fetch('/assets/emojis.json')
      .then(response => {
        return response.json();
      })
      .then(jsonData => {
        setEmojis(jsonData.emojis);
      });
  }, []);

  const queryEmojis = (query: any) => {
    if (query.length === 0) return;

    const matches = emojis
      .filter(emoji => {
        return emoji.shortname.indexOf((':' + query).toLowerCase()) > -1;
      })
      .slice(0, 10);
    return matches.map(({ emoji }) => ({ id: emoji }));
  };
  return (
    <>
      <form onSubmit={handleSubmit}>
        {!replyMsg?.deleted && (
          <ReplyMsg
            replyMsg={replyMsg}
            mb="-28px"
            mr="17px"
            background="#131E2D"
            mt="20px"
            showCloseIcon
            handleReplyMsgClick={handleReplyMsgClick}
          />
        )}
        <Box my={20} pr="17px">
          <MentionsInput
            disabled={!isAuthenticated}
            inputRef={inputRef}
            placeholder="Say something..."
            value={msg}
            onChange={onChange}
            onKeyDownCapture={handleKeyPress}
            forceSuggestionsAboveCursor
            className="mentions"
            classNames={classNames}
            onKeyDown={e => {
              if (e.code.endsWith('Enter') && !e.shiftKey && !isOnComposition)
                handleSubmit(e);
            }}
            onCompositionStartCapture={() => setIsOnComposition(true)}
            onCompositionEndCapture={() => setIsOnComposition(false)}
          >
            <Mention
              trigger="@"
              markup="@{__id__}"
              data={users.map(user => ({
                id: user.name,
                role: user.role,
                isUser: user.id === userId
              }))}
              className={classNames.mentions__mention}
              displayTransform={id => `@${formatUserName(id)}`}
              //@ts-ignore
              renderSuggestion={(
                suggestion: { id: number; role?: string; isUser: boolean },
                search,
                highlightedDisplay,
                index,
                focused
              ) => (
                <div
                  ref={ref => {
                    if (focused) {
                      setFocusIndex(index);
                      mentionItemRef.current = ref;
                    }
                  }}
                  className={`user ${focused ? 'focused' : ''}`}
                >
                  <DuelItemContainer>
                    <Span
                      fontWeight={500}
                      color={
                        suggestion.role === 'admin'
                          ? suggestion.isUser
                            ? '#DB00FF'
                            : '#9E00FF'
                          : suggestion.role === 'moderator'
                          ? suggestion.isUser
                            ? '#FFC700'
                            : '#FFA800'
                          : suggestion.role === 'ambassador'
                          ? suggestion.isUser
                            ? '#00FFFF'
                            : '#00C4FF'
                          : '#FFF'
                      }
                    >
                      {formatUserName(`${suggestion.id}`)}
                    </Span>
                  </DuelItemContainer>
                </div>
              )}
            />
            <Mention
              trigger="#"
              markup="__id__"
              regex={neverMatchingRegex}
              data={queryEmojis}
            />
            <Mention
              trigger=":"
              markup=":{__id__}"
              displayTransform={id => `:{${id}}`}
              data={duelEmojis.map(emoji => ({
                id: emoji.id
              }))}
              renderSuggestion={(
                suggestion,
                search,
                highlightedDisplay,
                index,
                focused
              ) => (
                <div
                  ref={ref => {
                    if (focused) {
                      setFocusIndex(index);
                      mentionItemRef.current = ref;
                    }
                  }}
                  className={`user ${focused ? 'focused' : ''}`}
                >
                  <DuelEmojiItem
                    emoji={
                      duelEmojis.find(emoji => emoji.id === suggestion.id)!
                    }
                    inputRef={inputRef}
                    showShortName
                  />
                </div>
              )}
            />
          </MentionsInput>
          {chatError === 'locked' ? (
            <LockedChat
              key={animateWarningKey}
              unlockAmount={wagerAmountToUnlockChat}
            />
          ) : chatError === 'banned' ? (
            <BannedUser />
          ) : chatError === 'muted' ? (
            <MutedUser />
          ) : chatError === 'notEnoughPeople' ? (
            <NotEnoughPeople />
          ) : chatError === 'notEnoughChip' ? (
            <NotEnoughChip />
          ) : null}
        </Box>
        <Flex justifyContent="space-between" mt="5px" pr="17px">
          <Flex alignItems="center">
            <Button
              type="button"
              background="transparent"
              onClick={onPresentEmojiModal}
            >
              <img src={frogImg} alt="frog" />
            </Button>

            <RuleBtn onClick={onPresentRulesModal}>RULES</RuleBtn>
          </Flex>
          <Flex alignItems="center" gap={10}>
            {!sendable && (
              <CountdownCircleTimer
                duration={chatCooldown}
                isPlaying
                colors={'#4F617B'}
                size={19}
                strokeWidth={2.5}
                trailColor={'#161D26'}
              >
                {({ remainingTime }) => (
                  <Span
                    fontSize="11.5px"
                    color="#4F617B"
                    fontWeight={500}
                    lineHeight={0}
                    mt="1px"
                  >
                    {remainingTime}
                  </Span>
                )}
              </CountdownCircleTimer>
            )}

            <Span
              color="#4F617B"
              fontWeight={500}
              width="32px"
              textAlign="center"
            >
              {remainLetters}
            </Span>
            <SendBtn
              type="submit"
              disabled={remainLetters < 0 || !isAuthenticated}
            >
              Send
            </SendBtn>
          </Flex>
        </Flex>
        <CommandModal
          inputRef={inputRef}
          show={showCommandModal}
          onDismiss={() => setShowCommandModal(false)}
        />
      </form>
    </>
  );
}
