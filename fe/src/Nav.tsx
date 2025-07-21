import React from "react";
import logo from "./logo.png"
import { RemoveTokenCookie } from "./common";
import { useHistory } from "react-router-dom";

export default function Nav(): React.ReactElement {
  const history = useHistory()

  function HandleLogout(e: React.MouseEvent<HTMLButtonElement>): void {
    RemoveTokenCookie()
    history.push("/login")
  }

  return (
    <nav className="navbar navbar-dark bg-primary">

      <div className="container-fluid">
        <a className="navbar-brand" href="/">

          <img src={logo} alt="Logo" width="30" height="30"
            className="d-inline-block align-text-center" />
          GoldenPay
        </a>

        <div>
          <button className={"btn btn-primary"} onClick={HandleLogout}>
            Logout
          </button>
        </div>
      </div>
    </nav>
  )
}