import React, { useState } from "react"
import "./Login.css"
import Common, { CSS, useRedirectToHomeIfAlreadyAuthenticated } from "./common"
import { Link, useHistory } from "react-router-dom";
import UserService from "./api/userService";

interface LoginFormData {
  email: string;
  password: string;
}

interface Code {
  id: number;
  msg: string;
}

export default function Login(): React.ReactElement | null {
  const [formData, setFormData] = useState<LoginFormData>({ email: "", password: "" })
  const [code, setCode] = useState<Code>({ id: 0, msg: "" })
  const history = useHistory()

  const { isCheckingAuth } = useRedirectToHomeIfAlreadyAuthenticated()

  function handleSubmit(e: React.FormEvent<HTMLFormElement>): void {
    e.preventDefault();
    console.log("🚀 Starting login process...")
    console.log("📧 Email:", formData.email)

    UserService()
      .Login(formData.email, formData.password)
      .then((c: Code) => {
        console.log("🔐 Login API response code:", c)
        if (c.id === 0) {
          console.log("✅ Login successful, waiting briefly for cookies to be set...")
          // Small delay to ensure cookies are set before redirect
          setTimeout(() => {
            console.log("🏠 Redirecting to home...")
            history.push('/')
          }, 100);
        } else {
          console.log("❌ Login failed with code:", c)
        }
        setCode({ id: c.id, msg: c.msg })
      })
      .catch((error) => {
        console.log("❌ Login API error:", error)
        setCode({ id: -1, msg: "Network Error" })
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
            <input className="form-control" id="email" name="email" type="email" onChange={handleChange}
              placeholder={""} value={formData.email} required />
            <label htmlFor="email">Email</label>
          </div>
          <div className="form-floating">
            <input className="form-control" name="password" type="password" onChange={handleChange} placeholder={""}
              value={formData.password} required />
            <label htmlFor="password">Password</label>
          </div>
          <button className="btn btn-primary" type="submit">Login</button>
          {code.id !== 0 && <div className="bs-callout bs-callout-danger">{code.msg}</div>}
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