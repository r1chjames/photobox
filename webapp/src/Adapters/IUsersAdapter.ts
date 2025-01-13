import {User} from "../Models/User";
import {Token} from "./UsersAdapter";

export interface IUsersAdapter {

  login(user: User): Promise<Token>;
  register(user: User): Promise<User>;
}
