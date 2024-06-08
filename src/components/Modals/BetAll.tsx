import { useMemo, useCallback, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { Modal, ModalProps, Flex, Chip, Button } from 'components';
import { useAppSelector } from 'state';
import { chipColors } from 'config/chip';
import { convertBalanceToChip } from 'utils/balance';

const StyledButton = styled(Button)`
  width: 152px;
  height: 46px;
  border-radius: 7px;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.16rem;
  background: transparent;
  color: white;

  display: flex;
  justify-content: center;
  align-items: center;

  &:hover {
    transform: translateY(-5px);
  }

  transition: all 0.3s;
`;

const ChipContainer = styled(Flex)<{ leftColor?: string }>`
  position: relative;
  border-radius: 5px;
  height: 38px;
  width: max-content;
  padding-left: 10px;
  padding-right: 10px;

  display: flex;
  justify-content: center;
  align-items: center;

  &::after {
    content: '';
    position: absolute;
    width: 2px;
    height: 18px;
    background-color: ${({ leftColor }) => (leftColor ? leftColor : '#ffe24b')};
    left: 0;
    top: 50%;
    transform: translate(0, -50%);
  }
`;

const Container = styled(Flex)`
  flex-direction: column;
  background: linear-gradient(180deg, #0f2035 0%, #1a293d 100%);
  border-radius: 10px;
  color: white;
  font-weight: 600;
  font-size: 22px;
  gap: 20px;
  padding: 31px 44px 40px;
`;

interface BetAllProps extends ModalProps {
  setValue: React.Dispatch<React.SetStateAction<string>>;
  maxAmount?: number;
  count?: number;
}

export default function BetAll({
  onDismiss,
  setValue,
  maxAmount,
  count = 1,
  ...props
}: BetAllProps) {
  const { balance, betBalanceType } = useAppSelector(state => state.user);

  const buttonRef = useRef<HTMLButtonElement>(null);

  const max = useMemo(() => {
    let bal = balance;
    bal = maxAmount ? Math.min(bal, maxAmount) : bal / count;
    return bal;
  }, [balance, count, maxAmount]);

  const handleClickOK = useCallback(
    (e: any) => {
      e.preventDefault();
      setValue((Math.floor(convertBalanceToChip(max) * 100) / 100).toFixed(2));
      onDismiss && onDismiss();
    },
    [max, setValue, onDismiss]
  );

  const handleClickCancel = useCallback(() => {
    onDismiss && onDismiss();
  }, [onDismiss]);

  useEffect(() => {
    setTimeout(() => {
      if (!buttonRef || !buttonRef.current) return;
      buttonRef.current.focus();
    }, 100);
  }, []);

  return (
    <Modal {...props} onDismiss={onDismiss} hideCloseButton>
      <Container>
        <form onSubmit={handleClickOK}>
          <Flex gap={17} justifyContent="center" alignItems="center">
            Bet all-in?
            <ChipContainer
              background={
                betBalanceType === 'coupon'
                  ? 'linear-gradient(90deg, #004150 0%, rgba(0, 65, 80, 0.25) 100%)'
                  : 'linear-gradient(90deg, #503b00 0%, #2f2814 100%)'
              }
              leftColor={
                chipColors[betBalanceType === 'coupon' ? 'coupon' : 'chip']
              }
            >
              <Chip
                chipType={betBalanceType === 'coupon' ? 'coupon' : 'chip'}
                fontSize="16px"
                fontWeight={600}
                price={(
                  Math.floor(convertBalanceToChip(max) * 100) / 100
                ).toFixed(2)}
              />
            </ChipContainer>
          </Flex>
          <Flex gap={22} mt="20px">
            <StyledButton
              ref={buttonRef}
              outlined
              borderColor="success"
              type="submit"
            >
              OK
            </StyledButton>
            <StyledButton
              outlined
              borderColor="warning"
              type="button"
              onClick={handleClickCancel}
            >
              CANCEL
            </StyledButton>
          </Flex>
        </form>
      </Container>
    </Modal>
  );
}
