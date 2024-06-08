import React, {
  createContext,
  FC,
  useRef,
  useState,
  useEffect,
  ReactNode
} from 'react';
import { AnimatePresence, domAnimation, LazyMotion, m } from 'framer-motion';
import styled, { css } from 'styled-components';
import { Handler } from './types';
import {
  animationHandler,
  animationMap,
  animationVariants,
  appearAnimation,
  disappearAnimation
} from 'utils/animationToolkit';
import { Overlay } from 'components/Overlay';

interface ModalsContext {
  isOpen: boolean;
  nodeId: string;
  modalNode: React.ReactNode;
  setModalNode: React.Dispatch<React.SetStateAction<React.ReactNode>>;
  onPresent: (
    node: React.ReactNode,
    newNodeId: string,
    closeOverlayClick: boolean,
    renderOverlay: boolean,
    asModalCss: boolean
  ) => void;
  onDismiss: Handler;
}

const ModalWrapper = styled(m.div)<{ asModal: boolean }>`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: fixed;
  right: 0;
  left: 0;
  ${({ asModal }) =>
    asModal
      ? css`
          top: 0;
          bottom: 0;
        `
      : css`
          top: 65px;
          bottom: 0;
        `}

  z-index: 250;
  will-change: opacity;
  opacity: 0;
  &.appear {
    animation: ${appearAnimation} 0.3s ease-in-out forwards;
  }
  &.disappear {
    animation: ${disappearAnimation} 0.3s ease-in-out forwards;
  }
`;

export const Context = createContext<ModalsContext>({
  isOpen: false,
  nodeId: '',
  modalNode: null,
  setModalNode: () => null,
  onPresent: () => null,
  onDismiss: () => null
});

const ModalProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [modalNode, setModalNode] = useState<React.ReactNode>();
  const [nodeId, setNodeId] = useState('');
  const [closeOnOverlayClick, setCloseOnOverlayClick] = useState(true);
  const [renderOverlay, setRenderOverlay] = useState(true);
  const [asModal, setAsModal] = useState(true);
  const animationRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const setViewportHeight = () => {
      const vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty('--vh', `${vh}px`);
    };
    setViewportHeight();
    window.addEventListener('resize', setViewportHeight);
    return () => window.removeEventListener('resize', setViewportHeight);
  }, []);

  const handlePresent = (
    node: React.ReactNode,
    newNodeId: string,
    closeOverlayClick: boolean,
    renderOverlay: boolean,
    asModalCss: boolean
  ) => {
    setModalNode(node);
    setIsOpen(true);
    setNodeId(newNodeId);
    setCloseOnOverlayClick(closeOverlayClick);
    setRenderOverlay(renderOverlay);
    setAsModal(asModalCss);
  };

  const handleDismiss = () => {
    setModalNode(undefined);
    setIsOpen(false);
    setNodeId('');
    setCloseOnOverlayClick(true);
    setRenderOverlay(true);
  };

  const handleOverlayDismiss = () => {
    if (closeOnOverlayClick) {
      handleDismiss();
    }
  };

  return (
    <Context.Provider
      value={{
        isOpen,
        nodeId,
        modalNode,
        setModalNode,
        onPresent: handlePresent,
        onDismiss: handleDismiss
      }}
    >
      <LazyMotion features={domAnimation}>
        <AnimatePresence>
          {isOpen && (
            <ModalWrapper
              ref={animationRef}
              onAnimationStart={() => animationHandler(animationRef.current)}
              {...animationMap}
              variants={animationVariants}
              transition={{ duration: 0.3 }}
              asModal={asModal}
            >
              {renderOverlay && <Overlay onClick={handleOverlayDismiss} />}
              {React.isValidElement(modalNode) &&
                React.cloneElement(modalNode as React.ReactElement<any>, {
                  onDismiss: modalNode.props.hasOwnProperty('onDismiss')
                    ? modalNode.props.onDismiss
                    : handleDismiss
                })}
            </ModalWrapper>
          )}
        </AnimatePresence>
      </LazyMotion>

      {children}
    </Context.Provider>
  );
};

export default ModalProvider;
