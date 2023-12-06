import React, {useState} from "react";
import Common, {CSS} from "./common"
import WalletService from "./api/walletService";
import {useHistory} from "react-router-dom";
import Nav from "./Nav";

export default function Topup() {
  const [formData, setFormData] = useState({amount: 0})
  const [code, setCode] = useState({id: 0, msg: ""})
  const history = useHistory()

  function handleChange(e) {
    Common().SetFormData(e, setFormData)
  }

  function handleSubmit(e) {
    e.preventDefault()
    WalletService().Topup(formData.amount).then(
      c => {
        if (c.id == 0) {
          history.push('/')
        }
        setCode(prev => {
          return {...c}
        })
      }
    )
  }

  return (
    <>
      <Nav></Nav>
      <div className="container">
        <div className="row justify-content-center">
          <form className={"form needs-validation" + CSS.FormCol} onSubmit={handleSubmit}>
            <div className="form-floating">
              <input className="form-control" id="amount" name="amount" type="number" onChange={handleChange}
                     placeholder={""}
                     value={formData.amount} required/>
              <label htmlFor="amount">Amount</label>
            </div>

            <button className="btn btn-primary" type="submit">Topup</button>
            {code.id != 0 && <div className="bs-callout bs-callout-danger">{code.msg}</div>}
          </form>
        </div>
      </div>
    </>
  )
}