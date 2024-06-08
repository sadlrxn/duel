import { FC } from 'react';
import styled from 'styled-components';

import { Badge } from 'components/Badge';
import { Flex } from 'components/Box';
import { useModal } from 'components/Modal';
import { ProfileModal } from 'components/Modals';
import { imageProxy } from 'config';

import StyledAvatar from './styles';
import AdminAvatar from './AdminAvatar';
import ModeratorAvatar from './ModeratorAvatar';
import AmbassadorAvatar from './AmbassadorAvatar';

const StyledBadge = styled(Badge)`
  position: absolute;
  bottom: -4px;
  left: 35%;
  background: #0a111a;
  color: white;
  border: 2px solid #ee00e1;
  border-radius: 13px;
  padding: 0px 7px;
`;

export interface AvatarProps {
  image?: string;
  level?: number;
  name?: string;
  size?: string;
  border?: string;
  borderRadius?: string;
  padding?: string;
  userId?: number;
  background?: string;
  filter?: string;
  role?: string;
  useProxy?: boolean;
}

const Avatar: FC<AvatarProps> = ({
  image,
  name = '',
  level,
  size = '42px',
  border = '0.84px solid #5E6E88',
  padding = '2px',
  borderRadius = '50%',
  userId = 0,
  background = 'transparent',
  filter,
  role,
  useProxy = true,
  ...props
}) => {
  const [onProfileModal] = useModal(
    <ProfileModal userId={userId} name={name} avatar={image} />,
    true
  );

  return (
    <Flex
      justifyContent="center"
      alignItems="center"
      position="relative"
      onClick={
        userId === 0 || name.toLowerCase() === 'duelbot'
          ? undefined
          : (e: any) => {
              e.stopPropagation();
              onProfileModal();
            }
      }
      cursor={
        userId === 0 || name.toLowerCase() === 'duelbot' ? 'auto' : 'pointer'
      }
    >
      {role === 'admin' ? (
        <AdminAvatar
          id={userId}
          image={useProxy ? imageProxy() + image : image}
          size={size}
        />
      ) : role === 'moderator' ? (
        <ModeratorAvatar
          id={userId}
          image={useProxy ? imageProxy() + image : image}
          size={size}
        />
      ) : role === 'ambassador' ? (
        <AmbassadorAvatar
          id={userId}
          image={useProxy ? imageProxy() + image : image}
          size={size}
        />
      ) : (
        <StyledAvatar
          src={useProxy ? imageProxy() + image : image}
          alt="avatar"
          size={size}
          border={border}
          padding={padding}
          borderradius={borderRadius}
          background={background}
          filter={filter}
          {...props}
        />
      )}

      {level && <StyledBadge>{level}</StyledBadge>}
    </Flex>
  );
};

export default Avatar;
