const config = {
  host: "http://localhost:5000/api/v1/"
}

const userService = {}
userService.Login = function (email, password) {
  const api = config.host + "login"
  const body = {
    email: email,
    password: password
  }

  fetch(api, {
      method: 'POST',
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(body)
    }
  ).then(() => {
    console.log("aaaaaaaaaaaaaaaaaa")
  }).catch(() => {
    console.log("bbbbbbbbbbbbbbbb")
  })
}

userService.Signup = function (email, password, name) {
  const api = config.host + "signup"
  const body = {
    email: email,
    password: password,
    name: name
  }

  fetch(api, {
      method: 'POST',
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(body)
    }
  ).then(() => {
    console.log("aaaaaaaaaaaaaaaaaa")
  }).catch(() => {
    console.log("bbbbbbbbbbbbbbbb")
  })

}
export default function UserService() {
  return userService
}
