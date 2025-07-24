import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Nav from "./Nav";
import { useRedirectToLoginIfNotAuthenticated } from "./common"
import RecentActivity from "./RecentActivity";
import WalletService from "./api/walletService";

interface WalletData {
  balance: number;
}

interface WalletResponse {
  data?: WalletData;
}

export default function Home(): React.ReactElement | null {
  const [walletResp, updateWalletResp] = useState<WalletResponse>({})
  const { isAuthenticated } = useRedirectToLoginIfNotAuthenticated()

  useEffect(() => {
    // Only make API call if authenticated
    if (isAuthenticated === true) {
      WalletService().GetUserWallet().then((resp: any) => {
        updateWalletResp(resp)
      }
      )
    }
  }, [isAuthenticated])

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
        {
          walletResp.data !== undefined &&
          <div className="row justify-content-center">

            <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12">
              <h5>Balance</h5>
              <h2>${walletResp.data.balance}</h2>
              <Link className="btn btn-primary" to={"/transfer"}>Transfer</Link>
              <Link className="btn btn-light" to={"/topup"}>Topup</Link>
              <Link className="btn btn-success" to={"/chat"}>Chat</Link>
            </div>

          </div>
        }

        <RecentActivity isAuthenticated={isAuthenticated}></RecentActivity>
      </div>
    </>
  )
}