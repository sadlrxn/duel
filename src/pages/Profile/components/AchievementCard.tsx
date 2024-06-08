import React, { FC } from "react";
import { Box, Flex } from "components/Box";
import { ReactComponent as PolygonIcon } from "assets/imgs/icons/polygon.svg";
import { Text } from "components/Text";

const AchievementCard: FC<{
  background: string;
  color: string;
  subject: string;
  content: string;
  current: string;
}> = ({ background, color, subject, content, current }) => {
  return (
    <Flex
      background={background}
      borderRadius="15px"
      p="20px"
      flex={"1 1 0px"}
      alignItems="center"
    >
      <PolygonIcon color={color} width="75px" height="75px" />

      <Box ml="15px">
        <Text color={color} fontSize="16px" fontWeight={600} mb="4px">
          {subject}
        </Text>
        <Text color="white" fontSize={"13px"}>
          {content}
        </Text>
        <Text color={"#96A8C2"}>current {current}</Text>
      </Box>
    </Flex>
  );
};

export default AchievementCard;
