import React from 'react';
import copy from 'copy-to-clipboard';
import { toast } from 'react-toastify';

import { ReactComponent as CopyIcon } from 'assets/imgs/icons/copy.svg';

import { Input, Text, Flex } from 'components';

import { CopyButton, InputBox, ChipIcon } from './styles';

const StatusItem = ({ status }: { status: string }) => {
  return (
    <Flex
      borderRadius="30px"
      px="30px"
      py="4px"
      background={status === 'Active' ? '#4FFF8B' : '#FFCE4F'}
      justifyContent="center"
      alignItems="center"
      fontSize="12px"
      fontWeight={700}
      color="black"
    >
      {status}
    </Flex>
  );
};

interface DetailProps {
  title?: string;
  readOnly?: boolean;
  enableCopy?: boolean;
  showChip?: boolean;
  text?: string;
  placeholder?: string;
  setText?: any;
  type?: string;
  status?: string;
}

const Detail: React.FC<DetailProps> = ({
  title = '',
  readOnly = false,
  showChip = false,
  text = '',
  placeholder = '',
  enableCopy = false,
  setText,
  type = 'text',
  status = ''
}) => {
  return (
    <Flex flexDirection="column" gap={8} mt={2}>
      <Text
        fontSize="0.8em"
        fontWeight={400}
        fontStyle="italic"
        color="#768BAD"
      >
        {title}
      </Text>

      <InputBox readOnly={readOnly}>
        {showChip === true && <ChipIcon width={14} height={14} />}
        <Input
          readOnly={readOnly}
          value={text}
          type={type}
          placeholder={placeholder}
          onChange={e => {
            setText(e.target.value);
          }}
        />
        {status !== '' && <StatusItem status={status} />}
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
