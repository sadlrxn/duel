import { FC, useCallback, useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { ModalProps, useModal } from 'components/Modal';
import { Box, Flex, NftBox } from 'components/Box';
import { Text } from 'components/Text';
import { InputBox } from 'components/InputBox';
import { Notification } from 'components/Badge';
import { ReactComponent as SearchIcon } from 'assets/imgs/icons/search.svg';

import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
import { Button } from 'components/Button';
import NFTCard from 'components/NFTCard';
import { useAppSelector } from 'state';
import useTrading from 'hooks/useTrading';
import Select from 'components/Select';
import { ReactComponent as CloseIcon } from 'assets/imgs/icons/close.svg';
import { imageProxy } from 'config';
import { convertBalanceToChip } from 'utils/balance';

const StyledNft = styled.img`
  width: 56px;
  height: 56px;
  border: 4px solid #203044;
  border-radius: 16px;

  &:not(:first-child) {
    margin-left: -25px;
  }
`;

const StyledFlex = styled(Flex)`
  margin-right: 20px;
  margin-bottom: 30px;
  gap: 25px;
  align-items: center;
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: flex;
  }
`;

const Button1 = styled(Button)`
  display: none;
  ${({ theme }) => theme.mediaQueries.md} {
    display: flex;
  }
`;
const Button2 = styled(Button)`
  display: flex;
  ${({ theme }) => theme.mediaQueries.md} {
    display: none;
  }
`;
const ListContainer = styled(Flex)`
  flex-direction: column;
  gap: 23px;
`;

export const StyledNotification = styled(Notification)`
  display: block;
  position: relative;
  transform: none;
  top: auto;
  right: auto;
  margin-left: -15px;
  height: min-content;
`;

const DepositWithdrawNftModal: FC<ModalProps> = ({ tabIndex, ..._props }) => {
  const { nfts } = useAppSelector(state => state.user);
  const [_, onDismiss] = useModal(<></>);
  const { depositNft, withdrawNft } = useTrading();
  const [selectedNfts, setSelectedNfts] = useState<any[]>([]);
  const [option, setOption] = useState({
    label: 'All collections',
    value: 'all'
  });

  const searchInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setTimeout(() => {
      if (!searchInputRef || !searchInputRef.current) return;
      searchInputRef.current.focus();
    }, 100);
  }, [option]);

  const handleSelectNft = (nft: any) => {
    const exist = selectedNfts.find(
      item => item.mintAddress === nft.mintAddress
    );

    if (exist) {
      setSelectedNfts(
        selectedNfts.filter(item => item.mintAddress !== nft.mintAddress)
      );
    } else setSelectedNfts([...selectedNfts, nft]);
  };

  const getOptions = useCallback(() => {
    let options: { label: string; value: string }[] = [
      { label: 'All collections', value: 'all' }
    ];

    // in case of deposit
    if (tabIndex === 0)
      options = options.concat(
        nfts.undeposited
          .map(value => value.collectionName)
          .filter((value, index, self) => self.indexOf(value) === index)
          .map(value => ({ label: value, value }))
      );
    else
      options = options.concat(
        nfts.deposited
          .map(value => value.collectionName)
          .filter((value, index, self) => self.indexOf(value) === index)
          .map(value => ({ label: value, value }))
      );

    return options;
  }, [nfts, tabIndex]);

  const getNfts = useCallback(() => {
    let nftList: any[];
    if (tabIndex === 0) nftList = nfts.undeposited;
    else nftList = nfts.deposited;

    if (option.value === 'all')
      return nftList.map(nft => (
        <NFTCard
          key={nft.name}
          price={nft.price}
          collectionName={nft.collectionName}
          name={nft.name}
          image={nft.image}
          selectable
          selected={
            selectedNfts.findIndex(
              item => item.mintAddress === nft.mintAddress
            ) !== -1
          }
          onClick={() => handleSelectNft(nft)}
        />
      ));

    return nftList
      .filter(value => value.collectionName === option.value)
      .map(nft => (
        <NFTCard
          key={nft.name}
          price={nft.price}
          collectionName={nft.collectionName}
          name={nft.name}
          image={nft.image}
          selectable
          selected={
            selectedNfts.findIndex(
              item => item.mintAddress === nft.mintAddress
            ) !== -1
          }
          onClick={() => handleSelectNft(nft)}
        />
      ));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [nfts, option, selectedNfts]);

  const handleClick = async (e: any) => {
    e.preventDefault();
    if (tabIndex === 0) {
      await depositNft(selectedNfts);
    } else await withdrawNft(selectedNfts);

    if (onDismiss) onDismiss();
  };

  return (
    <div className="container">
      <div className="box">
        <form onSubmit={handleClick}>
          <Flex flexDirection={'column'} flexGrow={1}>
            <StyledFlex>
              <Text
                color={'white'}
                fontWeight={500}
                fontSize={'20px'}
                letterSpacing="0.18em"
              >
                {tabIndex === 0 ? 'DEPOSIT NFTs' : 'Withdraw NFTs'}
              </Text>

              <InputBox gap={20} p="10px 20px">
                <SearchIcon width={25} height={25} />
                <input
                  ref={searchInputRef}
                  type={'text'}
                  name="search"
                  placeholder="Search NFTs"
                />
              </InputBox>

              <Select
                background="#192637"
                hoverBackground="#03060933"
                color="#B2D1FF"
                options={getOptions()}
                onChange={(selectedOption: any) => setOption(selectedOption)}
                // components={{
                //   Menu: CustomMenu,
                // }}
              />
              <CloseIcon color="#96A8C2" onClick={onDismiss} cursor="pointer" />
            </StyledFlex>

            <Flex flex={1}>
              <NftBox
                p="10px !important"
                // height={'auto'}
                marginTop="0px !important"
                maxHeight={[
                  'calc(100vh - 411px)',
                  'calc(100vh - 411px)',
                  'calc(100vh - 411px)',
                  'auto'
                ]}
                width="100%"
                overflow="auto"
              >
                {getNfts()}
              </NftBox>
            </Flex>
          </Flex>

          <ListContainer>
            <Flex
              gap={20}
              flexDirection={['column', 'column', 'column', 'row']}
              justifyContent={'space-between'}
              alignItems={['start', 'start', 'start', 'center']}
              borderTop="1px solid #4F617B"
              background={[
                'transparent',
                'transparent',
                'transparent',
                '#203044'
              ]}
              margin={'0px -30px -30px -30px'}
              px="30px"
              py={['20px', '20px', '20px', '5px']}
            >
              <Flex alignItems={'center'}>
                <Box
                  background={'#4FFF8B'}
                  borderRadius="13px"
                  px="8px"
                  py="0px"
                  mr={'10px'}
                >
                  <Text fontSize={'12px'} fontWeight={600} color="#0B141E">
                    {selectedNfts.length}
                  </Text>
                </Box>

                <Text
                  color={'#B2D1FF80'}
                  fontWeight={500}
                  fontSize={'12px'}
                  width="54px"
                  mr={'20px'}
                >
                  Selected NFTs
                </Text>

                <Flex>
                  {selectedNfts.slice(0, 3).map(value => (
                    <StyledNft
                      src={imageProxy(300) + value.image}
                      alt="nft"
                      key={value.mintAddress}
                    />
                  ))}
                  {selectedNfts.length > 3 && (
                    <StyledNotification>
                      {selectedNfts.length - 3}
                    </StyledNotification>
                  )}
                </Flex>
              </Flex>

              <Flex
                alignItems="center"
                justifyContent={[
                  'space-between',
                  'space-between',
                  'space-between',
                  'start'
                ]}
                width={['100%', '100%', '100%', 'auto']}
              >
                <Text
                  color={'#B2D1FF80'}
                  fontWeight={500}
                  fontSize={'12px'}
                  width="58px"
                  mr={'20px'}
                  ml="22px"
                >
                  Estimated Value
                </Text>

                <Flex
                  background={'#4FFF8B1A'}
                  p="6px"
                  alignItems={'center'}
                  borderRadius="9px 7px 7px 8px"
                  mr="30px"
                >
                  <CoinIcon />
                  <Text
                    fontSize={'14px'}
                    fontWeight={700}
                    color={'#4FFF8B'}
                    ml="4px"
                  >
                    {convertBalanceToChip(
                      selectedNfts.reduce(
                        (item1, item2) => ({
                          price: item1.price + item2.price
                        }),
                        { price: 0 }
                      ).price
                    )}
                  </Text>
                </Flex>
                <Button1
                  fontSize={'16px'}
                  fontWeight={600}
                  px="20px"
                  py={'12px'}
                  borderRadius="5px"
                  my="5px"
                  type="submit"
                  disabled={selectedNfts.length === 0}
                >
                  {tabIndex === 0 ? 'Deposit NFTs' : 'Withdraw NFTs'}
                </Button1>
              </Flex>
            </Flex>
            <Button2
              fontSize={'16px'}
              fontWeight={600}
              px="20px"
              py={'12px'}
              borderRadius="5px"
              my="5px"
              type="submit"
              disabled={selectedNfts.length === 0}
            >
              {tabIndex === 0 ? 'Deposit NFTs' : 'Withdraw NFTs'}
            </Button2>
          </ListContainer>
        </form>
      </div>
    </div>
  );
};

export default DepositWithdrawNftModal;
