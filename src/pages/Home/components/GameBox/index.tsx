import React, { FC } from 'react';
import styled from 'styled-components';
import { Link } from 'react-router-dom';
import { Box, Flex } from 'components/Box';
import { Text } from 'components/Text';

const StyledBox = styled(Box)`
  transition: 0.5s;
  &:hover {
    transform: translateY(-10px);
  }
`;

const GameBox: FC<{
  img: string;
  title?: string;
  to?: string;
  toGrayScale?: boolean;
}> = ({ img, title = 'Coin flip', to = '/', toGrayScale = false }) => {
  const component = (
    <StyledBox>
      <Flex alignItems={'baseline'} position="relative">
        <img
          src={img}
          alt="game card"
          width={'100%'}
          style={toGrayScale ? { filter: 'grayscale(100%)' } : {}}
        />
        <AbsoluteBox>
          <Text
            fontFamily={'Termina'}
            fontSize="6px"
            textTransform="uppercase"
            fontWeight={600}
            lineHeight="8px"
            color={'white'}
          >
            {title === '' ? '' : 'Duel ORIGINAL'}
          </Text>

          <Text
            fontFamily={'Termina'}
            fontSize="16.5px"
            textTransform="uppercase"
            fontWeight={800}
            lineHeight="20px"
            color={'white'}
          >
            {title}
          </Text>
        </AbsoluteBox>
      </Flex>
    </StyledBox>
  );

  if (toGrayScale) return component;
  else
    return (
      <Link to={to} style={{ flex: '1 1 0px' }}>
        {component}
      </Link>
    );
};

const AbsoluteBox = styled(Box)`
  position: absolute;
  left: 20px;
  bottom: 15px;
`;

export default GameBox;
