import { useMemo } from 'react';
import styled from 'styled-components';

import person from 'assets/imgs/icons/person.svg';

import { BetPlayer as Player } from 'api/types/jackpot';
import { Box, Flex, Grid, Span, Text, Chip } from 'components';
import { useAppSelector } from 'state';

import {
  UsdAmount,
  StackedImage,
  StyledNotification
} from 'pages/Jackpot/components/UserStatus/styles';
import UserStatus from 'pages/Jackpot/components/UserStatus';
import { NFT } from 'api/types/nft';
import Logo from 'components/Icon/Logo';
import { convertBalanceToChip } from 'utils/balance';

const Persons = styled(Flex)`
  font-size: 14px;
  line-height: 17px;
  gap: 3px;
  align-items: center;
`;

const ListContainer = styled(Flex)`
  flex-direction: column;
  gap: 18px;

  font-weight: 600;
  font-size: 14px;
  line-height: 17px;
  /* letter-spacing: 0.18em; */
  color: white;
  width: 100%;
`;

const DataContainer = styled(Flex)`
  width: 100%;
  justify-content: space-between;
  align-items: center;
  color: ${({ theme }) => theme.colors.success};
  margin-left: 2px;
  margin-bottom: 35px;
`;

const UsersContainer = styled(Flex)`
  flex-direction: column;
  gap: 15px;
`;

const AdminBet = styled(Grid)`
  grid-template-columns: auto;
  gap: 20px;
  background: linear-gradient(90deg, #1a3452 0%, #132131 56.14%);
  border-radius: 8px;
  padding: 26px 26px 22px 25px;
  .width_700 & {
    grid-template-columns: 1fr 1fr;
    gap: 4.5%;
  }
`;

const Most = styled(Flex)`
  /* grid-template-columns: auto; */
  width: 100%;
  flex-direction: column;
  gap: 24px;
  .width_700 & {
    /* grid-template-columns: 1fr 1fr; */
    flex-direction: row;
  }
`;

const Container = styled(Box)`
  background: #0f1a26;
  border: 0;
  border-radius: 13px;
  padding: 27px 25px 40px 25px;
`;

export interface GrandPlayersProps {
  handleShowNFT?: any;
  players: Player[];
  winnerId: number;
  nftsToShow?: number;
}

