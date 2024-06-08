import React from 'react';
import styled from 'styled-components';

import { Flex, FlexProps, Button, Span } from 'components';

import { useAppSelector } from 'state';
import { useCrash } from 'hooks';

import { Players } from '../Players';
import { Persons } from '../Persons';
import AutoButton from './AutoButton';

interface AutomaticButtonsProps extends FlexProps {
  gameHeight: number;
  gameWidth: number;
}

export default function AutomaticButtons({
  gameHeight,
  gameWidth,
  ...props
}: AutomaticButtonsProps) {
  const { status } = useAppSelector(state => state.crash);
  const { showStatus, setShowStatus } = useCrash();

  return (
    <Container {...props}>
      <CustomPlayers
        className={'crash_mobile_list' + (!showStatus ? '--hide' : '')}
        maxHeight={`calc(100vh - ${
          gameHeight - 20 - 20 + 16 + (status === 'bet' ? 104 : 0)
        }px)`}
        width={`${gameWidth - 16 - (status === 'bet' ? 18 : 90)}px`}
      />
      <ToggleButton
        variant="secondary"
        border="1px solid #2B4160"
        borderRadius="10px"
        background="#24354C"
        height="100%"
        width="50px"
        minWidth="50px"
        type="button"
        onClick={() => {
          setShowStatus(!showStatus);
        }}
      >
        {showStatus ? (
          <Span
            style={{
              transform: 'rotate(-90deg) scaleX(0.6)'
            }}
            fontWeight={600}
            color="#B2D1FF"
            fontSize="20px"
          >
            {'<'}
          </Span>
        ) : (
          <Persons />
        )}
      </ToggleButton>

      <AutoButton />
    </Container>
  );
}

const ToggleButton = styled(Button)``;

const CustomPlayers = styled(Players)`
  position: absolute;
  left: 16px;
  top: 0px;
  transform: translate(0px, -100%);
  border-bottom-left-radius: 0px;
  border-bottom-right-radius: 0px;

  transition: all 0.5s;

  .width_700 & {
    display: none;
  }
`;

const Container = styled(Flex)`
  gap: 12px;

  position: absolute;
  transform: translate(0px, -100%);
  padding: 17px 16px 15px;
  background: #131e2d;
  height: 74px;
  width: 100%;

  top: -35px;
  left: 0px;

  margin-top: 0px;

  .width_700 & {
    position: relative;
    transform: translate(0, 0);
    background: transparent;
    top: 0px;
    padding: 0px;
    height: 40px;
    min-height: 40px;
  }

  .crash_mobile_list--hide {
    display: none;
  }
`;
