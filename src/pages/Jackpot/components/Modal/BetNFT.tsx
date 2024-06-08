import { useMemo, useCallback, useState, useEffect, useRef } from 'react';
import Select, { StylesConfig } from 'react-select';
import ClipLoader from 'react-spinners/ClipLoader';
import { toast } from 'react-toastify';

import { ReactComponent as SearchIcon } from 'assets/imgs/icons/search.svg';
import { Modal, ModalProps, Box, Flex, Grid, Text, NftBox } from 'components';
import NFTCard from 'components/NFTCard';
import { NFT } from 'api/types/nft';

import state, { useAppSelector } from 'state';
import { setRequest } from 'state/jackpot/actions';
import { sendMessage } from 'state/socket';
import { updateBalance } from 'state/user/actions';

import { InputContainer } from './BetNFT.styles';
import { DepositButton } from './BetCash.styles';
import UserStatus, { UserStatusProps } from '../UserStatus';
import { convertBalanceToChip } from 'utils/balance';

const customStyles: StylesConfig = {
  control: provided => ({
    ...provided,
    background: '#03060933',
    border: 0,
    borderRadius: 11,
    cursor: 'pointer',
    boxShadow: 'none',
    '&:hover': {
      background: '#03060933'
    }
  }),
  option: provided => ({
    ...provided,
    background: 'transparent',
    fontFamily: 'Inter',
    fontWeight: '500',
    fontSize: '16px',
    lineHeight: '19px',
    color: '#B2D1FF',
    cursor: 'pointer',
    '&:hover': {
      background: '#03060933'
    }
  }),
  input: base => ({
    ...base
  }),

  singleValue: provided => ({
    ...provided,
    color: '#B2D1FF'
  }),
  indicatorSeparator: () => ({ display: 'none' }),
  dropdownIndicator: (provided, state) => ({
    ...provided,
    color: '#B2D1FF',
    transition: '0.5s',
    transform: state.selectProps.menuIsOpen ? 'scaleY(-1)' : 'scaleY(1)',
    '&:hover': {
      color: '#B2D1FF'
    }
  }),
  menu: provided => ({
    ...provided,
    background: 'transparent'
  }),
  menuList: provided => ({
    ...provided,
    background: '#03060933',
    borderRadius: '7px',
    transition: '0.5s'
  }),
  valueContainer: base => ({
    ...base,
    fontFamily: 'Inter',
    fontWeight: '500',
    fontSize: '16px',
    lineHeight: '19px',

    color: '#B2D1FF'
  })
};

interface BetNFTProps extends ModalProps {
  userData: UserStatusProps;
  status?: string;
}

