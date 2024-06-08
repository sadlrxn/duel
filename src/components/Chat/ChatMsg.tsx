import React, { FC, useMemo, useCallback } from 'react';
import { LazyLoadImage } from 'react-lazy-load-image-component';
import dayjs from 'dayjs';

import { Message, ChatUser } from 'api/types/chat';
import { useModal } from 'components/Modal';
import { RainModal } from 'components/Modals';
import Avatar from 'components/Avatar';
import { Box, Flex, Grid } from 'components/Box';
import { Text } from 'components/Text';
import { Badge } from 'components/Badge';
import { Chip } from 'components/Chip';
import { HeartIcon, ReplyIcon, RubbishBinIcon } from 'components/Icon';
import { duelEmojis, generateDuelEmojiUrl } from 'config';
import { useAppSelector, useAppDispatch } from 'state';
import { formatUserName } from 'utils/format';

import {
  notifyRegex,
  duelEmojiRegex,
  tipMsgRegex,
  muteRegex,
  unmuteRegex,
  banRegex,
  unbanRegex,
  setMaxLengthRegex,
  setWagerLimitRegex,
  rainRegex
} from './config';
import {
  BubbleBtn,
  BubbleContainer,
  HistoryTime,
  MsgBox,
  MsgSpan
} from './styles';
import {
  Notifier,
  TipTransfer,
  Mute,
  Unmute,
  Ban,
  Unban,
  MaxLength,
  WagerLimit,
  Rain
} from './utils';
import ReplyMsg from './ReplyMsg';
import { sendMessage } from 'state/socket';
import { setReply } from 'state/actions';
import { convertBalanceToChip } from 'utils/balance';

