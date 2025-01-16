import * as rpc from "vlens/rpc"

export interface AddUserRequest {
    Username: string
}

export interface UserListResponse {
    AllUsernames: string[]
}

export async function AddUser(data: AddUserRequest): Promise<rpc.Response<UserListResponse>> {
    return await rpc.call<UserListResponse>('AddUser', JSON.stringify(data));
}

