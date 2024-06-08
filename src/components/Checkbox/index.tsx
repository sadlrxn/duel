import React, { FC, HTMLProps } from 'react';
import styled from 'styled-components';
import { Span } from 'components/Text';

const StyledCheckbox = styled.label<{ available: boolean }>`
  display: flex;
  align-items: center;

  opacity: ${({ available }) => (available ? 1 : 0.3)};

  margin-right: 1rem;
  padding-left: 2rem;
  position: relative;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;

  input[type='checkbox'] {
    position: absolute;
    opacity: 0;

    &:checked ~ b {
      color: #4fff8b;
      background: #0b121b
        url("data:image/svg+xml,%3Csvg width='14' height='11' viewBox='0 0 14 11' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M12.0482 0.770421L12.0482 0.770478L12.0544 0.763902C12.2055 0.603565 12.4106 0.510224 12.6261 0.500444L12.6293 0.500349C12.9732 0.490388 13.294 0.696011 13.4317 1.03071L13.4321 1.03183C13.5712 1.36755 13.4917 1.7548 13.2384 2.00466L13.2383 2.0046L13.2318 2.01124C11.8155 3.46119 10.5923 4.77711 9.36763 6.09459L9.34564 6.11825C8.11317 7.44413 6.87785 8.77263 5.43508 10.2494C5.12589 10.5643 4.63736 10.5846 4.3038 10.2978C4.3037 10.2978 4.3036 10.2977 4.3035 10.2976L0.808973 7.27742L0.808979 7.27741L0.806535 7.27533C0.63418 7.12829 0.523728 6.9145 0.503344 6.67927C0.483133 6.44602 0.555211 6.21513 0.702142 6.03609C0.848267 5.85908 1.05415 5.75105 1.27356 5.73165C1.4951 5.71316 1.71586 5.78459 1.88605 5.93157L1.8862 5.93171L4.43209 8.1287L4.79124 8.43863L5.11961 8.09626C6.15596 7.01571 7.10361 5.99975 8.05206 4.98294C9.31265 3.6315 10.5746 2.27855 12.0482 0.770421Z' fill='%234FFF8B' stroke='%234FFF8B'/%3E%3C/svg%3E%0A")
        50% 50% no-repeat;
    }
  }

  b {
    border-radius: 6px;
    position: absolute;
    left: 0;
    top: 50%;
    width: 23px;
    height: 23px;
    background-color: #0b121b;

    transform: translate(0, -50%);

    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    cursor: pointer;
  }
`;

const Checkbox: FC<HTMLProps<HTMLInputElement>> = ({
  label,
  disabled,
  ...props
}) => {
  return (
    <StyledCheckbox available={!disabled}>
      <input type="checkbox" {...props} disabled={disabled} />
      <b />
      <Span
        color={'#768BAD'}
        fontSize="16px"
        fontWeight={500}
        style={{ cursor: 'pointer' }}
      >
        {label}
      </Span>
    </StyledCheckbox>
  );
};

export default Checkbox;
