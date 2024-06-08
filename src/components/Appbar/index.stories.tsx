import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { withRouter } from 'storybook-addon-react-router-v6';
import Provider from '../../Providers';
import Appbar from '.';

const withProvider = (story: any) => <Provider>{story()}</Provider>;

export default {
  title: 'Components/Appbar',
  component: Appbar,
  argTypes: {},
  decorators: [withRouter, withProvider]
} as ComponentMeta<typeof Appbar>;

export const Default: ComponentStory<typeof Appbar> = args => (
  <Appbar {...args} />
);

Default.args = {
  login: async () => {},
  logout: async () => {}
};
