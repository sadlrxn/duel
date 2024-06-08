import styled from 'styled-components';

export const ChipIcon = styled.div<{
  $size?: number;
  $background?: string;
  $border?: string;
}>`
  background: ${({ $background }) => $background};
  /* background-color: #ffe24b; */
  border: 2px solid ${({ $border }) => $border};
  border-radius: 9999px;
  width: ${({ $size }) => `${$size}px`};
  height: ${({ $size }) => `${$size}px`};
`;

ChipIcon.defaultProps = {
  $size: 12,
  $background: '#ffe24b',
  $border: '#ffb31f'
};

export const Container = styled.div`
  display: inline-flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 5px;
`;
