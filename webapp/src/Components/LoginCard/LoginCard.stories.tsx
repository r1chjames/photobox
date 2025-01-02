import type {Meta, StoryObj} from '@storybook/react';
import {LoginCard} from "./LoginCard";
import {MockUsersAdapter} from "../../Adapters/MockUsersAdapter";

const meta: Meta<typeof LoginCard> = {
  component: LoginCard,
};

export default meta;
type Story = StoryObj<typeof LoginCard>;

export const Default: Story = {
  args: {
      usersAdapter: new MockUsersAdapter(),
  }
};
