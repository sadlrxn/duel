import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { withRouter } from 'storybook-addon-react-router-v6';

import Sidebar from '.';

export default {
  title: 'Components/Sidebar',
  component: Sidebar,
  argTypes: {},
  decorators: [withRouter]
} as ComponentMeta<typeof Sidebar>;

export const Default: ComponentStory<typeof Sidebar> = args => (
  <Sidebar {...args} />
);

Default.args = {
  login: async () => {},
  logout: async () => {}
};
