import styled from "styled-components";
import { Flex, Span } from "components/index";

interface StatusButtonProps {
  text: string;
  children?: React.ReactNode;
}

const StatusButtonWrapper = styled(Flex)`
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.black};
  padding: 0 10px;
  margin: 7px 0;
`;

export default function StatusButton({ text, children }: StatusButtonProps) {
  return (
    <StatusButtonWrapper alignItems="center" gap={8}>
      {children}
      <Span color="text" fontSize={14}>
        {text}
      </Span>
    </StatusButtonWrapper>
  );
}
