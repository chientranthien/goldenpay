import React, { useEffect, useState } from "react";
import Activity from "./Activity";
import WalletService from "./api/walletService";

interface Transaction {
  id: string;
  from: {
    id: number;
    name: string;
  };
  to: {
    id: number;
    name: string;
  };
  amount: number;
  status: number;
  ctime: string;
}

interface TransactionResponse {
  code?: {
    id: number;
  };
  data?: {
    transactions: Transaction[];
    nextPagination: Pagination;
  };
}

interface Pagination {
  limit: number;
  hasMore: boolean;
}

interface RecentActivityProps {
  isAuthenticated: boolean | null;
}

export default function RecentActivity({ isAuthenticated }: RecentActivityProps): React.ReactElement {
  const [transactions, setTransactions] = useState<React.ReactElement[]>([])
  const [pagination, setPagination] = useState<Pagination>({ limit: 3, hasMore: false })

  function getUserTransactions(): void {
    WalletService().GetUserTransactions(pagination).then((resp: TransactionResponse) => {
      if (resp.code !== undefined && resp.code.id === 0 &&
        resp.data?.transactions !== null && resp.data?.transactions !== undefined) {
        setTransactions(prev => {
          const rs = [...prev]
          rs.push(...resp.data!.transactions.map(e => <Activity key={e.id} from={e.from} to={e.to} amount={e.amount}
            status={e.status} ctime={e.ctime} />))
          return rs
        })
        setPagination(resp.data.nextPagination)
      }
    }
    )
  }

  useEffect(() => {
    // Only make API call if authenticated
    if (isAuthenticated === true) {
      getUserTransactions()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated])

  function handleLoadMore(): void {
    getUserTransactions()
  }

  // Don't render anything if not authenticated
  if (isAuthenticated !== true) {
    return <></>
  }

  return (
    <>
      {
        transactions.length > 0 && <div className="row justify-content-center">

          <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12">
            <h5>Recent Activity</h5>
            {transactions}
            {pagination.hasMore === true &&
              <button className={"btn btn-primary"} onClick={handleLoadMore}>Load More</button>}
          </div>

        </div>
      }
    </>
  )
}