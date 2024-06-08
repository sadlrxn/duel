import { useMemo, memo, FC } from 'react';
import styled, { css } from 'styled-components';
import dayjs from 'dayjs';

import { Box, Text, Flex, BoxProps, Span } from 'components';
import NFTCard from '../NFTCard';
import { formatUserName } from 'utils/format';

import { NFT } from 'api/types/nft';
import { formatNumber } from 'utils/format';
import { useQuery } from 'hooks';
import { convertBalanceToChip } from 'utils/balance';

const ChipIcon = styled.div`
  background-color: #ffe24b;
  border: 0.2rem solid #ffb31f;
  border-radius: 100%;
  width: 0.6em;
  height: 0.6em;
`;

const WinnerText = styled(Text)`
  font-weight: 600;
  font-size: 44px;
  line-height: 53px;
  text-align: center;

  display: flex;
  flex-direction: column;
  gap: 5px;

  color: #fff6ca;

  text-shadow: 0px 0px 17px rgba(255, 226, 75, 0.8);
`;

const TotalChip = styled(Flex)`
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 13px 20px 13px 17px;
  font-weight: 600;
  font-size: 20px;
  color: #ffffff;
  border-radius: 9px;
  border: 4px solid #c0b264;
  background: linear-gradient(
    360deg,
    rgba(255, 226, 75, 0.2) 0%,
    rgba(255, 226, 75, 0) 100%
  );
`;

const ResultContainer = styled(Flex)`
  flex-direction: column;
  gap: 15px;
  color: #ffffff;
  font-size: 14px;
  letter-spacing: 0.18em;
  align-items: center;

  .width_700 & {
    align-items: end;
  }
`;

const TitleContainer = styled(Flex)`
  align-items: center;
  justify-content: space-between;
  gap: 26px;
  font-weight: 600;
  font-size: 14px;
  color: #fff6ca;
  flex-direction: column;

  .width_700 & {
    flex-direction: row;
  }
`;

const NFTContainer = styled(Flex)`
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
  max-width: 100%;
`;

interface ContainerProps extends BoxProps {
  variant: 'primary' | 'secondary';
  win: boolean;
}

const Container = styled(Box)<ContainerProps>`
  /* display: grid;
  grid-auto-rows: max-content; */
  display: flex;
  flex-direction: column;
  gap: 23px;
  /* position: absolute; */
  /* background-color: #12131ae0; */
  background: transparent;
  height: 80%;
  /* width: 90%; */

  ${({ variant }) => {
    if (variant === 'secondary')
      return css`
        position: fixed;
        left: 0;
        top: 0;
        width: 100vw;
        height: 100vh;
        align-items: center;
        background: #12131ae0;
        gap: 47px;
        z-index: 10;
        overflow: auto;
        padding: 100px 0;

        scrollbar-width: none;
        &::-webkit-scrollbar {
          display: none;
        }

        ${TitleContainer} {
          flex-direction: column;
        }
        ${TotalChip} {
          padding: 29px 36px 28px 45px;
          gap: 26px;
          font-size: 44px;
          border-radius: 20px;
        }
        ${ResultContainer} {
          font-size: 20px;
          align-items: center;
          gap: 28px;
          color: #fff6ca;
        }
        ${WinnerText} {
          font-size: 96px;
          line-height: 116px;
        }
        ${NFTContainer} {
          max-width: 70%;
        }
      `;
  }}

  ${({ win, variant }) => {
    if (!win)
      return css`
        ${TitleContainer} {
          justify-content: center;
        }
        ${WinnerText} {
          display: none;
        }

        ${ResultContainer} {
          font-size: ${variant === 'primary' ? 20 : 40}px;
          align-items: center;
          gap: ${variant === 'primary' ? 15 : 28}px;
          color: ${variant === 'primary' ? 'white' : '#fff6ca'};
        }
      `;
  }}
`;

export interface ContentProps {
  variant?: 'primary' | 'secondary';
  nftsToShow?: number;
  nfts?: NFT[];
  usdBetAmount?: number;
  nftBetAmount?: number;
  onClick?: any;
  onClose?: any;
  win?: boolean;
  winnerName?: string;
  usdProfit?: number;
  usdFee?: number;
  nftProfit?: NFT[];
  nftFee?: NFT[];
  roundId?: number;
  time?: number;
}

