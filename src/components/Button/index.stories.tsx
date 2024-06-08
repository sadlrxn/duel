import React from "react";
import { ComponentMeta } from "@storybook/react";

import { Text } from "components/Text";
import logout from "assets/imgs/icons/logout.svg";
import lock from "assets/imgs/icons/lock.svg";

import Button from "./Button";

export default {
  title: "Components/Button",
  component: Button,
  argTypes: {},
} as unknown as ComponentMeta<typeof Button>;

export const Buttons = () => {
  return (
    <>
      <Button>
        <Text fontSize="16px" mx="26px" my="17px">
          Create Game
        </Text>
      </Button>
      <Button
        variant="secondary"
        scale="xs"
        color="text"
        backgroundColor="#242F42"
      >
        <Text fontSize="14px" mx="23px">
          Withdraw
        </Text>
      </Button>
      <Button
        variant="secondary"
        scale="xs"
        color="success"
        backgroundColor="#1A5032"
      >
        <Text fontSize="14px" mx="23px">
          Deposit
        </Text>
      </Button>
      <Button
        variant="secondary"
        outlined
        scale="xs"
        backgroundColor="transparent"
        color="success"
        borderColor="success"
      >
        <Text fontSize="14px" mx="23px">
          Join
        </Text>
      </Button>
      <Button
        variant="secondary"
        outlined
        scale="lg"
        width="52px"
        backgroundColor="transparent"
        borderColor="#374355"
      >
        <img src={lock} alt="" width={18} height={23} />
      </Button>
      <Button variant="secondary" scale="xs" width="38px">
        <img src={logout} alt="" width={13} height={13} />
      </Button>
    </>
  );
};
