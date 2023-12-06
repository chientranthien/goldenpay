import React, {useState} from "react"
import "./Login.css"
import Common, {CSS, RedirectToHomeIfAlreadyAuthenticated} from "./common"
import {Link, useHistory} from "react-router-dom";
import UserService from "./api/userService";


export default function Login() {
  const [formData, setFormData] = useState({email: "", password: ""})
  const [code, setCode] = useState({id: 0, msg: ""})
  const history = useHistory()

  RedirectToHomeIfAlreadyAuthenticated()

  function handleSubmit(e) {
    e.preventDefault();
    UserService()
      .Login(formData.email, formData.password)
      .then(c => {
        if (c.id == 0) {
          history.push('/')
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
        <form className={"form needs-validation " + CSS.FormCol} onSubmit={handleSubmit}>
          <div className="form-floating">
            <input className="form-control" id="email" name="email" type="email" onChange={handleChange}
                   placeholder={""} value={formData.email} required/>
            <label htmlFor="email">Email</label>
          </div>
          <div className="form-floating">
            <input className="form-control" name="password" type="password" onChange={handleChange} placeholder={""}
                   value={formData.password} required/>
            <label htmlFor="password">Password</label>
          </div>
          <button className="btn btn-primary" type="submit">Login</button>
          {code.id != 0 && <div className="bs-callout bs-callout-danger">{code.msg}</div>}
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