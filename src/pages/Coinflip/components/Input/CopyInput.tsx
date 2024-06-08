import { Flex, Button } from "components/index";
import copy from "assets/imgs/icons/copy.svg";

import { Container, StyledLabel, Image, Input } from "./CopyInput.styles";

interface IInputProps {
  name: string;
  value: any;
  label?: String;
  placeholder?: string;
  readOnly?: boolean;
  onCopy?: any;
}
export default function CopyInput({
  label = "",
  onCopy,
  ...rest
}: IInputProps) {
  return (
    <Container>
      {label && <StyledLabel>{label}</StyledLabel>}
      <Flex flexDirection="row" alignItems="center" gap={20}>
        <Input type="text" disabled {...rest} />

        {onCopy && (
          <Button
            variant="secondary"
            backgroundColor="#2a3d57"
            scale="lg"
            width="52px"
            onClick={onCopy}
          >
            <Image src={copy} alt="" />
          </Button>
        )}
      </Flex>
    </Container>
  );
}
