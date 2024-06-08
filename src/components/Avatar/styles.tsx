import styled from 'styled-components';
import {
  LazyLoadImage,
  LazyLoadImageProps
} from 'react-lazy-load-image-component';

export interface StyledAvatarProps extends LazyLoadImageProps {
  size: string;
  border: string;
  borderradius: string;
  padding: string;
  background: string;
  role?: string;
  filter?: string;
}

const StyledAvatar = styled(LazyLoadImage)<StyledAvatarProps>`
  width: ${({ size }) => size};
  height: ${({ size }) => size};
  border: ${({ border }) => border};
  border-radius: ${({ borderradius }) => borderradius};
  padding: ${({ padding }) => padding};
  background-color: ${({ background }) => background};
  filter: ${({ filter }) => filter};
  flex: none;
`;

export default StyledAvatar;