const ChatMsg: FC<{
  msg: Message;
  inputRef: React.MutableRefObject<any>;
  handleReplyMsgClick?: any;
}> = ({ msg, inputRef, handleReplyMsgClick }) => {
  const dispatch = useAppDispatch();

  const {
    name: currentUserName,
    role: currentUserRole,
    id: currentUserId
  } = useAppSelector(state => state.user);
  const differ = useAppSelector(state => state.socket.differ);
  const msgs = useAppSelector(state => state.chat.msgs);
  const {
    author: { name, avatar, id, role },
    message,
    time,
    replyTo,
    sponsors,
    deleted,
    id: msgId
  } = msg;

  const replyMsg = useMemo(() => {
    if (!replyTo) return undefined;
    return msgs.find(msg => msg.id === replyTo);
  }, [replyTo, msgs]);

  const color = useMemo(() => {
    let color = '#FFFFFF';
    const isUser = id === currentUserId;
    switch (role) {
      case 'admin':
        color = isUser ? '#DB00FF' : '#9E00FF';
        break;
      case 'moderator':
        color = isUser ? '#FFC700' : '#FFA800';
        break;
      case 'ambassador':
        color = isUser ? '#00FFFF' : '#00C4FF';
        break;
      default:
        break;
    }
    return color;
  }, [currentUserId, id, role]);

  const isSponsor = useMemo(() => {
    return sponsors?.find(s => s === currentUserId);
  }, [currentUserId, sponsors]);

  const handleDelete = useCallback(() => {
    dispatch(
      sendMessage({
        type: 'delete',
        room: 'chat',
        content: JSON.stringify({
          id: msgId
        })
      })
    );
  }, [dispatch, msgId]);

  const handleReply = useCallback(
    (e: any) => {
      console.info(e);
      dispatch(setReply(msgId));
      inputRef.current && inputRef.current.focus();
    },
    [dispatch, inputRef, msgId]
  );

  const handleSponsor = useCallback(() => {
    dispatch(
      sendMessage({
        type: 'sponsor',
        room: 'chat',
        content: JSON.stringify({
          id: msgId
        })
      })
    );
  }, [dispatch, msgId]);

  const [isRain, amount, duelers, rainToUser] = useMemo(() => {
    const isRain = message.startsWith('/rain');
    const amount = message.split(' ')[2] ? +message.split(' ')[2] : 0;
    let duelers = [];
    if (isRain && message.split(' ')[3]) {
      try {
        duelers = JSON.parse(message.split(' ')[3]);
      } catch {
        duelers = [];
      }
    }
    const rainToUser =
      isRain &&
      duelers.findIndex((dueler: any) => dueler.id === currentUserId) !== -1;
    return [isRain, amount, duelers, rainToUser];
  }, [message, currentUserId]);

  const [onRainModal] = useModal(
    <RainModal duelers={duelers} amount={amount} rainer={msg.author} />,
    true
  );

  return (
    <>
      <Box>
        {!deleted && (
          <ReplyMsg
            replyMsg={replyMsg}
            handleReplyMsgClick={handleReplyMsgClick}
            mb="-8px"
          />
        )}
        <MsgBox
          className="chat_original_msg"
          as={Grid}
          gridTemplateColumns={deleted ? '1fr' : '40px 1fr'}
          onClick={
            deleted
              ? undefined
              : isRain
              ? (e: any) => {
                  e.stopPropagation();
                  onRainModal();
                }
              : undefined
          }
          cursor={isRain ? 'pointer' : 'auto'}
        >
          {deleted ? (
            <Text color="#B2D1FF" fontWeight={600} fontSize="14px">
              This message has been deleted.
            </Text>
          ) : (
            <>
              <Avatar userId={id} name={name} image={avatar} role={role} />
              <Flex flexDirection={'column'} alignItems="start" gap={7}>
                <Flex justifyContent="space-between" width="100%">
                  {role === 'admin' ? (
                    <Flex gap={8} alignItems="center">
                      <Text color={color} fontSize={13.5} fontWeight={600}>
                        {formatUserName(name)}
                      </Text>
                      <Badge color={color} background={'#3B156A'}>
                        Admin
                      </Badge>
                    </Flex>
                  ) : role === 'moderator' ? (
                    <Flex gap={8} alignItems="center">
                      <Text color={color} fontSize={13.5} fontWeight={600}>
                        {formatUserName(name)}
                      </Text>
                      <Badge color={color} background={'#6A4815'}>
                        Mod
                      </Badge>
                    </Flex>
                  ) : role === 'ambassador' ? (
                    <Flex gap={8} alignItems="center">
                      <Text color={color} fontSize={13.5} fontWeight={600}>
                        {formatUserName(name)}
                      </Text>
                      <Badge color={color} background={'#15566A'}>
                        Ambassador
                      </Badge>
                    </Flex>
                  ) : (
                    <Text
                      color={name === currentUserName ? '#4FFF8B' : '#fff'}
                      fontSize={13.5}
                      fontWeight={600}
                    >
                      {formatUserName(name)}
                    </Text>
                  )}
                  <Flex gap={4}>
                    {sponsors && sponsors.length > 0 && (
                      <Text
                        fontSize="11px"
                        color={isSponsor ? 'success' : '#788CA9'}
                        lineHeight={1.2}
                      >
                        {sponsors.length}
                      </Text>
                    )}
                    <HeartIcon
                      borderColor={isSponsor ? '#4FFF8B' : '#788CA9'}
                      background={isSponsor ? '#4FFF8B' : 'transparent'}
                    />
                  </Flex>
                </Flex>

                <Flex
                  flexDirection="row"
                  flexWrap="wrap"
                  alignItems="center"
                  gap={4}
                  fontSize={14}
                >
                  {message.split(notifyRegex).map((str, i) => {
                    if (str.startsWith('@')) {
                      const notifyUserName = str.substring(
                        str.indexOf('{') + 1,
                        str.indexOf('}')
                      );
                      return (
                        <Notifier
                          key={`${time}_${i}`}
                          notifier={notifyUserName}
                        />
                      );
                    }
                    if (str.match(tipMsgRegex)) {
                      // tipMsgRegex.test(str) I am not sure why this regex checking not render react component, weired
                      const amount = parseFloat(
                        str.substring(2, str.indexOf(' ', 3))
                      );
                      const to: ChatUser = JSON.parse(
                        str.substring(str.indexOf(' ', 3), str.length)
                      );
                      return (
                        <TipTransfer
                          key={`${time}_${i}`}
                          to={to}
                          amount={amount}
                        />
                      );
                    }
                    if (str.startsWith('/mute') && str.match(muteRegex)) {
                      const user = str.substring(
                        str.indexOf(' ') + 1,
                        str.length
                      );
                      return <Mute key={`${time}_${i}`} user={user} />;
                    }
                    if (str.startsWith('/unmute') && str.match(unmuteRegex)) {
                      const user = str.substring(
                        str.indexOf(' ') + 1,
                        str.length
                      );
                      return <Unmute key={`${time}_${i}`} user={user} />;
                    }
                    if (str.startsWith('/ban') && str.match(banRegex)) {
                      const user = str.substring(
                        str.indexOf(' ') + 1,
                        str.length
                      );
                      return <Ban key={`${time}_${i}`} user={user} />;
                    }
                    if (str.startsWith('/unban') && str.match(unbanRegex)) {
                      const user = str.substring(
                        str.indexOf(' ') + 1,
                        str.length
                      );
                      return <Unban key={`${time}_${i}`} user={user} />;
                    }
                    if (
                      str.startsWith('/setMaxLength') &&
                      str.match(setMaxLengthRegex)
                    ) {
                      const length = str.substring(
                        str.indexOf(' ') + 1,
                        str.length
                      );
                      return <MaxLength key={`${time}_${i}`} length={length} />;
                    }
                    if (
                      str.startsWith('/setWagerLimit') &&
                      str.match(setWagerLimitRegex)
                    ) {
                      const limit = str.substring(
                        str.indexOf(' ') + 1,
                        str.length
                      );
                      return <WagerLimit key={`${time}_${i}`} limit={limit} />;
                    }
                    if (str.startsWith('/rain') && str.match(rainRegex)) {
                      const split = str.split(' ');
                      return (
                        <Rain
                          key={`${time}_${i}`}
                          split={+split[1]}
                          amount={+split[2]}
                        />
                      );
                    }

                    let containsOnlyEmoji = true;
                    str.split(duelEmojiRegex).forEach(str => {
                      if (!str.startsWith(':{') && str.length > 1) {
                        containsOnlyEmoji = false;
                      }
                    });

                    return str.split(duelEmojiRegex).map((str2, j) => {
                      if (str2.startsWith(':{')) {
                        const emojiId = str2.substring(2, str2.length - 1);
                        const emoji = duelEmojis.find(
                          emoji => emoji.id === emojiId
                        );
                        if (emoji)
                          return (
                            <LazyLoadImage
                              key={`chat_message_${time}_emoji_${i}_${j}`}
                              width={containsOnlyEmoji ? 28 : 20}
                              height={containsOnlyEmoji ? 28 : 20}
                              src={generateDuelEmojiUrl(emoji)}
                              alt={emojiId}
                            />
                          );
                      }

                      return (
                        <MsgSpan
                          key={`chat_message_${time}_content_${i}_${j}`}
                          fontWeight={400}
                        >
                          {str2}
                        </MsgSpan>
                      );
                    });
                  })}
                </Flex>
                <HistoryTime sound={false}>
                  {dayjs(new Date(time + differ).getTime()).format('hh:mm A')}
                </HistoryTime>
                <BubbleContainer>
                  {(currentUserRole === 'admin' ||
                    currentUserRole === 'moderator') && (
                    <BubbleBtn type="button" fill onClick={handleDelete}>
                      <RubbishBinIcon />
                    </BubbleBtn>
                  )}
                  <BubbleBtn type="button" stroke onClick={handleReply}>
                    <ReplyIcon />
                  </BubbleBtn>
                  <BubbleBtn type="button" onClick={handleSponsor}>
                    <HeartIcon
                      size={16}
                      borderColor={isSponsor ? '#4FFF8B' : '#788CA9'}
                      background={isSponsor ? '#4FFF8B' : 'transparent'}
                    />
                  </BubbleBtn>
                </BubbleContainer>
              </Flex>
            </>
          )}
        </MsgBox>
      </Box>
      {!deleted && rainToUser && (
        <Box>
          <MsgBox as={Grid} gridTemplateColumns={deleted ? '1fr' : '40px 1fr'}>
            <Avatar userId={id} name={name} image={avatar} role={role} />
            <Flex flexDirection={'column'} alignItems="start" gap={7}>
              <Flex justifyContent="space-between" width="100%">
                {role === 'admin' ? (
                  <Flex gap={8} alignItems="center">
                    <Text color={color} fontSize={13.5} fontWeight={600}>
                      {formatUserName(name)}
                    </Text>
                    <Badge color={color} background={'#3B156A'}>
                      Admin
                    </Badge>
                  </Flex>
                ) : role === 'moderator' ? (
                  <Flex gap={8} alignItems="center">
                    <Text color={color} fontSize={13.5} fontWeight={600}>
                      {formatUserName(name)}
                    </Text>
                    <Badge color={color} background={'#6A4815'}>
                      Mod
                    </Badge>
                  </Flex>
                ) : role === 'ambassador' ? (
                  <Flex gap={8} alignItems="center">
                    <Text color={color} fontSize={13.5} fontWeight={600}>
                      {formatUserName(name)}
                    </Text>
                    <Badge color={color} background={'#15566A'}>
                      Ambassador
                    </Badge>
                  </Flex>
                ) : (
                  <Text
                    color={name === currentUserName ? '#4FFF8B' : '#fff'}
                    fontSize={13.5}
                    fontWeight={600}
                  >
                    {formatUserName(name)}
                  </Text>
                )}
              </Flex>

              <Flex
                flexDirection="row"
                flexWrap="wrap"
                alignItems="center"
                gap={4}
                fontSize={14}
              >
                <Badge
                  variant="secondary"
                  fontSize={14}
                  lineHeight="17px"
                  px="5px"
                  py="1px"
                >
                  <Chip price={convertBalanceToChip(amount) / duelers.length} />
                </Badge>{' '}
                Rained into your wallet.
                <Notifier notifier={currentUserName} />
              </Flex>
              <HistoryTime sound={false}>
                {dayjs(new Date(time + differ).getTime()).format('hh:mm A')}
              </HistoryTime>
            </Flex>
          </MsgBox>
        </Box>
      )}
    </>
  );
};

export default ChatMsg;
