import React, { FC } from 'react';

import { Box, Flex } from 'components/Box';
import { Text } from 'components/Text';
import { ReactComponent as DotIcon } from 'assets/imgs/icons/dot.svg';
import { ReactComponent as CoinIcon } from 'assets/imgs/coins/coin.svg';
import {
  StyledText,
  StyledContainer,
  StyledUserInfo,
  No,
  Avatar,
  StyleButton
} from './styles';
import { Options } from 'react-select';
import Select from 'components/Select';

export interface UserInfo {
  user: {
    avatar: string;
    name: string;
  };
  info: object;
}

const TopItem: FC<{
  title: string;
  options: Options<{
    label: string;
    value: string;
  }>;
  isWinComp: boolean;
  data: UserInfo[];
}> = ({ title, options, isWinComp }) => {
  return (
    <Box width={'100%'}>
      <Flex alignItems="center">
        <StyledText>{title}</StyledText>
        <Select
          options={options}
          background="#1a3d39"
          hoverBackground="#4fff8b26"
          color="#4FFF8B"
          fontSize="20px"
        />
      </Flex>

      <Box background={'#16202C'} borderRadius="10px" p={'20px 22px'} mt="30px">
        <StyledUserInfo>
          <StyledContainer isWinComp={isWinComp}>
            <Flex alignItems={'center'}>
              <No>1</No>
              <DotIcon width={'20px'} height="20px" />

              <Avatar>
                <img
                  src={'https://avatars.githubusercontent.com/u/97012368?v=4'}
                  alt="avatar"
                />
              </Avatar>
              <Text color={'#fff'}>Username 01</Text>
            </Flex>
            {isWinComp ? (
              <Flex alignItems={'center'}>
                <CoinIcon style={{ marginRight: '7px' }} />
                <Text color={'#FFF6CA'}>145,332</Text>
              </Flex>
            ) : (
              <Flex alignItems={'center'} gap={8}>
                <Text color={'#B2D1FF'} fontWeight={600}>
                  1,1M &nbsp;
                  <small>XP</small>
                </Text>
                <Text color={'#4FFF8B'}>Dueler IV</Text>
              </Flex>
            )}
          </StyledContainer>
        </StyledUserInfo>

        <Flex justifyContent={'center'} mt="15px">
          <StyleButton>See more</StyleButton>
        </Flex>
      </Box>
    </Box>
  );
};

export default TopItem;
