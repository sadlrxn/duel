import { Flex, Span, Button } from 'components';
import { FairnessIcon } from 'components/Icon';

export default function TopBar() {
  return (
    <Flex gap={30} flexWrap="wrap">
      <Span color={'#768BAD'} fontSize={23} fontWeight={600}>
        COIN FLIP
      </Span>

      <Flex gap={20} height="34px">
        <Button
          border="2px solid #4F617B"
          borderRadius="0px"
          borderWidth="0px 0px 0px 2px"
          background="#070C12"
          color="#4F617B"
          fontWeight={500}
          nonClickable={true}
        >
          <FairnessIcon />
          Fair Game
        </Button>

        <Button
          border="2px solid #4F617B"
          borderRadius="0px"
          borderWidth="0px 0px 0px 2px"
          background="#070C12"
          color="#4F617B"
          fontWeight={500}
          nonClickable={true}
        >
          <FairnessIcon />
          4% Fee
        </Button>
      </Flex>
    </Flex>
  );
}
