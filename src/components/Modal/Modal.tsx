import { FC, memo, useEffect, useRef } from 'react';
import { ReactComponent as CloseIcon } from 'assets/imgs/icons/close.svg';

import { ModalContainer } from './styles';
import { ModalProps } from './types';
import styled from 'styled-components';

const StyledModalContainer = styled(ModalContainer)`
  display: flex;
  flex-direction: column;

  height: calc(100vh - 65px);
  ${({ theme }) => theme.mediaQueries.md} {
    height: auto;
  }
`;

const Modal: FC<ModalProps> = ({
  onDismiss,
  onBack: _,
  children,
  hideCloseButton = false,
  ...props
}) => {
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!modalRef || !modalRef.current) return;
    modalRef.current.focus();
  }, []);

  return (
    <StyledModalContainer
      tabIndex={0}
      ref={modalRef}
      onKeyDown={e => {
        if (e.key === 'Escape' || e.key === 'Esc') {
          onDismiss && onDismiss();
        }
      }}
      {...props}
    >
      {!hideCloseButton && (
        <CloseIcon color="#96A8C2" onClick={onDismiss} className="close" />
      )}
      {children}
    </StyledModalContainer>
  );
};

export default memo(Modal);
