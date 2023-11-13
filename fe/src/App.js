import './App.css';
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";

import Login from './Login'
import Signup from './Signup'

export default function App() {
  return (
    <Router>
      <div>
        {/*<nav>*/}
        {/*  <ul>*/}
        {/*    <li>*/}
        {/*      <Link to="/">Home</Link>*/}
        {/*    </li>*/}
        {/*    <li>*/}
        {/*      <Link to="/about">About</Link>*/}
        {/*    </li>*/}
        {/*    <li>*/}
        {/*      <Link to="/users">Users</Link>*/}
        {/*    </li>*/}
        {/*  </ul>*/}
        {/*</nav>*/}

        {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
        <Switch>
          <Route path="/about">
            <About/>
          </Route>
          <Route path="/users">
            <Users/>
          </Route>
          <Route exact path="/">
            <Home/>
          </Route>
          <Route path="/login">
            <Login/>
          </Route>
          <Route path="/signup">
            <Signup/>
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

function Home() {
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
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">To: Tran Thien Chien</h5>
              <h6 class="card-subtitle mb-2 text-body-secondary">Oct 2023</h6>
              <p class="card-text">$100</p>
              <span class="badge text-bg-warning">Pending</span>
            </div>
          </div>
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">To: Tran Thien Chien</h5>
              <p class="card-text">$100</p>
              <span class="badge text-bg-success">Success</span>
            </div>
          </div>
        </div>

      </div>
    </div>
  )

}

function About() {
  return <h2>About</h2>;
}

function Users() {
  return <h2>Users</h2>;
}
