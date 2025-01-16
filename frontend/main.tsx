import * as vlens from "vlens";
import * as server from "@app/server";

// route will have a fetch and view function
// no state
async function main() {
    vlens.initRoutes([
        vlens.routeHandler("/users", () => import("@app/users")),
        vlens.routeHandler("/", () => import("@app/home")),
    ]);
}

main();

(window as any).server = server;