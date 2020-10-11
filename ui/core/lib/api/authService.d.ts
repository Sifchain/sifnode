import { AuthAccountInfo } from "../entities/Auth";
export declare const authService: {
    getAccount(address: string): Promise<AuthAccountInfo>;
};
