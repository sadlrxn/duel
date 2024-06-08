import {
  useState,
  useCallback,
  ChangeEvent,
  useEffect,
  useRef,
  useMemo
} from 'react';
import styled from 'styled-components';
import { toast } from 'react-toastify';

import { useAppSelector, useAppDispatch } from 'state';
import { updateBalance } from 'state/user/actions';

import coin from 'assets/imgs/coins/coin.svg';

import Avatar from 'components/Avatar';
import { Box, Flex } from 'components/Box';
import { Button } from 'components/Button';
import Checkbox from 'components/Checkbox';
import { Input } from 'components/Input';
import { Modal, ModalProps } from 'components/Modal';
import { Text } from 'components/Text';
import { useUserInfo } from 'hooks';
import { sendTip } from 'services';
import { formatUserName } from 'utils/format';
import { convertBalanceToChip, convertChipToBalance } from 'utils/balance';
export interface ChipModalProps extends ModalProps {
  userId?: number;
  name?: string;
  avatar?: string;
  amount?: number;
}

export default function TipModal({
  userId,
  name,
  avatar,
  amount = 0,
  onDismiss,
  ...props
}: ChipModalProps) {
  const dispatch = useAppDispatch();
  const {
    balances,
    // statistics: { private_profile: isPrivate },
    id
  } = useAppSelector(state => state.user);
  const { data, error } = useUserInfo({ userId, userName: name });
  const [value, setValue] = useState(amount || '');
  const [hide, setHide] = useState(false);
  const [errText, setErrText] = useState('');

  const balance = useMemo(() => {
    return balances.chip.balance;
  }, [balances.chip.balance]);

  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setTimeout(() => {
      if (!inputRef || !inputRef.current) return;
      inputRef.current.focus();
    }, 100);
  }, []);

  const handleSend = useCallback(
    async (e: any) => {
      e.preventDefault();
      const val = (+value).toFixed(2).toString();
      if (id === userId) {
        setErrText('You cannot send yourself a tip');
      } else if (+val < 0.01) {
        setErrText('Must be at least 0.01');
      } else if (+val > convertBalanceToChip(balance)) {
        setErrText('Insufficient funds');
      } else if (+val > 10000) {
        setErrText('Must be at most 10,000');
      } else {
        setErrText('');
        setValue(val);
        const amount = +convertChipToBalance(+value).toFixed(0);
        dispatch(updateBalance({ type: -1, usdAmount: amount }));

        try {
          const { recipient, amount: usdAmount } = await sendTip({
            userId: userId || data?.info.id,
            amount,
            showInChat: !hide
          });

          toast.success(
            `Send ${convertBalanceToChip(usdAmount)} chip${
              usdAmount !== 100000 ? 's' : ''
            } to ${recipient.name} success.`
          );
          setValue('');
          onDismiss && onDismiss();
        } catch (err) {
          dispatch(updateBalance({ type: 1, usdAmount: amount }));
          //@ts-ignore
          toast.error(err.response.data.status);
        }
      }
    },
    [value, balance, dispatch, userId, hide, onDismiss, data?.info.id, id]
  );

  const handleChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    if (+e.target.value < 0) return;
    setValue(e.target.value);
  }, []);

  useEffect(() => {
    if (errText === '') return;
    const timer = setTimeout(() => {
      setErrText('');
    }, 5000);

    return () => {
      clearTimeout(timer);
    };
  }, [errText]);

  return (
    <Modal {...props} onDismiss={onDismiss}>
      <form onSubmit={handleSend}>
        <Container>
          {error || !data ? (
            <Text
              fontSize="24px"
              fontWeight={600}
              lineHeight="29px"
              color="white"
              textAlign="center"
            >
              User not found.
            </Text>
          ) : data.info.name.toLowerCase() === 'hidden' ? (
            <Text
              fontSize="24px"
              fontWeight={600}
              lineHeight="29px"
              color="white"
              textAlign="center"
            >
              This user is hidden and tip is not available.
            </Text>
          ) : (
            <>
              <Flex gap={10} alignItems="center" mb="20px">
                <Avatar
                  image={avatar || data?.info.avatar}
                  size="50px"
                  padding="0px"
                />
                <Text color="chip" fontSize="25px" fontWeight={600}>
                  {formatUserName(data?.info.name || '')}
                </Text>
              </Flex>

              <Text color="#4F617B" fontSize="16px" fontWeight={400} mb="11px">
                Tip amount
              </Text>

              <InputContainer>
                <StyledInput
                  ref={inputRef}
                  placeholder="0.00"
                  value={value}
                  type="number"
                  onChange={handleChange}
                />
                <img src={coin} width={14} height={14} alt="" />
              </InputContainer>

              {errText !== '' && <ErrorText>{errText}</ErrorText>}

              <Box mt="14px">
                <Checkbox
                  name="hide"
                  value="hide"
                  label="DON'T SHOW TIP IN CHAT"
                  // disabled={isPrivate}
                  checked={hide}
                  // defaultChecked={false}
                  onChange={(e: ChangeEvent<HTMLInputElement>) => {
                    setHide(e.target.checked);
                  }}
                />
              </Box>
              <Box mt="50px">
                <Button
                  type="submit"
                  px="20px"
                  py="14px"
                  fontSize="16px"
                  fontWeight={600}
                  disabled={data === undefined}
                >
                  SEND TIP
                </Button>
              </Box>
            </>
          )}
        </Container>
      </form>
    </Modal>
  );
}

const Container = styled(Box)`
  /* background: linear-gradient(180deg, #030508 0%, #0b141e 100%); */
  background-color: #1a293d;
  border: 2px solid #43546c;
  border-radius: 15px;
  padding: 40px;

  min-width: 340px;
  max-height: 90vh;

  overflow: hidden auto;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
`;

export const InputContainer = styled(Flex)`
  gap: 41px;
  padding: 13px 15px;
  justify-content: space-between;
  align-items: center;
  border-radius: 9px;
  background: ${({ theme }) => theme.jackpot.input};

  margin-bottom: 7px;
`;

export const StyledInput = styled(Input)`
  font-size: 20px;
  font-weight: 400;
  line-height: 24px;
  color: ${({ theme }) => theme.coinflip.private};
`;

export const ErrorText = styled(Text)`
  font-size: 16px;
  line-height: 20px;
  font-weight: 600;
  color: #a31f4e;
`;
