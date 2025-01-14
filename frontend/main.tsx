import * as vlens from "vlens";

// route will have a fetch and view function
async function main() {
    vlens.initRoutes([
        vlens.routeHandler("/", () => import("@app/home")),
    ]);
}

main();