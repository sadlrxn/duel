import React from 'react';
import {
  LazyLoadImage,
  LazyLoadImageProps
} from 'react-lazy-load-image-component';
import styled from 'styled-components';
import { NFT } from 'api/types/nft';
import { Flex } from '../Box';
import { Span } from '../Text';
import { Badge } from '../Badge';
import { Chip } from '../Chip';
import { ReactComponent as PlusIcon } from 'assets/imgs/icons/plus.svg';
import { ReactComponent as CheckedIcon } from 'assets/imgs/icons/checked.svg';
import { Box } from 'components/Box';
import { imageProxy } from 'config';
import { Button } from 'components';
import ThreeDot from 'components/Icon/ThreeDot';
import { convertBalanceToChip } from 'utils/balance';

const StyledLazyLoadImage = styled(LazyLoadImage)<LazyLoadImageProps>`
  border-radius: 6px;
  height: 100%;
  width: 100%;
  margin: auto;
  object-fit: cover;
`;

export interface NFTCardProps extends NFT {
  selectable?: boolean;
  selected?: boolean;
  onClick?: () => void;
}

const StyledNFTCard = styled(Flex)<{ selectable: boolean; selected: boolean }>`
  border: 1.5px solid
    ${({ theme, selected }) =>
      selected ? `${theme.colors.success}` : `${theme.colors.text}80`};
  border-radius: 14px;
  flex-direction: column;
  gap: 4px;
  padding: 8px;

  width: 100%;

  /* &:hover {
    ${({ selectable }) =>
    selectable &&
    `
    box-shadow: 0 0 0 5px #768bad33;
  `}
  } */

  ${({ selectable }) =>
    selectable &&
    `
    cursor: pointer;
  `}

  ${({ selected }) =>
    selected &&
    `
    box-shadow: 0 0 0 5px #4fff8b33;
  `}
`;

const StyledPlusIcon = styled(PlusIcon)`
  position: absolute;
  top: 13px;
  right: 13px;
`;

const StyledCheckedIconIcon = styled(CheckedIcon)`
  position: absolute;
  top: 13px;
  right: 13px;
`;

export default function NFTCard({
  name,
  image,
  collectionName,
  price = 0,
  selectable = false,
  selected = false,
  onClick
}: NFTCardProps) {
  return (
    <StyledNFTCard
      selectable={selectable}
      selected={selected}
      onClick={onClick}
    >
      <Box position={'relative'} height={225} width="100%">
        <StyledLazyLoadImage
          width={200}
          height={200}
          src={imageProxy(300) + image}
          alt={name}
        />
        {selectable && !selected && <StyledPlusIcon />}
        {selectable && selected && <StyledCheckedIconIcon />}
      </Box>

      <Box background={'#1A2534'} borderRadius="6px" mt="8px" p="8px">
        <Flex justifyContent={'space-between'}>
          <Span color="white" fontSize={'14px'} fontWeight={700}>
            {collectionName}
          </Span>
          <Span color="white" fontSize={'14px'} fontWeight={700}>
            #{name?.split('#')[1]}
          </Span>
        </Flex>

        <Flex alignItems={'center'} justifyContent="space-between" mt="5px">
          <Badge>
            <Chip
              color="success"
              price={convertBalanceToChip(price)}
              fontWeight={700}
            />
          </Badge>

          {/* <Button size={'23px'} background="#242F42">
            <ThreeDot />
          </Button> */}
        </Flex>
      </Box>
    </StyledNFTCard>
  );
}
