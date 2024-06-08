import React, { FC } from 'react';
import { Tab, TabList } from 'react-tabs';
import { StyledTabs, StyledTabPanel } from './styles';
import { Modal, ModalProps, useModal } from 'components/Modal';
import { Flex } from 'components/Box';
import { Text } from 'components/Text';
import DepositSolTab from './DepositSolTab';
import WithdrawSolTab from './WithdrawSolTab';
import HistoryTab from './HistoryTab';
import 'react-tabs/style/react-tabs.css';
import DepositWithdrawNftTab from './NFTTab';
import CouponTab from './CouponTab';
import styled from 'styled-components';
import { ReactComponent as CloseIcon } from 'assets/imgs/icons/close.svg';

const StyledHeader = styled(Flex)`
  background: linear-gradient(180deg, #132031 0%, #1a293d 100%);
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  display: flex;
  ${({ theme }) => theme.mediaQueries.md} {
    display: none;
  }
`;

const DepositWithdrawModal: FC<ModalProps> = ({ tabIndex, ...props }) => {
  const [, onDismiss] = useModal(<></>, false);
  return (
    <Modal {...props}>
      <StyledHeader>
        <Text
          textTransform="uppercase"
          color={'#768BAD'}
          fontSize="22px"
          fontWeight={600}
        >
          Wallet
        </Text>

        <CloseIcon color="#96A8C2" onClick={onDismiss} cursor="pointer" />
      </StyledHeader>
      <StyledTabs defaultIndex={tabIndex}>
        <TabList>
          <Flex>
            <Tab>
              DEPOSIT
              <b />
            </Tab>
            <Tab>
              WITHDRAW
              <b />
            </Tab>
            <Tab>
              NFT DEPOSIT
              <b />
            </Tab>
            <Tab>
              NFT WITHDRAW
              <b />
            </Tab>
            <Tab>
              COUPONS
              <b />
            </Tab>
            <Tab>
              HISTORY
              <b />
            </Tab>
          </Flex>
        </TabList>

        <StyledTabPanel>
          <DepositSolTab />
        </StyledTabPanel>
        <StyledTabPanel>
          <WithdrawSolTab />
        </StyledTabPanel>
        <StyledTabPanel>
          <DepositWithdrawNftTab tabIndex={0} />
        </StyledTabPanel>
        <StyledTabPanel>
          <DepositWithdrawNftTab tabIndex={1} />
        </StyledTabPanel>
        <StyledTabPanel>
          <CouponTab />
        </StyledTabPanel>
        <StyledTabPanel>
          <HistoryTab />
        </StyledTabPanel>
      </StyledTabs>
    </Modal>
  );
};

export default DepositWithdrawModal;
