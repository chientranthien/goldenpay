import React, {useState} from "react"
import "./Login.css"
import UserService from "./api/userService.js"
import Common from "./common"

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
    <div className="container login-page">
      <div className="row justify-content-end">
        <form onSubmit={handleSubmit}>
          <input className="form-control" name="email" type="email" placeholder="Email" onChange={handleChange}
                 value={formData.email}/>
          <input className="form-control" name="password" type="password" placeholder="Password" onChange={handleChange}
                 value={formData.password}/>
          <button className="btn" type="submit">Login</button>
        </form>
      </div>
    </div>
  )
}