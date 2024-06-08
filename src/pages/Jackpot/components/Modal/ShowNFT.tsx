import React from "react";
import styled from "styled-components";

import { Modal, ModalProps, Box, NftBox, Flex, Text } from "components";
import NFTCard from "components/NFTCard";
import { NFT } from "api/types/nft";

const StyledText = styled(Text)`
  font-size: 16px;
  font-weight: 600;
  color: #0b141e;
  padding: 3px 8px;
  border-radius: 17px;
  width: max-content;
`;

const UserName = styled(Flex)`
  background: rgba(76, 105, 137, 0.2);
  border-radius: 8px;
  color: #b2d1ff;
  font-weight: 500;
  font-size: 14px;
  gap: 9px;
  align-items: center;
  padding: 7px 17px 8px 13px;
`;

interface ShowNFTProps extends ModalProps {
  name?: string;
  level?: number;
  nfts?: NFT[];
}

const ShowNFT: React.FC<ShowNFTProps> = ({
  name = "",
  nfts = [],
  ...props
}) => {
  return (
    <Modal {...props}>
      <Box
        background={"linear-gradient(180deg, #202F44 0%, #1B283A 100%)"}
        borderRadius="17px"
        maxWidth="80vw"
        width="1000px"
        px={"30px"}
        pt="30px"
        pb="15px"
      >
        <Flex mb="30px" gap={10} alignItems="center" px="20px">
          <StyledText backgroundColor={name === "" ? "#6D81A2" : "success"}>
            {nfts.length}
          </StyledText>
          <Text fontSize="20px" fontWeight={600} color="white">
            {name === "" ? "NFT IN THE JACKPOT" : "NFT BET BY "}
          </Text>
          {name !== "" && <UserName>{name}</UserName>}
        </Flex>
        <NftBox maxHeight="60vh">
          {nfts.map((nft) => {
            return (
              <NFTCard
                key={nft.mintAddress}
                price={nft.price}
                collectionName={nft.collectionName ?? ""}
                name={nft.name ?? ""}
                image={nft.image}
              />
            );
          })}
        </NftBox>
      </Box>
    </Modal>
  );
};

export default React.memo(ShowNFT);
