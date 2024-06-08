import { FC, memo } from 'react';

import { Duel as DuelIcon, Ana as AnaIcon, HatIcon, Box } from 'components';
import crown from 'assets/imgs/icons/crown.svg';
import { formatUserName } from 'utils/format';

import {
  Container,
  Image,
  Divider,
  Avatar,
  StyledText,
  DataContainer,
  DetailContainer,
  StyledChip,
  CrownImage
} from './Side.styles';
import { useAppSelector } from 'state';
import { convertBalanceToChip } from 'utils/balance';

interface SideProps {
  side: string;
  avatar: string;
  prize: number;
  name: string;
  winner?: boolean;
  end?: boolean;
  user?: boolean;
  userId?: number;
  paidBalanceType?: 'chip' | 'coupon';
}

const Side: FC<SideProps> = ({
  side = 'duel',
  name,
  avatar = 'https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://www.arweave.net/FEFMTQEgWHhDd33e2N2ldQZ93Bk0BSVKxp7TPdP-3ao',
  prize,
  winner = false,
  end = false,
  user = false,
  userId = 0,
  paidBalanceType = 'chip'
}) => {
  const isHoliday = useAppSelector(state => state.user.isHoliday);

  return (
    <Container side={side} active={name !== '' && !(end && !winner)}>
      <Image side={side} position="relative">
        {end && winner && isHoliday ? (
          <Box position="absolute" top="-27px" left="-4px">
            <HatIcon />
          </Box>
        ) : (
          <></>
        )}
        {side === 'duel' ? <DuelIcon size={28} /> : <AnaIcon size={28} />}
      </Image>
      {name !== '' && (
        <>
          <Divider side={side} />
          <Avatar
            userId={userId}
            name={name}
            padding="unset"
            border="0px"
            side={side}
            size="34px"
            image={avatar}
          />

          <DataContainer side={side}>
            <CrownImage src={crown} alt="" show={end && winner} />
            <DetailContainer side={side}>
              <StyledText color={end && winner && user ? 'success' : 'white'}>
                {formatUserName(name)}
              </StyledText>
              {!!end && !!winner && (
                <StyledChip
                  chipType={paidBalanceType ? paidBalanceType : 'chip'}
                  // background={balanceType ? '#4BE9FF' : '#ffe24b'}
                  // border={balanceType ? '#1FAEFF' : '#ffb31f'}
                  price={convertBalanceToChip(prize)}
                  color={user ? 'success' : 'chip'}
                />
              )}
            </DetailContainer>
          </DataContainer>
        </>
      )}
    </Container>
  );
};

export default memo(Side);
