import * as preact from "preact";
import * as rpc from "vlens/rpc";

type Data = {};

export async function fetch(route: string, prefix: string) {
    // has reference ID which is stable
    return rpc.ok<Data>({});
}

export function view(route: string, prefix: string, data: Data): preact.ComponentChild {
    return <div>
        <h2> Handcrafted Forum !!!</h2>
        <img src="/images/xav-and-i.jpg" />
    </div>
}