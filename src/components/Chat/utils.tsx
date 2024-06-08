import React from 'react';
import { ChatUser } from 'api/types/chat';
import { ReactComponent as MuteIcon } from 'assets/imgs/icons/mute.svg';
import { ReactComponent as UnmuteIcon } from 'assets/imgs/icons/unmute.svg';
import { ReactComponent as BanIcon } from 'assets/imgs/icons/ban.svg';
import { ReactComponent as UnbanIcon } from 'assets/imgs/icons/handshake.svg';
import { ReactComponent as LockIcon } from 'assets/imgs/icons/lock2.svg';
import { ReactComponent as ChatIcon } from 'assets/imgs/icons/chat.svg';
import { ReactComponent as WagerIcon } from 'assets/imgs/icons/wager.svg';
import { ReactComponent as RainIcon } from 'assets/imgs/icons/rain.svg';
import { Badge } from 'components/Badge';
import { Chip } from 'components/Chip';
import { useModal } from 'components/Modal';
import { ProfileModal } from 'components/Modals';
import { useAppSelector } from 'state';
import { NotifierWrapper, ChatWarningContainer } from './styles';
import { formatUserName } from 'utils/format';
import { Span } from 'components/Text';
import { convertBalanceToChip } from 'utils/balance';

export const Notifier = ({ notifier }: { notifier: string }) => {
  const { name: currentUserName } = useAppSelector(state => state.user);
  const [onProfileModal] = useModal(<ProfileModal name={notifier} />, true);
  return (
    <NotifierWrapper>
      {currentUserName === notifier ? (
        <Badge
          onClick={onProfileModal}
          fontSize={14}
          lineHeight="17px"
          px="5px"
          py="3px"
        >{`@${formatUserName(notifier)}`}</Badge>
      ) : (
        <strong onClick={onProfileModal}>{`@${formatUserName(
          notifier
        )}`}</strong>
      )}
    </NotifierWrapper>
  );
};

export const TipTransfer = ({
  amount,
  to
}: {
  amount: number;
  to: ChatUser;
}) => {
  return (
    <>
      <Badge
        variant="secondary"
        fontSize={14}
        lineHeight="17px"
        px="5px"
        py="1px"
      >
        <Chip price={convertBalanceToChip(amount)} />
      </Badge>{' '}
      tip for
      <Notifier notifier={to.name} />
    </>
  );
};

export const Mute = ({
  user,
  duration = 15
}: {
  user: string;
  duration?: number;
}) => {
  return (
    <>
      <MuteIcon /> Has muted
      <Notifier notifier={user} /> for {duration} minutes.
    </>
  );
};

export const Unmute = ({ user }: { user: string }) => {
  return (
    <>
      <UnmuteIcon /> Has unmuted
      <Notifier notifier={user} />
    </>
  );
};

export const Ban = ({ user }: { user: string }) => {
  return (
    <>
      <BanIcon /> Has banned
      <Notifier notifier={user} />
    </>
  );
};

export const Unban = ({ user }: { user: string }) => {
  return (
    <>
      <UnbanIcon /> Has unbanned
      <Notifier notifier={user} />
    </>
  );
};

export const MaxLength = ({ length }: { length: string }) => {
  return (
    <div style={{ display: 'inline' }}>
      <ChatIcon /> Has changed the chat character limit to {length} characters.
    </div>
  );
};

export const WagerLimit = ({ limit }: { limit: string }) => {
  return (
    <div style={{ display: 'inline' }}>
      <WagerIcon /> Has changed the chat minimum wager to{' '}
      {Number(limit).toFixed(2)} CHIPS.
    </div>
  );
};

export const Rain = ({ split, amount }: { split: number; amount: number }) => {
  return (
    <>
      <RainIcon />
      <Span color="success">Made it RAIN! </Span>
      <Badge
        variant="secondary"
        fontSize={14}
        lineHeight="17px"
        px="5px"
        py="1px"
      >
        <Chip price={convertBalanceToChip(amount)} />
      </Badge>{' '}
      was tipped to {split} random Duelers.
    </>
  );
};

export const LockedChat = ({ unlockAmount }: { unlockAmount: number }) => {
  return (
    <ChatWarningContainer className="animate__animated animate__headShake">
      <LockIcon />
      Wager <Chip price={convertBalanceToChip(unlockAmount)} /> more to unlock
      chat
    </ChatWarningContainer>
  );
};

export const BannedUser = () => {
  return (
    <ChatWarningContainer>
      <BanIcon />
      You are banned from the chat.
    </ChatWarningContainer>
  );
};

export const MutedUser = ({ duration = 15 }: { duration?: number }) => {
  return (
    <ChatWarningContainer>
      <MuteIcon />
      You have been muted for {duration} minutes.
    </ChatWarningContainer>
  );
};

export const NotEnoughPeople = () => {
  return (
    <ChatWarningContainer>
      There’s not enough people in chat to split your tip. Use /tip max amount
      instead to split it amongst everyone equally.
    </ChatWarningContainer>
  );
};

export const NotEnoughChip = () => {
  return (
    <ChatWarningContainer>
      You don’t have enough CHIPS to Make it Rain.
    </ChatWarningContainer>
  );
};
