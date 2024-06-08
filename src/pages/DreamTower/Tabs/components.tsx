import styled from 'styled-components';
import { Button } from 'components';

export const UpperBtn = styled(Button)`
  padding: 7.2px 5.5px;
  border-radius: 7px 7px 0px 0px;
  background: #24354d;
`;

export const DownBtn = styled(Button)`
  padding: 7.2px 5.5px;
  border-radius: 0px 0px 7px 7px;
  background: #24354d;
  svg {
    transform: rotate(180deg);
  }
`;

export const StyledInput = styled.input`
  font-size: 20px !important;
`;

export const ModeBtn = styled(Button)<{ selected: boolean }>`
  width: 100%;
  height: 38px;
  background: #242f42;
  font-size: 14px;
  font-weight: 600;
  color: #768bad;
  border-radius: 5px;

  ${({ selected }) => selected && `border: 1px solid #4FFF8B; color: #4FFF8B;`};
`;
