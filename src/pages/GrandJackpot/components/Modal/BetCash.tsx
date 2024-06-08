import { useCallback, useState, ChangeEvent, useEffect } from 'react';
import ClipLoader from 'react-spinners/ClipLoader';

import { toast } from 'utils/toast';
import state, { useAppSelector } from 'state';
import { setRequest } from 'state/grandJackpot/actions';
import { sendMessage } from 'state/socket';
import { updateBalance } from 'state/user/actions';

import { Modal, ModalProps, Label, Button } from 'components';
import { UserStatus } from 'pages/Jackpot/components';
import coin from 'assets/imgs/coins/coin.svg';

import {
  Container,
  InputContainer,
  StyledInput,
  MaxText,
  DepositButton
} from 'pages/Jackpot/components/Modal/BetCash.styles';
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
  const { game } = useAppSelector(state => state.grandJackpot);
  const meta = useAppSelector(state => state.meta.grandJackpot);
  const [value, setValue] = useState('');

  const handleMax = useCallback(() => {
    let amount = balance;
    setValue(convertBalanceToChip(amount).toFixed(2).toString());
  }, [balance]);

  const handleBet = useCallback(() => {
    const val = (+value).toFixed(2).toString();
    if (balance < meta.minBetAmount) {
      toast.error('Insufficient funds');
    } else if (+val < convertBalanceToChip(meta.minBetAmount)) {
      toast.warning(
        `Can't bet less than ${convertBalanceToChip(meta.minBetAmount)}.`
      );
    } else if (+val > convertBalanceToChip(balance)) {
      toast.error('Insufficient funds');
    } else {
      setValue(val);
      const amount = +convertChipToBalance(+value).toFixed(0);
      state.dispatch(updateBalance({ type: -1, usdAmount: amount }));
      const content = JSON.stringify({
        amount
      });
      state.dispatch(
        sendMessage({ type: 'event', room: 'grandJackpot', content })
      );
      state.dispatch(setRequest(true));
      onDismiss && onDismiss();
    }
  }, [balance, value, onDismiss, meta]);

  const handleChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    setValue(e.target.value);
  }, []);

  useEffect(() => {
    if (status !== 'started') onDismiss && onDismiss();
  }, [status, onDismiss]);

  return (
    <Modal {...props} onDismiss={onDismiss} title="BET CHIPS">
      <Container>
        <Label fontWeight={600} fontSize="20px" mb={30}>
          BET CHIPS
        </Label>
        <Label>Bet chips or NFT to win the jackpot.</Label>
        <InputContainer mb="36px">
          <img src={coin} width={14} height={14} alt="" />
          <StyledInput
            placeholder="0.00"
            value={value}
            type="number"
            onChange={handleChange}
          />
          <Button
            backgroundColor="#121D2A"
            variant="secondary"
            onClick={handleMax}
          >
            <MaxText>Max</MaxText>
          </Button>
        </InputContainer>
        <DepositButton
          onClick={game.request ? null : handleBet}
          disabled={+value <= 0}
        >
          {game.request ? (
            <ClipLoader color="#ffffff" loading={game.request} size={30} />
          ) : (
            'BET'
          )}
        </DepositButton>
        <UserStatus {...userData} variant="secondary" background="#1A293C" />
      </Container>
    </Modal>
  );
}
