import { GraphQLClient } from "graphql-request"
import { getSdk } from "./query"

export function getClient() {
  if (process.env.GRAPHQL_URL === undefined) {
    throw new Error("GRAPHQL_URL environment variable is undefined")
  }

  const headers: Record<string, string> = {}
  const client = new GraphQLClient(process.env.GRAPHQL_URL, {
    headers,
  })

  const sdk = getSdk(client)

  return {
    ...sdk,
    unwrap() {
      return client
    },
    setToken(token: string | undefined) {
      if (token === undefined) {
        delete headers["Authorization"]
      } else {
        headers["Authorization"] = token
      }
    },
  }
}
