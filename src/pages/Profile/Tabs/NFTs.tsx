import React, { useCallback, useState } from 'react';

import { Flex, Box, NftBox } from 'components/Box';
import { Text } from 'components/Text';
import { InputBox } from 'components/InputBox';
import { Button } from 'components/Button';

import { ReactComponent as SearchIcon } from 'assets/imgs/icons/search.svg';
import { useAppSelector } from 'state';
import NFTCard from 'components/NFTCard';
import useTrading from 'hooks/useTrading';

import { DepositWithdrawModal, useModal } from 'components';
import toast from 'utils/toast';

export default function NFTs() {
  const { deposited } = useAppSelector(state => state.user.nfts);
  const [selectedNfts, setSelectedNfts] = useState<any[]>([]);
  const { withdrawNft } = useTrading();
  const [onPresentDeposit] = useModal(
    <DepositWithdrawModal tabIndex={2} hideCloseButton={true} />,
    true
  );

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

  const getNFTs = useCallback(() => {
    return deposited.map(nft => (
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
  }, [deposited, selectedNfts]);

  const handleWithdraw = async () => {
    if (selectedNfts.length === 0) {
      toast.warn('select NFTs to withdraw!');
      return;
    }
    withdrawNft(selectedNfts);
  };

  return (
    <div className="container">
      <div className="box">
        <Flex
          alignItems={'center'}
          gap={15}
          flexDirection={['column', 'column', 'column', 'row']}
        >
          <Text
            display={['none', 'none', 'none', 'block']}
            color={'white'}
            fontSize="25px"
            fontWeight={500}
            mr="20px"
          >
            My NFTs
          </Text>

          <InputBox
            gap={20}
            p="10px 20px"
            padding={'10px 20px !important'}
            background="#142131 !important"
            width={['100%', '100%', '100%', 'auto']}
          >
            <SearchIcon />
            <input
              placeholder="Search NFTs, collections"
              name="search"
              color="#495D7E"
            />
          </InputBox>

          <Flex
            gap={15}
            justifyContent="space-between"
            width={['100%', '100%', '100%', 'auto']}
          >
            <Button
              background={'#242F42'}
              borderRadius={'5px'}
              color="#768BAD"
              py="12px"
              px="20px"
              fontSize={'14px'}
              fontWeight={600}
              width={['100%', '100%', '100%', 'auto']}
              onClick={onPresentDeposit}
            >
              Deposit NFTs
            </Button>
            <Button
              background={'#242F42'}
              borderRadius={'5px'}
              color="#768BAD"
              py="12px"
              px="20px"
              fontSize={'14px'}
              fontWeight={600}
              width={['100%', '100%', '100%', 'auto']}
              onClick={handleWithdraw}
            >
              Withdraw All
            </Button>
          </Flex>
        </Flex>
        <NftBox maxHeight={'auto'}>{getNFTs()}</NftBox>
      </div>
    </div>
  );
}
