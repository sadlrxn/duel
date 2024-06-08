import React, { useMemo } from 'react';
import styled from 'styled-components';

import person from 'assets/imgs/icons/person.svg';
import arrow from 'assets/imgs/icons/arrow-right.svg';
import { BetPlayer as Player } from 'api/types/jackpot';
import { Flex, Span, Text, Snow1, Snow2 } from 'components';
import { useAppSelector } from 'state';

import UserStatus from '../UserStatus';
import { getTotalNfts } from '../../utils';
import { convertBalanceToChip } from 'utils/balance';

const Persons = styled(Flex)`
  font-size: 14px;
  line-height: 17px;
  gap: 3px;
  align-items: center;
`;

const StyledText = styled(Text)`
  display: flex;
  justify-content: center;
  align-items: center;
  background: #6d81a2;
  color: #0b141e;
  border-radius: 17px;
  font-size: 16px;
  font-weight: 600;
  line-height: 19px;
  padding: 1px 7px;
`;

const NFTs = styled(Flex)`
  color: #6d81a2;
  align-items: center;
  gap: 5px;
`;

const DataContainer = styled(Flex)`
  width: 100%;
  justify-content: space-between;
  align-items: center;
  color: ${({ theme }) => theme.colors.success};
  margin-left: 2px;
  margin-bottom: 15px;
`;

const UsersContainer = styled(Flex)`
  flex-direction: column;
  gap: 9px;
  overflow: hidden auto;
  padding-right: 8px;
`;

const Container = styled(Flex)`
  position: relative;
  flex-direction: column;
  flex: none;
  background: #0f1a26;
  border: 0;
  border-radius: 13px;
  gap: 9px;
  padding: 20px 20px 13px 18px;
  max-height: 710px;
`;

export interface PlayersProps {
  handleShowNFT?: any;
  players: Player[];
  winnerId: number;
}

export default function Players({
  players,
  handleShowNFT,
  winnerId
}: PlayersProps) {
  const meta = useAppSelector(state => state.meta.jackpot[state.jackpot.room]);
  const { id, name, isHoliday } = useAppSelector(state => state.user);

  const totalUsdAmount = useMemo(() => {
    return players.map(p => p.usdAmount).reduce((partial, a) => partial + a, 0);
  }, [players]);

  const totalNftAmount = useMemo(() => {
    return players
      .map(p => p.nftAmount)
      .reduce((partial: number, a) => partial + (a ?? 0), 0);
  }, [players]);

  const totalBetAmount = useMemo(() => {
    return totalUsdAmount + totalNftAmount;
  }, [totalUsdAmount, totalNftAmount]);

  const userDatas = useMemo(() => {
    return players.map(player => {
      // TODO: Check
      const usd = convertBalanceToChip(player.usdAmount);
      const nft = convertBalanceToChip(player.nftAmount ?? 0);
      const total = usd + nft;
      return {
        user: {
          id: player.id,
          avatar: player.avatar,
          level: 1,
          name: player.id === id ? name : player.name,
          percent:
            ((player.usdAmount + (player.nftAmount ?? 0)) * 100) /
            totalBetAmount
        },
        nfts: player.nfts ?? [],
        amount: { usd, nft, total }
      };
    });
  }, [totalBetAmount, id, name, players]);

  const bettedNfts = useMemo(() => {
    return getTotalNfts(players);
  }, [players]);

  return (
    <Container>
      {isHoliday && (
        <>
          <Snow1 position="absolute" top={-17} left={-13} />
          <Snow2 position="absolute" top={-13} right={-7} />
        </>
      )}
      <DataContainer>
        <Persons>
          <Span fontWeight={600}>
            {players.length}/{meta.playerLimit}
          </Span>
          <img src={person} alt="" width={13} height={14} />
        </Persons>
        <NFTs>
          <StyledText>{bettedNfts.length}</StyledText>
          <Flex
            alignItems="center"
            gap={5}
            style={{ cursor: bettedNfts.length > 0 ? 'pointer' : 'auto' }}
            onClick={() =>
              handleShowNFT &&
              bettedNfts.length > 0 &&
              handleShowNFT({ nfts: bettedNfts, name: '', leve: 0 })
            }
          >
            <Span>NFT{bettedNfts.length !== 1 && 's'}</Span>
            <img src={arrow} alt="" width={12} height={12} />
          </Flex>
        </NFTs>
      </DataContainer>
      <UsersContainer>
        {userDatas.map(user => {
          return (
            <UserStatus
              {...user}
              handleShowNFT={handleShowNFT}
              key={user.user.id}
              background={
                user.user.id === winnerId
                  ? 'linear-gradient(90deg, rgba(255, 226, 75, 0.2) 0%, rgba(255, 226, 75, 0) 100%)'
                  : 'linear-gradient(90deg, #182738 0%, #182738 100%)'
              }
            />
          );
        })}
      </UsersContainer>
    </Container>
  );
}
