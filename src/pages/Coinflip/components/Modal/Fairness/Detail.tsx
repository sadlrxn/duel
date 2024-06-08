import React from 'react';
import copy from 'copy-to-clipboard';
import { toast } from 'react-toastify';

import { ReactComponent as CopyIcon } from 'assets/imgs/icons/copy.svg';

import { Input, Text, Flex } from 'components';

import { Coin } from '../../Coin';
import { CopyButton, InputBox } from './styles';

interface DetailProps {
  title?: string;
  readOnly?: boolean;
  side?: string;
  enableCopy?: boolean;
  text?: string;
  type?: string;
  setText?: any;
  placeholder?: string;
}

const Detail: React.FC<DetailProps> = ({
  title = '',
  readOnly = false,
  side = '',
  text = '',
  type = 'text',
  setText,
  placeholder = '',
  enableCopy = false
}) => {
  return (
    <Flex flexDirection="column" gap={8}>
      <Text
        fontSize="0.8em"
        fontWeight={400}
        fontStyle="italic"
        color="#768BAD"
      >
        {title}
      </Text>

      <InputBox readOnly={readOnly}>
        {side !== '' && <Coin side={side} size={28} />}
        <Input
          readOnly={readOnly}
          value={text}
          type={type}
          placeholder={placeholder}
          onChange={e => {
            setText(e.target.value);
          }}
        />
        {enableCopy === true && (
          <CopyButton
            onClick={() => {
              copy(text);
              toast.success('Copy to clipboard success.');
            }}
          >
            <CopyIcon width={12} height={12} />
          </CopyButton>
        )}
      </InputBox>
    </Flex>
  );
};

export default React.memo(Detail);
