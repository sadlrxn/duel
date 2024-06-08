import React from 'react';
import styled from 'styled-components';

import { Flex, FlexProps } from 'components/Box';
import { Span } from 'components/Text';
import { Badge } from 'components/Badge';
import { formatUserName } from 'utils/format';
import { imageProxy } from 'config';

const Image = styled.img`
  border-radius: 13px;
  background-color: #4c6989;
  width: 48px;
  height: 48px;
`;

const Container = styled(Flex)`
  width: 120px;
  height: 145px;

  flex-direction: column;
  /* justify-content: center; */
  align-items: center;
  padding: 13px 24px 28px;
  gap: 8px;

  font-size: 12px;
  line-height: 15px;
  color: white;
`;

export interface SpinCardProps extends FlexProps {
  name?: string;
  percent?: number;
  avatar?: string;
  count?: number;
}

export default function SpinCard({
  name = 'Username1',
  percent = 10,
  avatar = 'https://beta.api.solanalysis.com/images/400x400/filters:frames(,0)/https://www.arweave.net/FEFMTQEgWHhDd33e2N2ldQZ93Bk0BSVKxp7TPdP-3ao',
  ...props
}: SpinCardProps) {
  return (
    <Container {...props}>
      <Span fontSize="14px" mt="14px">
        {formatUserName(name)}
      </Span>
      <Badge>
        <Span fontWeight={600}>{percent.toFixed(2)}%</Span>
      </Badge>
      <Image src={imageProxy() + avatar} alt="" />
    </Container>
  );
}
