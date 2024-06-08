import React, { FC } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Flex } from 'components/Box';
import { StyledTabs, StyledTabPanel } from './styles';
import { Tab, TabList } from 'react-tabs';
import DailyTabs from './Tabs/Daily';
import WeeklyTabs from './Tabs/Weekly';

const Race: FC<{ tabIndex?: number }> = ({ tabIndex = 0 }) => {
  const navigate = useNavigate();
  const handleSelect = (index: number) => {
    const tabs = ['/daily-race', '/weekly-raffle'];
    navigate(tabs[index]);
  };

  return (
    <Box
      maxWidth={'1094px'}
      marginX="auto"
      padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}
    >
      <StyledTabs selectedIndex={tabIndex} onSelect={handleSelect}>
        <TabList>
          <Flex>
            <Tab>
              DAILY RACE <b />
            </Tab>
            <Tab>
              WEEKLY RAFFLE <b />
            </Tab>
          </Flex>
        </TabList>
        <StyledTabPanel>
          <DailyTabs />
        </StyledTabPanel>
        <StyledTabPanel>
          <WeeklyTabs />
        </StyledTabPanel>
      </StyledTabs>
    </Box>
  );
};

export default Race;
