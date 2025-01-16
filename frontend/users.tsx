import * as preact from "preact";
import * as server from "@app/server";

export async function fetch(route: string, prefix: string) {
    return server.ListUsers({})
}

export async function view(route: string, prefix: string, data: server.UserListResponse) : preact.ComponentChild {
    return <div>
        <h3> Users </h3>
        {data.AllUsernames.map(name => <div key={name}>{name}</div>)}
    </div>
}
