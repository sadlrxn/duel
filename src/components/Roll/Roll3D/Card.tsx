import React from 'react';
import styled from 'styled-components';
import { LazyLoadImage } from 'react-lazy-load-image-component';

import { Box, BoxProps, Flex } from 'components/Box';
import { Text, Span } from 'components/Text';
import { Badge } from 'components/Badge';
import { formatUserName } from 'utils/format';
import { imageProxy } from 'config';

export interface SpinCardProps extends BoxProps {
  name?: string;
  percent?: number;
  avatar?: string;
  count?: number;
}

export default function Card({
  name = 'Username1',
  percent = 10,
  avatar = 'https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://www.arweave.net/FEFMTQEgWHhDd33e2N2ldQZ93Bk0BSVKxp7TPdP-3ao',
  ...props
}: SpinCardProps) {
  return (
    <Container {...props}>
      <Image src={imageProxy() + avatar} alt="" />
      <Flex flexDirection="column" gap={4} width="100%" alignItems="center">
        <Text
          textAlign="center"
          fontSize="30px"
          style={{
            width: '190px',
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis'
          }}
        >
          {formatUserName(name)}
        </Text>
        <Badge>
          <Span fontWeight={600} fontSize="25px">
            {percent.toFixed(2)}%
          </Span>
        </Badge>
      </Flex>
    </Container>
  );
}

const Image = styled(LazyLoadImage)`
  width: 100%;
  border-radius: 50px;
  padding: 20px;
  padding-bottom: 0px;
  aspect-ratio: 1.1;
  background: transparent;
`;

const Container = styled(Box)`
  width: 200px;
  height: 300px;

  box-sizing: border-box;
  border-radius: 13px;
  line-height: 1.3em;
  color: white;
  background: #4c698933;

  position: absolute;
`;