export default function BetNFT({
  userData,
  onDismiss,
  status = 'available',
  ...props
}: BetNFTProps) {
  const room = useAppSelector(state => state.jackpot.room);
  const { game } = useAppSelector(state => state.jackpot[state.jackpot.room]);
  const meta = useAppSelector(state => state.meta.jackpot[state.jackpot.room]);
  const depositedNfts = useAppSelector(state => state.user.nfts.deposited);

  const inputRef = useRef<HTMLInputElement>(null);

  const [selectedNfts, setSelectedNfts] = useState<NFT[]>([]);
  const [option, setOption] = useState({
    label: 'All collections',
    value: 'all'
  });

  const minBetAmount = useMemo(() => {
    const betAmount: number = userData?.amount?.total ?? 0;
    return Math.max(convertBalanceToChip(meta.minBetAmount) - betAmount, 0);
  }, [meta.minBetAmount, userData]);

  const maxBetAmount = useMemo(() => {
    const betAmount: number = userData?.amount?.total ?? 0;
    return convertBalanceToChip(meta.maxBetAmount) - betAmount;
  }, [userData, meta.maxBetAmount]);

  useEffect(() => {
    setTimeout(() => {
      if (!inputRef || !inputRef.current) return;
      inputRef.current.focus();
    }, 100);
  }, []);

  useEffect(() => {
    if (
      !(status === 'available' || status === 'created' || status === 'started')
    )
      onDismiss && onDismiss();
  }, [status, onDismiss]);

  const [wager, nftAmount] = useMemo(() => {
    const nftAmount = selectedNfts.reduce(
      (sum: number, nft) => sum + nft.price,
      0
    );
    // const wager = Math.ceil((nftAmount * meta.fee) / (100 - meta.fee));
    const wager = 0;
    return [convertBalanceToChip(wager), convertBalanceToChip(nftAmount)];
  }, [selectedNfts]);

  const user = useMemo(() => {
    const usd = userData.amount.total;
    const nft = nftAmount;
    const total = usd + nft + wager;
    return {
      user: userData.user,
      nfts: [...userData.nfts, ...selectedNfts],
      nftsToShow: 4,
      amount: { usd, nft, total }
    };
  }, [userData, wager, nftAmount, selectedNfts]);

  const nfts = useMemo(() => {
    return depositedNfts.filter(nft => {
      return (
        userData.nfts.findIndex(
          //@ts-ignore
          item => item.mintAddress === nft.mintAddress
        ) === -1
      );
    });
  }, [depositedNfts, userData]);

  const handleSelectNft = useCallback(
    (nft: any) => {
      const exist = selectedNfts.find(
        item => item.mintAddress === nft.mintAddress
      );

      if (exist) {
        setSelectedNfts(
          selectedNfts.filter(item => item.mintAddress !== nft.mintAddress)
        );
      } else setSelectedNfts([...selectedNfts, nft]);
    },
    [selectedNfts]
  );

  const getOptions = useCallback(() => {
    let options: { label: string; value: string }[] = [
      { label: 'All collections', value: 'all' }
    ];

    options = options.concat(
      nfts
        .map(value => value.collectionName)
        .filter((value, index, self) => self.indexOf(value) === index)
        .map(value => ({ label: value, value }))
    );

    return options;
  }, [nfts]);

  const getNfts = useCallback(() => {
    let nftList = nfts;

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
        />
      ));
  }, [nfts, option, selectedNfts, handleSelectNft]);

  const handleBet = useCallback(
    (e: any) => {
      if (game.request) return;
      e.preventDefault();
      if (nftAmount < minBetAmount) {
        toast.warning(`Can't bet less than ${minBetAmount}.`);
      } else if (nftAmount > maxBetAmount) {
        toast.warning(`Can't bet more than ${maxBetAmount}.`);
      } else {
        state.dispatch(
          updateBalance({ type: -1, usdAmount: wager, nfts: selectedNfts })
        );
        const content = JSON.stringify({
          amount: wager,
          nfts: selectedNfts.map(nft => nft.mintAddress),
          nftAmount
        });
        state.dispatch(
          sendMessage({
            type: 'event',
            room: 'jackpot',
            level: room,
            content
          })
        );
        state.dispatch(setRequest({ room, request: true }));
        onDismiss && onDismiss();
      }
    },
    [
      game.request,
      nftAmount,
      minBetAmount,
      maxBetAmount,
      wager,
      selectedNfts,
      room,
      onDismiss
    ]
  );

  return (
    <Modal {...props} onDismiss={onDismiss}>
      <form onSubmit={handleBet}>
        <Box
          background={'linear-gradient(180deg, #202F44 0%, #1B283A 100%)'}
          borderRadius={['0px', '0px', '17px']}
          maxWidth={['100vw', '100vw', '90vw']}
          width={['100vw', '100vw', '1000px']}
          maxHeight={['calc(100vh - 65px)', 'calc(100vh - 65px)', '90vh']}
          height={['100vh', '100vh', 'max-content']}
          overflow="auto"
        >
          <Box px={'30px'} pt="30px" pb="15px">
            <Flex
              alignItems={['start', 'start', 'start', 'start', 'center']}
              gap={25}
              mr="20px"
              flexDirection={['column', 'column', 'column', 'column', 'row']}
            >
              <Text
                color={'white'}
                fontWeight={500}
                fontSize={'20px'}
                letterSpacing="0.18em"
              >
                BET NFT
              </Text>

              <Flex
                flexGrow={1}
                gap={20}
                width="100%"
                flexDirection={['column', 'column', 'row']}
              >
                <InputContainer flexGrow={1}>
                  <SearchIcon width={18} height={18} />
                  <input
                    type={'text'}
                    ref={inputRef}
                    name="search"
                    placeholder="Search NFTs"
                  />
                </InputContainer>

                <Select
                  isSearchable={false}
                  defaultValue={{ label: 'All collections', value: 'all' }}
                  styles={customStyles}
                  options={getOptions()}
                  onChange={(selectedOption: any) => setOption(selectedOption)}
                  // components={{
                  //   Menu: CustomMenu,
                  // }}
                />
              </Flex>
            </Flex>

            <NftBox>{getNfts()}</NftBox>
          </Box>
          <Grid pt="10px" pb="20px" pl="37px" pr="47px">
            <UserStatus {...user} background="#0F1A26" />
            <Flex gap={20} mt={35} justifyContent="end">
              {/* <WagerButton>
              <Duelana />
              <Box width="1px" background="#4f617b" height="100%" />
              <Flex flexDirection="column" gap={2} justifyContent="center">
                DUEL WAGER*
                <Chip price={wager} fontSize="14px" fontWeight={700} />
              </Flex>
            </WagerButton> */}
              <DepositButton type="submit" disabled={selectedNfts.length === 0}>
                {game.request ? (
                  <ClipLoader color="#ffffff" loading={game.request} />
                ) : (
                  'BET'
                )}
              </DepositButton>
            </Flex>
          </Grid>
        </Box>
      </form>
    </Modal>
  );
}
