import React from 'react';

import { Flex, FlexProps } from 'components/Box';
import { Span } from 'components/Text';
import { Button } from 'components/Button';
import { FairnessIcon } from 'components/Icon';

interface TopBarProps extends FlexProps {
  title?: string;
  fee?: number;
}

export default function Topbar({ title = '', fee, ...props }: TopBarProps) {
  return (
    <Flex gap={20} flexWrap="wrap" {...props}>
      <Span color={'#768BAD'} fontSize={23} fontWeight={600}>
        {title}
      </Span>

      <Flex gap={20} height="34px">
        <Button
          border="2px solid #4F617B"
          borderRadius={'0px'}
          borderWidth="0px 0px 0px 2px"
          background={'#070C12'}
          color="#4F617B"
          fontWeight={500}
          nonClickable={true}
        >
          <FairnessIcon />
          Fair Game
        </Button>

        {fee !== undefined && (
          <Button
            border="2px solid #4F617B"
            borderRadius={'0px'}
            borderWidth="0px 0px 0px 2px"
            background={'#070C12'}
            color="#4F617B"
            fontWeight={500}
            nonClickable={true}
          >
            <FairnessIcon />
            {fee!}% Fee
          </Button>
        )}
      </Flex>
    </Flex>
  );
}
