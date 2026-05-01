import {User} from "../Models/User";
import {Token} from "./UsersAdapter";

export interface IUsersAdapter {

  login(user: User): Promise<Token>;
  register(user: User): Promise<Token>;
  getAllUsers(): Promise<User[]>;
  updateUser(userId: string, updates: Partial<User>): Promise<User>;
  deleteUser(userId: string): Promise<void>;
}
