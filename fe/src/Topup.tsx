import React, { useState } from "react";
import Common, { CSS, useRedirectToLoginIfNotAuthenticated } from "./common"
import WalletService from "./api/walletService";
import { useHistory } from "react-router-dom";
import Nav from "./Nav";

interface TopupFormData {
  amount: number;
}

interface Code {
  id: number;
  msg: string;
}

export default function Topup(): React.ReactElement | null {
  const [formData, setFormData] = useState<TopupFormData>({ amount: 0 })
  const [code, setCode] = useState<Code>({ id: 0, msg: "" })
  const history = useHistory()

  const { isAuthenticated } = useRedirectToLoginIfNotAuthenticated()

  function handleChange(e: React.ChangeEvent<HTMLInputElement>): void {
    Common().SetFormData(e, setFormData)
  }

  function handleSubmit(e: React.FormEvent<HTMLFormElement>): void {
    e.preventDefault()
    WalletService().Topup(formData.amount.toString()).then(
      (c: Code) => {
        if (c.id === 0) {
          history.push('/')
        }
        setCode({ id: c.id, msg: c.msg })
      }
    )
  }

  // Don't render anything until authentication status is determined
  if (isAuthenticated === null) {
    return <div>Loading...</div>
  }

  // If not authenticated, this component shouldn't render (will redirect)
  if (isAuthenticated === false) {
    return null
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
                value={formData.amount} required />
              <label htmlFor="amount">Amount</label>
            </div>

            <button className="btn btn-primary" type="submit">Topup</button>
            {code.id !== 0 && <div className="bs-callout bs-callout-danger">{code.msg}</div>}
          </form>
        </div>
      </div>
    </>
  )
}