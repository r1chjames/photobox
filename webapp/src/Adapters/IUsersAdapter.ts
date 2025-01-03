import {User} from "../Models/User";

export interface IUsersAdapter {

  login(user: User): Promise<string>;
  register(user: User): Promise<any>;
}
