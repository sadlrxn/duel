import React, { useCallback } from 'react';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import { ReactComponent as CloseIcon } from 'assets/imgs/icons/close.svg';
import { ReactComponent as MuteIcon } from 'assets/imgs/icons/mute.svg';
import { ReactComponent as UnmuteIcon } from 'assets/imgs/icons/unmute.svg';
import { ReactComponent as BanIcon } from 'assets/imgs/icons/ban.svg';
import { ReactComponent as UnbanIcon } from 'assets/imgs/icons/handshake.svg';

import { Message, ChatUser } from 'api/types/chat';
import Avatar from 'components/Avatar';
import { BoxProps, Flex, Grid } from 'components/Box';
import { Chip } from 'components/Chip';
import { Span } from 'components/Text';
import { ReplyIcon } from 'components/Icon';
import { duelEmojis, generateDuelEmojiUrl } from 'config';
import { useAppDispatch } from 'state';
import { setReply } from 'state/actions';

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
import { MsgSpan, ReplyMsgBox } from './styles';
import { formatUserName } from 'utils/format';
import { MaxLength, WagerLimit, Rain } from './utils';
import { convertBalanceToChip } from 'utils/balance';

export interface ReplyMsgProps extends BoxProps {
  replyMsg?: Message;
  showCloseIcon?: boolean;
  handleReplyMsgClick?: any;
}

export default function ReplyMsg({
  replyMsg,
  showCloseIcon = false,
  handleReplyMsgClick,
  ...props
}: ReplyMsgProps) {
  const dispatch = useAppDispatch();

  const handleReplyCancel = useCallback(() => {
    dispatch(setReply());
  }, [dispatch]);

  if (!replyMsg) return <></>;

  const {
    author: { name, avatar, id, role },
    message,
    time,
    deleted,
    id: msgId
  } = replyMsg;

  if (deleted)
    return (
      <ReplyMsgBox
        as={Grid}
        background="#1B283A"
        px="13px"
        pt="5px"
        pb="13px"
        mr="10px"
        {...props}
      >
        <Span fontStyle="italic" fontSize="13px">
          Original message was deleted.
        </Span>
      </ReplyMsgBox>
    );

  return (
    <>
      <ReplyMsgBox
        onClick={() => handleReplyMsgClick(msgId)}
        cursor="pointer"
        as={Grid}
        gridTemplateColumns="40px 1fr"
        background="#1B283A"
        px="13px"
        pt="5px"
        pb="13px"
        mr="10px"
        {...props}
      >
        <Flex justifyContent="space-between" alignItems="center">
          <ReplyIcon size={10} />
          <Avatar
            userId={id}
            name={name}
            image={avatar}
            role={role}
            size="20px"
          />
        </Flex>
        <Flex
          flexDirection="row"
          flexWrap="wrap"
          alignItems="center"
          gap={4}
          fontSize={14}
          fontWeight={400}
          mr="15px"
        >
          <Span fontWeight={500}>
            @{formatUserName(replyMsg.author.name)}:{' '}
          </Span>
          {message
            .replaceAll('\n', ' ')
            .slice(0, message[0] === '$' ? undefined : 50)
            .concat(message.length >= 50 && message[0] !== '$' ? '...' : '')
            .split(notifyRegex)
            .map((str, i) => {
              if (str.startsWith('@')) {
                const notifyUserName = str.substring(
                  str.indexOf('{') + 1,
                  str.indexOf('}')
                );
                return (
                  <Span fontWeight={600}>
                    @{formatUserName(notifyUserName)}
                  </Span>
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
                  <>
                    <Chip price={convertBalanceToChip(amount)} />
                    tip for @{formatUserName(to.name)}
                  </>
                );
              }
              if (str.startsWith('/mute') && str.match(muteRegex)) {
                const user = str.substring(str.indexOf(' ') + 1, str.length);
                return (
                  <>
                    <MuteIcon /> Has muted @{user} for 15 minutes.
                  </>
                );
              }
              if (str.startsWith('/unmute') && str.match(unmuteRegex)) {
                const user = str.substring(str.indexOf(' ') + 1, str.length);
                return (
                  <>
                    <UnmuteIcon /> Has unmuted @{user}
                  </>
                );
              }
              if (str.startsWith('/ban') && str.match(banRegex)) {
                const user = str.substring(str.indexOf(' ') + 1, str.length);
                return (
                  <>
                    <BanIcon /> Has banned @{user}
                  </>
                );
              }
              if (str.startsWith('/unban') && str.match(unbanRegex)) {
                const user = str.substring(str.indexOf(' ') + 1, str.length);
                return (
                  <>
                    <UnbanIcon /> Has unbanned @{user}
                  </>
                );
              }
              if (
                str.startsWith('/setMaxLength') &&
                str.match(setMaxLengthRegex)
              ) {
                const length = str.substring(str.indexOf(' ') + 1, str.length);
                return <MaxLength key={`${time}_${i}`} length={length} />;
              }
              if (
                str.startsWith('/setWagerLimit') &&
                str.match(setWagerLimitRegex)
              ) {
                const limit = str.substring(str.indexOf(' ') + 1, str.length);
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
                  const emoji = duelEmojis.find(emoji => emoji.id === emojiId);
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
        {showCloseIcon && (
          <CloseIcon
            cursor="pointer"
            style={{ position: 'absolute', right: '13px', top: '10px' }}
            width={10}
            height={10}
            onClick={handleReplyCancel}
          />
        )}
      </ReplyMsgBox>
    </>
  );
}
