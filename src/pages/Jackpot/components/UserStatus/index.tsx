import { FC, memo, useRef, useEffect, useMemo, useState } from 'react';
import { useContainerQuery } from 'react-container-query';
import classnames from 'classnames';

import { NFT } from 'api/types/nft';
import { formatUserName } from 'utils/format';
import { imageProxy } from 'config';

import { Chip, Flex, Text, Avatar, BoxProps, Span, Badge } from 'components';
import {
  Divider,
  DividerHorizontal,
  UsdAmount,
  UserContainer,
  User,
  UserInfo,
  NftContainer,
  StyledIntro,
  StackedImage,
  StyledNotification,
  Nft,
  AmountContainer,
  Container,
  FlexBox,
  NftAmount
} from './styles';

const query = {
  user_status_max_640: {
    maxWidth: 639
  }
};

export interface UserStatusProps extends BoxProps {
  user: {
    id: number;
    name: string;
    avatar: string;
    level: number;
    percent: number;
  };
  nfts: NFT[];
  amount: {
    usd: number;
    nft: number;
    total: number;
  };
  nftsToShow?: number;
  handleShowNFT?: any;
}

const UserStatus: FC<UserStatusProps> = ({
  user,
  nfts,
  amount,
  // nftsToShow = 3,
  handleShowNFT,
  ...props
}) => {
  const { id: userId, name, avatar, level, percent } = user;
  const { usd: usdAmount, nft: nftAmount, total: totalAmount } = amount;

  const nftRef = useRef<any>(null);
  const [nftsToShow, setNftsToShow] = useState(3);

  const [params, containerRef] = useContainerQuery(query, { width: 640 });

  const resizeObserver = useMemo(
    () =>
      new ResizeObserver(() => {
        if (!nftRef || !nftRef.current) setNftsToShow(3);
        else {
          const width = nftRef.current.offsetWidth as number;
          setNftsToShow(1 + Math.max(Math.floor((width - 45) / 29), 0));
        }
      }),
    []
  );

  useEffect(() => {
    if (!nftRef) return;
    resizeObserver.observe(nftRef.current);
  }, [resizeObserver]);

  return (
    <Container ref={containerRef} className={classnames(params)} {...props}>
      <UserContainer>
        <User>
          <Avatar
            userId={userId}
            name={name}
            image={avatar}
            border="none"
            borderRadius="12px"
            padding="0px"
          />
          <UserInfo>
            <Span
              style={{
                width: '105px',
                overflow: 'hidden',
                whiteSpace: 'nowrap',
                textOverflow: 'ellipsis'
              }}
            >
              {formatUserName(name)}
            </Span>
            <Flex flexDirection="row" gap={1} alignItems="center">
              <Span color={'#768BAD'} fontSize={12}>
                Win:&nbsp;
              </Span>
              <Text color={percent > 0 ? 'success' : 'text'} fontSize={12}>
                {Number(percent.toFixed(2)) === 0
                  ? percent.toFixed(4)
                  : percent.toFixed(2)}
                %
              </Text>
            </Flex>
          </UserInfo>
        </User>
        <UsdAmount>
          <Chip price={usdAmount} color="chip" />
        </UsdAmount>
      </UserContainer>
      <Divider />
      {nfts.length > 0 && <DividerHorizontal />}
      <NftContainer className={nfts.length > 0 ? '' : 'hide'}>
        <FlexBox flexDirection={'column'} mr="30px">
          <Text color={'#768BAD'} fontWeight={500} fontSize={14}>
            Cash Bet
          </Text>
          <UsdAmount>
            <Chip price={usdAmount} color="chip" />
          </UsdAmount>
        </FlexBox>
        <Divider />
        <FlexBox flexDirection={'column'} ml="30px">
          <Text color={'#768BAD'} fontWeight={500} fontSize={14}>
            NFT Bet
          </Text>
          <Badge>
            <Chip price={nftAmount} color="success" />
          </Badge>
        </FlexBox>

        <Nft>
          <StyledIntro>
            NFT
            <br /> Bet
          </StyledIntro>
          <Flex
            flexDirection="row"
            gap={4}
            alignItems="center"
            marginLeft={20}
            marginRight={20}
            position="relative"
            height={45}
            style={{ cursor: nfts.length > 0 ? 'pointer' : 'auto' }}
            onClick={() =>
              handleShowNFT &&
              nfts.length > 0 &&
              handleShowNFT({ nfts, name, level })
            }
            ref={nftRef}
          >
            {nfts.length > 0 &&
              nfts
                .slice(0, Math.min(nftsToShow, nfts.length))
                .map(nft => (
                  <StackedImage
                    key={nft.image}
                    src={imageProxy(300) + nft.image}
                    alt={''}
                  />
                ))}
            {nfts.length > nftsToShow && (
              <StyledNotification>
                {nfts.length - nftsToShow}
              </StyledNotification>
            )}
          </Flex>
          <NftAmount>
            <Chip price={nftAmount} color="success" />
          </NftAmount>
        </Nft>
      </NftContainer>
      <Divider />
      <AmountContainer>
        <Flex flexDirection="column" gap={4} color="text" alignItems="center">
          <Text fontSize={12} color="text" fontWeight={500}>
            Total Bet
          </Text>
          <Chip price={totalAmount} fontSize={12} />
        </Flex>
      </AmountContainer>
    </Container>
  );
};

export default memo(UserStatus);
