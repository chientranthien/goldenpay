import React, {useState} from 'react'
import UserService from "./api/userService.js"
import Common from "./common"
import "./common.css"
import {Link} from "react-router-dom";

export default function Signup() {
  const [formData, setFormData] = useState({email: "", name: "", password: "", confirmPassword: ""})

  function handleSubmit(e) {
    e.preventDefault()
    UserService().Login(formData.email, formData.password)
  }

  function handleChange(e) {
    Common().SetFormData(e, setFormData)
  }

  function validateEmail(e, email) {
    if (e.target.name === "email") {
      return ""
    }

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(email)) {
      return "invalid email"
    }

    return ""
  }

  function validatePassword(e, password, confirmPassword) {
    if (e.target.name === "password") {
      return ""
    }

    if (password === "") {
      return ""
    }

    if (password.length <= 8) {
      return "password must have more than 7 characters"
    }

    if (confirmPassword === "") {
      return ""
    }

    if (password !== confirmPassword) {
      return "confirm password doesn't match"
    }

    return ""
  }

  return (
    <div className="container form-container">
      <div className="row justify-content-center">
        <form onSubmit={handleSubmit} className="form col-6">
          <input className="form-control" type="email" name="email" placeholder="Email" onChange={handleChange}
                 value={formData.email}/>
          <input className="form-control" type="text" name="name" placeholder="Name" value={formData.name}
                 onChange={handleChange}/>
          <input className="form-control" type="password" name="password" placeholder="Password"
                 value={formData.password}
                 onChange={handleChange}/>
          <input className="form-control" type="password" name="confirmPassword" placeholder="Confirm Password"
                 value={formData.confirmPassword}
                 onChange={handleChange}/>

          <button className="btn btn-primary bold" type="submit">Create Account</button>
        </form>

      </div>
      <div className="row justify-content-center">
        <div className="info col-6">
          <label htmlFor="login-link">Have an Account?</label>
          <Link to="/login" id="login-link">Login</Link>
        </div>
      </div>

    </div>
  )
}