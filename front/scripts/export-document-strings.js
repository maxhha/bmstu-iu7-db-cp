/* Export *DocumentString constants */
import { readFileSync, writeFileSync } from "fs"

const file = process.argv[2]
let data = readFileSync(file, "utf8")

data = data.replace(
    /const (.*?)DocumentString =/g,
    "export const $1DocumentString ="
)

writeFileSync(file, data)
