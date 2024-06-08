import { useMemo } from 'react';
import styled, { css } from 'styled-components';

import { Flex, Grid, Span, Avatar, Chip } from 'components';
import { formatUserName } from 'utils/format';
import { CrashGameStatus } from 'api/types/crash';
import { PaidBalanceType } from 'api/types/chip';
import { convertBalanceToChip } from 'utils/balance';

export interface UserStatusProps {
  user: {
    id: number;
    name: string;
    avatar: string;
  };
  betAmount?: number;
  multiplier?: number;
  payoutMultiplier?: number;
  status?: CrashGameStatus;
  isUser?: boolean;
  cashOutAt?: number;
  maxCashOut: number;
  balanceType: PaidBalanceType;
  profit?: number;
}

export default function UserStatus({
  user,
  betAmount = 0,
  multiplier = 1,
  payoutMultiplier,
  status = 'bet',
  maxCashOut,
  balanceType,
  profit,
  isUser = false
}: UserStatusProps) {
  const [winType, color] = useMemo(() => {
    let winType = -1;
    let color = 'warning';
    if (status === 'bet' || status === 'ready') {
      winType = 0;
      color = 'white';
    } else if (status === 'play') {
      if (payoutMultiplier) {
        winType = 1;
        color = 'success';
      } else {
        winType = 0;
        color = 'white';
      }
    } else if (payoutMultiplier) {
      winType = 1;
      color = 'success';
    }
    return [winType, color];
  }, [payoutMultiplier, status]);

  return (
    <Container isUser={isUser}>
      <Flex gap={6} alignItems="center">
        <Avatar
          userId={user.id}
          name={user.name}
          image={user.avatar}
          border="none"
          borderRadius="100%"
          padding="0px"
          size="22px"
        />
        <Span
          fontWeight={500}
          fontSize="12px"
          color="white"
          style={{
            width: '115px',
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis'
          }}
        >
          {formatUserName(user.name)}
        </Span>
      </Flex>
      <Span color={color} fontWeight={500} fontSize="12px">
        {winType >= 0 && payoutMultiplier
          ? `${payoutMultiplier.toFixed(2)}x`
          : '-'}
      </Span>
      <Chip
        decimal={2}
        prefix={winType > 0 ? '+' : winType < 0 ? '-' : ''}
        chipType={balanceType}
        price={
          winType === 1
            ? convertBalanceToChip(Math.min(profit!, maxCashOut)).toFixed(2)
            : winType === 0
            ? convertBalanceToChip(
                Math.min(betAmount * multiplier, maxCashOut)
              ).toFixed(2)
            : convertBalanceToChip(betAmount).toFixed(2)
        }
        fontSize="12px"
        size={12}
        fontWeight={700}
        color={color}
      />
    </Container>
  );
}

const Container = styled(Grid)<{ isUser: boolean }>`
  grid-template-columns: 115px 1fr max-content;
  background: #182738;
  border-radius: 8px;
  align-items: center;
  padding: 7px 13px 7px 15px;

  ${({ isUser }) =>
    isUser
      ? css`
          border: 1px solid #4fff8b;
          box-shadow: 0px 0px 10px rgba(79, 255, 139, 0.25);
        `
      : css``}
`;
