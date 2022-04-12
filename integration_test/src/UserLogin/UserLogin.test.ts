import { UserLogin } from "./UserLogin"

describe("UserLogin", () => {
  const test = new UserLogin()
  beforeAll(() => test.before())
  afterAll(() => test.after())

  test.run()
})
