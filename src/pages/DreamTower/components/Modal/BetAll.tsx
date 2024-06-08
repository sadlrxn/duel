import { useMemo, useCallback } from 'react';
import styled from 'styled-components';
import { Modal, ModalProps, Flex, Chip, Button } from 'components';
import { useAppSelector } from 'state';
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

const ChipContainer = styled(Flex)`
  position: relative;
  background: linear-gradient(90deg, #503b00 0%, #2f2814 100%);
  border-radius: 5px;
  height: 38px;
  width: 102px;

  display: flex;
  justify-content: center;
  align-items: center;

  &::after {
    content: '';
    position: absolute;
    width: 2px;
    height: 18px;
    background-color: #ffe24b;
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
}

export default function BetAll({ onDismiss, setValue, ...props }: BetAllProps) {
  const { balance } = useAppSelector(state => state.user);

  const max = useMemo(() => {
    return balance;
  }, [balance]);

  const handleClickOK = useCallback(() => {
    setValue(convertBalanceToChip(max).toFixed(2).toString());
    onDismiss && onDismiss();
  }, [max, setValue, onDismiss]);

  const handleClickCancel = useCallback(() => {
    onDismiss && onDismiss();
  }, [onDismiss]);

  return (
    <Modal {...props} onDismiss={onDismiss} hideCloseButton>
      <Container>
        <Flex gap={17} justifyContent="center" alignItems="center">
          Bet all-in?
          <ChipContainer>
            <Chip
              fontSize="16px"
              fontWeight={600}
              price={convertBalanceToChip(max).toFixed(2)}
            />
          </ChipContainer>
        </Flex>
        <Flex gap={22}>
          <StyledButton outlined borderColor="success" onClick={handleClickOK}>
            OK
          </StyledButton>
          <StyledButton
            outlined
            borderColor="warning"
            onClick={handleClickCancel}
          >
            CANCEL
          </StyledButton>
        </Flex>
      </Container>
    </Modal>
  );
}
