import * as preact from "preact"
import * as vlens from "vlens"
import * as events from "vlens/events"
import * as server from "@app/server"

// stable state across rendering cycles
type Form = {
    // all names
    data: server.UserListResponse
    // latest name
    name: string
    error: string
} 

// outputs stable hook to input
const useForm = vlens.declareHook((data: server.UserListResponse): Form => ({
    data, name: "", error: ""
}))

async function onAddUserClicked(form: Form) {
    let [r, e] = await server.AddUser({Username: form.name})
    if (r) {
        form.name = ""
        form.data = r
        form.error = ""
    } else {
        form.error = e
    }
    vlens.scheduleRedraw()
}

export async function fetch(route: string, prefix: string) {
    return server.ListUsers({})
}

// this function called on redraws. therefore, cannot have closures if want stability
// data is immutable from first fetch(), so can be used for cache key
// TODO: input sanitisation, e.g. no empty, only alphanumeric etc.
export function view(route: string, prefix: string, data: server.UserListResponse) : preact.ComponentChild {
    let form = useForm(data)
    return <div>
        <h3> Users </h3>
        {form.data.AllUsernames.map(name => <div key={name}>{name}</div>)}
        <h3> Add User </h3>
        <input type="text" {...events.inputAttrs(vlens.ref(form, "name"))} />
        <button onClick={vlens.cachePartial(onAddUserClicked, form)}>Add</button>
        {form.name && <div> You're inputting: <code>{form.name}</code></div>}
    </div>
}
