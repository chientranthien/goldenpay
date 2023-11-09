import React, {useState} from "react";
import "./Login.css"
import UserService from "./api/userService.js"

export default function Login() {
  const [formData, setFormData] = useState({username: "", password: ""})

  function handleSubmit(e) {
    e.preventDefault()
    UserService().Login(formData.username, formData.password)
  }

  function handleOnChange(e) {
    setFormData(prev => {
      return {
        ...prev,
        [e.target.name]: e.target.value
      }
    })
  }

  return (
    <div className="container login-page">
      <div className="row justify-content-end">
        <form onSubmit={handleSubmit}>
          <input className="form-control" name="username" type="username" placeholder="Username"
                 onChange={handleOnChange} value={formData.username}/>
          <input className="form-control" name="password" type="password" placeholder="Password"
                 onChange={handleOnChange} value={formData.password}/>
          <button className="btn" type="submit">Login</button>
        </form>
      </div>
    </div>
  )
}