export default function GrandPlayers({
  players,
  handleShowNFT,
  winnerId,
  nftsToShow = 3
}: GrandPlayersProps) {
  const { id, name } = useAppSelector(state => state.user);

  const [users, admins] = useMemo(() => {
    const users: Player[] = [];
    const admins: Player[] = [];
    players.forEach(p => {
      if (p.role === 'admin') admins.push(p);
      else users.push(p);
    });
    return [users, admins];
  }, [players]);

  const [adminUsd, adminNfts] = useMemo(() => {
    let adminUsd = 0;
    let adminNfts: NFT[] = [];
    admins.forEach(p => {
      adminUsd += p.usdAmount;
      adminNfts = [...adminNfts, ...(p.nfts ?? [])];
    });
    return [adminUsd, adminNfts];
  }, [admins]);

  const [userDatas, whales, snipers] = useMemo(() => {
    let userDatas = users
      .filter(p => p.role !== 'admin')
      .map(player => {
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
            percent: player.percent!
          },
          nfts: player.nfts ?? [],
          amount: { usd, nft, total }
        };
      });
    const playerLength = users.length > 0 ? users.length : 1;
    const whalePercent = 100 / playerLength;
    const sniperPercent = 100 / playerLength / 2;
    const whales =
      userDatas.length >= 6
        ? userDatas.filter(item => item.user.percent > whalePercent).slice(0, 3)
        : [];
    const snipers =
      userDatas.length >= 6
        ? userDatas
            .filter(item => item.user.percent < sniperPercent)
            .slice(-3)
            .reverse()
        : [];
    userDatas = userDatas.filter(p => {
      if (whales.find(whale => p.user.id === whale.user.id)) return false;
      if (snipers.find(whale => p.user.id === whale.user.id)) return false;
      return true;
    });
    return [userDatas, whales, snipers];
  }, [id, name, users]);

  return (
    <Container>
      <DataContainer>
        <Persons>
          <Span fontWeight={600}>{users.length}/∞</Span>
          <img src={person} alt="" width={13} height={14} />
        </Persons>
      </DataContainer>
      <Flex flexDirection="column" gap={34} px="5px">
        <AdminBet>
          <Flex flexDirection="column" gap={15}>
            <Logo />
            <Text
              color="textWhite"
              fontWeight={500}
              fontSize="14px"
              lineHeight="17px"
            >
              NFTs and CHIPS submitted by the Duel team and admins have a 0%
              chance of winning the Grand Jackpot when the time is up.
            </Text>
            <div>
              <Span
                color="text"
                fontWeight={500}
                fontSize="12px"
                lineHeight="14px"
              >
                Win:{' '}
              </Span>
              <Span
                color="success"
                fontWeight={500}
                fontSize="12px"
                lineHeight="14px"
              >
                0.00%
              </Span>
            </div>
          </Flex>
          {admins.length > 0 && (
            <Flex mt="8px" justifyContent="space-between">
              <Flex flexDirection="column" alignItems="center" gap={20}>
                {adminNfts.length > 0 && (
                  <>
                    <Text
                      color="text"
                      fontWeight={500}
                      fontSize="14px"
                      lineHeight="17px"
                    >
                      NFTs by Duel
                    </Text>
                    <Flex
                      flexDirection="row"
                      gap={4}
                      alignItems="center"
                      marginLeft={20}
                      marginRight={20}
                      position="relative"
                      height={45}
                      style={{
                        cursor: adminNfts.length > 0 ? 'pointer' : 'auto'
                      }}
                      onClick={() =>
                        handleShowNFT &&
                        adminNfts.length > 0 &&
                        handleShowNFT({
                          nfts: adminNfts,
                          name: 'Admins',
                          level: 0
                        })
                      }
                    >
                      {adminNfts
                        .slice(0, Math.min(nftsToShow, adminNfts.length))
                        .map(nft => (
                          <StackedImage
                            key={nft.image}
                            src={nft.image}
                            alt={nft.image}
                          />
                        ))}
                      {adminNfts.length > nftsToShow && (
                        <StyledNotification>
                          {adminNfts.length - nftsToShow}
                        </StyledNotification>
                      )}
                    </Flex>
                  </>
                )}
              </Flex>

              {adminUsd > 0 && (
                <Flex flexDirection="column" alignItems="center" gap={20}>
                  <Text
                    color="text"
                    fontWeight={500}
                    fontSize="14px"
                    lineHeight="17px"
                  >
                    CHIPs by Duel
                  </Text>
                  <UsdAmount>
                    <Chip
                      price={convertBalanceToChip(adminUsd)}
                      fontSize="17px"
                      color="chip"
                    />
                  </UsdAmount>
                </Flex>
              )}
            </Flex>
          )}
        </AdminBet>
        {(whales.length > 0 || snipers.length > 0) && (
          <Most>
            {whales.length > 0 && (
              <ListContainer>
                WHALES
                <UsersContainer>
                  {whales.map(user => {
                    return (
                      <UserStatus
                        {...user}
                        handleShowNFT={handleShowNFT}
                        key={'whales' + user.user.id}
                        background={
                          user.user.id === winnerId
                            ? 'linear-gradient(90deg, rgba(255, 226, 75, 0.2) 0%, rgba(255, 226, 75, 0) 100%)'
                            : 'linear-gradient(90deg, #182738 0%, #182738 100%)'
                        }
                      />
                    );
                  })}
                </UsersContainer>
              </ListContainer>
            )}
            {snipers.length > 0 && (
              <ListContainer>
                SNIPERS
                <UsersContainer>
                  {snipers.map(user => {
                    return (
                      <UserStatus
                        {...user}
                        handleShowNFT={handleShowNFT}
                        key={'snipers' + user.user.id}
                        background={
                          user.user.id === winnerId
                            ? 'linear-gradient(90deg, rgba(255, 226, 75, 0.2) 0%, rgba(255, 226, 75, 0) 100%)'
                            : 'linear-gradient(90deg, #182738 0%, #182738 100%)'
                        }
                      />
                    );
                  })}
                </UsersContainer>
              </ListContainer>
            )}
          </Most>
        )}
        {users.length > 0 && (
          <ListContainer>
            DUELERS
            <UsersContainer>
              {userDatas.map(user => {
                return (
                  <UserStatus
                    {...user}
                    handleShowNFT={handleShowNFT}
                    key={'duelers' + user.user.id}
                    background={
                      user.user.id === winnerId
                        ? 'linear-gradient(90deg, rgba(255, 226, 75, 0.2) 0%, rgba(255, 226, 75, 0) 100%)'
                        : 'linear-gradient(90deg, #182738 0%, #182738 100%)'
                    }
                  />
                );
              })}
            </UsersContainer>
          </ListContainer>
        )}
      </Flex>
    </Container>
  );
}
