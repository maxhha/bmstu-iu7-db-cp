import { Register } from "./Register"

describe("Register", () => {
  const test = new Register()
  beforeAll(() => test.before())
  afterAll(() => test.after())

  test.run()
})
