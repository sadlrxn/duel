import React from 'react';
import dayjs from 'dayjs';
import { Box, Text, Flex } from 'components';
import RaffleLogo from '../Icons/RaffleLogo';

interface RaffleTicketProps {
  ticketId?: number;
  date?: string;
}
export default function RaffleTicket({
  ticketId = 0,
  date = '2023-03-01T08:26:01.2552809Z'
}: RaffleTicketProps) {
  return (
    <Flex>
      <Box p="10px" background={'#182738'} borderRadius="8px 4px 4px 8px">
        <RaffleLogo />
      </Box>
      <Box
        p="10px 20px"
        background={'#182738'}
        borderRadius="4px 8px 8px 4px"
        borderLeft={'2px dashed #FFFFFF1A'}
        width="100%"
        textAlign="center"
      >
        <Text color="#D7D7D7" fontSize={'20px'} fontWeight={700}>
          {ticketId.toString().padStart(4, '0')}
        </Text>
        <Text color="#FFFFFF80" fontSize={'10px'} fontWeight={500}>
          {dayjs(new Date(date)).format('MM/DD/YYYY')}
        </Text>
      </Box>
    </Flex>
  );
}
