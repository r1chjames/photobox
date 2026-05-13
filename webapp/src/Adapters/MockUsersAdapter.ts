import {IUsersAdapter} from "./IUsersAdapter";
import {User} from "../Models/User";


export class MockUsersAdapter implements IUsersAdapter {

  public login = async (user: User) => {
    console.log("MockUsersAdapter.login: " + JSON.stringify(user));
    return { token: "v2.local.wf808vjXGvdPkKGDMi35vQaONl_Oc2Ep73alaX1UEHAzIDo2e40lzyvDmPL6R58tuqDX8-y0BmxPSN2c64Js-cPeff94H2AO6zMsewCp3BDHC_ALcrdiQ77A2clIgVwzU4YRH7HGxkrxKpQPCO8BgfiA9ApYxpUVl8zVwKVXzesS5LP1NWrpY5zcX61qgDeHlpJ-BUv48WccAAm41XdwD5cMmEERxWJ6fO53KGGOK0eWMLNsd6Q8vaURjNUA.bnVsbA" };
  }

  public register = async (user: User) => {
    console.log("MockUsersAdapter.register: " + JSON.stringify(user));
    return { token: "mock-registration-token" };
  }

  public getAllUsers = async (): Promise<User[]> => {
    return [];
  }

  public updateUser = async (_userId: string, updates: Partial<User>): Promise<User> => {
    return { username: 'mock', email: 'mock@example.com', ...updates };
  }

  public deleteUser = async (_userId: string): Promise<void> => {
    return;
  }
}
