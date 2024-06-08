import React, { FC, useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { StyledTabs, StyledTabPanel } from './styles';
import { Flex } from 'components/Box';
import { Tab, TabList } from 'react-tabs';
import { Box } from 'components/Box';
import { useQuery } from 'hooks';

import ProfileTab from './Tabs/Profile';
import NFTsTab from './Tabs/NFTs';
import ReferralTab from './Tabs/Referral';
import RewardsTab from './Tabs/Rewards';
import SelfExclusionTab from './Tabs/SelfExclusion';

const Profile: FC = () => {
  const query = useQuery();
  const navigate = useNavigate();

  const [tabIndex, setTabIndex] = useState(0);

  useEffect(() => {
    const tab = query.get('tab');
    switch (tab) {
      case 'profile':
        setTabIndex(0);
        break;
      case 'nfts':
        setTabIndex(1);
        break;
      case 'rewards':
        setTabIndex(2);
        break;
      case 'referral':
        setTabIndex(3);
        break;
      case 'self-exclusion':
        setTabIndex(4);
        break;

      default:
        setTabIndex(0);
        break;
    }
  }, [query]);

  const handleSelect = (index: number) => {
    const tabs = ['profile', 'nfts', 'rewards', 'referral', 'self-exclusion'];
    navigate(`/profile?tab=${tabs[index]}`);
  };

  return (
    <Box padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}>
      <StyledTabs selectedIndex={tabIndex} onSelect={handleSelect}>
        <TabList>
          <Flex>
            <Tab>
              PROFILE
              <b />
            </Tab>
            <Tab>
              NFTS
              <b />
            </Tab>
            <Tab>
              REWARDS
              <b />
            </Tab>
            <Tab>
              Referral
              <b />
            </Tab>
            <Tab>
              Self exclusion
              <b />
            </Tab>
          </Flex>
        </TabList>
        <StyledTabPanel>
          <ProfileTab />
        </StyledTabPanel>
        <StyledTabPanel>
          <NFTsTab />
        </StyledTabPanel>
        <StyledTabPanel>
          <RewardsTab />
        </StyledTabPanel>
        <StyledTabPanel>
          <ReferralTab />
        </StyledTabPanel>
        <StyledTabPanel>
          <SelfExclusionTab />
        </StyledTabPanel>
      </StyledTabs>
    </Box>
  );
};
export default Profile;
