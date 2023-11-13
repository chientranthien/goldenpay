import React, {useState} from "react"
import "./Login.css"
import UserService from "./api/userService.js"
import Common from "./common"
import {Link} from "react-router-dom";

export default function Login() {
  const [formData, setFormData] = useState({email: "", password: ""})

  function handleSubmit(e) {
    e.preventDefault()
    UserService().Login(formData.email, formData.password)
  }

  function handleChange(e) {
    Common().SetFormData(e, setFormData)
  }

  return (
    <div className="container form-container">
      <div className="row justify-content-center">
        <form className="form needs-validation col-lg-6 col-md-8 col-sm-12  col-xs-12" onSubmit={handleSubmit}>
          <div className="form-floating">
            <input className="form-control" name="email" type="email" placeholder="Email" onChange={handleChange}
                   value={formData.email}/>
            <label htmlFor="email">Email</label>
          </div>
          <div className="form-floating">
            <input className="form-control" name="password" type="password" placeholder="Password"
                   onChange={handleChange} value={formData.password}/>
            <label htmlFor="password">Password</label>
          </div>
          <button className="btn btn-primary" type="submit">Login</button>
        </form>

      </div>

      <div className="row justify-content-center">
        <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12 ">
          <label htmlFor="signup-link">Don't have an account yet?</label>
          <Link to="/signup" id="signup-link">Signup</Link>
        </div>
      </div>
    </div>
  )
}