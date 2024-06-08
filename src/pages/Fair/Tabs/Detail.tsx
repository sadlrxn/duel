import React from 'react';
import styled from 'styled-components';
import copy from 'copy-to-clipboard';
import { toast } from 'react-toastify';
import { ClipLoader } from 'react-spinners';

import { ReactComponent as CopyIcon } from 'assets/imgs/icons/copy.svg';

import { Input, Flex, Text, Button, BoxProps } from 'components';

interface TableItemProps {
  text?: string;
  enableCopy?: boolean;
}

export function TableItem({ text = '', enableCopy = false }: TableItemProps) {
  return (
    <>
      <Flex alignItems="center">
        <Text fontSize="19px" lineHeight="1.2em" color="white" width={180}>
          {text.length > 14 ? text.slice(0, 8) + '...' + text.slice(-4) : text}
        </Text>
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
      </Flex>
    </>
  );
}

interface DetailProps extends BoxProps {
  title?: string;
  readOnly?: boolean;
  placeholder?: string;
  text?: string;
  setText?: any;
  type?: string;
  enableCopy?: boolean;
  buttonText?: string;
  buttonClick?: any;
  isLoading?: boolean;
}

export default function Detail({
  title = '',
  readOnly = false,
  placeholder = '',
  type = 'text',
  setText,
  text = '',
  enableCopy = false,
  buttonClick,
  buttonText = '',
  isLoading = false,
  ...props
}: DetailProps) {
  return (
    <>
      <Flex
        flexDirection="column"
        fontSize="16px"
        lineHeight="1.2em"
        gap={8}
        {...props}
      >
        <Text color="#768BAD">{title}</Text>
        <RowContainer>
          <InputBox readOnly={readOnly}>
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
          {buttonText !== '' && (
            <StyledButton onClick={isLoading ? undefined : buttonClick}>
              {isLoading ? <ClipLoader color="#fff" size={20} /> : buttonText}
            </StyledButton>
          )}
        </RowContainer>
      </Flex>
    </>
  );
}

const InputBox = styled(Flex)<{ readOnly?: boolean }>`
  align-items: center;
  gap: 20px;
  padding: 10px 20px;
  background: ${({ readOnly }) => (readOnly ? '#121C2A' : '#03060999')};
  border-radius: 11px;
  flex-grow: 1;
  height: 42px;

  width: 100%;
  .width_700 & {
    width: auto;
  }

  input {
    color: white;
    flex-grow: 1;
    border: none;
    outline: none;
    background: transparent;
    font-size: 16px;
    line-height: 19px;

    &::placeholder {
      color: #ffffff80;
    }

    ::-webkit-outer-spin-button,
    ::-webkit-inner-spin-button {
      -webkit-appearance: none;
      margin: 0;
    }
    -moz-appearance: textfield;
  }
`;

const CopyButton = styled(Button)`
  min-width: 28px;
  min-height: 28px;

  svg {
    min-width: 12px;
    min-height: 12px;
  }
`;

CopyButton.defaultProps = { variant: 'secondary' };

const StyledButton = styled(Button)`
  border: 2px solid ${({ theme }) => theme.colors.success};
  background: linear-gradient(180deg, #070b10 0%, rgba(7, 11, 16, 0.3) 100%);
  border-radius: 6.75px;

  font-size: 14px;
  font-weight: 600;
  line-height: 17px;
  letter-spacing: 0.16em;
  color: white;

  padding: 9px 13px;
  gap: 10px;

  text-transform: uppercase;
  min-width: max-content;

  width: 100%;
  .width_700 & {
    width: max-content;
  }
`;

const RowContainer = styled(Flex)`
  gap: 10px;
  justify-content: space-between;
  align-items: center;
  flex-direction: column;

  .width_700 & {
    gap: 25px;
    flex-direction: row;
  }
`;
