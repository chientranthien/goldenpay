import React, {useState} from "react";
import {CSS} from "./common"

export default function Transfer() {
  const [formData, setFormDate] = useState({email: "", amount: 0})
  const [code, setCode] = useState({id: 0, msg: ""})

  function handleChange(e) {

  }

  return (
    <div className="container form-container">
      <div className="row justify-content-center">
        <form className={"form needs-validation" + CSS.FormCol}>
          <div className="form-floating">
            <input className="form-control" id="email" name="email" type="email" onChange={handleChange}
                   value={formData.email} required/>
            <label htmlFor="email">To Email</label>
          </div>
          <div className="form-floating">
            <input className="form-control" id="amount" name="amount" type="number" onChange={handleChange}
                   value={formData.amount} required/>
            <label htmlFor="amount">Amount</label>
          </div>

          <button className="btn btn-primary" type="submit">Transfer</button>
          {code.id != 0 && <div className="bs-callout bs-callout-danger">{code.msg}</div>}
        </form>
      </div>
    </div>
  )
}