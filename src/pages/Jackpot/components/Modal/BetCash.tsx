import {
  useCallback,
  useState,
  ChangeEvent,
  useEffect,
  useRef,
  useMemo
} from 'react';
import ClipLoader from 'react-spinners/ClipLoader';

import { toast } from 'utils/toast';
import state, { useAppSelector } from 'state';
import { setRequest } from 'state/jackpot/actions';
import { sendMessage } from 'state/socket';
import { updateBalance } from 'state/user/actions';

import { Modal, ModalProps, Label, Button } from 'components';
import UserStatus from '../UserStatus';
import coin from 'assets/imgs/coins/coin.svg';

import {
  Container,
  InputContainer,
  StyledInput,
  MaxText,
  DepositButton
} from './BetCash.styles';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';

interface BetCashProps extends ModalProps {
  userData: any;
  balance?: number;
  status?: string;
}

export default function BetCash({
  userData,
  balance = 0,
  onDismiss,
  status = 'available',
  ...props
}: BetCashProps) {
  const room = useAppSelector(state => state.jackpot.room);
  const { game } = useAppSelector(state => state.jackpot[state.jackpot.room]);
  const meta = useAppSelector(state => state.meta.jackpot[state.jackpot.room]);
  const [value, setValue] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  const minBetAmount = useMemo(() => {
    const betAmount: number = userData?.amount?.total ?? 0;
    return Math.max(convertBalanceToChip(meta.minBetAmount) - betAmount, 0);
  }, [meta.minBetAmount, userData]);

  const maxBetAmount = useMemo(() => {
    const betAmount: number = userData?.amount?.total ?? 0;
    return Math.min(
      convertBalanceToChip(meta.maxBetAmount) - betAmount,
      convertBalanceToChip(balance)
    );
  }, [userData, meta.maxBetAmount, balance]);

  const handleMin = useCallback(() => {
    setValue(minBetAmount.toFixed(2).toString());
  }, [minBetAmount]);

  const handleMax = useCallback(() => {
    setValue(maxBetAmount.toFixed(2).toString());
  }, [maxBetAmount]);

  const handleBet = useCallback(
    (e: any) => {
      if (game.request) return;
      e.preventDefault();
      const val = (Math.floor(+value * 100) / 100).toFixed(2);
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
        const amount = convertChipToBalance(Math.floor(+value * 100) / 100);
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
        onDismiss && onDismiss();
      }
    },
    [balance, value, onDismiss, room, minBetAmount, maxBetAmount, game.request]
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
    if (
      !(status === 'available' || status === 'created' || status === 'started')
    )
      onDismiss && onDismiss();
  }, [status, onDismiss]);

  useEffect(() => {
    setTimeout(() => {
      if (!inputRef || !inputRef.current) return;
      inputRef.current.focus();
    }, 100);
  }, []);

  return (
    <Modal {...props} onDismiss={onDismiss} title="CUSTOM BET">
      <form onSubmit={handleBet}>
        <Container>
          <Label fontWeight={600} fontSize="20px" mb={30}>
            CUSTOM BET
          </Label>
          <Label color="text" mb="-5px">
            Custom Bet
          </Label>
          <InputContainer mb="36px">
            <img src={coin} width={14} height={14} alt="" />
            <StyledInput
              ref={inputRef}
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
            <Button
              type="button"
              backgroundColor="#121D2A"
              variant="secondary"
              onClick={handleMax}
            >
              <MaxText>Max</MaxText>
            </Button>
          </InputContainer>
          <DepositButton type="submit" disabled={+value <= 0}>
            {game.request ? (
              <ClipLoader color="#ffffff" loading={game.request} size={30} />
            ) : (
              'BET'
            )}
          </DepositButton>
          <UserStatus
            {...userData}
            variant="secondary"
            background="#1A293C"
            mb="30px"
          />
        </Container>
      </form>
    </Modal>
  );
}
