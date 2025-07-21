import React, { useState } from 'react'
import UserService from "./api/userService"
import Common, { useRedirectToHomeIfAlreadyAuthenticated, CSS } from "./common"
import "./common.css"
import { Link, useHistory } from "react-router-dom";

interface FormData {
  email: string;
  name: string;
  password: string;
  confirmPassword: string;
}

interface Code {
  id: number;
  msg: string;
}

export default function Signup(): React.ReactElement | null {
  const [formData, setFormData] = useState<FormData>({ email: "", name: "", password: "", confirmPassword: "" })
  const [code, setCode] = useState<Code>({ id: 0, msg: "" })
  const history = useHistory()

  const { isCheckingAuth } = useRedirectToHomeIfAlreadyAuthenticated()

  function handleSubmit(e: React.FormEvent<HTMLFormElement>): void {
    e.preventDefault()

    UserService()
      .Signup(formData.email, formData.password, formData.name)
      .then((c: Code) => {
        if (c.id === 0) {
          history.push('/login')
        }
        setCode({ id: c.id, msg: c.msg })
      })
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement>): void {
    Common().SetFormData(e, setFormData)
  }

  // Don't render form while checking authentication
  if (isCheckingAuth) {
    return <div>Loading...</div>
  }

  return (
    <div className="container">
      <div className="row justify-content-center">
        <form className={"form needs-validation " + CSS.FormCol} onSubmit={handleSubmit}>
          <div className="form-floating">
            <input className="form-control" id="email" type="email" name="email" placeholder="Email"
              onChange={handleChange} value={formData.email} required />
            <label htmlFor="email">Email</label>
          </div>

          <div className="form-floating">
            <input className="form-control" id="name" type="text" name="name" placeholder="Name" value={formData.name}
              onChange={handleChange} required />
            <label htmlFor="name">Name</label>
          </div>

          <div className="form-floating">
            <input className="form-control" id="password" type="password" name="password" placeholder="Password"
              value={formData.password} onChange={handleChange} required />
            <label htmlFor="password">Password</label>
          </div>

          <div className="form-floating">
            <input className="form-control" id="confirmPassword" type="password" name="confirmPassword"
              placeholder="Confirm Password" value={formData.confirmPassword} onChange={handleChange} required />
            <label htmlFor="confirmPassword">Confirm Password</label>
          </div>

          <button className="btn btn-primary" type="submit">Create Account</button>
          {code.id !== 0 && <div className="bs-callout bs-callout-danger">{code.msg}</div>}
        </form>
      </div>

      <div className="row justify-content-center">
        <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12 ">
          <label htmlFor="login-link">Have an account?</label>
          <Link to="/login" id="login-link">Go to login</Link>
        </div>
      </div>
    </div>
  )
}