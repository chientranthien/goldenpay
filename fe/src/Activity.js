import React from "react";
import {GetUserIdFromCookie} from "./common";

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

  let card
  if (props.from.id == GetUserIdFromCookie()) {
     card = <div className="card-body">
      <h5 className="card-title">To: {props.to.name}</h5>
      <div>
        <span className={"badge text-bg-" + badge}>{statusText}</span>
        <span className="text-body-secondary">{d}</span>
      </div>
      <p className="card-text text-danger">-${props.amount}</p>
    </div>
  } else {
    card = <div className="card-body">
      <h5 className="card-title">From: {props.from.name}</h5>
      <div>
        <span className={"badge text-bg-" + badge}>{statusText}</span>
        <span className="text-body-secondary">{d}</span>
      </div>
      <p className="card-text text-success">+${props.amount}</p>
    </div>

  }
  return (
    <div className="card">
      {card}
    </div>
  )
}