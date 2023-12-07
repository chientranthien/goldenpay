import {config, FromJSON, ToJSON} from "./common"

const walletService = {}

walletService.Transfer = async function (toEmail, amount) {
  const api = config.host + "users/transactions"
  const body = {
    toEmail: toEmail,
    amount: parseInt(amount)
  }

  let c = {}
  await fetch(api, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
      body: ToJSON(body)
    }
  ).then(resp => {
    return resp.json()
  }).then(json => {
    c = json.code
  }).catch(() => {
    console.log("failed to transfer")
  })

  return c

}

walletService.Topup = async function (amount) {
  const api = config.host + "users/topups"
  const body = {
    amount: parseInt(amount)
  }

  let c = {}
  await fetch(api, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
      body: ToJSON(body)
    }
  ).then(resp => {
    return resp.json()
  }).then(json => {
    c = json.code
  }).catch(() => {
    console.log("failed to topup")
  })

  return c

}

walletService.GetUserTransactions = async function (pagination) {
  const api = config.host + "users/transactions/_query"
  let resp = {}
  const body = {
    pagination: pagination,
  }

  await fetch(api, {
      method: 'POST',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
      body: ToJSON(body)
    }
  ).then(resp => {
    return resp.text()
  }).then(json => {
    resp = FromJSON(json)
    console.log(resp)
  }).catch((e) => {
    console.log("failed to Get User Transactions", e)
  })

  return resp
}

walletService.GetUserWallet = async function () {
  const api = config.host + "users/wallets"
  let resp = {}
  await fetch(api, {
      method: 'GET',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
    }
  ).then(resp => {
    return resp.json()
  }).then(json => {
    resp = json
  }).catch((e) => {
    console.log("failed to Get User Wallet", e)
  })

  return resp
}
export default function WalletService() {
  return walletService
}