const Content: FC<ContentProps> = ({
  variant = 'primary',
  nftsToShow = 2,
  nfts = [],
  usdBetAmount = 0,
  nftBetAmount = 0,
  onClick,
  onClose,
  win = false,
  winnerName = '',
  usdProfit,
  usdFee,
  nftProfit,
  nftFee,
  roundId = 0,
  time = Date.now()
}) => {
  const query = useQuery();
  const showTime = useMemo(() => query.get('roundId'), [query]);

  const totalBetAmount = useMemo(
    () => usdBetAmount + nftBetAmount,
    [usdBetAmount, nftBetAmount]
  );

  return (
    <Container variant={variant} win={win} onClick={onClose}>
      <TitleContainer mb={'15px'}>
        <WinnerText>
          {showTime !== null && (
            <Flex
              gap={19}
              flexWrap="wrap"
              fontSize="14px"
              lineHeight="1.3em"
              letterSpacing="0.18em"
            >
              <Span color="white" fontWeight={600}>
                #{roundId}
              </Span>
              <Span color="white" fontWeight={600}>
                {dayjs(time).format('MMM DD, hh:mm A')}
              </Span>
            </Flex>
          )}
          Won By {formatUserName(winnerName)}
        </WinnerText>
        <ResultContainer>
          {variant !== 'primary' && 'TOTAL '}JACKPOT VALUE
          <TotalChip>
            <ChipIcon />
            {formatNumber(totalBetAmount)}
          </TotalChip>
        </ResultContainer>
      </TitleContainer>
      {nftProfit !== undefined && nftProfit !== undefined ? (
        <>
          <NFTContainer>
            <NFTCard
              price={convertBalanceToChip(usdProfit!)}
              type="chip"
              clickable={false}
            />
            {nftProfit!
              .slice(0, variant === 'primary' ? nftsToShow : nftProfit!.length)
              .map(nft => {
                return (
                  <NFTCard
                    image={nft.image}
                    price={convertBalanceToChip(nft.price)}
                    type="nft"
                    key={nft.mintAddress}
                    clickable={false}
                  />
                );
              })}
            {variant === 'primary' && nftProfit!.length > nftsToShow && (
              <NFTCard
                price={
                  nftProfit!.length - nftsToShow > 0
                    ? nftProfit!.length - nftsToShow
                    : 0
                }
                type="more"
                onClick={onClick}
                clickable={false}
              />
            )}
          </NFTContainer>
          <Text
            fontSize={14}
            fontWeight={600}
            color="white"
            textAlign={'center'}
            mt="15px"
          >
            HOUSE FEES
          </Text>
          <NFTContainer>
            <NFTCard
              price={convertBalanceToChip(usdFee!)}
              type="chip"
              size={46}
              clickable={false}
            />
            {nftFee !== undefined &&
              nftFee !== undefined &&
              nftFee!
                .slice(0, variant === 'primary' ? nftsToShow : nftFee!.length)
                .map(nft => {
                  return (
                    <NFTCard
                      image={nft.image}
                      price={convertBalanceToChip(nft.price)}
                      type="nft"
                      key={nft.mintAddress}
                      size={46}
                      clickable={false}
                    />
                  );
                })}
            {variant === 'primary' &&
              nftFee !== undefined &&
              nftFee!.length > nftsToShow && (
                <NFTCard
                  price={
                    nftFee!.length - nftsToShow > 0
                      ? nftFee!.length - nftsToShow
                      : 0
                  }
                  size={46}
                  type="more"
                  onClick={onClick}
                  clickable={false}
                />
              )}
          </NFTContainer>
        </>
      ) : (
        <NFTContainer>
          <NFTCard price={usdBetAmount} type="chip" clickable={false} />
          {nfts.length > 0 &&
            nfts
              .slice(0, variant === 'primary' ? nftsToShow : nfts.length)
              .map(nft => {
                return (
                  <NFTCard
                    image={nft.image}
                    price={convertBalanceToChip(nft.price)}
                    type="nft"
                    key={nft.mintAddress}
                    clickable={false}
                  />
                );
              })}
          {variant === 'primary' && nfts.length > nftsToShow && (
            <NFTCard
              price={
                nfts.length - nftsToShow > 0 ? nfts.length - nftsToShow : 0
              }
              type="more"
              onClick={onClick}
              clickable={false}
            />
          )}
        </NFTContainer>
      )}
    </Container>
  );
};

export default memo(Content);
