import React, {useState} from 'react'
import UserService from "./api/userService.js"
import Common, {RedirectToHomeIfAlreadyAuthenticated} from "./common"
import "./common.css"
import {Link, useHistory} from "react-router-dom";

const config = {
  host: "http://localhost:5000/api/v1/"
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

export default function Signup() {
  const [formData, setFormData] = useState({email: "", name: "", password: "", confirmPassword: ""})
  const [code, setCode] = useState({id: 0, msg: ""})
  const history = useHistory()

  RedirectToHomeIfAlreadyAuthenticated()

  function handleSubmit(e) {
    e.preventDefault()

    UserService()
      .Signup(formData.email, formData.password, formData.name)
      .then(c => {
        if (c.id == 0) {
          history.push('/login')
        }
        setCode(prev => {
          return {...c}
        })
      })
  }

  function handleChange(e) {
    Common().SetFormData(e, setFormData)
  }

  return (

    <div className="container">
      <div className="row justify-content-center">
        <form onSubmit={handleSubmit} className="form needs-validation col-lg-6 col-md-8 col-sm-12 col-xs-12">
          <div className="form-floating">
            <input className="form-control" id="email" type="email" name="email" placeholder="Email"
                   onChange={handleChange} value={formData.email} required/>
            <label htmlFor="email">Email</label>
          </div>

          <div className="form-floating">
            <input className="form-control" id="name" type="text" name="name" placeholder="Name" value={formData.name}
                   onChange={handleChange} required/>
            <label htmlFor="name">Name</label>
          </div>
          <div className="form-floating">
            <input className="form-control" id="password" type="password" name="password" placeholder="Password"
                   value={formData.password} onChange={handleChange} required/>
            <label htmlFor="password">Password</label>
          </div>
          <div className="form-floating">
            <input className="form-control" id="confirmPassword" type="password" name="confirmPassword"
                   placeholder="Confirm Password" value={formData.confirmPassword} onChange={handleChange} required/>
            <label htmlFor="confirmPassword">Confirm Password</label>
          </div>

          <button className="btn btn-primary bold" type="submit">Create Account</button>
          {code.id != 0 && <div className="bs-callout bs-callout-danger">{code.msg}</div>}
        </form>

      </div>


      <div className="row justify-content-center">
        <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12 ">
          <label htmlFor="login-link">Have an Account?</label>
          <Link to="/login" id="login-link">Login</Link>
        </div>
      </div>

    </div>
  )
}