import {
  useRef,
  useState,
  useMemo,
  useCallback,
  ChangeEvent,
  useEffect
} from 'react';

import coin from 'assets/imgs/coins/coin.svg';

import { Flex, Text, Modal, ModalProps, Label, Button } from 'components';
import state, { useAppSelector } from 'state';
import { toast } from 'utils/toast';
import { setAutoBet } from 'state/jackpot/actions';
import { setRequest } from 'state/jackpot/actions';
import { updateBalance } from 'state/user/actions';
import { sendMessage } from 'state/socket';
import { imageProxy } from 'config';

import {
  Container,
  // NftAutoBetButton,
  StackedImage,
  StyledNotification
} from './AutoBet.styles';
import {
  InputContainer,
  StyledInput,
  MaxText,
  DepositButton
} from './BetCash.styles';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';

interface AutoBetProps extends ModalProps {
  nftsToShow?: number;
}

export default function AutoBet({
  nftsToShow = 3,
  onDismiss,
  ...props
}: AutoBetProps) {
  const room = useAppSelector(state => state.jackpot.room);
  const { game: currentGame, autoBet } = useAppSelector(
    state => state.jackpot[state.jackpot.room]
  );
  const meta = useAppSelector(state => state.meta.jackpot[state.jackpot.room]);
  const { balance, id: userId } = useAppSelector(state => state.user);
  const [value, setValue] = useState('');
  const inputRef = useRef(null);

  const [chip, nfts] = useMemo(() => {
    return autoBet ? [autoBet.chip ?? 0, autoBet.nfts ?? []] : [0, []];
  }, [autoBet]);

  const minBetAmount = useMemo(() => {
    return convertBalanceToChip(meta.minBetAmount);
  }, [meta.minBetAmount]);

  const maxBetAmount = useMemo(() => {
    return convertBalanceToChip(meta.maxBetAmount);
  }, [meta.maxBetAmount]);

  const handleMin = useCallback(() => {
    setValue(minBetAmount.toFixed(2).toString());
  }, [minBetAmount]);

  const handleBet = useCallback(
    (e: any) => {
      e.preventDefault();
      const val = (+value).toFixed(2).toString();
      if (convertBalanceToChip(balance) < minBetAmount) {
        toast.error('Insufficient funds');
      } else if (+val < minBetAmount) {
        toast.warning(`Can't bet less than ${minBetAmount}.`);
      } else if (+val > maxBetAmount) {
        toast.warning(`Can't bet more than ${maxBetAmount}.`);
      } else if (+val > convertBalanceToChip(balance)) {
        toast.error('Insufficient funds');
      } else {
        setValue(val);
        const amount = +convertChipToBalance(+value).toFixed(0);
        state.dispatch(setAutoBet({ room, autoBet: { chip: amount } }));

        if (
          currentGame.players.findIndex(player => player.id === userId) === -1
        ) {
          state.dispatch(updateBalance({ type: -1, usdAmount: amount }));
          const content = JSON.stringify({
            amount
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
        }
        onDismiss && onDismiss();
      }
    },
    [
      value,
      balance,
      minBetAmount,
      maxBetAmount,
      room,
      currentGame.players,
      onDismiss,
      userId
    ]
  );

  const handleChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      if (+e.target.value < 0) return;
      if (+e.target.value > maxBetAmount) return;
      setValue(e.target.value);
    },
    [maxBetAmount]
  );

  useEffect(() => {
    if (!inputRef || !inputRef.current) return;
    //@ts-ignore
    inputRef.current.focus();
  }, []);

  return (
    <Modal {...props} onDismiss={onDismiss} title="AUTO BET">
      <form onSubmit={handleBet}>
        <Container>
          <Flex flexDirection="column">
            <Label fontWeight={600} fontSize="20px" mb={21}>
              AUTO BET
            </Label>
            <Label
              fontSize="15px"
              fontWeight={400}
              lineHeight="18px"
              color="#B2D1FF"
              mb="27px"
            >
              While auto bet is enabled, your bet will be placed every round
              until you disable it or close your browser. Your bet will continue
              when you play other games.
            </Label>
            <Label fontWeight={400} color="text" mb="-5px">
              Auto Bet CHIPs
            </Label>
          </Flex>
          <InputContainer mb="42px">
            <img src={coin} width={14} height={14} alt="" />
            <StyledInput
              ref={inputRef}
              defaultValue={chip.toFixed(2)}
              placeholder="0.00"
              value={value}
              type="number"
              onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) => {
                if (
                  e.key === 'e' ||
                  e.key === 'E' ||
                  e.key === '+' ||
                  e.key === '-'
                )
                  e.preventDefault();
              }}
              step="any"
              onChange={handleChange}
            />
            <Button
              type="button"
              backgroundColor="#121D2A"
              variant="secondary"
              onClick={handleMin}
            >
              <MaxText>Min</MaxText>
            </Button>
          </InputContainer>
          {/* <NftAutoBetButton type="button">Auto Bet NFTs</NftAutoBetButton> */}
          {/* <Flex my="20px" alignItems="center" gap={14}> */}
          {nfts.length > 0 && (
            <>
              <Text
                backgroundColor="success"
                px="10px"
                py="3px"
                color="black"
                fontSize="12px"
                fontWeight={700}
                lineHeight="14.52px"
                borderRadius="15px"
                minWidth="max-content"
              >
                {nfts.length}
              </Text>
              <Text fontSize="12px" fontWeight={500} color="#768bad">
                Selected
                <br /> NFTs
              </Text>
              <Flex
                flexDirection="row"
                gap={4}
                alignItems="center"
                marginLeft={20}
                marginRight={20}
                position="relative"
                height={45}
                style={{ cursor: nfts.length > 0 ? 'pointer' : 'auto' }}
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
            </>
          )}
          {/* </Flex> */}
          <DepositButton type="submit">Start Auto Bet</DepositButton>
        </Container>
      </form>
    </Modal>
  );
}
