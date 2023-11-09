import React, {useState} from 'react'

export default function Signup() {
  const [formData, setFromData] = useState({username: "", password: "", confirmPassword: ""})

  function handleOnChange(e) {
    setFromData(prev => {
      return {
        ...prev,
        [e.target.name]: e.target.value
      }
    })
  }

  return (
    <div className="container">
      <div className="row">
        <form>
          <input type="username" placeholder="Username" value={formData.username} onChange={handleOnChange}/>
          <input type="password" placeholder="Password" value={formData.password} onChange={handleOnChange}/>
          <input type="confirmPassword" laceholder="Confirm Password" value={formData.confirmPassword}
                 onChange={handleOnChange}/>
          <button className="btn" type="submit">Create Account</button>
        </form>
      </div>
    </div>
  )
}