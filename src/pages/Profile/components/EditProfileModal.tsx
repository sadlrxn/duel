import { ChangeEvent, FC, useCallback, useState, useEffect } from 'react';
import { ClipLoader } from 'react-spinners';
import { useDropzone } from 'react-dropzone';
import { useAppSelector } from 'state';
import { Modal, ModalProps } from 'components/Modal';
import { Box, Flex } from 'components/Box';
import { Text } from 'components/Text';
import Avatar from 'components/Avatar';
import { Button } from 'components/Button';
import { InputBox } from 'components/InputBox';
import Checkbox from 'components/Checkbox';
import { ReactComponent as TrashIcon } from 'assets/imgs/icons/trash.svg';
import api from 'utils/api';
import { useAppDispatch } from 'state';
import { saveProfile } from 'state/user/actions';
import { toast } from 'utils/toast';
import styled from 'styled-components';

const EditProfileModal: FC<ModalProps> = ({ onDismiss, ...props }) => {
  const dispatch = useAppDispatch();
  const { name, avatar, statistics } = useAppSelector(state => state.user);

  const [image, setImage] = useState<any>(undefined as any);
  const [imageLoadError, setImageLoadError] = useState(false);
  const onDrop = (acceptedFiles: any[]) => {
    try {
      setImage(
        Object.assign(acceptedFiles[0], {
          preview: URL.createObjectURL(acceptedFiles[0])
        })
      );
      setImageLoadError(false);
    } catch {
      setImageLoadError(true);
    }
  };

  const { getRootProps, getInputProps, open } = useDropzone({
    noClick: true,
    noKeyboard: true,
    noDrag: true,
    multiple: false,
    accept: { 'image/*': [] },
    maxSize: 819200,
    onDrop
  });

  const [request, setRequest] = useState(false);
  const [formData, setFormData] = useState({
    name: name,
    private: statistics.private_profile
  });

  const handleChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const format = /^[A-Za-z]\w*$/;
      const value = e.target.value;
      if (format.test(value) || value === '') {
        setFormData({ name: value, private: formData.private });
      }
    },
    [formData.private]
  );

  const handlePrivateChange = useCallback(() => {
    setFormData(prev => ({ name: prev.name, private: !prev.private }));
  }, []);

  const handleSaveProfile = async () => {
    const config = {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    };

    let data = new FormData();
    data.append('name', formData.name);
    data.append('image', image);
    data.append('isPrivate', formData.private ? '1' : '0');

    setRequest(true);
    try {
      if (formData.name.toLowerCase() === 'hidden') {
        toast.warning("Can't set name as 'Hidden'");
        return;
      }
      const response = await api.post('/user/update', data, config);
      dispatch(
        saveProfile({
          avatar: response.data.avatar ? response.data.avatar : avatar,
          name: response.data.name,
          statistics: {
            ...statistics,
            private_profile: response.data.isPrivate
          }
        })
      );
      toast.success('Successfully saved.');
      if (onDismiss) onDismiss();
    } catch (error: any) {
      if (error.response.status === 429) {
        toast.error(error.response.data.message);
      } else toast.error(error.response.data.status);
    }
    setRequest(false);
  };

  useEffect(() => {
    setFormData(prev => ({
      name: prev.name,
      private: statistics.private_profile
    }));
  }, [statistics.private_profile]);

  useEffect(() => {
    let timeout: NodeJS.Timeout;
    if (imageLoadError) {
      timeout = setTimeout(() => {
        setImageLoadError(false);
      }, 5000);
    }

    return () => {
      clearTimeout(timeout);
    };
  }, [imageLoadError]);

  return (
    <Modal onDismiss={onDismiss} {...props}>
      <Container>
        <Text
          color="white"
          fontSize="20px"
          letterSpacing={'2px'}
          fontWeight={600}
          textTransform="uppercase"
        >
          Edit Profile
        </Text>

        <Box mt="30px" {...getRootProps()}>
          <input {...getInputProps()} />
          {image ? (
            <Avatar image={image.preview} size="145px" useProxy={false} />
          ) : (
            <Avatar image={avatar} size="145px" />
          )}
        </Box>

        <Box>
          <Flex mt="20px" justifyContent={'center'}>
            <Button
              borderRadius={'5px'}
              background="#242F42"
              fontWeight={600}
              color="#768BAD"
              px="30px"
              py="10px"
              onClick={open}
            >
              Upload Profile Picture
            </Button>

            <Button
              borderRadius={'5px'}
              background="#242F42"
              p="10px"
              ml="12px"
            >
              <TrashIcon />
            </Button>
          </Flex>

          <Text
            color={'#B2D1FF'}
            fontSize="15px"
            textAlign={'center'}
            mt="10px"
          >
            JPEG or PNG. Max size 800 KB.
          </Text>
        </Box>

        {imageLoadError && (
          <Text
            color="#a31f4e"
            fontSize="15px"
            fontWeight={600}
            textAlign="center"
            mt="5px"
          >
            Image file is too large.
          </Text>
        )}

        <Text
          color={'#768BAD'}
          fontSize="16px"
          fontWeight={500}
          mb={'8px'}
          mt="40px"
          width={['100%', '100%', '100%', 'auto']}
        >
          Username
        </Text>
        <InputBox
          gap={20}
          p="8px 20px"
          width={['100%', '100%', '100%', 'auto']}
        >
          <input
            name="name"
            placeholder="input name"
            value={formData.name}
            onChange={handleChange}
            maxLength={16}
          />
        </InputBox>

        <Box mt="30px">
          {/* <Box>
              <Checkbox
                name="private"
                value="private"
                label="Private profile"
              />
              <Span color={"#4F617B"}>(Hide statistics + your name)</Span>
            </Box> */}
          <Box mt="10px">
            <Checkbox
              name="hide"
              value="hide"
              label="Private Profile (username and stats will be hidden from other users.)"
              checked={formData.private}
              onChange={handlePrivateChange}
            />
          </Box>
        </Box>

        <Flex justifyContent={'center'} mt="20px">
          <Button
            px="35px"
            py="10px"
            background={'#1A5032'}
            color="#4FFF8B"
            fontSize="16px"
            fontWeight={600}
            width="100%"
            borderRadius={'5px'}
            onClick={request ? undefined : handleSaveProfile}
          >
            {request ? <ClipLoader size={20} color="#fff" /> : 'Save Profile'}
          </Button>
        </Flex>
      </Container>
    </Modal>
  );
};

const Container = styled(Flex)`
  flex-direction: column;
  flex: 1;
  background: linear-gradient(180deg, #132031 0%, #1a293c 100%);

  padding: 40px 21px;

  min-width: 350px;

  ${({ theme }) => theme.mediaQueries.md} {
    border: 2px solid #43546c;
    border-radius: 15px;

    padding: 40px 40px;
  }

  overflow: hidden auto;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    display: none;
  }
`;

export default EditProfileModal;
