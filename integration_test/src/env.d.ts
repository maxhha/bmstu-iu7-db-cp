declare global {
  namespace NodeJS {
    interface ProcessEnv {
      GRAPHQL_URL: string
      REDIS_CONNECTION_URL: string
    }
  }
}
