import React from "react";

export default function Activity(props) {
  let badge = "warning"
  let statusText = "Pending"
  if (props.status == 2) {
    badge = "danger"
    statusText = "Rejected"
  } else if (props.status == 3) {
    badge = "success"
    statusText = "Success"
  }

  const DATE_OPTIONS = {weekday: 'short', year: 'numeric', month: 'short', day: 'numeric'};
  let d = new Date(props.ctime).toLocaleDateString('en-US', DATE_OPTIONS)

  return (
    <div className="card">
      <div className="card-body">
        <h5 className="card-title">To: {props.name}</h5>
        <div>
          <span className={"badge text-bg-" + badge}>{statusText}</span>
          <span className="text-body-secondary">{d}</span>
        </div>
        <p className="card-text">$ {props.amount}</p>
      </div>
    </div>
  )
}