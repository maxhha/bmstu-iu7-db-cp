/* Fix type-only imports */
import { readFileSync, writeFileSync } from "fs"

const file = process.argv[2]
let data = readFileSync(file, "utf8")

data = data
    .replace(
        'import { GraphQLClient } from "graphql-request"',
        'import type { GraphQLClient } from "graphql-request"'
    )
    .replace(
        'import * as Dom from "graphql-request/dist/types.dom"',
        'import type * as Dom from "graphql-request/dist/types.dom"'
    )
    .replace(
        'import { GraphQLError } from "graphql-request/dist/types"',
        'import type { GraphQLError } from "graphql-request/dist/types"'
    )

writeFileSync(file, data)
