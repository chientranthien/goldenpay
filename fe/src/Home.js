import React from "react";

export default  function Home(){
  return (
    <div className="container form-container">
      <div className="row justify-content-center">

        <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12">
          <h5>Balance</h5>
          <h2>$100</h2>
          <button className="btn btn-primary">Transfer</button>
        </div>

      </div>
      <div className="row justify-content-center">

        <div className="info col-lg-6 col-md-8 col-sm-12  col-xs-12">
          <h5>Recent Activity</h5>
          <div className="card">
            <div className="card-body">
              <h5 className="card-title">To: Tran Thien Chien</h5>
              <h6 className="card-subtitle mb-2 text-body-secondary">Oct 2023</h6>
              <p className="card-text">$100</p>
              <span className="badge text-bg-warning">Pending</span>
            </div>
          </div>
          <div className="card">
            <div className="card-body">
              <h5 className="card-title">To: Tran Thien Chien</h5>
              <p className="card-text">$100</p>
              <span className="badge text-bg-success">Success</span>
            </div>
          </div>
        </div>

      </div>
    </div>
  )
}