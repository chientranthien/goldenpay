import { config, ToJSON } from './common'

// API Response interfaces
interface ApiResponse {
  code: Code;
  data?: any;
  message?: string;
}

interface Code {
  id: number;
  msg: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

interface SignupRequest {
  email: string;
  password: string;
  name: string;
}

interface IUserService {
  Login: (email: string, password: string) => Promise<Code>;
  Authz: () => Promise<Code>;
  Signup: (email: string, password: string, name: string) => Promise<Code>;
}

const userService: IUserService = {
  Login: async function (email: string, password: string): Promise<Code> {
    const api = config.host + "login"
    const body: LoginRequest = {
      email: email,
      password: password
    }

    let code: Code = { id: 0, msg: "" }
    await fetch(api, {
      method: 'POST',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
      body: ToJSON(body)
    }
    ).then(resp => {
      return resp.json()
    }).then((json: ApiResponse) => {
      code = json.code
    }).catch(() => {
      console.log("failed to login")
    })

    return code
  },

  Authz: async function (): Promise<Code> {
    const api = config.host + "authz"
    const body = {}

    let code: Code = { id: 0, msg: "" }
    await fetch(api, {
      method: 'POST',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
      body: ToJSON(body)
    }
    ).then(resp => {
      return resp.json()
    }).then((json: ApiResponse) => {
      code = json.code
    }).catch(() => {
      console.log("failed to authz")
    })

    return code
  },

  Signup: async function (email: string, password: string, name: string): Promise<Code> {
    const api = config.host + "signup"
    const body: SignupRequest = {
      email: email,
      password: password,
      name: name
    }

    let code: Code = { id: 0, msg: "" }
    await fetch(api, {
      method: 'POST',
      headers: {
        "Content-Type": "application/json"
      },
      body: ToJSON(body)
    }
    ).then(resp => resp.json())
      .then((json: ApiResponse) => {
        code = json.code
      }).catch(() => {
        console.error("failed to signup")
      })

    return code
  }
}

export default function UserService(): IUserService {
  return userService
}
