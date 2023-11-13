import React, {useState} from 'react'
import UserService from "./api/userService.js"
import Common from "./common"
import "./common.css"
import {Link} from "react-router-dom";

export default function Signup() {
  const [formData, setFormData] = useState({email: "", name: "", password: "", confirmPassword: ""})

  function handleSubmit(e) {
    e.preventDefault()

    UserService().Signup(formData.email, formData.password)
  }

  function handleChange(e) {
    Common().SetFormData(e, setFormData)
  }

  // TODO(tom): use validation
  function validateEmail(e, email) {
    if (e.target.name === "email") {
      return ""
    }

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(email)) {
      console.log("invalid email " + email)
      return "invalid email"
    }

    return ""
  }

  // TODO(tom): use validation
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
        <form onSubmit={handleSubmit} className="form col-6 needs-validation">
          <div class="form-floating mb-3">
            <input className="form-control" id="email" type="email" name="email" placeholder="Email"
                   onChange={handleChange} value={formData.email} required/>
            <label for="email">Email</label>
          </div>

          <div class="form-floating">
            <input className="form-control" id="name" type="text" name="name" placeholder="Name" value={formData.name}
                   onChange={handleChange} required/>
            <label for="name">Name</label>
          </div>
          <div class="form-floating">
            <input className="form-control" id="password" type="password" name="password" placeholder="Password"
                   value={formData.password} onChange={handleChange} required/>
            <label for="password">Password</label>
          </div>
          <div class="form-floating">
            <input className="form-control" id="confirmPassword" type="password" name="confirmPassword"
                   placeholder="Confirm Password" value={formData.confirmPassword} onChange={handleChange} required/>
            <label for="confirmPassword">Confirm Password</label>
          </div>

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