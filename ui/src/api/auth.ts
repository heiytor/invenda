import { AuthApi, AuthUserRequest } from "@/lib/client";
import { CONFIG } from "./config";

let api = new AuthApi(CONFIG)

export async function authUser({ identifier, password }: AuthUserRequest) {
    try {
        await api.authUser({ identifier, password });
    } catch(error: any) {
        return;
    }
}
