import { useMemo, memo, FC } from 'react';
import styled from 'styled-components';

import { Box, Text, Flex, Button } from 'components';
import { NFTCard } from 'pages/Jackpot/components';
import { formatUserName } from 'utils/format';

import { NFT } from 'api/types/nft';
import { formatNumber } from 'utils/format';
import { convertBalanceToChip } from 'utils/balance';

const CashWrapper = styled(Flex)`
  flex-direction: column;
  align-items: center;
  .width_800 & {
    align-items: start;
  }
`;

const TopNFTContainer = styled(Flex)`
  display: none;
  .width_800 & {
    display: flex;
  }
`;

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
  border: 2px solid #c0b264;
  background: linear-gradient(
    360deg,
    rgba(255, 226, 75, 0.2) 0%,
    rgba(255, 226, 75, 0) 100%
  );
`;

const ResultContainer = styled(Flex)`
  flex-direction: column;
  gap: 15px;
  align-items: center;
  color: #ffffff;
  font-size: 14px;
  letter-spacing: 0.18em;

  .width_800 & {
    align-items: end;
  }
`;

const TitleContainer = styled(Box)`
  display: grid;
  gap: 26px;
  font-weight: 600;
  font-size: 14px;
  color: #fff6ca;

  .width_800 & {
    grid-template-columns: max-content auto max-content;
  }
`;

const NFTPrizeContainer = styled(Box)`
  justify-content: center;
  gap: 24px;
  max-width: 100%;

  overflow-x: scroll;
`;

const Container = styled(Box)`
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
  handleShowNFT?: any;
}

const GrandContent: FC<ContentProps> = ({
  nfts = [],
  usdBetAmount = 0,
  nftBetAmount = 0,
  onClose,
  win = false,
  winnerName = '',
  handleShowNFT
}) => {
  const totalBetAmount = useMemo(
    () => usdBetAmount + nftBetAmount,
    [usdBetAmount, nftBetAmount]
  );

  return (
    <Container onClick={onClose}>
      <TitleContainer mb={'34px'}>
        <CashWrapper gap={19}>
          <Text
            fontWeight={600}
            fontSize={14}
            color="white"
            letterSpacing="0.18em"
          >
            CASH PRIZE
          </Text>
          <NFTCard
            size={125}
            price={usdBetAmount}
            type="chip"
            clickable={false}
          />
        </CashWrapper>
        <div>
          {win ? (
            <WinnerText>Won By {formatUserName(winnerName)}</WinnerText>
          ) : (
            <TopNFTContainer
              flexDirection="column"
              gap={18}
              alignItems="center"
            >
              <Text
                fontWeight={600}
                fontSize={14}
                color="white"
                letterSpacing="0.18em"
              >
                TOP NFTS
              </Text>
              <Flex gap={16}>
                <>
                  {nfts
                    .slice()
                    .sort((nft1, nft2) => nft2.price - nft1.price)
                    .slice(0, Math.min(3, nfts.length))
                    .map(nft => {
                      return (
                        <NFTCard
                          size={125}
                          image={nft.image}
                          price={convertBalanceToChip(nft.price)}
                          type="nft"
                          key={'top_nfts_' + nft.mintAddress}
                          clickable={false}
                        />
                      );
                    })}
                </>
              </Flex>
            </TopNFTContainer>
          )}
        </div>
        <ResultContainer>
          JACKPOT VALUE
          <TotalChip>
            <ChipIcon />
            {formatNumber(totalBetAmount)}
          </TotalChip>
        </ResultContainer>
      </TitleContainer>
      {nfts.length >= 0 && (
        <Flex flexDirection="column" gap={20} alignItems="center">
          <Text
            fontWeight={600}
            fontSize={14}
            color="white"
            letterSpacing="0.18em"
          >
            NFT PRIZES
          </Text>
          <NFTPrizeContainer>
            <Flex
              gap={24}
              justifyContent="center"
              width="min-content"
              pb="15px"
            >
              {nfts.map(nft => {
                return (
                  <NFTCard
                    size={110}
                    image={nft.image}
                    price={convertBalanceToChip(nft.price)}
                    type="nft"
                    key={nft.mintAddress}
                    clickable={false}
                  />
                );
              })}
            </Flex>
          </NFTPrizeContainer>
          <Button
            variant="secondary"
            onClick={() => {
              handleShowNFT && handleShowNFT({ nfts, name: '', level: 0 });
            }}
          >
            <Text
              color="text"
              fontWeight={600}
              fontSize={14}
              px="25px"
              py="10px"
            >
              Show All NFTs
            </Text>
          </Button>
        </Flex>
      )}
    </Container>
  );
};

export default memo(GrandContent);
