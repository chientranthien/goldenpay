import Activity from "./Activity";
import WalletService from "./api/walletService";
import {useEffect, useState} from "react";

export default function RecentActivity() {
  const [transactions, setTransactions] = useState([])
  const [pagination, setPagination] = useState({limit: 3, hasMore: false})

  function getUserTransactions() {
    WalletService().GetUserTransactions(pagination).then(resp => {
        if (resp.code != undefined && resp.code.id == 0 &&
          resp.data.transactions != null && resp.data.transactions != undefined) {
          console.log(resp.data)
          setTransactions(prev => {
            const rs = [...prev]
            rs.push(...resp.data.transactions.map(e => <Activity key={e.id} from={e.from} to={e.to} amount={e.amount}
                                                                 status={e.status} ctime={e.ctime}/>))
            return rs
          })
          setPagination(prev => resp.data.nextPagination)
        }
      }
    )
  }

  useEffect(() => {
    getUserTransactions()
  }, [])

  function handleLoadMore() {
    getUserTransactions()
  }

  return (
    <>
      {
        transactions.length > 0 && <div className="row justify-content-center">

          <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12">
            <h5>Recent Activity</h5>
            {transactions}
            {pagination.hasMore == true &&
            <button className={"btn btn-primary"} onClick={handleLoadMore}>Load More</button>}
          </div>

        </div>
      }
    </>
  )
}