import React from 'react';
import { ComponentMeta } from '@storybook/react';

import Modal from './Modal';
import { ModalProps } from './types';
import { Button } from 'components/Button';
import useModal from './useModal';
import DepositWithdrawModal from 'components/Modals/DepositWithdrawModal';

import LogModal from 'components/Modals/LogModal';

export default {
  title: 'Components/Modal',
  component: Modal,
  argTypes: {}
} as ComponentMeta<typeof Modal>;

const CustomModal: React.FC<ModalProps> = ({ title, onDismiss, ...props }) => (
  <Modal title={title} onDismiss={onDismiss} {...props}>
    <h1>{title}</h1>
    <Button>This button Does nothing</Button>
  </Modal>
);

export const Primary = () => {
  const [onPresent] = useModal(<CustomModal title="Modal 1" />, false);
  return (
    <>
      <button type="button" onClick={onPresent}>
        Open Modal
      </button>
    </>
  );
};

export const DepositWithdraw = () => {
  const [onPresent] = useModal(<DepositWithdrawModal />, false);

  return (
    <>
      <button type="button" onClick={onPresent}>
        Open Modal
      </button>
    </>
  );
};

export const Log = () => {
  const [onPresent] = useModal(<LogModal />, false);

  return (
    <>
      <button type="button" onClick={onPresent}>
        Open Modal
      </button>
    </>
  );
};
