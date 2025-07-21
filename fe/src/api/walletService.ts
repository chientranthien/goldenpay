import { config, FromJSON, ToJSON } from "./common"

// API Response interfaces
interface ApiResponse {
  code: Code;
  data?: any;
  message?: string;
}

interface TransferRequest {
  toEmail: string;
  amount: number;
}

interface TopupRequest {
  amount: number;
}

interface GetTransactionsRequest {
  pagination: any; // You might want to define a proper Pagination interface
}

interface Code {
  id: number;
  msg: string;
}

interface IWalletService {
  Transfer: (toEmail: string, amount: string) => Promise<Code>;
  Topup: (amount: string) => Promise<Code>;
  GetUserTransactions: (pagination: any) => Promise<any>;
  GetUserWallet: () => Promise<any>;
}

const walletService: IWalletService = {
  Transfer: async function (toEmail: string, amount: string): Promise<Code> {
    const api = config.host + "users/transactions"
    const body: TransferRequest = {
      toEmail: toEmail,
      amount: parseInt(amount)
    }

    let c: Code = { id: 0, msg: "" }
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
    }).then((json: ApiResponse) => {
      c = json.code
    }).catch(() => {
      console.log("failed to transfer")
    })

    return c
  },

  Topup: async function (amount: string): Promise<Code> {
    const api = config.host + "users/topups"
    const body: TopupRequest = {
      amount: parseInt(amount)
    }

    let c: Code = { id: 0, msg: "" }
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
    }).then((json: ApiResponse) => {
      c = json.code
    }).catch(() => {
      console.log("failed to topup")
    })

    return c
  },

  GetUserTransactions: async function (pagination: any): Promise<any> {
    const api = config.host + "users/transactions/_query"
    let resp: any = {}
    const body: GetTransactionsRequest = {
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
    }).then((json: string) => {
      resp = FromJSON(json)
    }).catch((e) => {
      console.log("failed to Get User Transactions", e)
    })

    return resp
  },

  GetUserWallet: async function (): Promise<any> {
    const api = config.host + "users/wallets"
    let resp: any = {}
    await fetch(api, {
      method: 'GET',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
    }
    ).then(resp => {
      return resp.json()
    }).then((json: any) => {
      resp = json
    }).catch((e) => {
      console.log("failed to Get User Wallet", e)
    })

    return resp
  }
}

export default function WalletService(): IWalletService {
  return walletService
}