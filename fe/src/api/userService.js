import React from 'react'
import {config, ToJSON} from './common'

const userService = {}
userService.Login = async function (email, password) {
  const api = config.host + "login"
  const body = {
    email: email,
    password: password
  }

  let c = {}
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
  }).then(json => {
    c = json.code
  }).catch(() => {
    console.log("failed to login")
  })

  return c
}

userService.Authz = async function () {
  const api = config.host + "authz"
  const body = {}

  let c = {}
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
  }).then(json => {
    c = json.code
  }).catch(() => {
    console.log("failed to authz")
  })

  return c
}


userService.Signup = async function (email, password, name) {
  const api = config.host + "signup"
  const body = {
    email: email,
    password: password,
    name: name
  }

  let code = {}
  await fetch(api, {
      method: 'POST',
      headers: {
        "Content-Type": "application/json"
      },
      body: ToJSON(body)
    }
  ).then(resp => resp.json())
    .then(json => {
      code = json.code
    }).catch(() => {
      console.error("failed to signup")
    })

  return code
}

export default function UserService() {
  return userService
}
