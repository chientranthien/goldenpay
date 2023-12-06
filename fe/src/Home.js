import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";
import Nav from "./Nav";
import {RedirectToLoginIfNotAuthenticated} from "./common"
import RecentActivity from "./RecentActivity";
import WalletService from "./api/walletService";

export default function Home() {
  const [walletResp, updateWalletResp] = useState({})

  RedirectToLoginIfNotAuthenticated()

  useEffect(() => {
    WalletService().GetUserWallet().then(resp => {
        updateWalletResp(prev => resp)
      }
    )
  }, [])

  return (
    <>
      <Nav></Nav>

      <div className="container">
        {
          walletResp.data != undefined &&
          <div className="row justify-content-center">

            <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12">
              <h5>Balance</h5>
              <h2>${walletResp.data.balance}</h2>
              <Link className="btn btn-primary" to={"/transfer"}>Transfer</Link>
              <Link className="btn btn-light" to={"/Topup"}>Topup</Link>
            </div>

          </div>
        }

        <RecentActivity></RecentActivity>
      </div>
    </>
  )
}