const config = {
  host: "http://localhost:5000/api/v1/"
}

const userService = {}
userService.Login = function (username, password) {
  const api = config.host + "login"
  const body = {
    username: username,
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

userService.Signup = function (email, password, user) {
  const api = config.host + "login"
  const body = {
    username: username,
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
export default function UserService() {
  return userService
}
